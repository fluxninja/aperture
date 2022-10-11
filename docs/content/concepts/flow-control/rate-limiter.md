---
title: Rate Limiter
sidebar_position: 5
---

:::info

See also [Rate Limiter reference][reference]

:::

Rate Limiter is a versatile tool that could be used to protect your services,
but also can be used for different purposes, like giving fairness to your users.
Some potential usecases are:

- protecting the service from bot traffic,
- ratelimiting users,
- preventing from reaching external API quotas.

Rate Limiter is configured as a [policy][policies] component.

## Distributed Counters

For each configured [Rate Limiter component][reference], every matching Aperture
Agent instantiates a copy of Rate Limiter. They're all sharing counters though,
so conceptually they work as a single Rate Limiter. That's possible thanks to
distributed counters, powered by [Agent-to-Agent peer-to-peer
network][agent-group].

### Lazy Syncing

If lazy syncing is enabled, rate-limiting counters are stored in-memory and
lazily-syced between Agent instances. Thanks to this, rate-limiting decisions
can be made without latency overhead at the slight cost of accuracy in
edge-cases.

## Limits

Rate-limiter accepts or rejects incoming flow based on per-label limits
(configured as number of requests per given period of time). Rate limiting label
is chosen from [flow-label][flow-label] with a given key. Eg. you can configure
each user to have a separate limit.

:::tip

The value of the limit is accepted as a circuit signal and can be changed (or
even disabled) on runtime.

:::

### Overrides

A limit for particular value of a label can be increased via the override
mechanisms. Eg. you might want to increase the limit for the admin user. See
[reference][reference] for more details.

[reference]: /references/configuration/policy.md#v1-rate-limiter
[agent-group]: /concepts/service.md#agent-group
[policies]: /concepts/policy/policy.md
[flow-label]: /concepts/flow-control/flow-label.md
