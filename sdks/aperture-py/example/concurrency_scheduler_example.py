import asyncio

from aperture_sdk.client import ApertureClient, FlowParams

user_tiers = {
    "platinum": 8,
    "gold": 4,
    "silver": 2,
    "free": 1,
}
requests_per_batch = 10
batch_interval = 1

agent_address = "agents.us-central1.gcp.latest.dev.fluxninja.com:443"
api_key = "6428f436ddf647e9ab6c94c391750f39"


async def initialize_aperture_client():
    aperture_client = ApertureClient.new_client(address=agent_address, api_key=api_key)
    return aperture_client


async def send_request_for_tier(aperture_client, tier, priority):
    print(f"[{tier} Tier] Sending request with priority {priority}...")
    flow_params = FlowParams(
        explicit_labels={
            "user_id": "some_user_id",
            "priority": str(priority),
            "workload": tier,
        },
    )
    flow = aperture_client.start_flow(
        control_point="concurrency-scheduling-feature",
        params=flow_params,
    )

    if flow.should_run():
        print(f"[{tier} Tier] Request accepted with priority {priority}.")
        await asyncio.sleep(5)
    else:
        print(f"[{tier} Tier] Request rejected. Priority was {priority}.")

    flow.end()


async def schedule_requests(aperture_client):
    while True:
        print("Sending new batch of requests...")
        tasks = []
        for tier, priority in user_tiers.items():
            for _ in range(requests_per_batch):
                task = asyncio.create_task(
                    send_request_for_tier(aperture_client, tier, priority)
                )
                tasks.append(task)

        await asyncio.gather(*tasks)
        await asyncio.sleep(batch_interval)


async def main():
    aperture_client = await initialize_aperture_client()
    await schedule_requests(aperture_client)


if __name__ == "__main__":
    asyncio.run(main())
