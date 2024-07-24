import asyncio

from aperture_sdk.client import ApertureClient, FlowParams

requests_per_second = 5
duration_in_seconds = 10

agent_address = ""
api_key = ""


async def initialize_aperture_client():
    aperture_client = ApertureClient.new_client(address=agent_address, api_key=api_key)
    return aperture_client


async def send_request(aperture_client):
    print("Sending request...")
    flow_params = FlowParams(
        explicit_labels={"user_id": "some_user_id"},
    )
    flow = aperture_client.start_flow(
        control_point="concurrency-limiting-feature",
        params=flow_params,
    )

    if flow.should_run():
        print("Request accepted. Processing...")
        # Simulate processing
        await asyncio.sleep(5)
    else:
        print("Request rejected due to concurrency limit. Try again later.")

    await flow.end()


async def handle_concurrency_limit(aperture_client):
    for _ in range(duration_in_seconds):
        tasks = [
            asyncio.create_task(send_request(aperture_client))
            for _ in range(requests_per_second)
        ]
        await asyncio.gather(*tasks)
        await asyncio.sleep(1)


async def main():
    aperture_client = await initialize_aperture_client()
    await handle_concurrency_limit(aperture_client)


if __name__ == "__main__":
    asyncio.run(main())
