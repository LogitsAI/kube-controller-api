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


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n$kube_controller_api/reconciler.proto\x12\x13kube_controller_api\x1a\x1egoogle/protobuf/duration.proto\x1a$kube_controller_api/controller.proto\"\xce\x01\n\x0fReconcileResult\x12\x12\n\x05\x65rror\x18\x01 \x01(\tH\x00\x88\x01\x01\x12\x0f\n\x07requeue\x18\x02 \x01(\x08\x12\x30\n\rrequeue_after\x18\x03 \x01(\x0b\x32\x19.google.protobuf.Duration\x12\x13\n\x06status\x18\x04 \x01(\x0cH\x01\x88\x01\x01\x12:\n\x08\x63hildren\x18\x05 \x03(\x0b\x32(.kube_controller_api.DesiredChildObjectsB\x08\n\x06_errorB\t\n\x07_status\"$\n\tWorkQueue\x12\x17\n\x0f\x63ontroller_name\x18\x01 \x01(\t\"\xe0\x01\n\x13\x44\x65siredChildObjects\x12\x41\n\x12group_version_kind\x18\x01 \x01(\x0b\x32%.kube_controller_api.GroupVersionKind\x12Q\n\rapply_configs\x18\x02 \x03(\x0b\x32:.kube_controller_api.DesiredChildObjects.ApplyConfigsEntry\x1a\x33\n\x11\x41pplyConfigsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x0c:\x02\x38\x01\"\xeb\x01\n\x14ObservedChildObjects\x12\x41\n\x12group_version_kind\x18\x01 \x01(\x0b\x32%.kube_controller_api.GroupVersionKind\x12X\n\x10observed_objects\x18\x02 \x03(\x0b\x32>.kube_controller_api.ObservedChildObjects.ObservedObjectsEntry\x1a\x36\n\x14ObservedObjectsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x0c:\x02\x38\x01\x42\x36Z4github.com/LogitsAI/kube-controller-api/controllerpbb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'kube_controller_api.reconciler_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z4github.com/LogitsAI/kube-controller-api/controllerpb'
  _DESIREDCHILDOBJECTS_APPLYCONFIGSENTRY._options = None
  _DESIREDCHILDOBJECTS_APPLYCONFIGSENTRY._serialized_options = b'8\001'
  _OBSERVEDCHILDOBJECTS_OBSERVEDOBJECTSENTRY._options = None
  _OBSERVEDCHILDOBJECTS_OBSERVEDOBJECTSENTRY._serialized_options = b'8\001'
  _globals['_RECONCILERESULT']._serialized_start=132
  _globals['_RECONCILERESULT']._serialized_end=338
  _globals['_WORKQUEUE']._serialized_start=340
  _globals['_WORKQUEUE']._serialized_end=376
  _globals['_DESIREDCHILDOBJECTS']._serialized_start=379
  _globals['_DESIREDCHILDOBJECTS']._serialized_end=603
  _globals['_DESIREDCHILDOBJECTS_APPLYCONFIGSENTRY']._serialized_start=552
  _globals['_DESIREDCHILDOBJECTS_APPLYCONFIGSENTRY']._serialized_end=603
  _globals['_OBSERVEDCHILDOBJECTS']._serialized_start=606
  _globals['_OBSERVEDCHILDOBJECTS']._serialized_end=841
  _globals['_OBSERVEDCHILDOBJECTS_OBSERVEDOBJECTSENTRY']._serialized_start=787
  _globals['_OBSERVEDCHILDOBJECTS_OBSERVEDOBJECTSENTRY']._serialized_end=841
# @@protoc_insertion_point(module_scope)
