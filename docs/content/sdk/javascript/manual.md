---
title: Manually setting feature control points
slug: manually-setting-feature-control-points-using-javascript-sdk
sidebar_position: 1
keywords:
  - js
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

[Aperture JavaScript SDK](https://www.npmjs.com/package/@fluxninja/aperture-js)
can be used to manually set feature control points within a JavaScript service.

To do so, first create an instance of ApertureClient:

:::info Agent API Key

You can create an Agent API key for your project in the Aperture Cloud UI. For
more information, refer to
[Define Control Points](/get-started/define-control-points.md).

:::

```javascript
import { ApertureClient, Flow, FlowStatusEnum } from "@fluxninja/aperture-js";

export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  agentAPIKey: "AGENT_API_KEY",
});
```

The created instance can then be used to start a flow:

```javascript
async function handleRequest(req, res) {
  const flow = await apertureClient.StartFlow("feature-name", {
    labels: {
      label_key: "some_user_id",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300, // ms
    },
  });

  if (flow.ShouldRun()) {
    // Do Actual Work
  } else {
    // Handle flow rejection
    flow.SetStatus(FlowStatusEnum.Error);
  }

  if (flow) {
    flow.End();
  }
}
```

For more context on using the Aperture JavaScript SDK to set feature control
points, refer to the [example app][example] available in the repository.

[example]: https://github.com/fluxninja/aperture-js/tree/main/example
