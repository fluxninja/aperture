import {ApertureClient} from "@fluxninja/aperture-js";
import grpc from "@grpc/grpc-js";

async function initializeApertureClient() {
  const address = process.env.APERTURE_AGENT_ADDRESS || "localhost:8080";
  const apiKey = process.env.APERTURE_API_KEY || "";

  const apertureClient = new ApertureClient({
    address: address,
    apiKey: apiKey,
    channelCredentials: grpc.credentials.createInsecure(),
  });

  return apertureClient;
}

// Define user tiers and associated priorities
const userTiers = {
  platinum: 8,
  gold: 4,
  silver: 2,
  free: 1,
};

async function sendRequestForTier(
  apertureClient: ApertureClient,
  tier: string,
  priority: number,
) {
  console.log(`[${tier} Tier] Sending request with priority ${priority}...`);
  const flow = await apertureClient.startFlow(
    "concurrency-scheduling-feature",
    {
      labels: {
        user_id: "some_user_id",
        priority: priority.toString(),
        tier: tier,
      },
      grpcCallOptions: {
        deadline: Date.now() + 120000, // ms
      },
    },
  );

  if (flow.shouldRun()) {
    console.log(`[${tier} Tier] Request accepted with priority ${priority}.`);
    // sleep for 5 seconds to simulate a long-running request
    await new Promise((resolve) => setTimeout(resolve, 5000));
  } else {
    console.log(`[${tier} Tier] Request rejected. Priority was ${priority}.`);
  }

  await flow.end();
}

// Launch each batch in parallel
async function scheduleRequests(apertureClient: ApertureClient) {
  const requestsPerBatch = 10;
  const batchInterval = 1000; // ms

  while (true) {
    console.log("Sending new batch of requests...");
    // Send requests for each tier
    const promises = Object.entries(userTiers).flatMap(([tier, priority]) => {
      return Array(requestsPerBatch)
        .fill(null)
        .map(() => sendRequestForTier(apertureClient, tier, priority));
    });

    await Promise.all(promises);
    await new Promise((resolve) => setTimeout(resolve, batchInterval));
  }
}

async function main() {
  const apertureClient = await initializeApertureClient();
  scheduleRequests(apertureClient);
}

main().catch((e) => console.error(e));
