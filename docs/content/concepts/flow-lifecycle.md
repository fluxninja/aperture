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

This diagram elucidates the stages a flow undergoes within the Aperture. The
stages encompass:

- **Selectors**: Act as filters, determining the flow's path based on scoping
  rules. They identify and direct the flow towards relevant components in line
  with these rules.
- **Classifiers**: Responsible for categorizing the flow, utilizing the flow's
  metadata as their basis.
- **FluxMeters**: Critical tools that measure the flow's metrics. They convert
  fluxes into a Prometheus histogram format for an understandable data
  visualization.
- **Sampler**: Manages the flow's rate according to the service's health and
  capacity. The sampler must approve the flow before advancing it to the
  subsequent stage.
- **Rate-Limiters**: Proactively guard against recurrent overloads by regulating
  excessive requests in accordance with per-label limits.
- **Schedulers**: Govern the efficient processing of requests, favoring critical
  application features over background workloads. There are two types of
  schedulers, quota and load, which operate concurrently at the same stage.
  - The **Load Scheduler** manages the queue of flows before reaching the
    service, ensuring active service protection and controlling the incoming
    tokens per second.
  - The **Quota Scheduler** utilizes a global token bucket as a ledger to manage
    the distribution of tokens across all agents. It allows for strategic
    prioritization of requests when hitting quota limits, and is especially
    effective in environments with strict global rate limits.

Once the flow has been processed, the decision is sent back to the originator.
