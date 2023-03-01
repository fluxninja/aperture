---
title: Concurrency Limiter
keywords:
  - scheduler
  - tokens
  - priority
  - fairness
  - queuing
  - actuators
sidebar_position: 6
---

:::info

See also
[AIMD Concurrency Controller reference](/reference/policies/spec.md#a-i-m-d-concurrency-controller)
and
[Concurrency Limiter reference](/reference/policies/spec.md#concurrency-limiter).

:::

The Concurrency Limiter is a powerful tool designed to protect your services
from overloads and prevent cascading failures. Its primary goal is to limit the
number of concurrent requests to a service, ensuring that the service can handle
the incoming workload. By defining workloads with different priorities and
weights, the Concurrency Limiter can prioritize certain requests over others,
enabling graceful degradation of service during times of high traffic.

As with other components of the Aperture platform, the Concurrency Limiter is
configured using a [policy][policies] component.

## Scheduler {#scheduler}

Scheduler prioritizes requests based on their priority and size. Each Aperture
Agent instantiates a
[Weighted Fair Queueing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
based scheduler as a way to prioritize requests. The controller applies a _Load
Multiplier_ that the scheduler uses to compute the fill rate of
[tokens](#tokens) per second, which it tries to maintain between each update
from the controller.

If the rate of tokens in requests entering the scheduler exceeds the desired
rate, requests are queued in the scheduler. If a flow can't be scheduled within
its specified timeout, it is rejected.

The scheduler helps ensure that requests are handled in a fair and efficient
manner, even during periods of high load or overload. By prioritizing critical
application features over background workloads, the scheduler helps maximize
user experience or revenue.

### Workload {#workload}

Workloads are groups of flows based on common [_Flow Labels_](../flow-label.md).
Workloads are expressed by [label matcher][label-matcher] rules in Aperture.
Aperture Agents schedule workloads based on their priorities and by (auto)
estimating their [tokens](#tokens).

### Priority {#priority}

Priority represents the importance of a request with respect to other requests
in the queue.

:::note

Priority levels are in the range `0 to 255`. `0` is the lowest priority and
`255` is the highest priority.

Priority levels have non-linear effect on the scheduler. The following formula
is used to compute the position in the queue based on the concept of (virtual
finish times)[https://en.wikipedia.org/wiki/Weighted_fair_queueing#Algorithm]:

`virtual_finish_time = virtual_time + (tokens * (256 - priority))`

This means that requests with a priority level `255` will see double the
acceptance rate compared to requests with a priority level `254`, if they have
same number of tokens.

:::

### Tokens {#tokens}

Tokens represent the unit of cost for processing a flow in the system.
Typically, tokens are based on the estimated response time of a flow. Estimating
the number of tokens for each request within a workload is critical for making
effective flow control decisions. The concept of tokens is aligned with
[Little's Law](https://en.wikipedia.org/wiki/Little%27s_law), which relates
response times, arrival rate, and the number of requests in the system
(concurrency).

In certain cases, tokens can be represented as the number of requests instead of
response times. For example, when applying flow control to external APIs that
have strict rate limits.

By default, Aperture can automatically estimate the tokens for each workload.
See the `auto_tokens` [configuration](/reference/policies/spec.md#scheduler)
configuration for more details.

### Token bucket {#token-bucket}

The Aperture Agents utilize a modified version of the
[token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) to regulate
the flow of incoming requests. In this algorithm, every flow is required to
obtain tokens from the bucket before a specified deadline in order to gain
admission to the system.

### Timeout Factor {#timeout-factor}

The timeout factor parameter determines the duration a request in the workload
can wait for tokens. A larger timeout factor results in a higher chance of the
request being scheduled, improving fairness. The timeout is computed as
`timeout = timeout_factor * tokens`.

:::info

It is recommended to configure the timeouts to be in the same order of magnitude
as the normal latency of the workload requests. This helps prevent retry storms
during overload scenarios.

:::

[label-matcher]: ../flow-selector.md#label-matcher
[policies]: /concepts/policy/policy.md
