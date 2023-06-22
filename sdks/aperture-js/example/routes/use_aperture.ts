import express from "express";

import { ApertureClient } from "../../sdk/client.js";
import { FlowStatus } from "../../sdk/flow.js";

// Create aperture client
export const apertureClient = new ApertureClient();

export const apertureRoute = express.Router();
apertureRoute.get("/", function (_: express.Request, res: express.Response) {
  // do some business logic to collect labels
  var labelsMap = new Map<string, string>().set("user", "kenobi");

  // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
  apertureClient
    .StartFlow("awesome-feature", labelsMap)
    .then((flow) => {
      // See whether flow was accepted by Aperture Agent.
      if (flow.ShouldRun()) {
        // Simulate work being done
        sleep(2000).then(() => {
          console.log("Work done!");
        });
        res.sendStatus(202);
      } else {
        // Flow has been rejected by Aperture Agent.
        flow.SetStatus(FlowStatus.Error);
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
