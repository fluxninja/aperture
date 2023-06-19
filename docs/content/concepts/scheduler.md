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

:::info

Additional information can be found in the
[Scheduler Reference](/reference/configuration/spec.md#scheduler)

:::

### Load Scheduler {#load-scheduler}

The _Load Scheduler_ is a primary actuation component layered atop the base
Scheduler. It provides active service protection by forming a queue for flows
before they reach the service. Through provided match rules, flows are organized
into Workloads, each determining the priority and cost of admitting its
associated flows. The _Load Scheduler_ controls the incoming tokens per second,
which, in concurrency limiting cases, can be interpreted as (avg. latency \*
in-flight requests) according to Little's Law. The proportion of incoming tokens
that are admitted is dictated by the signal at the `load_multiplier` port.

:::info See Also

Load Scheduler [Reference](/reference/configuration/spec.md#load-scheduler)

:::

### Adaptive Load Scheduler {#adaptive-load-scheduler}

Building upon the foundation of the _Load Scheduler_, the _Adaptive Load
Scheduler_ fine-tunes the request management process. It incorporates a gradient
controller and integrator for computing the load multiplier, which adjusts based
on a signal, setpoint, and overload signals. The key functionality of the
_Adaptive Load Scheduler_ lies in its ability to adjust the accepted token rate
based on the deviation of the input signal from the setpoint. This adaptability
allows for greater control and precision in managing load.

:::info See Also

Adaptive Load Scheduler
[Reference](/reference/configuration/spec.md#adaptive-load-scheduler)

:::

### Quota Scheduler {#quota-scheduler}

The Quota Scheduler uses a global-level token bucket that functions as a ledger.
This ledger maintains an account of the total available tokens that can be
allocated across all agents. In a scenario where the total load or the total
number of tokens is known upfront, the Quota Scheduler can be particularly
effective. These Tokens can be thought of as a fixed quota that is distributed
among the agents. Each agent has access to this global ledger and deducts tokens
from it when processing requests. If the ledger runs out of tokens, new requests
are queued or rejected until more tokens become available. The Quota Scheduler
is useful in scenarios where there are strict global rate limits, for example,
when scheduling access to external APIs that have a fixed rate limit.

Workloads, a property of the scheduler, can be defined within the Quota
Scheduler too, allowing for strategic prioritization of requests when hitting
quota limits.

:::info See Also

Quota Scheduler [Reference](/reference/configuration/spec.md#quota-scheduler)

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

The scheduler adeptly manages incoming requests. These are classified based on
their urgency, utilizing parameters such as priority and tokens, which are
applicable to flows within a workload.

Priority varies from 0 to an unlimited positive integer, indicating the urgency
level, with higher numbers denoting higher priority. The position of a flow in
the queue is computed based on its virtual finish time using the following
formula:

$$
inverted\_priority = {\frac {\operatorname{lcm} \left(priorities\right)} {priority}}
$$

$$
virtual\_finish\_time = virtual\_time + \left(tokens \cdot inverted\_priority\right)
$$

To manage prioritized requests, the scheduler seeks tokens from the token
bucket. If tokens are available, they are used for scheduling and route requests
to the suitable service. In cases where tokens are not readily available,
requests are queued, waiting either until tokens become accessible or until a
timeout occurs - the latter being dependent on the workload or `check()` call
timeout.

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

The Dynamic Token Bucket, primarily utilized for service protection, employs a
Load Multiplier to control the [tokens](#tokens) refill rate per second in the
Scheduler. The objective is to maintain a stable refill rate between each
controller update. If the incoming request rate surpasses the anticipated rate,
the Scheduler queues these requests. Any request that fails to be scheduled
within its designated timeout is rejected. This strategy enables the adjustment
of the token rate corresponding to varying load signals, thereby providing an
active defense mechanism for the service.

Each agent operates an individual token bucket, utilizing lazy synchronization
for coordination among agents. This dynamic resizing functionality ensures
effective request management while maintaining optimal levels of service
protection.

#### Fixed Token Bucket {#fixed-bucket}

In contrast to the Dynamic Token Bucket, the Fixed Token Bucket operates with a
predetermined
[number of tokens](/reference/configuration/spec.md#rate-limiter-ins),
establishing a constant token rate. This type of bucket is particularly
advantageous when the load or the number of tokens is known upfront.

The Fixed Token Bucket can be employed for service protection but is especially
useful in scenarios where there are strict rate limits. For instance, when
scheduling access to external APIs with a global quota, the Fixed Token Bucket
ensures that the API's rate limits are adhered to.

Both the Dynamic and Fixed Token Buckets serve to manage and control the flow of
requests, their specific usage and operation depend on the requirements and
constraints of the system, providing versatility and adaptability in handling
workloads.

### Queue Timeout {#queue-timeout}

The queue timeout is determined by the gRPC timeout provided on the `check()`
call. When a request is made, it includes a timeout value that specifies the
maximum duration the request can wait in the queue. If the request receives the
necessary tokens within this timeout duration, it is admitted. Otherwise, if the
timeout expires before the tokens are available, the request is rejected.

The gRPC timeout on the `check()` call is set in the Envoy filter and the SDK
during initialization. It serves as an upper bound on the queue timeout,
preventing requests from waiting excessively long.

[label-matcher]: ./selector.md#label-matcher
