package manager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/LogitsAI/kube-controller-api/controllerpb"
	"github.com/LogitsAI/kube-controller-api/names"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// childKeyAnnotation is the annotation in which we store the map key that
	// the user-written reconciler uses to uniquely identify a given child
	// object within the context of its parent.
	childKeyAnnotation = "kube-controller-api.logits.ai/child-key"
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
	slog           *slog.Logger
	controllerName string
	parentGVK      schema.GroupVersionKind
	childGVKs      []schema.GroupVersionKind
	kubeClient     client.Client
	eventRecorder  record.EventRecorder
	dynClient      *DynamicClient
	requests       chan reconcileRequest
}

func newReconcilerAdapter(mgr manager.Manager, dynClient *DynamicClient, controllerName string, parentGVK schema.GroupVersionKind, childGVKs []schema.GroupVersionKind) *reconcilerAdapter {
	return &reconcilerAdapter{
		slog:           slog.With("parentGVK", parentGVK, "controller", controllerName),
		controllerName: controllerName,
		parentGVK:      parentGVK,
		childGVKs:      childGVKs,
		kubeClient:     mgr.GetClient(),
		eventRecorder:  mgr.GetEventRecorderFor(controllerName),
		dynClient:      dynClient,
		requests:       make(chan reconcileRequest),
	}
}

func (r *reconcilerAdapter) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	slog := r.slog.With("parentName", req.NamespacedName)

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
	observed := map[schema.GroupVersionKind][]unstructured.Unstructured{}
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
		observed[childGVK] = controlled
	}

	// Send a request and wait for a worker to accept it.
	slog.DebugContext(ctx, "Sending reconcile request")
	request := newReconcileRequest(ctx, parent, observed)

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
			newStatus["observedGeneration"] = parent.GetGeneration()

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

		// Decode apply configs for desired children.
		applyConfigs := map[schema.GroupVersionKind]map[string]unstructured.Unstructured{}
		for _, desired := range result.Children {
			gvk := desired.GroupVersionKind.GroupVersionKind()
			configs := make(map[string]unstructured.Unstructured, len(desired.ApplyConfigs))
			for childKey, configJSON := range desired.ApplyConfigs {
				out := unstructured.Unstructured{}
				if err := json.Unmarshal(configJSON, &out); err != nil {
					return reconcile.Result{}, fmt.Errorf("failed to unmarshal apply config for %v: %v", gvk, err)
				}
				out.SetGroupVersionKind(gvk)

				// Generate a unique, deterministic name if one was not provided.
				if out.GetNamespace() == "" {
					out.SetNamespace(parent.GetNamespace())
				}
				if out.GetName() == "" {
					setDefaultChildName(&out, parent, childKey)
				}

				// Set the controller reference.
				ref := metav1.OwnerReference{
					APIVersion:         r.parentGVK.GroupVersion().String(),
					Kind:               r.parentGVK.Kind,
					Name:               parent.GetName(),
					UID:                parent.GetUID(),
					BlockOwnerDeletion: pointer.Bool(true),
					Controller:         pointer.Bool(true),
				}
				out.SetOwnerReferences([]metav1.OwnerReference{ref})

				// Add the child key annotation.
				annotations := out.GetAnnotations()
				if annotations == nil {
					annotations = map[string]string{}
				}
				annotations[childKeyAnnotation] = childKey
				out.SetAnnotations(annotations)

				configs[childKey] = out
			}
			applyConfigs[gvk] = configs
		}

		// Reconcile children.
		for _, childGVK := range r.childGVKs {
			if err := r.reconcileChildren(ctx, slog, parent, childGVK, observed[childGVK], applyConfigs[childGVK]); err != nil {
				return reconcile.Result{}, err
			}
		}

		slog.DebugContext(ctx, "Reconcile succeeded")
		return result.ReconcileResult(), nil
	}
}

func (r *reconcilerAdapter) reconcileChildren(ctx context.Context, slog *slog.Logger, parent *unstructured.Unstructured, childGVK schema.GroupVersionKind, observed []unstructured.Unstructured, desired map[string]unstructured.Unstructured) error {
	slog = slog.With("childGVK", childGVK)

	childClient := r.dynClient.Kind(childGVK)

	// Call Server-Side Apply for each desired child.
	// TODO: Use the controller-runtime client for this if they ever add support for SSA.
	for _, des := range desired {
		slog := slog.With("childName", fmt.Sprintf("%s/%s", des.GetNamespace(), des.GetName()))

		slog.DebugContext(ctx, "Applying child")
		opts := metav1.ApplyOptions{Force: true, FieldManager: r.controllerName}
		if _, err := childClient.Namespace(des.GetNamespace()).Apply(ctx, des.GetName(), &des, opts); err != nil {
			return fmt.Errorf("failed to apply child %v %v/%v: %v", des.GetKind(), des.GetNamespace(), des.GetName(), err)
		}
	}

	// Delete any observed children that are no longer desired.
	for _, obs := range observed {
		slog := slog.With("childName", fmt.Sprintf("%s/%s", obs.GetNamespace(), obs.GetName()))

		childKey := obs.GetAnnotations()[childKeyAnnotation]
		if childKey == "" {
			slog.DebugContext(ctx, "Ignoring child with missing child key annotation")
			continue
		}

		if _, ok := desired[childKey]; !ok {
			slog.DebugContext(ctx, "Deleting child")

			// Only delete the child if the UID and ResourceVersion are the same
			// as what we last observed. Otherwise, we should let the controller
			// re-reconcile to decide if it should still be deleted.
			uid := obs.GetUID()
			rv := obs.GetResourceVersion()
			cond := &client.Preconditions{UID: &uid, ResourceVersion: &rv}

			if err := r.kubeClient.Delete(ctx, &obs, cond); err != nil {
				return fmt.Errorf("failed to delete child %v %v/%v: %v", obs.GetKind(), obs.GetNamespace(), obs.GetName(), err)
			}
			r.eventRecorder.Eventf(parent, "Normal", "DeletedChild", "Deleted child %s %s/%s", childGVK.Kind, obs.GetNamespace(), obs.GetName())
		}
	}

	return nil
}

func setDefaultChildName(child *unstructured.Unstructured, parent *unstructured.Unstructured, childKey string) {
	parentGVK := parent.GetObjectKind().GroupVersionKind()
	if parentGVK.Kind == "" {
		panic("parent does not have a Kind")
	}
	childGVK := child.GetObjectKind().GroupVersionKind()
	if childGVK.Kind == "" {
		panic("child does not have a Kind")
	}

	cons := names.DefaultConstraints
	if childGVK.Group == "" && childGVK.Kind == "Service" {
		cons = names.ServiceConstraints
	}

	// WARNING: DO NOT change the contents of either the visible or hidden parts
	// arrays. Doing so would break determinism for controllers when the
	// kube-controller-api server is upgraded, causing objects to be
	// unintentionally deleted and recreated with new names.
	visibleParts := []string{parent.GetName(), childKey}
	hiddenParts := []string{parentGVK.Group, parentGVK.Kind}
	child.SetName(names.JoinSalt(cons, hiddenParts, visibleParts...))
}
