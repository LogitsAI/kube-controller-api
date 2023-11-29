from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ControllerConfig(_message.Message):
    __slots__ = ["name", "parent", "children"]
    NAME_FIELD_NUMBER: _ClassVar[int]
    PARENT_FIELD_NUMBER: _ClassVar[int]
    CHILDREN_FIELD_NUMBER: _ClassVar[int]
    name: str
    parent: GroupVersionKind
    children: _containers.RepeatedCompositeFieldContainer[GroupVersionKind]
    def __init__(self, name: _Optional[str] = ..., parent: _Optional[_Union[GroupVersionKind, _Mapping]] = ..., children: _Optional[_Iterable[_Union[GroupVersionKind, _Mapping]]] = ...) -> None: ...

class GroupVersionKind(_message.Message):
    __slots__ = ["group", "version", "kind"]
    GROUP_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    group: str
    version: str
    kind: str
    def __init__(self, group: _Optional[str] = ..., version: _Optional[str] = ..., kind: _Optional[str] = ...) -> None: ...
