---
title: Prioritization
keywords:
  - policies
  - scheduler
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Overview

When dealing with services in resource-limited scenarios, it becomes paramount
to prioritize key user experiences and business-critical features over less
crucial tasks or background workloads. For instance, in an e-commerce platform,
the checkout process must take precedence over functionalities like personalized
recommendations, especially during resource shortage or high traffic. Aperture's
[Weighted Fair Queuing Scheduler (WFQ)](/concepts/flow-control/components/load-scheduler.md#scheduler)
enables such prioritization of flows over others based on their labels, ensuring
user experience or revenue is maximized during overloads or other failure
scenarios.

## Configuration

This policy defines service protection on
**`cart-service.prod.svc.cluster.local`** using a load scheduler and overload
detection is based on average latency similar to the
[protection policy](protection.md). In addition, workload prioritization is
specified in the load scheduler based on user types accessing the service. User
types are identified based on the value of a header label
`http.request.header.user_type`. Requests matching label value **`guest`** are
assigned a priority of 50, whereas those matching **`subscriber`** are given a
priority of 200.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/prioritization/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/prioritization/policy.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/prioritization/graph.mmd.svg) for this policy.

:::

## Policy in Action

The traffic generator in this scenario is configured to generate similar traffic
pattern (number of concurrent users) for 2 types of users - subscribers and
guests.

The below dashboard shows that, during overload periods, requests from
`subscriber` users have a higher acceptance rate than those from `guest` users.

<Zoom>

![Workload Prioritization](./assets/prioritization/dashboard.png)

</Zoom>

### Demo Video

The below demo video shows the basic concurrency limiter and workload
prioritization policy in action within Aperture Playground.

[![Demo Video](https://img.youtube.com/vi/m070bAvrDHM/0.jpg)](https://www.youtube.com/watch?v=m070bAvrDHM)