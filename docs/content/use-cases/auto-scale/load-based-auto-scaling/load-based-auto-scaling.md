---
title: Load-based Auto Scaling
sidebar_position: 1
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
  - HPA
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

:::note

The following policy is based on the
[Service Protection with Average Latency Feedback](/reference/blueprints/policies/service-protection/average-latency.md)
blueprint.

:::

## Overview

Responding to fluctuating service demand is a common challenge for maintaining
stable and responsive services. This policy introduces a mechanism to
dynamically scale service resources based on observed load, optimizing resource
allocation and ensuring that the service remains responsive even under high
load.

This policy employs two key strategies: service protection and auto-scaling.

1. Service Protection: Based on the trend of observed latency, the service gets
   protected from sudden traffic spikes using a
   [_Load Scheduler_](/concepts/flow-control/components/load-scheduler.md)
   component. Load on the service is throttled when the observed latency exceeds
   the long-term trend by a certain percentage threshold. This ensures the
   service stays responsive even under high load.
2. Auto-Scaling: The auto-scaling strategy is based on the throttling behavior
   of the service protection policy. An
   [_Auto Scaler_](/concepts/auto-scale/components/auto-scaler.md) component is
   used to dynamically adjust the number of service instances in response to
   changes in load. This load-based auto-scaling is enacted by a scale-out
   Controller that reads Load Scheduler signals. The service replicas are scaled
   out when the load is being throttled, effectively scaling resources to match
   the demand. During periods of low load, the policy attempts to scale in after
   periodic intervals to reduce excess replicas.

By combining service protection with auto-scaling, this policy ensures that the
number of service replicas is adjusted to match persistent changes in demand,
maintaining service stability and responsiveness.

## Configuration

This policy, provides protection against overloads at the
**`search-service.prod.svc.cluster.local`** service. Auto-scaling is applied to
the Deployment `search-service` with a minimum of `1` and a maximum of `10`
replicas.

To prevent frequent fluctuation in replicas, scale-in and scale-out cooldown
periods are set to `40` and `30` seconds, respectively. A periodic scale-in
interval of `60` seconds is also set to reduce excess replicas during periods of
low load.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/policy.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/graph.mmd.svg) for this policy.

:::

### Policy in Action

During transient load spikes, the response latency on the service increases. The
service protection policy queues a proportion of the incoming requests. The
_Auto Scaler_ makes a scale-out decision as the `OBSERVED_LOAD_MULTIPLIER` falls
below 1. This triggers the auto-scale policy, which scales up the deployment.
With the additional replicas in the deployment, the service is now better
equipped to handle the increased load. The `OBSERVED_LOAD_MULTIPLIER` rises
above 1, enabling the service to meet the heightened demand. As a result, the
response latency returns to a normal range, and the Load Scheduler ceases
throttling.

After the scale-out cooldown period, the periodic scale-in function is
triggered, which reduces the number of replicas in response to decreased load.

<Zoom>

![Auto Scale](./assets/dashboard.png)

</Zoom>
