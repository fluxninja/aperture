import asyncio
from datetime import timedelta

from aperture_sdk.client_async import ApertureClientAsync, FlowParams
from aperture_sdk.flow_common import FlowStatus

user_tiers = {
    "free": 1,
    "silver": 2,
    "gold": 4,
    "platinum": 8,
}
requests_per_batch = 50
batch_interval = 1

agent_address = ""
api_key = ""


def initialize_aperture_client():
    aperture_client = ApertureClientAsync.new_client(
        address=agent_address, api_key=api_key
    )
    return aperture_client


async def send_request_for_tier(i, aperture_client, tier, priority):
    print(f"[{tier} Tier] Sending request with priority {priority}...")
    flow_params = FlowParams(
        explicit_labels={
            "priority": str(priority),
            "workload": tier,
        },
        check_timeout=timedelta(seconds=10),
    )

    flow = await aperture_client.start_flow(
        control_point="concurrency-scheduling-feature",
        params=flow_params,
    )

    if flow.should_run():
        await asyncio.sleep(0.1)
        print(f"[{tier} Tier] Request accepted. Priority was {priority}.")
    else:
        flow.set_status(FlowStatus.Error)
        print(f"[{tier} Tier] Request rejected. Priority was {priority}.")

    await flow.end()


async def schedule_requests(aperture_client):
    while True:
        print("Sending new batch of requests...")

        tasks = []
        for tier, priority in user_tiers.items():
            for i in range(requests_per_batch):
                task = asyncio.create_task(
                    send_request_for_tier(i, aperture_client, tier, priority)
                )
                tasks.append(task)

        await asyncio.gather(*tasks)
        await asyncio.sleep(batch_interval)


async def main():
    aperture_client = initialize_aperture_client()
    await schedule_requests(aperture_client)


if __name__ == "__main__":
    asyncio.run(main())
