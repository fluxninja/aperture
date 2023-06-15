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

The Scheduler is a crucial component in Aperture, prioritizing requests based on
their priority and size. It operates based on a
[Weighted Fair Queue-based](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
methodology, which ensures the fair handling of workloads, particularly during
periods of high or overload traffic. By prioritizing critical application
features over low priority workloads, the scheduler helps optimize user
experience and maximize good throughput.

Each Agent instantiates an independent copy of the _Scheduler_, but output
signals for accepted and incoming token rate are aggregated across all agents.
The Scheduler utilizes a Load Multiplier to determine the [tokens](#tokens)
refill rate per second, striving to maintain this rate stable between each
controller update. If the incoming token rate in requests surpasses the desired
rate, requests are queued within the Scheduler. Any request that cannot be
scheduled within its designated timeout is rejected.

As with other components of the Aperture platform, the _Load Scheduler_ is
configured using a [policy][policies] component.

:::info

Additional information can be found in the
[Scheduler Reference](/reference/configuration/spec.md#scheduler)

:::

### Workload {#workload}

Workloads are groups of flows based on common [_Flow Labels_](./flow-label.md).
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

![Token Bucket](./assets/img/token-bucket-light.svg#gh-light-mode-only)
![Token Bucket](./assets/img/token-bucket-dark.svg#gh-dark-mode-only)

The diagram at hand illustrates the use of the token bucket algorithm for
workload prioritization by the scheduler. It demonstrates how the Aperture
Controller periodically adjusts the token bucket's load multiplier and refill
rate by broadcasting signals to the Aperture Agent.

In the scope of the Aperture Agent, the scheduler efficiently handles incoming
requests, categorizing them according to their urgency. This classification
scale ranges from 0 to 255, where 0 indicates the lowest and 255 the highest
priority.

To manage a prioritized request, the scheduler obtains tokens from the token
bucket, which either returns tokens, if available, or returns an await signal.
Then, tokens are used to schedule and route requests to the appropriate service.
Any requests exceeding the available token limit are strategically discarded to
maintain optimal system performance.

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
[configuration](/reference/configuration/spec.md#load-scheduler) for more
details.

Alternatively, tokens can also be represented as the number of requests instead
of response times. For example, when scheduling access to external APIs that
have strict rate limits (global quota). In this case, the number of tokens
represents the number of requests that can be made to the API within a given
time window.

Tokens are determined in the following order of precedence:

- Specified in the flow labels.
- Specified in the `Workload.tokens` setting.
- Estimated tokens (see
  [`workload_latency_based_tokens`](/reference/configuration/spec.md#load-scheduler)
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

The Scheduler employs two distinct types of token buckets: the Dynamic Token
Bucket and the Fixed Token Bucket. These buckets, each with a unique method of
operation, enhance the system's capacity to effectively adapt to various
scheduling scenarios.

#### Dynamic Token Bucket {#dynamic-bucket}

The Dynamic Token Bucket is primarily used for service protection. It adjusts
the token rate in response to the changing load signal, thereby ensuring an
active defense for the service. Each agent maintains an individual token bucket,
and synchronization between agents happens in lazy sync. This dynamic resizing
enables efficient management of requests and maintains optimal service
protection levels.

#### Fixed Token Bucket {#fixed-bucket}

In contrast to the Dynamic Token Bucket, the Fixed Token Bucket operates with a
predetermined number of tokens, establishing a constant token rate. This type of
bucket is particularly advantageous when the load or the number of tokens is
known upfront.

The Fixed Token Bucket can be employed for service protection but is especially
useful in scenarios where there are strict rate limits. For instance, when
scheduling access to external APIs with a global quota, the Fixed Token Bucket
ensures that the API's rate limits are adhered to.

Both the Dynamic and Fixed Token Buckets serve to manage and control the flow of
requests, their specific usage and operation depend on the requirements and
constraints of the system, providing versatility and adaptability in handling
workloads.

### Queue timeout {#queue-factor}

The timeout factor parameter determines the duration a request in the workload
can wait for tokens. A larger timeout factor results in a higher chance of the
request being scheduled, improving fairness. The timeout is computed as
`timeout = timeout_factor * tokens`.

:::info

It's recommended to configure the timeouts to be in the same order of magnitude
as the normal latency of the workload requests. This helps prevent retry storms
during overload scenarios.

:::

[label-matcher]: ./selector.md#label-matcher
[policies]: /concepts/advanced/policy.md
