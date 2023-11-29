#!/usr/bin/env python3
import asyncio
import logging
import os
from datetime import timedelta
from typing import Optional

import grpc
from aperture_sdk import ApertureClient, FlowParams, FlowStatus
from quart import Quart

default_agent_address = "localhost:8089"
app = Quart(__name__)

agent_address = os.getenv("APERTURE_AGENT_ADDRESS", default_agent_address)
api_key = os.getenv("APERTURE_API_KEY", "")
insecure = os.getenv("APERTURE_AGENT_INSECURE", "true").lower() == "true"

aperture_client = ApertureClient.new_client(
    address=agent_address, insecure=insecure, api_key=api_key
)

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger("aperture-manual-flow-example")


@app.get("/super")
async def super_handler():
    # START: manualFlow
    # business logic produces labels
    labels = {
        "key": "value",
    }
    flow_params = FlowParams(
        check_timeout=timedelta(seconds=200),
        explicit_labels=labels,
    )
    # start_flow performs a flowcontrol.v1.Check call to Aperture Agent.
    # It returns a Flow or raises an error if any.
    flow = aperture_client.start_flow(
        control_point="AwesomeFeature",
        params=flow_params,
    )

    # Check if flow check was successful.
    if not flow.success:
        logger.info("Flow check failed - will fail-open")

    # See whether flow was accepted by Aperture Agent.
    if flow.should_run():
        # do actual work
        pass
    else:
        # handle flow rejection by Aperture Agent
        flow.set_status(FlowStatus.Error)
    flow.end()
    # Simulate work being done
    await asyncio.sleep(2)
    return "", 202
    # END: manualFlow


@app.get("/super2")
async def super2_handler():
    # business logic produces labels
    labels = {
        "key": "value",
    }
    # START: contextManagerFlow

    flow_params = FlowParams(
        explicit_labels=labels,
        check_timeout=timedelta(seconds=200),
    )

    with aperture_client.start_flow(
        control_point="AwesomeFeature",
        params=flow_params,
    ) as flow:
        if flow.should_run():
            # do actual work
            # if you do not call flow.end() explicitly, it will be called automatically
            # when the context manager exits - with the status of the flow
            # depending on whether an error was raised or not
            pass
    # END: contextManagerFlow


@app.get("/connected")
async def connected_handler():
    state: Optional[grpc.ChannelConnectivity] = None

    def subscribe_callback(connectivity: grpc.ChannelConnectivity):
        nonlocal state
        state = connectivity

    # gRPC does not expose a way to get the current state of the channel
    # so we subscribe to the channel and unsubscribe immediately,
    # as callback is triggered immediately when subscribing
    grpc_channel = aperture_client.grpc_channel
    grpc_channel.subscribe(subscribe_callback, try_to_connect=True)
    grpc_channel.unsubscribe(subscribe_callback)
    print(f"State: {state}")
    http_code = 200 if state == grpc.ChannelConnectivity.READY else 503
    return str(state), http_code


@app.get("/health")
async def health_handler():
    return "Healthy", 200


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
