#!/usr/bin/env python3
import asyncio
import logging
import os
from typing import Optional

import grpc
from aperture_sdk import ApertureClient
from quart import Quart

default_agent_address = "localhost:8089"

agent_address = os.getenv("APERTURE_AGENT_ADDRESS", default_agent_address)
api_key = os.getenv("APERTURE_API_KEY", "")
insecure = os.getenv("APERTURE_AGENT_INSECURE", "true").lower() == "true"

app = Quart(__name__)
aperture_client = ApertureClient.new_client(
    address=agent_address, insecure=insecure, api_key=api_key
)

logging.basicConfig(level=logging.DEBUG)


@app.get("/super")
@aperture_client.decorate(
    "awesomeFeature", on_reject=lambda: ("Flow was rejected", 503)
)
async def super_handler():
    # Simulate work being done
    await asyncio.sleep(2)
    return "", 202


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
