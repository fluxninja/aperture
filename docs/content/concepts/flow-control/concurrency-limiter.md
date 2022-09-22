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
[Concurrency Limiter reference](/references/configuration/policies.md#v1-concurrency-limiter).

:::

Concurrency Limiter is about protecting your services from overload. Its goal is
to limit number of concurrent requests to service to a level the service can
handle. It's implemented by configurable load-shedding. Thanks to the ability to
define workloads of different priorities and weights, it allows to shed some
“less useful” flows, while not affecting the more important ones.

Concurrency Limiter is configured as a [policy][policies] component.

## Scheduler

Each Aperture Agent instantiates a Scheduler as a way to run the Concurrency
Limiter. Concurrency Limiter applies a load-shed-factor that the Scheduler uses
to compute a level of [tokens](#tokens) per second, which it tries to maintain.

If rate of tokens in flows entering the scheduler exceeds the desired rate,
flows are queued in the scheduler. If a flow can't be scheduled within its
specified timeout, it will be rejected.

### Workload

Workloads are groups of flows based on common attributes. Workloads are
expressed by [label matcher][label-matcher] rules in Aperture. Aperture Agents
schedule workloads based on their priorities and by estimating their
[tokens](#tokens).

### Tokens

Tokens represent the cost of admitting a flow in the system. Most commonly,
tokens are estimated based on milliseconds of response time observed when a flow
is processed. Token estimation of flows within a workload is crucial when making
flow control decisions. The concept of tokens is aligned with
[Little's Law](https://en.wikipedia.org/wiki/Little%27s_law), which defines the
relationship between response times, arrival rate and the number of requests
currently in the system (concurrency).

In some cases, tokens can be represented as the number of requests instead of
response times, e.g. when performing flow control on external APIs that have
hard rate-limits.

### Token bucket

Aperture Agents use a variant of a
[token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) to control
the flows entering the system. Each flow has to acquire tokens from the bucket
within a deadline period in order to be admitted.

### Timeouts

The timeout parameter decides how long a request in the workload can wait for
tokens. This value impacts fairness because the larger the timeout the higher
the chance a request has to get scheduled.

[label-matcher]: /concepts/flow-control/selector.md#label-matcher
[policies]: /concepts/policy/policy.md
