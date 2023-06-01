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

External Rate Limiting is a technique to limit the number of outgoing requests
from services to an external API server. Turning your apps into spend aware and
keeping them within quota limits to avoid cost overages. However, not all
workloads are on same priority, based on application, their priority can be
different. While doing external rate limiting, it is important to ensure
prioritized access for critical workloads. This policy builds upon the
[`Quota Scheduler`](/reference/policies/bundled-blueprints/policies/quota-scheduler.md)
Blueprint, which comprises components like the token bucket rate limiting to
ensure quota limits and a
[Weighted Fair Queuing (WFQ)](/concepts/flow-control/components/load-scheduler.md#scheduler)
based Workload Scheduler to assure prioritized access for critical workloads.

## Policy Key Concepts

At a high-level, this policy consists of:

- [Selector](../../concepts/flow-control/selector.md): Selectors are the traffic
  signal managers for flow control and observability components in the Aperture
  Agents. They lay down the traffic rules determining how these components
  should select flows for their operations.
- [Control Point](../../concepts/flow-control/selector.md): Think of Control
  Points as designated checkpoints in your code or data plane. They're the
  strategic points where flow control decisions are applied. Developers define
  these using SDKs or during API Gateways or Service Meshes integration.
- [Rate Limiter](../../concepts/flow-control/components/rate-limiter.md):
  Implemented on a token bucket algorithm, the rate limiter is an effective tool
  used to avoid recurring heavy traffic. This parking meter is flexible; it can
  be configured to work with different labels and limits.
- [Scheduler](../../concepts/flow-control/components/load-scheduler.md): The
  Scheduler ensures that requests are served based on their priority and size.
  It employs a Weighted Fair Queue-based system to serve the requests, using a
  Load Multiplier to calculate the token refill rate. If the incoming token rate
  exceeds the desired rate, it queues the requests and helps provide a smooth
  user experience, prioritizing critical orders over less urgent ones to ensure
  customer satisfaction.

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
