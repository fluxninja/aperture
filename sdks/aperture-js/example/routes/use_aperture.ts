import express from "express";

import { ApertureClient, FlowStatusEnum } from "@fluxninja/aperture-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: process.env.APERTURE_AGENT_ADDRESS !== undefined ? process.env.APERTURE_AGENT_ADDRESS : "localhost:8089",
  agentAPIKey: process.env.APERTURE_AGENT_API_KEY || undefined,
  isInsecure: process.env.APERTURE_AGENT_INSECURE === "true",
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
    })
    .then((flow) => {
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
      // Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
      // Status set using SetStatus() informs whether the feature captured by the Flow was successful or resulted in an error.
      flow.End();
    })
    .catch((e) => {
      console.log(e);
      res.send(`Error occurred: ${e}`);
    });
});

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
