// START: clientConstructor

import { ApertureClient } from "@fluxninja/aperture-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  apiKey: "API_KEY",
});
// END: clientConstructor

// START: handleRequestWithCache
import { FlowStatus, LookupStatus } from "@fluxninja/aperture-js";
import { Request, Response } from "express";

async function handleRequest(req: Request, res: Response) {
  const flow = await apertureClient.startFlow("archimedes-service", {
    labels: {
      user: "user1",
      tier: "premium",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300, // 300ms deadline
    },
    resultCacheKey: "cache", // optional, in case caching is needed
  });

  if (flow.shouldRun()) {
    // Check if the response is cached in Aperture from a previous request
    if (flow.resultCache().getLookupStatus() === LookupStatus.Hit) {
      res.send({ message: flow.resultCache().getValue()?.toString() });
    } else {
      // Do Actual Work
      // After completing the work, you can store the response in the cache and return it, for example
      const resString = "foo";

      // create a new buffer
      const buffer = Buffer.from(resString);

      // START: setResultCache
      // set cache value
      const setResp = await flow.setResultCache({
        value: buffer,
        ttl: {
          seconds: 30,
          nanos: 0,
        },
      });
      if (setResp.getError()) {
        console.log(`Error setting cache value: ${setResp.getError()}`);
      }
      // END: setResultCache

      res.send({ message: resString });
    }
  } else {
    // Aperture has rejected the request due to a rate limiting policy
    res.status(429).send({ message: "Too many requests" });
    // Handle flow rejection
    flow.setStatus(FlowStatus.Error);
  }

  flow.end();
}
// END: handleRequestWithCache

// START: handleRequestRateLimit

async function handleRequestRateLimit(req: Request, res: Response) {
  const flow = await apertureClient.startFlow("awesomeFeature", {
    labels: {
      limit_key: "some_user_id",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300, // ms
    },
  });

  if (flow.shouldRun()) {
    // Add business logic to process incoming request
    console.log("Request accepted. Processing...");
    const resString = "foo";
    res.send({ message: resString });

  } else {
    console.log("Request rate-limited. Try again later.");
    // Handle flow rejection
    flow.setStatus(FlowStatus.Error);
    res.status(429).send({ message: "Too many requests" });
  }

  flow.end();
}
// END: handleRequestRateLimit

// START: handleQuotaScheduler
async function handleQuotaScheduler(apertureClient: ApertureClient, tier: string, priority: number) {
    // START: QSStartFlow
    const flow = await apertureClient.startFlow("quota-scheduling-feature", {
        labels: {
            limit_key: "some_user_id",
            priority: priority.toString(),
            workload: `${tier} user`,
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });
    console.log(`Request sent for ${tier} tier with priority ${priority}.`);
    flow.end();
    // END: QSStartFlow
}

function scheduleRequests(apertureClient: ApertureClient) {
    Object.entries(userTiers).forEach(([tier, priority]) => {
        setInterval(() => {
            handleQuotaScheduler(apertureClient, tier, priority);
        }, 1000);
    });
}
// END: handleQuotaScheduler

// START: handleConcurrencyLimit
async function sendRequest(apertureClient: ApertureClient) {
  const flow = await apertureClient.startFlow("concurrency-limiting-feature", {
    labels: {
      limit_key: "some_user_id",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300,
    },
  });

  if (flow.shouldRun()) {
    console.log("Request accepted. Processing..." + flow.checkResponse());
  } else {
    console.log("Request rejected due to concurrency limit. Try again later.");
  }

  await flow.end();
}

async function handleConcurrencyLimit(apertureClient: ApertureClient) {
  const requestsPerSecond = 10;
  const durationInSeconds = 200;

  for (let i = 0; i < durationInSeconds; i++) {
    const requests = Array.from({ length: requestsPerSecond }, () =>
      sendRequest(apertureClient),
    );
    // sending requests in parallel to simulate concurrency
    await Promise.all(requests);

    // Wait 1 second before sending the next batch of requests
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
}
// END: handleConcurrencyLimit

// Define user tiers and associated priorities
const userTiers = {
  platinum: 8,
  gold: 4,
  silver: 2,
  free: 1,
};

// START: handleConcurrencyScheduler
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
        limit_key: "some_user_id",
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
async function handleConcurrencyScheduler(apertureClient: ApertureClient) {
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
// END: handleConcurrencyScheduler
