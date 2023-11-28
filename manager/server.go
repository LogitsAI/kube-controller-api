package manager

import (
	"context"
	"io"
	"log/slog"
	"sync"

	"github.com/LogitsAI/kube-controller-api/controllerpb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

type managerEntry struct {
	id      string
	name    string
	manager manager.Manager

	controllers map[string]controllerEntry
}

type controllerEntry struct {
	name       string
	reconciler *reconcilerAdapter
}

type ControllerManagerServer struct {
	controllerpb.UnimplementedControllerManagerServer

	kubeConfig    *rest.Config
	signalHandler context.Context

	mu       sync.RWMutex
	managers map[string]managerEntry
}

func NewControllerManagerServer(kubeConfig *rest.Config) *ControllerManagerServer {
	return &ControllerManagerServer{
		kubeConfig:    kubeConfig,
		signalHandler: signals.SetupSignalHandler(),

		managers: make(map[string]managerEntry),
	}
}

func (s *ControllerManagerServer) addManager(name string, mgr manager.Manager, controllers map[string]controllerEntry) string {
	id := uuid.New().String()

	s.mu.Lock()
	s.managers[id] = managerEntry{
		id:          id,
		name:        name,
		manager:     mgr,
		controllers: controllers,
	}
	s.mu.Unlock()

	return id
}

func (s *ControllerManagerServer) CreateManager(ctx context.Context, in *controllerpb.CreateManagerRequest) (*controllerpb.CreateManagerResponse, error) {
	mgr, err := manager.New(s.kubeConfig, manager.Options{})
	if err != nil {
		return nil, err
	}

	controllers := map[string]controllerEntry{}
	for _, controller := range in.Controllers {
		bldr := builder.ControllerManagedBy(mgr)

		if controller.Parent != nil {
			obj := &unstructured.Unstructured{}
			obj.SetGroupVersionKind(controller.Parent.GroupVersionKind())
			bldr = bldr.For(obj)
		}
		for _, child := range controller.Children {
			obj := &unstructured.Unstructured{}
			obj.SetGroupVersionKind(child.GroupVersionKind())
			bldr = bldr.Owns(obj)
		}

		reconciler := newReconcilerAdapter(mgr.GetClient(), controller.Parent.GroupVersionKind(), controller.Name)
		if err := bldr.Complete(reconciler); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to create controller %q: %v", controller.Name, err)
		}
		controllers[controller.Name] = controllerEntry{
			name:       controller.Name,
			reconciler: reconciler,
		}
		slog.Debug("Created controller", "name", controller.Name)
	}

	go func() {
		if err := mgr.Start(s.signalHandler); err != nil {
			slog.Error("failed to start manager", "name", in.Name, "error", err)
		}
	}()

	slog.Info("Started manager", "name", in.Name)
	id := s.addManager(in.Name, mgr, controllers)

	return &controllerpb.CreateManagerResponse{
		Id: id,
	}, nil
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
	// Find the specified manager.
	mgrEntry, ok := s.managers[subMsg.Subscribe.GetManagerId()]
	if !ok {
		return status.Errorf(codes.NotFound, "manager not found")
	}
	// Find the specified controller within that manager.
	controller, ok := mgrEntry.controllers[subMsg.Subscribe.GetControllerName()]
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
		data, err := json.Marshal(request.object)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal object: %v", err)
		}
		resp := &controllerpb.ReconcileLoopResponse{
			Object: data,
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
