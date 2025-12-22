import asyncio
from datetime import timedelta

from aperture_sdk.client import ApertureClient, FlowParams, LookupStatus

agent_address = ""
api_key = ""


async def initialize_aperture_client():
    aperture_client = ApertureClient.new_client(address=agent_address, api_key=api_key)
    return aperture_client


async def monitor_cache_and_update(aperture_client):
    labels = {
        "userId": "some_user_id",
    }
    flow_params = FlowParams(
        check_timeout=timedelta(seconds=400),
        explicit_labels=labels,
        result_cache_key="cache",
    )
    while True:
        flow = aperture_client.start_flow(
            control_point="caching-example",
            params=flow_params,
        )

        cache_response = flow.result_cache()
        if cache_response.get_lookup_status() == LookupStatus.HIT:
            print("Cache hit: Value =", cache_response.get_value().decode())
        else:
            print("Cache miss, setting cache value")
            cache_set_response = flow.set_result_cache(
                value="Hello, world!",
                ttl=timedelta(seconds=10),
            )
            if cache_set_response.get_error():
                print(f"Error setting cache: {cache_set_response.get_error()}")

        flow.end()
        await asyncio.sleep(1)


# Main function
async def main():
    aperture_client = await initialize_aperture_client()
    await monitor_cache_and_update(aperture_client)


# Run the main function
if __name__ == "__main__":
    asyncio.run(main())
