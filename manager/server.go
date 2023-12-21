package manager

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/LogitsAI/kube-controller-api/controllerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

type controllerEntry struct {
	name       string
	reconciler *reconcilerAdapter
}

type ControllerManagerServer struct {
	controllerpb.UnimplementedControllerManagerServer

	kubeConfig    *rest.Config
	signalHandler context.Context

	mu            sync.RWMutex
	manager       manager.Manager
	managerConfig *controllerpb.StartRequest
	controllers   map[string]controllerEntry
}

func NewControllerManagerServer(kubeConfig *rest.Config) *ControllerManagerServer {
	return &ControllerManagerServer{
		kubeConfig:    kubeConfig,
		signalHandler: signals.SetupSignalHandler(),
	}
}

func (s *ControllerManagerServer) getController(name string) (controllerEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	controller, ok := s.controllers[name]
	return controller, ok
}

func (s *ControllerManagerServer) Start(ctx context.Context, in *controllerpb.StartRequest) (*controllerpb.StartResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.manager != nil {
		// The manager was already started. Check if the config is the same.
		if !proto.Equal(in, s.managerConfig) {
			return nil, status.Errorf(codes.AlreadyExists, "manager already started with different config")
		}
		return &controllerpb.StartResponse{}, nil
	}

	mgr, err := manager.New(s.kubeConfig, manager.Options{})
	if err != nil {
		return nil, err
	}
	s.manager = mgr
	s.managerConfig = in

	s.controllers = map[string]controllerEntry{}
	childGVKs := []schema.GroupVersionKind{}
	for _, controller := range in.Controllers {
		bldr := builder.ControllerManagedBy(mgr)

		if controller.Parent != nil {
			obj := &unstructured.Unstructured{}
			obj.SetGroupVersionKind(controller.Parent.GroupVersionKind())
			bldr = bldr.For(obj)
		}
		for _, child := range controller.Children {
			childGVK := child.GroupVersionKind()
			childGVKs = append(childGVKs, childGVK)

			obj := &unstructured.Unstructured{}
			obj.SetGroupVersionKind(childGVK)
			bldr = bldr.Owns(obj)
		}

		reconciler := newReconcilerAdapter(
			mgr.GetClient(),
			controller.Name,
			controller.Parent.GroupVersionKind(),
			childGVKs,
		)
		if err := bldr.Complete(reconciler); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to create controller %q: %v", controller.Name, err)
		}
		s.controllers[controller.Name] = controllerEntry{
			name:       controller.Name,
			reconciler: reconciler,
		}
		slog.Debug("Created controller", "name", controller.Name)
	}

	slog.Info("Starting manager...")
	go func() {
		if err := mgr.Start(s.signalHandler); err != nil {
			slog.Error("failed to start manager", "error", err)
		}
	}()

	return &controllerpb.StartResponse{}, nil
}

func (s *ControllerManagerServer) ReconcileLoop(stream controllerpb.ControllerManager_ReconcileLoopServer) error {
	// Look for the first message, which should be a Subscribe.
	req, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	subMsg, ok := req.Msg.(*controllerpb.ReconcileLoopRequest_Subscribe)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "first message should be a subscribe message")
	}
	if s.manager == nil {
		return status.Errorf(codes.FailedPrecondition, "manager has not been started")
	}
	// Find the specified controller within that manager.
	controller, ok := s.getController(subMsg.Subscribe.GetControllerName())
	if !ok {
		return status.Errorf(codes.NotFound, "controller not found")
	}

	var request *reconcileRequest

	// If we return for any reason while a request is in-flight, return an error.
	defer func() {
		if request != nil {
			request.result <- &controllerpb.ReconcileResult{
				Error: pointer.String("ReconcileLoop() stream ended unexpectedly"),
			}
		}
	}()

	for {
		// Pick the next object from the work queue.
		select {
		case req := <-controller.reconciler.requests:
			request = &req
		case <-stream.Context().Done():
			return stream.Context().Err()
		}

		// Push it down to the worker by returning it as a streaming response.
		parentJSON, err := json.Marshal(request.parent)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal object: %v", err)
		}
		children, err := childObjectsProto(request.children)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal child objects: %v", err)
		}
		resp := &controllerpb.ReconcileLoopResponse{
			Parent:   parentJSON,
			Children: children,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}

		// Wait for an acknowledgement of success or failure.
		// TODO: Impose a timeout on the worker responding.
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		ackMsg, ok := req.Msg.(*controllerpb.ReconcileLoopRequest_Acknowledge)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "expected acknowledge message")
		}

		// Send the result back to the reconcilerAdapter.
		request.result <- ackMsg.Acknowledge
		// We no longer have a request in-flight.
		request = nil
	}
}

func childObjectsProto(children map[schema.GroupVersionKind][]unstructured.Unstructured) ([]*controllerpb.ObservedChildObjects, error) {
	var childrenProto []*controllerpb.ObservedChildObjects
	for gvk, objs := range children {
		observed := make(map[string][]byte, len(objs))
		for _, obj := range objs {
			objJSON, err := json.Marshal(obj)
			if err != nil {
				return nil, err
			}
			key := obj.GetAnnotations()[childKeyAnnotation]
			if key == "" {
				return nil, fmt.Errorf("child %s %s/%s is missing child key annotation (%q)", obj.GetKind(), obj.GetNamespace(), obj.GetName(), childKeyAnnotation)
			}
			observed[key] = objJSON
		}
		childrenProto = append(childrenProto, &controllerpb.ObservedChildObjects{
			GroupVersionKind: &controllerpb.GroupVersionKind{
				Group:   gvk.Group,
				Version: gvk.Version,
				Kind:    gvk.Kind,
			},
			ObservedObjects: observed,
		})
	}
	return childrenProto, nil
}
