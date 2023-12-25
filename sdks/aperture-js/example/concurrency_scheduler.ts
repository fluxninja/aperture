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
  } else {
    console.log(`[${tier} Tier] Request rejected. Priority was ${priority}.`);
  }

  const flowEndResponse = await flow.end();
}

async function scheduleRequests(apertureClient: ApertureClient) {
  const requestsPerBatch = 10;
  const batchInterval = 1000; // ms

  setInterval(() => {
    console.log("Sending new batch of requests...");
    Object.entries(userTiers).forEach(([tier, priority]) => {
      for (let i = 0; i < requestsPerBatch; i++) {
        sendRequestForTier(apertureClient, tier, priority);
      }
    });
  }, batchInterval);
}

async function main() {
  const apertureClient = await initializeApertureClient();
  scheduleRequests(apertureClient);
}

main().catch((e) => console.error(e));
