from google.protobuf import duration_pb2 as _duration_pb2
from kube_controller_api import controller_pb2 as _controller_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ReconcileResult(_message.Message):
    __slots__ = ["error", "requeue", "requeue_after", "status", "children"]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    REQUEUE_FIELD_NUMBER: _ClassVar[int]
    REQUEUE_AFTER_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CHILDREN_FIELD_NUMBER: _ClassVar[int]
    error: str
    requeue: bool
    requeue_after: _duration_pb2.Duration
    status: bytes
    children: _containers.RepeatedCompositeFieldContainer[DesiredChildObjects]
    def __init__(self, error: _Optional[str] = ..., requeue: bool = ..., requeue_after: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., status: _Optional[bytes] = ..., children: _Optional[_Iterable[_Union[DesiredChildObjects, _Mapping]]] = ...) -> None: ...

class WorkQueue(_message.Message):
    __slots__ = ["controller_name"]
    CONTROLLER_NAME_FIELD_NUMBER: _ClassVar[int]
    controller_name: str
    def __init__(self, controller_name: _Optional[str] = ...) -> None: ...

class DesiredChildObjects(_message.Message):
    __slots__ = ["group_version_kind", "apply_configs"]
    class ApplyConfigsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: bytes
        def __init__(self, key: _Optional[str] = ..., value: _Optional[bytes] = ...) -> None: ...
    GROUP_VERSION_KIND_FIELD_NUMBER: _ClassVar[int]
    APPLY_CONFIGS_FIELD_NUMBER: _ClassVar[int]
    group_version_kind: _controller_pb2.GroupVersionKind
    apply_configs: _containers.ScalarMap[str, bytes]
    def __init__(self, group_version_kind: _Optional[_Union[_controller_pb2.GroupVersionKind, _Mapping]] = ..., apply_configs: _Optional[_Mapping[str, bytes]] = ...) -> None: ...

class ObservedChildObjects(_message.Message):
    __slots__ = ["group_version_kind", "observed_objects"]
    class ObservedObjectsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: bytes
        def __init__(self, key: _Optional[str] = ..., value: _Optional[bytes] = ...) -> None: ...
    GROUP_VERSION_KIND_FIELD_NUMBER: _ClassVar[int]
    OBSERVED_OBJECTS_FIELD_NUMBER: _ClassVar[int]
    group_version_kind: _controller_pb2.GroupVersionKind
    observed_objects: _containers.ScalarMap[str, bytes]
    def __init__(self, group_version_kind: _Optional[_Union[_controller_pb2.GroupVersionKind, _Mapping]] = ..., observed_objects: _Optional[_Mapping[str, bytes]] = ...) -> None: ...
