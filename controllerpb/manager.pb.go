// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: kube_controller_api/manager.proto

package controllerpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Controllers []*ControllerConfig `protobuf:"bytes,1,rep,name=controllers,proto3" json:"controllers,omitempty"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_controller_api_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kube_controller_api_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file_kube_controller_api_manager_proto_rawDescGZIP(), []int{0}
}

func (x *StartRequest) GetControllers() []*ControllerConfig {
	if x != nil {
		return x.Controllers
	}
	return nil
}

type StartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartResponse) Reset() {
	*x = StartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_controller_api_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartResponse) ProtoMessage() {}

func (x *StartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kube_controller_api_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartResponse.ProtoReflect.Descriptor instead.
func (*StartResponse) Descriptor() ([]byte, []int) {
	return file_kube_controller_api_manager_proto_rawDescGZIP(), []int{1}
}

// ReconcileLoop() is a bidirectional streaming call.
// Each request can contain one of several types of sub-messages.
type ReconcileLoopRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Msg:
	//
	//	*ReconcileLoopRequest_Subscribe
	//	*ReconcileLoopRequest_Acknowledge
	Msg isReconcileLoopRequest_Msg `protobuf_oneof:"msg"`
}

func (x *ReconcileLoopRequest) Reset() {
	*x = ReconcileLoopRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_controller_api_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconcileLoopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconcileLoopRequest) ProtoMessage() {}

func (x *ReconcileLoopRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kube_controller_api_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconcileLoopRequest.ProtoReflect.Descriptor instead.
func (*ReconcileLoopRequest) Descriptor() ([]byte, []int) {
	return file_kube_controller_api_manager_proto_rawDescGZIP(), []int{2}
}

func (m *ReconcileLoopRequest) GetMsg() isReconcileLoopRequest_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *ReconcileLoopRequest) GetSubscribe() *WorkQueue {
	if x, ok := x.GetMsg().(*ReconcileLoopRequest_Subscribe); ok {
		return x.Subscribe
	}
	return nil
}

func (x *ReconcileLoopRequest) GetAcknowledge() *ReconcileResult {
	if x, ok := x.GetMsg().(*ReconcileLoopRequest_Acknowledge); ok {
		return x.Acknowledge
	}
	return nil
}

type isReconcileLoopRequest_Msg interface {
	isReconcileLoopRequest_Msg()
}

type ReconcileLoopRequest_Subscribe struct {
	// Subscribe should be sent once at the beginning of the stream,
	// to indicate what controller the client is interested in.
	Subscribe *WorkQueue `protobuf:"bytes,1,opt,name=subscribe,proto3,oneof"`
}

type ReconcileLoopRequest_Acknowledge struct {
	// Acknowledge should be sent after each object is processed,
	// whether it was successful or not. The server will not send
	// another object until the last outstanding response has been
	// acknowledged.
	Acknowledge *ReconcileResult `protobuf:"bytes,2,opt,name=acknowledge,proto3,oneof"`
}

func (*ReconcileLoopRequest_Subscribe) isReconcileLoopRequest_Msg() {}

func (*ReconcileLoopRequest_Acknowledge) isReconcileLoopRequest_Msg() {}

// ReconcileLoopResponse contains a single object to be reconciled.
type ReconcileLoopResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Object is a JSON-encoded Kubernetes object.
	Object []byte `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
}

func (x *ReconcileLoopResponse) Reset() {
	*x = ReconcileLoopResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_controller_api_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconcileLoopResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconcileLoopResponse) ProtoMessage() {}

