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

## Policy Overview

Responding to fluctuating service demand is a common challenge for maintaining
stable and responsive services. This policy, based on the Service Protection
with Load-based Pod Auto-Scaler
[blueprint](/reference/policies/bundled-blueprints/policies/service-protection-with-load-based-pod-auto-scaler/average-latency.md),
presents an evolved strategy to tackle this. It introduces a mechanism to
dynamically scale the service resources based on observed load, thereby
optimizing resource allocation and maintaining a balanced system.

This policy employs two key strategies: it protects the service from sudden
traffic spikes, and it ensures the service scales proportionally to accommodate
sustained load changes. The policy uses an
[_Auto Scaler_](/concepts/auto-scale/components/auto-scaler.md) component to
dynamically adjust the number of service instances in response to changes in
load and CPU utilization. This load-based auto-scaling is enacted by a scale-out
Controller that takes input from the
[_Load Scheduler_](/concepts/flow-control/components/load-scheduler.md) signal,
effectively throttling traffic into a queue and scaling resources to match the
demand.

## Policy Key Concepts

This policy integrates a suite of concepts and components to enable a dynamic,
load-responsive service operation:

- Service Protection Core: This component employs an [`adaptive_load_scheduler`]
  to manage incoming traffic and prevent chaotic load situations. It uses a Load
  Scheduler to limit concurrent requests, assigning different priorities and
  weights to workloads to ensure high-priority requests are served first during
  heavy traffic.
- [`latency_baseliner`]: This subsystem includes a [`flux_meter`] that measures
  the scope of latency by converting a flux of flows matching a flow selector
  into a Prometheus histogram. By default, it tracks the workload duration of a
  flow, but it can flexibly track any metric from OpenTelemetry attributes
  depending on the insertion method.
- Auto Scaling: A crucial part of this policy, this component facilitates
  automatic scaling of service instances based on the current load. It includes
  a Kubernetes replica scaling backend that adjusts the number of replicas of
  the Kubernetes Deployment for the service, ensuring it matches the current
  demand. The [`dry_run`] parameter can be used to simulate scaling actions
  without actually performing them, which is useful for testing and
  verification.
- Scaling Parameters: These are crucial for controlling the behavior of the
  [`auto scaling`] component. Parameters such as [`scale_in_cooldown`] and
  [`scale_out_cooldown`] define the minimum amount of time between consecutive
  scale-in and scale-out actions, preventing overactive scaling.

## Policy Configuration

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

### Playground

When the above policy is loaded in Aperture's
[Playground](https://github.com/fluxninja/aperture/blob/main/playground/README.md),
it can be observed that as the response latency increases, the service
protection policy queues a proportion of requests. The _Auto Scaler_ makes a
scale-out decision as the `OBSERVED_LOAD_MULTIPLIER` becomes less than 1. As
replicas get added to the deployment, the `OBSERVED_LOAD_MULTIPLIER` increases
to more than 1, allowing the service to meet increased demand. The response
latency returns to a normal range, and the _Load Scheduler_ won't throttle any
traffic.

After the scale-out cooldown period, the scale-in based on CPU utilization gets
triggered, which will cause the replicas to decrease. Once the traffic ramps up
again, the above cycle continues.

<Zoom>

![Auto Scale](./assets/dashboard.png)

</Zoom>
