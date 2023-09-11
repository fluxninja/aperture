---
title: Scheduler
keywords:
  - scheduler
  - tokens
  - priority
  - queuing
  - actuators
sidebar_position: 12
---

:::info See Also

Scheduler [Reference](/reference/configuration/spec.md#scheduler)

:::

## Overview {#overview}

Schedulers provide a mechanism for throttling and scheduling requests based on
importance when service resources are limited. The throttling is achieved
through [token buckets](https://en.wikipedia.org/wiki/Token_bucket). To gain
admittance, each request must obtain tokens from the bucket. When tokens are
depleted, incoming requests enter a queue, awaiting admittance based on a
[weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
algorithm. This algorithm ensures equitable resource allocation across
workloads, factoring in the priority and weight (tokens) of each request.

This diagram illustrates the working of a scheduler for workload prioritization.

![Scheduler](./assets/img/scheduler-light.svg#gh-light-mode-only)
![Scheduler](./assets/img/scheduler-dark.svg#gh-dark-mode-only)

Aperture offers two variants of scheduler:
[_Load Scheduler_](./load-scheduler.md) and
[_Quota Scheduler_](./quota-scheduler.md). While both use the same weighted fair
queuing-based scheduling algorithm, they differ in the throttling mechanism by
employing distinct types of token buckets. The _Load Scheduler_ uses a token
bucket local to each agent, which gets adjusted based on the past token rate at
the agent. This is useful for service protection scenarios since it provides a
robust mechanism to relatively adjust the token rate. The _Quota Scheduler_,
uses a centralized token bucket within an [agent group](../agent-group.md). This
is useful for scenarios involving known limits, like third-party API rate limits
or inter-service API quotas.

### Workload {#workload}

Workloads are groups of requests based on common
[_Flow Labels_](../flow-label.md). Workloads are expressed by [label
matcher][label-matcher] rules in the _Scheduler_ definition. Aperture Agents
schedule workloads based on their [priorities](#priority) and [tokens](#tokens).

### Priority {#priority}

Priority represents the importance of a request compared to the other requests
in the queue. It varies from 0 to an unlimited positive integer, indicating the
urgency level, with higher numbers denoting higher priority. The position of a
flow in the queue is computed based on its virtual finish time using the
following formula:

$$
inverted\_priority = {\frac {1} {priority}}
$$

$$
virtual\_finish\_time = virtual\_time + \left(tokens \cdot inverted\_priority\right)
$$

To manage prioritized requests, the scheduler seeks tokens from the token
bucket. If tokens are available, the request gets admitted. In cases where
tokens are not readily available, requests are queued, waiting either until
tokens become accessible or until a timeout occurs - the latter being dependent
on the workload or [`flowcontrol.v1.Check`][flowcontrol-proto] call timeout.

### Tokens {#tokens}

Tokens represent the cost for admitting a specific request. Typically, tokens
are based on the estimated response time of a request. Estimating the number of
tokens for each request within a workload is critical for making effective flow
control decisions.

Aperture can automatically estimate the tokens for each workload based on
historical latency measurements. See the `workload_latency_based_tokens`
[configuration](/reference/configuration/spec.md#load-scheduler-parameters) for
more details. The latency based token calculation is aligned with
[Little's Law](https://en.wikipedia.org/wiki/Little%27s_law), which relates
response times, arrival rate, and the system concurrency (number of in-flight
requests).

Alternatively, tokens can also be represented as the number of requests instead
of response times. For example, when scheduling access to external APIs that
have strict rate limits (global quota). In this case, the number of tokens
represents the number of requests that can be made to the API within a given
time window.

Tokens are determined in the following order of precedence:

- Specified in the flow labels.
- Estimated tokens (see
  [`workload_latency_based_tokens`](/reference/configuration/spec.md#load-scheduler)
  setting).
- Specified in the `Workload.tokens` setting.

### Queue Timeout {#queue-timeout}

The queue timeout is determined by the gRPC timeout provided on the
[`flowcontrol.v1.Check`][flowcontrol-proto] call. When a request is made, it
includes a timeout value that specifies the maximum duration the request can
wait in the queue. If the request receives the necessary tokens within this
timeout duration, it is admitted. Otherwise, if the timeout expires before the
tokens are available, the request is rejected.

The gRPC timeout on the [`flowcontrol.v1.Check`][flowcontrol-proto] call is set
in the Envoy filter and the SDK during initialization. It serves as an upper
bound on the queue timeout, preventing requests from waiting excessively long.

[label-matcher]: ../selector.md#label-matcher
[flowcontrol-proto]:
  https://buf.build/fluxninja/aperture/docs/main:aperture.flowcontrol.check.v1