func (x *ReconcileLoopResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kube_controller_api_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconcileLoopResponse.ProtoReflect.Descriptor instead.
func (*ReconcileLoopResponse) Descriptor() ([]byte, []int) {
	return file_kube_controller_api_manager_proto_rawDescGZIP(), []int{3}
}

func (x *ReconcileLoopResponse) GetObject() []byte {
	if x != nil {
		return x.Object
	}
	return nil
}

var File_kube_controller_api_manager_proto protoreflect.FileDescriptor

var file_kube_controller_api_manager_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65,
	0x72, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x13, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x1a, 0x24, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24,
	0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f,
	0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x57, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x47, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x2e,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x73, 0x22, 0x0f, 0x0a,
	0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xa7,
	0x01, 0x0a, 0x14, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x4c, 0x6f, 0x6f, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3e, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6b, 0x75, 0x62,
	0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69,
	0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x51, 0x75, 0x65, 0x75, 0x65, 0x48, 0x00, 0x52, 0x09, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x48, 0x0a, 0x0b, 0x61, 0x63, 0x6b, 0x6e, 0x6f,
	0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6b,
	0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x61, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x2f, 0x0a, 0x15, 0x52, 0x65, 0x63, 0x6f,
	0x6e, 0x63, 0x69, 0x6c, 0x65, 0x4c, 0x6f, 0x6f, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x32, 0xd3, 0x01, 0x0a, 0x11, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12,
	0x50, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x21, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6b, 0x75,
	0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x6c, 0x0a, 0x0d, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x4c, 0x6f,
	0x6f, 0x70, 0x12, 0x29, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69,
	0x6c, 0x65, 0x4c, 0x6f, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e,
	0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x63, 0x69, 0x6c, 0x65, 0x4c, 0x6f, 0x6f,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42,
	0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4c, 0x6f,
	0x67, 0x69, 0x74, 0x73, 0x41, 0x49, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x2d, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kube_controller_api_manager_proto_rawDescOnce sync.Once
	file_kube_controller_api_manager_proto_rawDescData = file_kube_controller_api_manager_proto_rawDesc
)

func file_kube_controller_api_manager_proto_rawDescGZIP() []byte {
	file_kube_controller_api_manager_proto_rawDescOnce.Do(func() {
		file_kube_controller_api_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_kube_controller_api_manager_proto_rawDescData)
	})
	return file_kube_controller_api_manager_proto_rawDescData
}

var file_kube_controller_api_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_kube_controller_api_manager_proto_goTypes = []interface{}{
	(*StartRequest)(nil),          // 0: kube_controller_api.StartRequest
	(*StartResponse)(nil),         // 1: kube_controller_api.StartResponse
	(*ReconcileLoopRequest)(nil),  // 2: kube_controller_api.ReconcileLoopRequest
	(*ReconcileLoopResponse)(nil), // 3: kube_controller_api.ReconcileLoopResponse
	(*ControllerConfig)(nil),      // 4: kube_controller_api.ControllerConfig
	(*WorkQueue)(nil),             // 5: kube_controller_api.WorkQueue
	(*ReconcileResult)(nil),       // 6: kube_controller_api.ReconcileResult
}
var file_kube_controller_api_manager_proto_depIdxs = []int32{
	4, // 0: kube_controller_api.StartRequest.controllers:type_name -> kube_controller_api.ControllerConfig
	5, // 1: kube_controller_api.ReconcileLoopRequest.subscribe:type_name -> kube_controller_api.WorkQueue
	6, // 2: kube_controller_api.ReconcileLoopRequest.acknowledge:type_name -> kube_controller_api.ReconcileResult
	0, // 3: kube_controller_api.ControllerManager.Start:input_type -> kube_controller_api.StartRequest
	2, // 4: kube_controller_api.ControllerManager.ReconcileLoop:input_type -> kube_controller_api.ReconcileLoopRequest
	1, // 5: kube_controller_api.ControllerManager.Start:output_type -> kube_controller_api.StartResponse
	3, // 6: kube_controller_api.ControllerManager.ReconcileLoop:output_type -> kube_controller_api.ReconcileLoopResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_kube_controller_api_manager_proto_init() }
func file_kube_controller_api_manager_proto_init() {
	if File_kube_controller_api_manager_proto != nil {
		return
	}
	file_kube_controller_api_controller_proto_init()
	file_kube_controller_api_reconciler_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_kube_controller_api_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_controller_api_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_controller_api_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconcileLoopRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_controller_api_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconcileLoopResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_kube_controller_api_manager_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*ReconcileLoopRequest_Subscribe)(nil),
		(*ReconcileLoopRequest_Acknowledge)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kube_controller_api_manager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kube_controller_api_manager_proto_goTypes,
		DependencyIndexes: file_kube_controller_api_manager_proto_depIdxs,
		MessageInfos:      file_kube_controller_api_manager_proto_msgTypes,
	}.Build()
	File_kube_controller_api_manager_proto = out.File
	file_kube_controller_api_manager_proto_rawDesc = nil
	file_kube_controller_api_manager_proto_goTypes = nil
	file_kube_controller_api_manager_proto_depIdxs = nil
}
