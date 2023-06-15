---
title: Flow Lifecycle
sidebar_label: Flow Lifecycle
sidebar_position: 2
keywords:
  - flows
  - services
  - discovery
  - labels
---

## Flow Lifecycle

To better understand how flows work in Aperture, it is essential to understand
the lifecycle of a flow. The following diagram shows the lifecycle of a flow in
Aperture.

<Zoom>

```mermaid
{@include: ./assets/gen/service/flow-lifecycle.mmd}
```

</Zoom>

This diagram shows the steps that a flow goes through in Aperture. The steps
are:

- Selectors: These define scoping rules, identifying and forwarding the flow to
  the relevant components based on the rules.
- Classifiers: used to classify the flow based on the flow's metadata.
- FluxMeters: These are critical instruments that quantify the flow's metrics,
  translating fluxes into a Prometheus histogram for clear data visualization.
- Sampler: It regulates the flow's rate, based on service health and capacity,
  and need to accept the flow before forwarding it to the next step.
- Rate-Limiters: They proactively mitigate recurring overloads by regulating
  heavy-hitters according to per-label limits.
- Schedulers: These ensure efficient handling of requests, prioritizing critical
  application features over background workloads.

Note that all components have the authority to reject a flow, which will stop
the flow from being processed further.

Once the flow has been processed, the decision is sent back to the originator.

<!-- vale on -->
