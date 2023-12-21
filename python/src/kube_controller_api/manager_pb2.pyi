from kube_controller_api import controller_pb2 as _controller_pb2
from kube_controller_api import reconciler_pb2 as _reconciler_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class StartRequest(_message.Message):
    __slots__ = ["controllers"]
    CONTROLLERS_FIELD_NUMBER: _ClassVar[int]
    controllers: _containers.RepeatedCompositeFieldContainer[_controller_pb2.ControllerConfig]
    def __init__(self, controllers: _Optional[_Iterable[_Union[_controller_pb2.ControllerConfig, _Mapping]]] = ...) -> None: ...

class StartResponse(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class ReconcileLoopRequest(_message.Message):
    __slots__ = ["subscribe", "acknowledge"]
    SUBSCRIBE_FIELD_NUMBER: _ClassVar[int]
    ACKNOWLEDGE_FIELD_NUMBER: _ClassVar[int]
    subscribe: _reconciler_pb2.WorkQueue
    acknowledge: _reconciler_pb2.ReconcileResult
    def __init__(self, subscribe: _Optional[_Union[_reconciler_pb2.WorkQueue, _Mapping]] = ..., acknowledge: _Optional[_Union[_reconciler_pb2.ReconcileResult, _Mapping]] = ...) -> None: ...

class ReconcileLoopResponse(_message.Message):
    __slots__ = ["parent", "children"]
    PARENT_FIELD_NUMBER: _ClassVar[int]
    CHILDREN_FIELD_NUMBER: _ClassVar[int]
    parent: bytes
    children: _containers.RepeatedCompositeFieldContainer[_reconciler_pb2.ObservedChildObjects]
    def __init__(self, parent: _Optional[bytes] = ..., children: _Optional[_Iterable[_Union[_reconciler_pb2.ObservedChildObjects, _Mapping]]] = ...) -> None: ...
