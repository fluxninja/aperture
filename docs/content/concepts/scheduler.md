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

The capacity of a service is divided among the workloads based on their
priorities. The higher the priority, the larger the share of its capacity. Note
that priorities are relative to each other and not absolute. For example, if
there are two workloads with priorities 1 and 2, the second workload gets twice
the capacity of the first workload. If there are three workloads with priorities
1, 2 and 3, the third workload gets three times the capacity of the first
workload.

If a certain workload does not have enough requests, the capacity is shared
among the other workloads. There is no upfront allocation of capacity to
workloads.

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

Tokens determine the relative allocation of capacity among requests. For
example, a request with 4 tokens is the same as 4 requests with 1 token each.

### Queue Timeout {#queue-timeout}

The queue timeout is determined by the gRPC timeout provided on the
[`flowcontrol.v1.Check`][flowcontrol-proto] call. When a request is made, it
includes a timeout value that specifies the maximum duration the request can
wait in the queue. If the request receives the necessary tokens within this
timeout duration, it is admitted. Otherwise, the timeout expires before the
tokens are available, the request is rejected. Thus, the timeout prevents
requests from waiting excessively long.

Developers can set the gRPC timeout on each `startFlow` call. In the case of
middlewares, gRPC timeout is configured statically specific to each middleware
integration, e.g. through Envoy filter.

The timeout can also be configured using the `queue_timeout` parameter in the
[workload parameters](/reference/configuration/spec#scheduler-workload-parameters).
The smaller of the two timeouts is used.

### Fairness {#fairness}

The requests in a workload at the same priority level are served in a
first-in-first-out manner. But that may not be desirable in multi-tenant
environments, where a few users or tenants may take up most of the capacity.

Aperture's scheduler can ensure fair usage among users within a workload using a
stochastic fair queuing strategy. Developers can send the fairness key
identifying users or tenants in their app as a flow label. The label for
fairness key is defined in the
[Scheduler specification](/reference/configuration/spec#scheduler) as
`fairness_label_key`.

Note that priorities determine relative allocation of capacity among workloads.
Fairness determines relative allocation of capacity among users within a
workload.

[label-matcher]: ./selector.md#label-matcher
[flowcontrol-proto]:
  https://buf.build/fluxninja/aperture/docs/main:aperture.flowcontrol.check.v1
