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
      user_id: "some_user_id",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300, // ms
    },
  });

  if (flow.shouldRun()) {
    // Add business logic to process incoming request
    console.log("Request accepted. Processing...");
  } else {
    console.log("Request rate-limited. Try again later.");
    // Handle flow rejection
    flow.setStatus(FlowStatus.Error);
  }

  if (flow) {
    flow.end();
  }
}
// END: handleRequestRateLimit
