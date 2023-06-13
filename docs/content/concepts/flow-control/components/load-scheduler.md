---
title: Load Scheduler
keywords:
  - scheduler
  - tokens
  - priority
  - queuing
  - actuators
sidebar_position: 6
---

:::info

See also
[_Adaptive Load Scheduler_ reference](/reference/policies/spec.md#adaptive-load-scheduler)
and
[_Load Scheduler_ reference](/reference/policies/spec.md#adaptive-load-scheduler).

:::

The _Load Scheduler_ is a powerful tool designed to protect your services from
overloads and prevent cascading failures. Its **primary goal is to limit** the
number of concurrent requests to a service, ensuring that the service can handle
the incoming workload. By defining workloads with different priorities and
weights, the _Load Scheduler_ can prioritize certain requests over others,
enabling graceful degradation of service during times of high traffic.

As with other components of the Aperture platform, the _Load Scheduler_ is
configured using a [policy][policies] component.

## Scheduler {#scheduler}

The scheduler prioritizes requests based on their priority and size. Each
Aperture Agent instantiates a
[Weighted Fair Queue-based](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
scheduler as a way to prioritize requests. The controller applies a _Load
Multiplier_ that the scheduler uses to compute the refill rate of
[tokens](#tokens) per second, which it tries to maintain between each update
from the controller.

If the rate of tokens in requests entering the scheduler exceeds the desired
rate, requests are queued in the scheduler. If a flow cannot be scheduled within
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

Priority represents the importance of a request compared to the other requests
in the queue.

:::note

Priority levels are in the range `0 to 255`. `0` is the lowest priority and
`255` is the highest priority.

Priority levels have a non-linear effect on the scheduler. The following formula
is used to compute the position in the queue based on the concept of
[virtual finish times](https://en.wikipedia.org/wiki/Weighted_fair_queueing#Algorithm):

`virtual_finish_time = virtual_time + (tokens * (256 - priority))`

This means that requests with a priority level `255` will see double the
acceptance rate compared to requests with a priority level `254`, if they have
the same number of tokens.

:::

### Tokens {#tokens}

Tokens represent the unit of cost for accepting a certain flow. Typically,
tokens are based on the estimated response time of a flow. Estimating the number
of tokens for each request within a workload is critical for making effective
flow control decisions. The concept of tokens is aligned with
[Little's Law](https://en.wikipedia.org/wiki/Little%27s_law), which relates
response times, arrival rate, and the number of requests in the system
(concurrency). Aperture can automatically estimate the tokens for each workload
based on historical latency measurements. See the
`workload_latency_based_tokens`
[configuration](/reference/policies/spec.md#load-scheduler) for more details.

Alternatively, tokens can also be represented as the number of requests instead
of response times. For example, when scheduling access to external APIs that
have strict rate limits (global quota). In this case, the number of tokens
represents the number of requests that can be made to the API within a given
time window.

Tokens are determined in the following order of precedence:

- Specified in the flow labels.
- Specified in the `Workload.tokens` setting.
- Estimated tokens (see
  [`workload_latency_based_tokens`](/reference/policies/spec.md#load-scheduler)
  setting).

### Token rate {#token-rate}

The Scheduler provides token rates for both incoming and accepted (admitted)
requests as output signals.

When tokens represent the number of requests, with each request counting as 1
token, the token rate is simply the number of requests per second.

However, when using the auto-tokens setting, tokens correspond to seconds of
response latency (work-seconds). In this case, the token rate represents the
work-seconds completed per unit of time (that is, work-seconds per second),
which is a measure of the system's concurrency. Concurrency is a dimensionless
metric that indicates the average number of flows being processed concurrently
by the system.

### Token bucket {#token-bucket}

The Aperture Agents use a modified version of the
[token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) to regulate
the flow of incoming requests. In this algorithm, every flow is required to
obtain tokens from the bucket before a specified deadline to gain admission to
the system.

### Timeout Factor {#timeout-factor}

The timeout factor parameter determines the duration a request in the workload
can wait for tokens. A larger timeout factor results in a higher chance of the
request being scheduled, improving fairness. The timeout is computed as
`timeout = timeout_factor * tokens`.

:::info

It's recommended to configure the timeouts to be in the same order of magnitude
as the normal latency of the workload requests. This helps prevent retry storms
during overload scenarios.

:::

[label-matcher]: ../selector.md#label-matcher
[policies]: /concepts/policy/policy.md
