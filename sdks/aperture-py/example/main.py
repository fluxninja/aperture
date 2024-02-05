#!/usr/bin/env python3
import asyncio
import logging
import os
from datetime import timedelta
from typing import Optional

import grpc
from quart import Quart

default_agent_address = "localhost:8089"
app = Quart(__name__)

# START: clientConstructor
from aperture_sdk.client import ApertureClient, FlowParams

agent_address = os.getenv("APERTURE_AGENT_ADDRESS", default_agent_address)
api_key = os.getenv("APERTURE_API_KEY", "")
insecure = os.getenv("APERTURE_AGENT_INSECURE", "true").lower() == "true"

aperture_client = ApertureClient.new_client(
    address=agent_address, insecure=insecure, api_key=api_key
)
# END: clientConstructor

logging.basicConfig(level=logging.DEBUG)
default_agent_address = "localhost:8089"
app = Quart(__name__)

# START: apertureDecorator
flow_params = FlowParams(
    check_timeout=timedelta(seconds=200),
)


@app.get("/super")
@aperture_client.decorate(
    "awesomeFeature", params=flow_params, on_reject=lambda: ("Flow was rejected", 503)
)
async def super_handler():
    # Simulate work being done
    await asyncio.sleep(2)
    return "", 202


# END: apertureDecorator


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


async def syncSnippet(aperture_client):
    # START: syncFlow
    labels = {
        "userId": "some_user_id",
        "userTier": "gold",
        "priority": "100",
    }
    flow_params = FlowParams(
        check_timeout=timedelta(seconds=400),
        explicit_labels=labels,
    )
    flow = aperture_client.start_flow(
        control_point="rate-limiting-feature",
        params=flow_params,
    )

    if flow.should_run():
        print("Request accepted. Processing...")
    else:
        print("Request rate-limited. Try again later.")
        flow.set_status()

    flow.end()
    # Simulate work being done
    await asyncio.sleep(1)
    # END: syncFlow
