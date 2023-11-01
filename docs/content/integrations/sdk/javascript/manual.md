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
[Agent API Keys](/get-started/aperture-cloud/agent-api-keys.md).

:::

```javascript
export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  agentAPIKey: "AGENT_API_KEY",
});
```

The created instance can then be used to start a flow:

```javascript
// do some business logic to collect labels
var labelsMap = new Map().set("key", "value");
var rampMode = false;

apertureClient
  .StartFlow("feature-name", labelsMap, rampMode)
  .then((flow) => {
    if (flow.ShouldRun()) {
      // Do actual work
    } else {
      // handle flow rejection by Aperture Agent
      flow.SetStatus(FlowStatus.Error);
    }
    flow.End();
  })
  .catch((e) => {
    console.log(e);
    res.send(`Error occurred: ${e}`);
  });
```

For more context on using the Aperture JavaScript SDK to set feature control
points, refer to the [example app][example] available in the repository.

[example]: https://github.com/fluxninja/aperture-js/tree/main/example
