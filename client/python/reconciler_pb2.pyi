from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ReconcileResult(_message.Message):
    __slots__ = ["error", "requeue", "requeue_after"]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    REQUEUE_FIELD_NUMBER: _ClassVar[int]
    REQUEUE_AFTER_FIELD_NUMBER: _ClassVar[int]
    error: str
    requeue: bool
    requeue_after: _duration_pb2.Duration
    def __init__(self, error: _Optional[str] = ..., requeue: bool = ..., requeue_after: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class WorkQueue(_message.Message):
    __slots__ = ["manager_id", "controller_name"]
    MANAGER_ID_FIELD_NUMBER: _ClassVar[int]
    CONTROLLER_NAME_FIELD_NUMBER: _ClassVar[int]
    manager_id: str
    controller_name: str
    def __init__(self, manager_id: _Optional[str] = ..., controller_name: _Optional[str] = ...) -> None: ...
