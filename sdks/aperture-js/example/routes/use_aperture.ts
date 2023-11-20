import express from "express";
import { ApertureClient, Flow, FlowStatusEnum, LookupStatus } from "@fluxninja/aperture-js";
import grpc from "@grpc/grpc-js";


// Create aperture client
export const apertureClient = new ApertureClient({
  address:
    process.env.APERTURE_AGENT_ADDRESS !== undefined
      ? process.env.APERTURE_AGENT_ADDRESS
      : "localhost:8089",
  agentAPIKey: process.env.APERTURE_AGENT_API_KEY || undefined,
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
    flow = await apertureClient.StartFlow("awesomeFeature", {
      labels: labels,
      grpcCallOptions: {
        deadline: Date.now() + 30000,
      },
      rampMode: false,
      cacheKey: "cache",
    });

    const endTimestamp = Date.now();
    console.log(`Flow took ${endTimestamp - startTimestamp}ms`);

    if (flow.ShouldRun()) {
      await sleep(200);

      if (flow.CachedValue().GetLookupStatus() === LookupStatus.Hit) {
        console.log("Cache hit:", flow.CachedValue().GetValue()?.toString());
      } else {
        console.log("Cache miss:", flow.CachedValue().GetOperationStatus(), flow.CachedValue().GetError()?.message);
        const resString = "awesomeString";

        // create a new buffer
        const buffer = Buffer.from(resString);

        // set cache value
        const setResult = await flow.SetCachedValue(buffer, { seconds: 30, nanos: 0 })
        if (setResult?.error) {
          console.log(`Error setting cache value: ${setResult.error}`);
        }
      }

      res.sendStatus(202);
    } else {
      flow.SetStatus(FlowStatusEnum.Error);
      res.sendStatus(403);
    }
  } catch (e) {
    console.log(e);
    res.status(500).send(`Error occurred: ${e}`);
  } finally {
    if (flow) {
      flow.End();
    }
  }
});

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
