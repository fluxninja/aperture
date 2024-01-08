---
title: Flow Lifecycle
sidebar_label: Flow Lifecycle
sidebar_position: 3
keywords:
  - flows
  - services
  - discovery
  - labels
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

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
{@include: ./assets/gen/flow-lifecycle/flow-lifecycle.mmd}
```

</Zoom>

Outlined below is the sequence of stages in detail, along with their roles:

:::note

Remember, a flow can bypass certain or all stages if there are no matching
components for that stage.

:::

### Selection, Classification, and Telemetry

- [**Selectors**](./selector.md) are the criteria used to determine the
  components that will be applied to a flow in the subsequent stages.
- [**Classifiers**](./advanced/classifier.md/) perform the task of assigning
  labels to the flow based on the HTTP attributes of the request. However,
  classifiers are only pertinent for HTTP or gRPC _Control Points_ and do not
  apply to flows associated with feature _Control Points_.
- [**FluxMeters**](./advanced/flux-meter.md) are employed to meter the flows,
  generating metrics such as latency, request counts, or other arbitrary
  telemetry based on access logs. They transform request flux that matches
  certain criteria into Prometheus histograms, enabling enhanced observability
  and control.

### Rate limiting (fast rejection)

- [**Samplers**](./advanced/load-ramp.md#sampler) manage load by permitting a
  portion of flows to be accepted, while immediately dropping the remainder with
  a forbidden status code. They are particularly useful in scenarios such as
  feature rollouts.
- [**Rate-Limiters**](./rate-limiter.md) proactively guard against abuse by
  regulating excessive requests in accordance with per-label limits.
- [**Concurrency-Limiters**](./concurrency-limiter.md) enforce in-flight request
  quotas to prevent overloads. They can also be used to enforce limits per
  entity such as a user to ensure fair access across users.

### Request Prioritization and Cache Lookup

[**Schedulers**](./scheduler.md) offer on-demand queuing based on a limit
enforced through a token bucket or a concurrency counter, and prioritize
requests using weighted fair queuing. Multiple matching schedulers can evaluate
concurrently, with each having the power to drop a flow. There are three
variants running at various stages of the flow lifecycle:

- The
  [**Concurrency Scheduler**](./request-prioritization/concurrency-scheduler.md)
  uses a global concurrency counter as a ledger, managing the concurrency across
  all Agents. It proves especially effective in environments with strict global
  concurrency limits, as it allows for strategic prioritization of requests when
  reaching concurrency limits.
- [**Caches**](./cache.md) Look of response and global caches occur at this
  stage. If a response cache hit occurs, the flow is not sent to the Concurrency
  and Load Scheduling stages, resulting in an early acceptance.
- The [**Quota Scheduler**](./request-prioritization/quota-scheduler.md) uses a
  global token bucket as a ledger, managing the token distribution across all
  Agents. It proves especially effective in environments with strict global rate
  limits, as it allows for strategic prioritization of requests when reaching
  quota limits.
- The [**Load Scheduler**](./request-prioritization/load-scheduler.md) oversees
  the current token rate in relation to the past token rate, adjusting as
  required based on health signals from a service. This scheduler type
  facilitates active service protection.

After traversing these stages, the flow's decision is sent back to the
initiating service.
