package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/LogitsAI/kube-controller-api/controllerpb"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type reconcileRequest struct {
	ctx      context.Context
	parent   *unstructured.Unstructured
	children map[schema.GroupVersionKind][]unstructured.Unstructured
	result   chan *controllerpb.ReconcileResult
}

func newReconcileRequest(ctx context.Context, parent *unstructured.Unstructured, children map[schema.GroupVersionKind][]unstructured.Unstructured) reconcileRequest {
	return reconcileRequest{
		ctx:      ctx,
		parent:   parent,
		children: children,
		// Use a buffered channel so we don't need to worry about blocking the sender
		// in case the receiver has stopped listening.
		result: make(chan *controllerpb.ReconcileResult, 1),
	}
}

type reconcilerAdapter struct {
	slog       *slog.Logger
	parentGVK  schema.GroupVersionKind
	childGVKs  []schema.GroupVersionKind
	kubeClient client.Client
	requests   chan reconcileRequest
}

func newReconcilerAdapter(kubeClient client.Client, controllerName string, parentGVK schema.GroupVersionKind, childGVKs []schema.GroupVersionKind) *reconcilerAdapter {
	return &reconcilerAdapter{
		slog:       slog.With("parentGVK", parentGVK, "controller", controllerName),
		parentGVK:  parentGVK,
		childGVKs:  childGVKs,
		kubeClient: kubeClient,
		requests:   make(chan reconcileRequest),
	}
}

func (r *reconcilerAdapter) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	slog := r.slog.With("name", req.NamespacedName)

	// Try to fetch the object.
	parent := &unstructured.Unstructured{}
	parent.SetGroupVersionKind(r.parentGVK)
	if err := r.kubeClient.Get(ctx, req.NamespacedName, parent); err != nil {
		if k8serrors.IsNotFound(err) {
			slog.DebugContext(ctx, "Ignoring reconcile request for object that no longer exists")
			return reconcile.Result{}, nil
		}
		slog.ErrorContext(ctx, "Failed to get object", "error", err)
		return reconcile.Result{}, err
	}

	// Fetch the child objects.
	children := map[schema.GroupVersionKind][]unstructured.Unstructured{}
	for _, childGVK := range r.childGVKs {
		listGVK := childGVK
		listGVK.Kind = listGVK.Kind + "List"

		childList := &unstructured.UnstructuredList{}
		childList.SetGroupVersionKind(listGVK)
		// TODO: Consider using a label selector.
		if err := r.kubeClient.List(ctx, childList, client.InNamespace(parent.GetNamespace())); err != nil {
			slog.ErrorContext(ctx, "Failed to list child objects", "error", err)
			return reconcile.Result{}, err
		}

		// Filter to only include children that are controlled by the parent.
		var controlled []unstructured.Unstructured
		for _, child := range childList.Items {
			if metav1.IsControlledBy(&child, parent) {
				controlled = append(controlled, child)
			}
		}
		children[childGVK] = controlled
	}

	// Send a request and wait for a worker to accept it.
	slog.DebugContext(ctx, "Sending reconcile request")
	request := newReconcileRequest(ctx, parent, children)

	select {
	case r.requests <- request:
	case <-ctx.Done():
		slog.ErrorContext(ctx, "Context was cancelled while waiting for reconcile request to be accepted by a worker", "error", ctx.Err())
		return reconcile.Result{}, ctx.Err()
	}

	// Wait for the result to be sent back.
	slog.DebugContext(ctx, "Waiting for reconcile result")
	select {
	case <-ctx.Done():
		slog.ErrorContext(ctx, "Context was cancelled while waiting for reconcile result", "error", ctx.Err())
		return reconcile.Result{}, ctx.Err()
	case result := <-request.result:
		// Set status if one was provided. We do this before checking the error
		// because the controller may still want to update status.
		if result.Status != nil {
			newStatus := map[string]any{}
			if err := json.Unmarshal(result.Status, &newStatus); err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to unmarshal status in ReconcileResult: %v", err)
			}

			// Only send a status update if something changed.
			if !equality.Semantic.DeepEqual(parent.Object["status"], newStatus) {
				parent.Object["status"] = newStatus

				if err := r.kubeClient.Status().Update(ctx, parent); err != nil {
					// We ignore conflict (optimistic concurrency) errors on status
					// updates because they are common during normal operation.
					// We will retry the update on the next reconcile.
					if k8serrors.IsConflict(err) {
						slog.DebugContext(ctx, "Ignoring conflict error while updating status", "error", err)
					} else {
						return reconcile.Result{}, fmt.Errorf("failed to update status: %v", err)
					}
				}
				slog.DebugContext(ctx, "Updated status")
			}
		}

		if result.Error != nil {
			slog.DebugContext(ctx, "Reconcile failed", "error", *result.Error)
			return reconcile.Result{}, errors.New(*result.Error)
		}
		slog.DebugContext(ctx, "Reconcile succeeded")
		return result.ReconcileResult(), nil
	}
}
