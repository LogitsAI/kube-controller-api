# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: kube_controller_api/reconciler.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from kube_controller_api import controller_pb2 as kube__controller__api_dot_controller__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n$kube_controller_api/reconciler.proto\x12\x13kube_controller_api\x1a\x1egoogle/protobuf/duration.proto\x1a$kube_controller_api/controller.proto\"\x92\x01\n\x0fReconcileResult\x12\x12\n\x05\x65rror\x18\x01 \x01(\tH\x00\x88\x01\x01\x12\x0f\n\x07requeue\x18\x02 \x01(\x08\x12\x30\n\rrequeue_after\x18\x03 \x01(\x0b\x32\x19.google.protobuf.Duration\x12\x13\n\x06status\x18\x04 \x01(\x0cH\x01\x88\x01\x01\x42\x08\n\x06_errorB\t\n\x07_status\"$\n\tWorkQueue\x12\x17\n\x0f\x63ontroller_name\x18\x01 \x01(\t\"b\n\x0c\x43hildObjects\x12\x41\n\x12group_version_kind\x18\x01 \x01(\x0b\x32%.kube_controller_api.GroupVersionKind\x12\x0f\n\x07objects\x18\x02 \x03(\x0c\x42\x36Z4github.com/LogitsAI/kube-controller-api/controllerpbb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'kube_controller_api.reconciler_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z4github.com/LogitsAI/kube-controller-api/controllerpb'
  _globals['_RECONCILERESULT']._serialized_start=132
  _globals['_RECONCILERESULT']._serialized_end=278
  _globals['_WORKQUEUE']._serialized_start=280
  _globals['_WORKQUEUE']._serialized_end=316
  _globals['_CHILDOBJECTS']._serialized_start=318
  _globals['_CHILDOBJECTS']._serialized_end=416
# @@protoc_insertion_point(module_scope)
