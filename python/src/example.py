import sys
import asyncio

sys.path.append("./python/src")

from kube_controller_api.client import (
    Controller,
    ControllerManager,
    ReconcileRequest,
    ReconcileResult,
    GroupVersionKind,
)

CONFIG_MAP_GVK = GroupVersionKind("", "v1", "ConfigMap")

class ExampleController(Controller,
        name="example.com/example-controller",
        parent=GroupVersionKind("example.com", "v1", "Example"),
        children=[CONFIG_MAP_GVK],
        ):

    async def reconcile(self, request: ReconcileRequest) -> ReconcileResult:
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


if __name__ == "__main__":
    controllers = [
        ExampleController(),
    ]

    manager = ControllerManager(address="localhost:8090", controllers=controllers)

    asyncio.run(manager.start())
