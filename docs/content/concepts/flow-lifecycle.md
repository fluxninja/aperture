---
title: Flow Lifecycle
sidebar_label: Flow Lifecycle
sidebar_position: 5
keywords:
  - flows
  - services
  - discovery
  - labels
---

## Flow Lifecycle

The lifecycle of a flow begins when a service initiates it, requesting a
decision from the Aperture Agent. As the flow enters the Aperture Agent, it
embarks on a journey through multiple stages before a final decision is made.
The following diagram illustrates these stages and the sequence they follow.

While some stages primarily serve to augment the flow with additional labels or
provide telemetry, others actively participate in the decision to accept or
reject the flow. Should a flow be rejected at any such active stage, its journey
through the subsequent stages is halted, and an immediate decision is dispatched
back to the service.

<Zoom>

```mermaid
{@include: ./assets/gen/service/flow-lifecycle.mmd}
```

</Zoom>

Outlined below is the sequence of stages in detail, along with their respective
roles:

:::note

Remember, a flow may bypass certain or all stages if there are no matching
components for that stage.

:::

- **Selectors** are the criteria used to determine the components that will be
  applied to a flow in the subsequent stages.
- **Classifiers** perform the task of assigning labels to the flow based on the
  HTTP attributes of the request. However, classifiers are only pertinent for
  HTTP or gRPC _Control Points_ and do not apply to flows associated with
  feature _Control Points_.
- **FluxMeters** are employed to meter the flows, generating metrics like
  latency, request counts, or other arbitrary telemetry based on access logs.
  They transform request flux that matches certain criteria into Prometheus
  histograms, enabling enhanced observability and control.
- **Samplers** manage load by permitting a portion of flows to be accepted,
  while immediately dropping the remainder with a forbidden status code. They
  are particularly useful in scenarios such as feature rollouts.
- **Rate-Limiters** proactively guard against abuse by regulating excessive
  requests in accordance with per-label limits.
- **Schedulers** offer on-demand queuing based on a token bucket algorithm, and
  prioritize requests using weighted fair queuing. Multiple matching schedulers
  can evaluate concurrently, with each having the power to drop a flow. There
  are two variants:
  - The **Load Scheduler** oversees the current token rate in relation to the
    past token rate, adjusting as required based on health signals from a
    service. This scheduler type facilitates active service protection.
  - The **Quota Scheduler** uses a global token bucket as a ledger, managing the
    token distribution across all agents. It proves especially effective in
    environments with strict global rate limits, as it allows for strategic
    prioritization of requests when reaching quota limits.

After traversing these stages, the flow's decision is sent back to the
initiating service.
