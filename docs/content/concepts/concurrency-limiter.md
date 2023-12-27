---
title: Concurrency Limiter
sidebar_position: 6
---

:::info See also

[_Concurrency Limiter_ reference][reference]

:::

The _Concurrency Limiter_ component is used to enforce in-flight request quotas
to prevent overloads. It can also be used to enforce limits per entity such as a
user to ensure fair access across users. Essentially, providing an added layer
of protection in additional to per-user rate limits.

_Concurrency Limiter_ can limit the number of concurrent requests to a control
point or certain labels that match within the control point. It achieves this by
maintaining a ledger of in-flight requests. If the number of in-flight requests
exceeds the configured limit, the _Concurrency Limiter_ rejects new requests
until the number of in-flight requests drops below the limit. The in-flight
requests are maintained by the Agents based on the flow start and end calls made
from the SDKs. Alternatively, for proxy integrations, the flow end is inferred
as the access log stream is received from the underlying middleware or proxy.

## Distributed Request Ledgers {#distributed-request-ledgers}

For each configured [_Concurrency Limiter Component_][reference], every matching
Aperture Agent instantiates a copy of the _Concurrency Limiter_. Although each
agent has its own copy of the component, they all share the in-flight request
ledger through a distributed cache. This means that they work together as a
single _Concurrency Limiter_, providing seamless coordination and control across
Agents. The Agents within an [agent group][agent-group] constantly share state
and detect failures using a gossip protocol.

## Lifecycle of a Request {#lifecycle-of-a-request}

The _Concurrency Limiter_ maintains a ledger of in-flight requests. The ledger
is updated by the Agents based on the flow start and end calls made from the
SDKs. Alternatively, for proxy integrations, the flow end is inferred as the
access log stream is received from the underlying middleware or proxy.

### Max In-flight Duration {#max-in-flight-duration}

In case of failures at the SDK or middleware/proxy, the flow end call might not
be made. To prevent stale entries in the ledger, the _Concurrency Limiter_
allows the definition of a maximum in-flight duration. This can be set according
to the maximum time a request is expected to take. If the request exceeds the
configured duration, it is automatically removed from the ledger by the
_Concurrency Limiter_.

[reference]: /reference/configuration/spec.md#concurrency-limiter
[agent-group]: /concepts/selector.md#agent-group
