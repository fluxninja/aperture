---
title: Rate Limiter
sidebar_position: 5
---

:::info

See also [_Rate Limiter_ reference][reference]

:::

The _Rate Limiter_ is a powerful tool that can be used to prevent recurring
overloads by proactively regulating heavy-hitters. It achieves this by accepting
or rejecting incoming flows based on per-label limits, which are configured as
the number of requests per given period of time. The _Rate Limiter_ is a
component of Aperture's [policy][policies] system, and it can be configured to
work with different labels and limits depending on the needs of your
application.

## Distributed Counters {#distributed-counters}

For each configured [_Rate Limiter Component_][reference], every matching
Aperture Agent instantiates a copy of the _Rate Limiter_. Although each agent
has its own copy of the component, they all share counters through a distributed
counter system. This means that they work together as a single _Rate Limiter_,
providing seamless coordination and control across agents. The distributed
counters are powered by the [Agent-to-Agent peer-to-peer network][agent-group]
network, which ensures reliable and efficient communication between agents.

### Lazy Syncing {#lazy-syncing}

When lazy syncing is enabled, rate-limiting counters are stored in-memory and
are only synchronized between Aperture Agent instances on-demand. This allows
for fast and low-latency rate-limiting decisions, at the cost of slight
inaccuracy in edge-cases.

## Limits {#limits}

The _Rate Limiter_ component accepts or rejects incoming flow based on per-label
limits, configured as the maximum number of requests per a given period of time.
The rate limiting label is chosen from the [flow-label][flow-label] with a
specific key, enabling you to configure separate limits for different users or
flows.

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
[agent-group]: /concepts/integrations/flow-control/flow-selector.md#agent-group
[policies]: /concepts/policy/policy.md
[flow-label]: /concepts/integrations/flow-control/flow-label.md
