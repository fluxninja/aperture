---
title: Concurrency Limiter
sidebar_position: 6
---

:::info See also

[_Concurrency Limiter_ reference][reference]

:::

The _Concurrency Limiter_ component is used to enforce in-flight request
concurrency to prevent overloads. It can also be used to enforce limits per
entity such as a user to ensure fair access across users. Requests are allowed
or denied based on whether the in-flight requests are within the configured
limit. Instead of measuring the number of requests, the _Concurrency Limiter_
can also be configured to measure the number of tokens associated with a
request. Tokens can be sent as flow labels using Aperture SDKs.

:::note

Only accepted requests are counted towards the in-flight concurrency.

:::

## Lifecycle of a Request {#lifecycle-of-a-request}

The _Concurrency Limiter_ maintains a ledger of in-flight requests. The ledger
is updated by the Agents based on the flow start and end calls made from the
SDKs. Alternatively, for proxy integrations, the flow end is inferred as the
access log stream is received from the underlying middleware or proxy.

## Distributed Request Ledgers {#distributed-request-ledgers}

For each configured [_Concurrency Limiter Component_][reference], every matching
Aperture Agent instantiates a copy of the _Concurrency Limiter_. Although each
agent has its own copy of the component, they all share the in-flight request
ledger through a distributed cache. This means that they work together as a
single _Concurrency Limiter_, providing seamless coordination and control across
Agents. The Agents within an [agent group][agent-group] constantly share state
and detect failures using a gossip protocol.

## Max In-flight Duration {#max-in-flight-duration}

In case of failures at the SDK or middleware/proxy, the flow end call might not
be made. To prevent stale entries in the ledger, the _Concurrency Limiter_
allows the definition of a maximum in-flight duration. This can be set according
to the maximum time a request is expected to take. If the request exceeds the
configured duration, it is automatically removed from the ledger by the
_Concurrency Limiter_.

[reference]: /reference/configuration/spec.md#concurrency-limiter
[agent-group]: /concepts/selector.md#agent-group
