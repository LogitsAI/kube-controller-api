import json
import contextlib
from datetime import timedelta
from typing import Callable, Coroutine, Any
from dataclasses import dataclass, field

import grpc.aio

from . import (
    manager_pb2_grpc,
    manager_pb2,
    controller_pb2,
    reconciler_pb2,
)

@dataclass
class ControllerConfig:
    name: str
    parent: controller_pb2.GroupVersionKind = None

    def set_parent(self, group, version, kind):
        self.parent = controller_pb2.GroupVersionKind(group=group, version=version, kind=kind)

    def to_proto(self):
        controller = controller_pb2.ControllerConfig(name=self.name)
        if self.parent is not None:
            controller.parent.CopyFrom(self.parent)
        return controller


@dataclass
class ControllerManagerConfig:
    controllers: list[ControllerConfig] = field(default_factory=list)

    def to_proto(self):
        request = manager_pb2.StartRequest()
        for controller in self.controllers:
            request.controllers.add().CopyFrom(controller.to_proto())
        return request


@dataclass
class ReconcileRequest:
    object: Any


@dataclass
class ReconcileResult:
    requeue: bool = False
    requeue_after: timedelta | None = None

    def to_proto(self):
        result = reconciler_pb2.ReconcileResult()
        result.requeue = self.requeue
        if self.requeue_after is not None:
            result.requeue_after.FromTimedelta(self.requeue_after)
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
            request = ReconcileRequest(object=json.loads(response.object))

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
