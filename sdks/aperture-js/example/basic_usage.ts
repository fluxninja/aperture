// START: clientConstructor

import { ApertureClient } from "@fluxninja/aperture-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  apiKey: "API_KEY",
});
// END: clientConstructor

// START: handleRequest
import { FlowStatusEnum, LookupStatus } from "@fluxninja/aperture-js";
import { Request, Response } from "express";

async function handleRequest(req: Request, res: Response) {
  const flow = await apertureClient.StartFlow("archimedes-service", {
    labels: {
      api_key: "some_api_key",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300, // 300ms deadline
    },
    cacheKey: "cache", // optional, in case caching is needed
  });

  if (flow.ShouldRun()) {
    // Check if the response is cached in Aperture from a previous request
    if (flow.CachedValue().GetLookupStatus() === LookupStatus.Hit) {
      res.send({ message: flow.CachedValue().GetValue()?.toString() });
    } else {
      // Do Actual Work
      // After completing the work, you can return store the response in cache and return it, for example:
      const resString = "foo";

      // create a new buffer
      const buffer = Buffer.from(resString);

      // START: setCache
      // set cache value
      const setResult = await flow.SetCachedValue(buffer, {
        seconds: 30,
        nanos: 0,
      });
      if (setResult?.error) {
        console.log(`Error setting cache value: ${setResult.error}`);
      }
      // END: setCache

      res.send({ message: resString });
    }
  } else {
    // Aperture has rejected the request due to a rate limiting policy
    res.status(429).send({ message: "Too many requests" });
    // Handle flow rejection
    flow.SetStatus(FlowStatusEnum.Error);
  }

  flow.End();
}
// END: handleRequest
