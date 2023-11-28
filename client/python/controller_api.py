import json
import contextlib
from datetime import timedelta
from typing import Callable, Coroutine, Any

import grpc.aio

import manager_pb2_grpc
import manager_pb2
import controller_pb2
import reconciler_pb2

class ControllerManagerConfig:
    def __init__(self, name):
        self.name = name
        self.controllers = []

    def to_proto(self):
        request = manager_pb2.CreateManagerRequest(name=self.name)
        for controller in self.controllers:
            request.controllers.add().CopyFrom(controller.to_proto())
        return request


class ControllerConfig:
    def __init__(self, name):
        self.name = name
        self.parent = None

    def set_parent(self, group, version, kind):
        self.parent = controller_pb2.GroupVersionKind(group=group, version=version, kind=kind)

    def to_proto(self):
        controller = controller_pb2.ControllerConfig(name=self.name)
        if self.parent is not None:
            controller.parent.CopyFrom(self.parent)
        return controller


class ReconcileRequest:
    def __init__(self, object):
        self.object = object


class ReconcileResult:
    def __init__(self, requeue: bool = False, requeue_after: timedelta | None = None):
        self.requeue = requeue
        self.requeue_after = requeue_after

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
    
    async def create_manager(self, config: ControllerManagerConfig):
        manager = await self.stub.CreateManager(config.to_proto())
        return manager

    async def reconcile_loop(self, manager_id, controller_name, reconcile_func: Callable[[ReconcileRequest], Coroutine[Any, Any, ReconcileResult]]):
        stream = self.stub.ReconcileLoop()
        
        # Specify the manager ID and controller name.
        await stream.write(manager_pb2.ReconcileLoopRequest(
            subscribe=reconciler_pb2.WorkQueue(manager_id=manager_id, controller_name=controller_name),
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
