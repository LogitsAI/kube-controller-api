// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: kube_controller_api/manager.proto

package controllerpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ControllerManagerClient is the client API for ControllerManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControllerManagerClient interface {
	// CreateManager creates a Controller Manager with the given config.
	CreateManager(ctx context.Context, in *CreateManagerRequest, opts ...grpc.CallOption) (*CreateManagerResponse, error)
	// ReconcileLoop returns the next object to be processed by the controller.
	ReconcileLoop(ctx context.Context, opts ...grpc.CallOption) (ControllerManager_ReconcileLoopClient, error)
}

type controllerManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewControllerManagerClient(cc grpc.ClientConnInterface) ControllerManagerClient {
	return &controllerManagerClient{cc}
}

func (c *controllerManagerClient) CreateManager(ctx context.Context, in *CreateManagerRequest, opts ...grpc.CallOption) (*CreateManagerResponse, error) {
	out := new(CreateManagerResponse)
	err := c.cc.Invoke(ctx, "/kube_controller_api.ControllerManager/CreateManager", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerManagerClient) ReconcileLoop(ctx context.Context, opts ...grpc.CallOption) (ControllerManager_ReconcileLoopClient, error) {
	stream, err := c.cc.NewStream(ctx, &ControllerManager_ServiceDesc.Streams[0], "/kube_controller_api.ControllerManager/ReconcileLoop", opts...)
	if err != nil {
		return nil, err
	}
	x := &controllerManagerReconcileLoopClient{stream}
	return x, nil
}

type ControllerManager_ReconcileLoopClient interface {
	Send(*ReconcileLoopRequest) error
	Recv() (*ReconcileLoopResponse, error)
	grpc.ClientStream
}

type controllerManagerReconcileLoopClient struct {
	grpc.ClientStream
}

func (x *controllerManagerReconcileLoopClient) Send(m *ReconcileLoopRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *controllerManagerReconcileLoopClient) Recv() (*ReconcileLoopResponse, error) {
	m := new(ReconcileLoopResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ControllerManagerServer is the server API for ControllerManager service.
// All implementations must embed UnimplementedControllerManagerServer
// for forward compatibility
type ControllerManagerServer interface {
	// CreateManager creates a Controller Manager with the given config.
	CreateManager(context.Context, *CreateManagerRequest) (*CreateManagerResponse, error)
	// ReconcileLoop returns the next object to be processed by the controller.
	ReconcileLoop(ControllerManager_ReconcileLoopServer) error
	mustEmbedUnimplementedControllerManagerServer()
}

// UnimplementedControllerManagerServer must be embedded to have forward compatible implementations.
type UnimplementedControllerManagerServer struct {
}

func (UnimplementedControllerManagerServer) CreateManager(context.Context, *CreateManagerRequest) (*CreateManagerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateManager not implemented")
}
func (UnimplementedControllerManagerServer) ReconcileLoop(ControllerManager_ReconcileLoopServer) error {
	return status.Errorf(codes.Unimplemented, "method ReconcileLoop not implemented")
}
func (UnimplementedControllerManagerServer) mustEmbedUnimplementedControllerManagerServer() {}

// UnsafeControllerManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControllerManagerServer will
// result in compilation errors.
type UnsafeControllerManagerServer interface {
	mustEmbedUnimplementedControllerManagerServer()
}

func RegisterControllerManagerServer(s grpc.ServiceRegistrar, srv ControllerManagerServer) {
	s.RegisterService(&ControllerManager_ServiceDesc, srv)
}

func _ControllerManager_CreateManager_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateManagerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerManagerServer).CreateManager(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kube_controller_api.ControllerManager/CreateManager",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerManagerServer).CreateManager(ctx, req.(*CreateManagerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControllerManager_ReconcileLoop_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ControllerManagerServer).ReconcileLoop(&controllerManagerReconcileLoopServer{stream})
}

type ControllerManager_ReconcileLoopServer interface {
	Send(*ReconcileLoopResponse) error
	Recv() (*ReconcileLoopRequest, error)
	grpc.ServerStream
}

type controllerManagerReconcileLoopServer struct {
	grpc.ServerStream
}

func (x *controllerManagerReconcileLoopServer) Send(m *ReconcileLoopResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *controllerManagerReconcileLoopServer) Recv() (*ReconcileLoopRequest, error) {
	m := new(ReconcileLoopRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ControllerManager_ServiceDesc is the grpc.ServiceDesc for ControllerManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ControllerManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kube_controller_api.ControllerManager",
	HandlerType: (*ControllerManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateManager",
			Handler:    _ControllerManager_CreateManager_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReconcileLoop",
			Handler:       _ControllerManager_ReconcileLoop_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "kube_controller_api/manager.proto",
}