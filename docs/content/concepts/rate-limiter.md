---
title: Rate Limiter
sidebar_position: 11
---

:::info See also

[_Rate Limiter_ reference][reference]

:::

The _Rate Limiter_ component can be used to prevent recurring overloads by
proactively regulating heavy-hitters. It achieves this by accepting or rejecting
incoming flows based on per-label limits, which are configured using the
[token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket).

The _Rate Limiter_ is a component of Aperture's [policy][policies] system, and
it can be configured to work with different labels and limits depending on the
needs of an application.

The following example creates a _Rate Limiter_ at the `ingress` control point
for service `checkout.default.svc.cluster.local`. A rate limit of `2` requests
per second with a burst capacity of `40` is applied per unique value of
`http.request.header.user_id` flow label:

```yaml
circuit:
components:
  - flow_control:
      rate_limiter:
        in_ports:
          bucket_capacity:
            constant_signal:
              value: 40
          fill_amount:
            constant_signal:
              value: 2
        parameters:
          interval: 1s
          label_key: http.request.header.user_id
        selectors:
          - control_point: ingress
            service: checkout.default.svc.cluster.local
```

## Distributed Counters {#distributed-counters}

For each configured [_Rate Limiter Component_][reference], every matching
Aperture agent instantiates a copy of the _Rate Limiter_. Although each agent
has its own copy of the component, they all share counters through a distributed
cache. This means that they work together as a single _Rate Limiter_, providing
seamless coordination and control across agents. The agents within an [agent
group][agent-group] constantly share state and detect failures using a gossip
protocol.

### Token Bucket Algorithm {#token-bucket-algorithm}

This algorithm allows users to execute a substantial number of requests in
bursts, and then continue at a steady rate. Here are the key points to
understand about the token bucket algorithm:

- Each user (or any flow label) has access to a bucket, which can hold, say, 60
  "tokens".
- Every second, a token is added to the bucket (if there's room). In this way,
  the bucket is steadily refilled over time.
- Each API request requires the user to remove a token from the bucket.
- If the bucket is empty, the user gets an error and has to wait for new tokens
  to be added to the bucket before making more requests.

This model ensures that apps that handle API calls judiciously will always have
a supply of tokens for a burst of requests when necessary. For example, if users
average 20 requests ("tokens") per second but suddenly need to make 30 requests
at once, users can do so if they have accumulated enough tokens.

### Lazy Syncing {#lazy-syncing}

When lazy syncing is enabled, rate-limiting counters are stored in-memory and
are only synchronized between Aperture agent instances on-demand. This allows
for fast and low-latency rate-limiting decisions, at the cost of slight
inaccuracy within a (small) time window (sync interval).

## Limits {#limits}

The _Rate Limiter_ component accepts or rejects incoming flows based on
per-label limits, configured as the maximum number of requests per a given
period of time. The rate-limiting label is chosen from the
[flow-label][flow-label] with a specific key, enabling distinct limits per user
as identified by unique values of the label.

:::tip

The limit value is provided as a signal within the circuit. It can be set
dynamically based on the circuit's logic.

:::

[reference]: /reference/configuration/spec.md#rate-limiter
[agent-group]: /concepts/selector.md#agent-group
[policies]: /concepts/advanced/policy.md
[flow-label]: /concepts/flow-label.md
