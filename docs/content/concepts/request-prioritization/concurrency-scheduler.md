---
title: Concurrency Scheduler
keywords:
  - scheduler
  - concurrency
  - queuing
sidebar_position: 2
---

:::info See Also

Concurrency Scheduler
[Reference](/reference/configuration/spec.md#concurrency-scheduler)

:::

The _Concurrency Scheduler_ is used to schedule requests based on importance
while ensuring that the application adheres to concurrency limits.

The _Concurrency Scheduler_ can be thought of as a combination of a
[_Scheduler_](../scheduler.md) and a
[_Concurrency Limiter_](../concurrency-limiter.md). It essentially provides
scheduling capabilities atop a _Concurrency Limiter_. Similar to the
_Concurrency Limiter_, this component takes `max_concurrency` as an input port
which determines the maximum number of in-flight requests in the global request
ledger.

The global request ledger is shared among Agents in an
[agent group](../advanced/agent-group.md). This ledger records the total number
of in-flight requests across the Agents. If the ledger exceeds the configured
`max_concurrency`, new requests are queued until the number of in-flight
requests drops below the limit or
[until timeout](../scheduler.md#queue-timeout).

:::note

Only accepted requests are counted towards the in-flight concurrency.

:::

In a scenario where the maximum concurrency is known upfront, the _Concurrency
Scheduler_ becomes particularly beneficial to enforce concurrency limits on a
per-service basis.

The _Concurrency Scheduler_ also allows the definition of
[workloads](../scheduler.md#workload), a property of the scheduler, which allows
for strategic prioritization of requests when faced with concurrency
constraints. As a result, the _Concurrency Scheduler_ ensures adherence to the
concurrency limits and simultaneously offers a mechanism to prioritize requests
based on their importance.
