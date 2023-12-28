---
title: Quota Scheduler
keywords:
  - scheduler
  - service protection
  - queuing
  - quota
sidebar_position: 1
---

:::info See Also

Quota Scheduler [Reference](/reference/configuration/spec.md#quota-scheduler)

:::

The _Quota Scheduler_ is used to schedule requests based on importance while
ensuring that the application adheres to third-party API rate limits or
inter-service API quotas.

The _Quota Scheduler_ can be thought of as a combination of a
[_Scheduler_](../scheduler.md) and a [_Rate Limiter_](../rate-limiter.md). It
essentially provides scheduling capabilities atop a _Rate Limiter_. In the
policy circuit, this component takes the same input ports as a _Rate Limiter_,
namely `fill_rate` and `bucket_capacity`. These ports facilitate adjustment of
the global token bucket, which can be used to model an API quota or rate limit.

The token bucket represents a fixed quota that is divided among the Agents. It
is used as a shared ledger for Agents in an
[agent group](../advanced/agent-group.md). This ledger records the total
available tokens that can be distributed across the Agents. Tokens are consumed
from it when admitting requests. If the ledger runs out of tokens, new requests
are queued until more tokens become available or
[until timeout](../scheduler.md#queue-timeout).

In a scenario where the token fill rate and bucket capacity (API quota) is known
upfront, the _Quota Scheduler_ becomes particularly beneficial to enforce
client-side rate limits.

The _Quota Scheduler_ also allows the definition of
[workloads](../scheduler.md#workload), a property of the scheduler, which allows
for strategic prioritization of requests when faced with quota constraints. As a
result, the _Quota Scheduler_ ensures adherence to the API's rate limits and
simultaneously offers a mechanism to prioritize requests based on their
importance.

:::info

Refer to the [API Quota Management guide][guide] for more information on how to
use the _Quota Scheduler_ using [aperture-js][aperture-js] SDK.

:::

[guide]: /guides/api-quota-management.md
[aperture-js]: https://github.com/fluxninja/aperture-js
