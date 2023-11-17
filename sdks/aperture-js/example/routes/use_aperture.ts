import express from "express";
import { ApertureClient, Flow, FlowStatusEnum } from "@fluxninja/aperture-js";
import grpc from "@grpc/grpc-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address:
    process.env.APERTURE_AGENT_ADDRESS !== undefined
      ? process.env.APERTURE_AGENT_ADDRESS
      : "localhost:8080",
  agentAPIKey: process.env.APERTURE_AGENT_API_KEY || undefined,
  // if process.env.APERTURE_AGENT_INSECURE set channelCredentials to insecure
  channelCredentials:
    process.env.APERTURE_AGENT_INSECURE !== undefined
      ? grpc.credentials.createInsecure()
      : grpc.credentials.createSsl(),
});

export const apertureRoute = express.Router();
apertureRoute.get("/", function (_: express.Request, res: express.Response) {
  // do some business logic to collect labels
  const labels: Record<string, string> = {
    user: "kenobi",
  };

  const startTimestamp = Date.now();

  // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
  apertureClient
    .StartFlow("awesomeFeature", {
      labels: labels,
      grpcCallOptions: {
        deadline: Date.now() + 30000,
      },
      rampMode: false,
      cacheKey: "cache",
    })
    .then(async (flow: Flow) => {
      const endTimestamp = Date.now();
      console.log(`Flow took ${endTimestamp - startTimestamp}ms`);
      // See whether flow was accepted by Aperture Agent.
      if (flow.ShouldRun()) {
        // Simulate work being done
        sleep(200).then(() => {
          console.log("Work done!");
          res.sendStatus(202);
        });
      } else {
        // Flow has been rejected by Aperture Agent.
        flow.SetStatus(FlowStatusEnum.Error);
        res.sendStatus(403);
      }
      // create a new buffer
      const buffer = Buffer.from("awesomeString");

      // set a cache value
      const setResult = await flow.SetCachedValue(buffer, { seconds: 30, nanos: 0 })
      if (setResult?.error) {
        console.log(`Error setting cache value: ${setResult.error}`);
      }

      const deleteResult = await flow.DeleteCachedValue()
      if (deleteResult?.error) {
        console.log(`Error deleting cache value: ${deleteResult.error}`);
      }

      // Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
      // Status set using SetStatus() informs whether the feature captured by the Flow was successful or resulted in an error.
      flow.End();
    })
    .catch((e: any) => {
      console.log(e);
      res.status(500).send(`Error occurred: ${e}`);
    });
});

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
