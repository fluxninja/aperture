---
title: Scheduler
keywords:
  - scheduler
  - tokens
  - priority
  - queuing
  - actuators
sidebar_position: 7
---

:::info See Also

Scheduler [Reference](/reference/configuration/spec.md#scheduler)

:::

## Overview {#overview}

Schedulers provide a mechanism for prioritizing requests based on importance
when the service capacity is limited. Scheduling of requests is based on a
[weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
algorithm that ensures equitable resource allocation across workloads, factoring
in the [priorities](#priority) and [tokens](#tokens) (weight) of each request.

Service capacity limits can be determined based on one of the following
techniques:

1. [Quota Scheduling](./request-prioritization/quota-scheduler.md): Global
   [token buckets](https://en.wikipedia.org/wiki/Token_bucket) are used to track
   the request rate quota. The limit is based on a known limit, such as
   third-party API rate limits or inter-service API quotas.
2. [Concurrency Scheduling](./request-prioritization/load-scheduler.md): Global
   token counters are used to track the concurrency. The limit is set based on
   the concurrent processing capacity of the service.
3. [Load Scheduling](./request-prioritization/load-scheduler.md): Uses a token
   bucket local to each agent, which gets adjusted based on the past token rate
   at the agent. The limit is adjusted based on load at the service, such as
   CPU, queue length or response latency.

Each request to the scheduler seeks tokens from the underlying token bucket. If
tokens are available, the request gets admitted. If tokens are not readily
available, requests are queued, waiting either until tokens become accessible or
until a [timeout](#queue-timeout) occurs.

This diagram illustrates the working of a scheduler for workload prioritization.

![Scheduler](./assets/img/scheduler-light.svg#gh-light-mode-only)
![Scheduler](./assets/img/scheduler-dark.svg#gh-dark-mode-only)

### Workload {#workload}

Workloads group requests based on common [_Flow Labels_](./flow-label.md).
Developers can send the workload parameters as flow labels using Aperture SDKs.
The label keys used to identify workloads are configurable. See
[\*\_label_key parameters](/reference/configuration/spec.md#scheduler)

Alternately, a list of
[Workloads](/reference/configuration/spec.md#scheduler-workload) can be defined
inside the Scheduler specification using [label matcher][label-matcher] rules.

### Priority {#priority}

Priority represents the importance of a request compared to the other requests
in the queue. It varies from 0 to any positive number, indicating the urgency
level, with higher numbers denoting higher priority. The position of a flow in
the queue is computed based on its virtual finish time using the following
formula:

$$
inverted\_priority = {\frac {1} {priority}}
$$

$$
virtual\_finish\_time = virtual\_time + \left(tokens \cdot inverted\_priority\right)
$$

### Tokens {#tokens}

Tokens represent the cost for admitting a specific request. They can be defined
either through request labels or through workload definition inside a policy. If
not specified, the default token value is assumed to be 1 for each request, thus
representing the number of requests. Estimating tokens accurately for each
request helps fairer flow control decisions.

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
timeout duration, it is admitted. Otherwise, the timeout expires before the
tokens are available, the request is rejected. Thus, the timeout prevents
requests from waiting excessively long.

Developers can set the gRPC timeout on each `startFlow` call. In case of
middlewares, gRPC timeout is configured statically specific to each middleware
integration, e.g. through Envoy filter.

The timeout can also be configured using the `queue_timeout` parameter in the
[workload parameters](/reference/configuration/spec#scheduler-workload-parameters).
The smaller of the two timeouts is used.

[label-matcher]: ./selector.md#label-matcher
[flowcontrol-proto]:
  https://buf.build/fluxninja/aperture/docs/main:aperture.flowcontrol.check.v1
