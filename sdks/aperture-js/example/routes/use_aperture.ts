import {
  ApertureClient,
  Flow,
  FlowStatus,
  LookupStatus,
} from "@fluxninja/aperture-js";
import grpc from "@grpc/grpc-js";
import express from "express";

// Create aperture client
export const apertureClient = new ApertureClient({
  address:
    process.env.APERTURE_AGENT_ADDRESS !== undefined
      ? process.env.APERTURE_AGENT_ADDRESS
      : "localhost:8089",
  apiKey: process.env.APERTURE_API_KEY || undefined,
  // if process.env.APERTURE_AGENT_INSECURE set channelCredentials to insecure
  channelCredentials:
    process.env.APERTURE_AGENT_INSECURE !== undefined
      ? grpc.credentials.createInsecure()
      : grpc.credentials.createSsl(),
});

export const apertureRoute = express.Router();

apertureRoute.get("/", async (_: express.Request, res: express.Response) => {
  const labels: Record<string, string> = { user: "kenobi" };
  const startTimestamp = Date.now();
  let flow: Flow | undefined = undefined;

  try {
    flow = await apertureClient.startFlow("awesomeFeature", {
      labels: labels,
      grpcCallOptions: {
        deadline: Date.now() + 30000,
      },
      rampMode: false,
      resultCacheKey: "cache",
    });

    const endTimestamp = Date.now();
    console.log(`Flow took ${endTimestamp - startTimestamp}ms`);

    if (flow.shouldRun()) {
      await sleep(200);

      if (flow.resultCache().getLookupStatus() === LookupStatus.Hit) {
        console.log("Cache hit:", flow.resultCache().getValue()?.toString());
      } else {
        console.log("Cache miss:", flow.resultCache().getError()?.message);
        const resString = "awesomeString";

        // create a new buffer
        const buffer = Buffer.from(resString);

        // set cache value
        const setResult = await flow.setResultCache({
          value: buffer,
          ttl: {
            seconds: 30,
            nanos: 0,
          },
        });
        if (setResult.getError()) {
          console.log(`Error setting cache value: ${setResult.getError()}`);
        }
      }

      res.sendStatus(202);
    } else {
      flow.setStatus(FlowStatus.Error);
      res.sendStatus(403);
    }
  } catch (e) {
    console.log(e);
    res.status(500).send(`Error occurred: ${e}`);
  } finally {
    if (flow) {
      flow.end();
    }
  }
});

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
