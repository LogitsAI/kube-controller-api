import controller_pb2 as _controller_pb2
import reconciler_pb2 as _reconciler_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CreateManagerRequest(_message.Message):
    __slots__ = ["name", "controllers"]
    NAME_FIELD_NUMBER: _ClassVar[int]
    CONTROLLERS_FIELD_NUMBER: _ClassVar[int]
    name: str
    controllers: _containers.RepeatedCompositeFieldContainer[_controller_pb2.ControllerConfig]
    def __init__(self, name: _Optional[str] = ..., controllers: _Optional[_Iterable[_Union[_controller_pb2.ControllerConfig, _Mapping]]] = ...) -> None: ...

class CreateManagerResponse(_message.Message):
    __slots__ = ["id"]
    ID_FIELD_NUMBER: _ClassVar[int]
    id: str
    def __init__(self, id: _Optional[str] = ...) -> None: ...

class ReconcileLoopRequest(_message.Message):
    __slots__ = ["subscribe", "acknowledge"]
    SUBSCRIBE_FIELD_NUMBER: _ClassVar[int]
    ACKNOWLEDGE_FIELD_NUMBER: _ClassVar[int]
    subscribe: _reconciler_pb2.WorkQueue
    acknowledge: _reconciler_pb2.ReconcileResult
    def __init__(self, subscribe: _Optional[_Union[_reconciler_pb2.WorkQueue, _Mapping]] = ..., acknowledge: _Optional[_Union[_reconciler_pb2.ReconcileResult, _Mapping]] = ...) -> None: ...

class ReconcileLoopResponse(_message.Message):
    __slots__ = ["object"]
    OBJECT_FIELD_NUMBER: _ClassVar[int]
    object: bytes
    def __init__(self, object: _Optional[bytes] = ...) -> None: ...
