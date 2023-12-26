import sys
import asyncio

sys.path.append("./python/src")

from kube_controller_api.client import (
    Connection,
    ControllerManagerConfig,
    ControllerConfig,
    ReconcileRequest,
    ReconcileResult,
    GroupVersionKind,
)

CONFIG_MAP_GVK = GroupVersionKind("", "v1", "ConfigMap")

async def reconcile_example(request: ReconcileRequest) -> ReconcileResult:
    namespace = request.parent["metadata"]["namespace"]
    name = request.parent["metadata"]["name"]
    print(f"Example {namespace}/{name} reconciled")

    observed_config_maps = request.children[CONFIG_MAP_GVK]
    desired_config_maps = {
        "child-1": {
            "data": {
                "key": "value 1",
            }
        },
        "child-2": {
            "data": {
                "key": "value 2",
            }
        }
    }

    return ReconcileResult(
        status={
            "output": request.parent["spec"]["input"] + " output",
            "configMaps": list(observed_config_maps),
        },
        children={
            CONFIG_MAP_GVK: desired_config_maps,
        }
    )

async def main():
    # Create a gRPC channel bound to the server address.
    async with Connection("localhost:8090") as conn:
        config = ControllerManagerConfig()

        controller = ControllerConfig(
            name="example-controller",
            parent=GroupVersionKind("example.com", "v1", "Example"),
            children=[
                GroupVersionKind("", "v1", "ConfigMap"),
                ],
            )
        config.controllers.append(controller)

        # Create a remote ControllerManager instance on the server.
        await conn.start_manager(config)

        # Start processing reconcile requests from the server's work queue.
        await conn.reconcile_loop(controller.name, reconcile_example)

if __name__ == "__main__":
    asyncio.run(main())
