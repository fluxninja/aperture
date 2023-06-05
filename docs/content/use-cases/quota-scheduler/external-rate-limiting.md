---
title: External Rate Limiting
keywords:
  - policies
  - quota
  - prioritization
  - external-api
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Policy Overview

The Quota Scheduler Policy is a sophisticated solution designed to manage and
limit outgoing requests from services to an external API server. This policy
makes applications cost-aware, ensuring that they operate within assigned quota
limits to prevent cost overruns. Workload priorities might differ based on the
application, and maintaining prioritized access for critical workloads during
external rate limiting is of paramount importance. The policy leverages the
[`Quota Scheduler`](/reference/policies/bundled-blueprints/policies/quota-scheduler.md)
Blueprint, which brings together the token bucket rate limiting and a
[Weighted Fair Queuing (WFQ)](/concepts/flow-control/components/load-scheduler.md#scheduler)
based Workload Scheduler to balance quota limits and priority-based access
efficiently.

## Policy Key Concepts

The policy operates around a set of core components within the
[`quota_scheduler`], each serving a specific function in the overall rate
limiting process.

- [`selectors`](../../concepts/flow-control/selector.md) define the rules that
  decide how these components should select flows for processing.
- [`control point`](../../concepts/flow-control/selector.md) can be considered
  as a critical checkpoint in code or data plane, a strategically placed spot
  where flow control decisions are applied. Developers define these points
  during the integration of API Gateways or Service Meshes or by using Aperture
  SDKs.
- [`rate_limiter`](../../concepts/flow-control/components/rate-limiter.md)
  prevents heavy traffic recurrence and its flexibility allows it to adapt to
  different labels, offering dynamic control over traffic flow.
- [`scheduler`](../../concepts/flow-control/components/load-scheduler.md)
  ensures that requests are serviced based on their priority and size. It
  employs a Weighted Fair Queue-based system, calculating the token refill rate
  using a Load Multiplier, effectively managing and prioritizing requests.

## Policy Configuration

In this policy,
[Quota Scheduler](/reference/policies/bundled-blueprints/policies/quota-scheduler.md#policy-quota-scheduler)
component is configured with `bucket_capacity`, and rate limiting is configured
based on label key `api_key` extracted from the request header. While the lazy
sync of between the agent is set to false.

WFQ Scheduler is configured two workloads priorities; `guest` and `subscriber`
with 50 and 200 respectively. Matching labels using `user_type` value from the
request header.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/with-external-api-calls-prioritization/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>

```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/with-external-api-calls-prioritization/policy.yaml}
```

</p>
</details>

## Playground

The above policy can be loaded using the `quota-scheduler` scenario in
[Playground](https://github.com/fluxninja/aperture/blob/main/playground/README.md)

:::info

[Circuit Diagram](./assets/with-external-api-calls-prioritization/graph.mmd.svg)
for this policy.

:::

<Zoom>

![External Rate Limiting With Prioritization ](./assets/with-external-api-calls-prioritization/dashboard.png)

</Zoom>
