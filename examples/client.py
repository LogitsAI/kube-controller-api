import sys
import asyncio

sys.path.append("./client/python")

from controller_api import (
    Connection,
    ControllerManagerConfig,
    ControllerConfig,
    ReconcileRequest,
    ReconcileResult,
)

async def reconcile_pod(request: ReconcileRequest) -> ReconcileResult:
    namespace = request.object["metadata"]["namespace"]
    name = request.object["metadata"]["name"]
    print(f"{namespace}/{name} reconciled")

    return ReconcileResult()

async def main():
    # Create a gRPC channel bound to the server address.
    async with Connection("localhost:8090") as conn:
        config = ControllerManagerConfig(name="my-manager")

        controller = ControllerConfig(name="my-controller")
        controller.set_parent("", "v1", "Pod")
        config.controllers.append(controller)

        # Create a remote ControllerManager instance on the server.
        manager = await conn.create_manager(config)

        # Start processing reconcile requests from the server's work queue.
        await conn.reconcile_loop(manager.id, controller.name, reconcile_pod)

if __name__ == "__main__":
    asyncio.run(main())
