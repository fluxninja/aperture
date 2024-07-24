import asyncio
from datetime import timedelta

from aperture_sdk.client import ApertureClient, FlowParams

user_tiers = {
    "platinum": 8,
    "gold": 4,
    "silver": 2,
    "free": 1,
}
interval_time = 1


async def initialize_aperture_client():
    agent_address = ""
    api_key = ""
    aperture_client = ApertureClient.new_client(address=agent_address, api_key=api_key)
    return aperture_client


async def send_request_for_tier(aperture_client, tier, priority):
    labels = {
        "user_id": "some_user_id",
        "priority": str(priority),
        "workload": f"{tier} user",
    }
    flow_params = FlowParams(
        check_timeout=timedelta(seconds=400),
        explicit_labels=labels,
    )

    flow = await aperture_client.start_flow(
        control_point="quota-scheduling-feature",
        params=flow_params,
    )
    print(f"Request sent for {tier} tier with priority {priority}.")
    flow.end()


async def schedule_requests(aperture_client):
    while True:
        for tier, priority in user_tiers.items():
            await send_request_for_tier(aperture_client, tier, priority)
            await asyncio.sleep(interval_time)


async def main():
    aperture_client = await initialize_aperture_client()
    await schedule_requests(aperture_client)


if __name__ == "__main__":
    asyncio.run(main())
