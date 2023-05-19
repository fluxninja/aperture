---
title: Rate Limiter
sidebar_position: 5
---

:::info

See also [_Rate Limiter_ reference][reference]

:::

The _Rate Limiter_ is a powerful component that can be used to prevent recurring
overloads by proactively regulating heavy-hitters. It achieves this by accepting
or rejecting incoming flows based on per-label limits, which are configured
using the [leaky bucket algorithm](https://en.wikipedia.org/wiki/Leaky_bucket).
The _Rate Limiter_ is a component of Aperture's [policy][policies] system, and
it can be configured to work with different labels and limits depending on the
needs of your application.

## Distributed Counters {#distributed-counters}

For each configured [_Rate Limiter Component_][reference], every matching
Aperture Agent instantiates a copy of the _Rate Limiter_. Although each Agent
has its own copy of the component, they all share counters through a distributed
cache. This means that they work together as a single _Rate Limiter_, providing
seamless coordination and control across agents. The agents within an [agent
group][agent-group] constantly share state and detect failures using a gossip
protocol.

### Leaky Bucket Algorithm {#leaky-bucket-algorithm}

This algorithm lets users make an unlimited number of requests in infrequent
bursts over time. The main points to understand about the leaky bucket metaphor
are as follows:

- Each user (or any flow label) has access to a bucket. It can hold, say, 60
  “marbles”.
- Each second, a marble is removed from the bucket (if there are any). That way
  there’s always more room.
- Each API request requires the user to toss a marble in the bucket.
- If the bucket gets full, the user gets an error and has to wait for room to
  become available in the bucket.

This model ensures that apps that manage API calls responsibly will always have
room in their buckets to make a burst of requests if needed. For example, if
users average 20 requests (“marbles”) per second but suddenly need to make 30
requests all at once, users can still do so without hitting their rate limit.
The basic principles of the leaky bucket algorithm apply to all our rate limits,
regardless of the specific methods used to apply them.

### Lazy Syncing {#lazy-syncing}

When lazy syncing is enabled, rate-limiting counters are stored in-memory and
are only synchronized between Aperture Agent instances on-demand. This allows
for fast and low-latency rate-limiting decisions, at the cost of slight
inaccuracy within a (small) time window (sync interval).

## Limits {#limits}

The _Rate Limiter_ component accepts or rejects incoming flows based on
per-label limits, configured as the maximum number of requests per a given
period of time. The rate-limiting label is chosen from the
[flow-label][flow-label] with a specific key, enabling you to configure separate
limits for different users or flows.

:::tip

The limit value is treated as a signal within the circuit and can be dynamically
modified or disabled at runtime.

:::

### Overrides {#overrides}

The override mechanism allows for increasing the limit for a particular value of
a label. For instance, you might want to increase the limit for the admin user.
Please refer to the [reference][reference] for more details on how to use this
feature.

[reference]: /reference/policies/spec.md#rate-limiter
[agent-group]: /concepts/flow-control/selector.md#agent-group
[policies]: /concepts/policy/policy.md
[flow-label]: /concepts/flow-control/flow-label.md
