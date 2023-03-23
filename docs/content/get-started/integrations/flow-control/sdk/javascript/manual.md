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

<a href={`https://www.npmjs.com/package/aperture-sdk`}>Aperture Javascript
SDK</a> can be used to manually set feature Control Points within a Javascript
service.

To do so, first create an instance of ApertureClient. Agent host and port will
be read from env variables `FN_AGENT_HOST` and `FN_AGENT_PORT`, defaulting to
localhost:8089.

```javascript
export const apertureClient = new ApertureClient();
```

The created instance can then be used to start a flow:

```javascript
// do some business logic to collect labels
var labelsMap = new Map().set("key", "value");

apertureClient
  .StartFlow("feature-name", labelsMap)
  .then((flow) => {
    if (flow.Accepted()) {
      // Do actual work
      flow.End(FlowStatus.Ok);
    } else {
      // handle flow rejection by Aperture Agent
      flow.End(FlowStatus.Error);
    }
  })
  .catch((e) => {
    console.log(e);
    res.send(`Error occurred: ${e}`);
  });
```

For more context on how to use Aperture Javascript SDK to set feature Control
Points, you can take a look at the [example app][example] available in our
repository.

[example]: https://github.com/fluxninja/aperture-js/tree/main/example
