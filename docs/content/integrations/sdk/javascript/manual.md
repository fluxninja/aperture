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

<a href={`https://www.npmjs.com/package/@fluxninja/aperture-js`}>Aperture
JavaScript SDK</a> can be used to manually set feature control points within a
JavaScript service.

To do so, first create an instance of ApertureClient. Agent host and port will
be read from environment variables `APERTURE_AGENT_HOST` and
`APERTURE_AGENT_PORT`, defaulting to localhost:8089.

```javascript
export const apertureClient = new ApertureClient();
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
