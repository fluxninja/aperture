import express from "express";

import { ApertureClient, FlowStatusEnum } from "@fluxninja/aperture-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: "localhost:8089",
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
    .StartFlow("awesome-feature", {
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
