import json
import contextlib
from datetime import timedelta
from typing import Callable, Coroutine, Any, TypeAlias
from dataclasses import dataclass, field

import grpc.aio

from . import (
    manager_pb2_grpc,
    manager_pb2,
    controller_pb2,
    reconciler_pb2,
)

@dataclass(frozen=True)
class GroupVersionKind:
    group: str
    version: str
    kind: str

    def __init__(self, group: str, version: str, kind: str):
        """This constructor allows us to use positional args."""

        # Bypass frozenness for construction, like the default dataclass constructor would.
        object.__setattr__(self, "group", group)
        object.__setattr__(self, "version", version)
        object.__setattr__(self, "kind", kind)

    def to_proto(self):
        return controller_pb2.GroupVersionKind(group=self.group, version=self.version, kind=self.kind)

    def api_version(self):
        if self.group == "":
            return self.version
        else:
            return f"{self.group}/{self.version}"

    @staticmethod
    def from_proto(gvk):
        return GroupVersionKind(group=gvk.group, version=gvk.version, kind=gvk.kind)

@dataclass
class ControllerConfig:
    name: str
    parent: GroupVersionKind
    children: list[GroupVersionKind] = field(default_factory=list)

    def to_proto(self):
        controller = controller_pb2.ControllerConfig(name=self.name)
        controller.parent.CopyFrom(self.parent.to_proto())
        for child in self.children:
            controller.children.add().CopyFrom(child.to_proto())
        return controller


@dataclass
class ControllerManagerConfig:
    controllers: list[ControllerConfig] = field(default_factory=list)

    def to_proto(self):
        request = manager_pb2.StartRequest()
        for controller in self.controllers:
            request.controllers.add().CopyFrom(controller.to_proto())
        return request


Object: TypeAlias = dict[str, Any]
ObjectMap: TypeAlias = dict[str, Object]


@dataclass
class ReconcileRequest:
    conn: 'Connection'
    parent: dict[str, Any]
    children: dict[GroupVersionKind, ObjectMap] = field(default_factory=dict)

    @staticmethod
    def from_proto(conn: 'Connection', response: manager_pb2.ReconcileLoopResponse):
        """
        Convert a ReconcileLoopResponse from the server into a ReconcileRequest.
        Note that it comes to us as an RPC response because we make an outgoing
        call to the server to subscribe to a stream of reconcile requests.
        """

        # Decode child objects.
        children = {}
        for child_objects in response.children:
            gvk = GroupVersionKind.from_proto(child_objects.group_version_kind)
            children[gvk] = {
                key: json.loads(obj) for key, obj in child_objects.observed_objects.items()
            }

        return ReconcileRequest(
                conn=conn,
                parent=json.loads(response.parent),
                children=children,
            )

@dataclass
class ReconcileResult:
    requeue: bool = False
    requeue_after: timedelta | None = None
    status: Any | None = None
    children: dict[GroupVersionKind, ObjectMap] = field(default_factory=dict)

    def to_proto(self):
        result = reconciler_pb2.ReconcileResult()
        result.requeue = self.requeue
        if self.requeue_after is not None:
            result.requeue_after.FromTimedelta(self.requeue_after)
        if self.status is not None:
            result.status = json.dumps(self.status).encode("utf-8")
        for gvk, objects in self.children.items():
            child_objects = result.children.add()
            child_objects.group_version_kind.CopyFrom(gvk.to_proto())
            for key, obj in objects.items():
                obj["apiVersion"] = gvk.api_version()
                obj["kind"] = gvk.kind
                child_objects.apply_configs[key] = json.dumps(obj).encode("utf-8")
        return result


class Connection(contextlib.AbstractAsyncContextManager):
    def __init__(self, address):
        self.address = address

    async def __aenter__(self):
        self.channel = grpc.aio.insecure_channel(self.address)
        await self.channel.__aenter__()
        self.stub = manager_pb2_grpc.ControllerManagerStub(self.channel)
        return self

    async def __aexit__(self, exc_type, exc_value, traceback):
        await self.channel.__aexit__(exc_type, exc_value, traceback)
    
    async def start_manager(self, config: ControllerManagerConfig):
        manager = await self.stub.Start(config.to_proto())
        return manager

    async def reconcile_loop(self, controller_name,
                             reconcile_func: Callable[[ReconcileRequest], Coroutine[Any, Any, ReconcileResult]]):
        stream = self.stub.ReconcileLoop()
        
        # Specify the manager ID and controller name.
        await stream.write(manager_pb2.ReconcileLoopRequest(
            subscribe=reconciler_pb2.WorkQueue(controller_name=controller_name),
        ))

        # Process objects from the work queue.
        async for response in stream:
            request = ReconcileRequest.from_proto(self, response)

            try:
                result = await reconcile_func(request)

                # Report the result to the server.
                await stream.write(manager_pb2.ReconcileLoopRequest(
                    acknowledge=result.to_proto(),
                ))
            except Exception as e:
                # Report the error to the server.
                await stream.write(manager_pb2.ReconcileLoopRequest(
                    acknowledge=reconciler_pb2.ReconcileResult(error=str(e)),
                ))
