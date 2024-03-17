import asyncio
from datetime import timedelta

from aperture_sdk.client import ApertureClient, FlowParams

agent_address = ""
api_key = ""


async def initialize_aperture_client():
    aperture_client = ApertureClient.new_client(address=agent_address, api_key=api_key)
    return aperture_client


async def handle_request_rate_limit(aperture_client):
    labels = {
        "userId": "some_user_id",
    }
    flow_params = FlowParams(
        check_timeout=timedelta(seconds=400),
        explicit_labels=labels,
    )
    while True:
        flow = await aperture_client.start_flow(
            control_point="rate-limiting-feature",
            params=flow_params,
        )

        if flow.should_run():
            print("Request accepted. Processing...")
        else:
            print("Request rate-limited. Try again later.")
            flow.set_status()

        res = flow.end()
        if res.get_error():
            pass
        elif res.get_flow_end_response():
            pass

        await asyncio.sleep(0.5)


# Main function
async def main():
    aperture_client = await initialize_aperture_client()
    await handle_request_rate_limit(aperture_client)


# Run the main function
if __name__ == "__main__":
    asyncio.run(main())
