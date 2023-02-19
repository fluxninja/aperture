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
[Concurrency Limiter reference](/reference/policies/spec.md#concurrency-limiter).

:::

Concurrency Limiter is about protecting your services from overload. Its goal is
to limit number of concurrent requests to service to a level the service can
handle. With the ability to define workloads of different priorities and
weights, it allows to shed some “less useful” requests, while minimally
impacting the "more important" ones.

Concurrency Limiter is configured as a [policy][policies] component.

## Scheduler {#scheduler}

Each Aperture Agent instantiates a
[Weighted Fair Queueing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
based scheduler as a way to prioritize requests based on their weights
(priority) and size(tokens). Concurrency Limiter applies a multiplier that the
scheduler uses to compute fill rate of [tokens](#tokens) per second, which it
tries to maintain.

If rate of tokens in requests entering the scheduler exceeds the desired rate,
requests are queued in the scheduler. If a flow can't be scheduled within its
specified timeout, it is rejected.

### Workload {#workload}

Workloads are groups of requests based on common attributes. Workloads are
expressed by [label matcher][label-matcher] rules in Aperture. Aperture Agents
schedule workloads based on their priorities and by estimating their
[tokens](#tokens).

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

Tokens represent the cost of admitting a flow in the system. Most commonly,
tokens are estimated based on milliseconds of response time observed when a flow
is processed. Token estimation of requests within a workload is crucial when
making flow control decisions. The concept of tokens is aligned with
[Little's Law](https://en.wikipedia.org/wiki/Little%27s_law), which defines the
relationship between response times, arrival rate and the number of requests
currently in the system (concurrency).

In some cases, tokens can be represented as the number of requests instead of
response times, e.g. when performing flow control on external APIs that have
hard rate-limits.

Aperture can be configured to automatically estimate the tokens for each
workload. See `auto-tokens`
[configuration](/reference/policies/spec.md#scheduler).

### Token bucket {#token-bucket}

Aperture Agents use a variant of a
[token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) to control
the requests entering the system. Each flow has to acquire tokens from the
bucket within a deadline period in order to be admitted.

### Timeout Factor {#timeout-factor}

The timeout factor parameter decides how long a request in the workload can wait
for tokens. This value impacts fairness because the larger the timeout the
higher the chance a request has to get scheduled.

The timeout is calculated as `timeout = timeout_factor * tokens`.

:::info

It's advisable to configure the timeouts in the same order of magnitude as the
normal latency of the workload requests in order to protect from retry storms
during overload scenarios.

:::

[label-matcher]: ../flow-selector.md#label-matcher
[policies]: /concepts/policy/policy.md
