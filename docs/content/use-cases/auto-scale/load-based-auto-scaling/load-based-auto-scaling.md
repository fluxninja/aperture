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

This policy builds upon the _Service Protection with Load-based Pod Auto-Scaler_
[blueprint](/reference/policies/bundled-blueprints/policies/service-protection-with-load-based-pod-auto-scaler/average-latency.md)
to add an escalation for auto-scaling. The basic service protection policy
protects the service from sudden traffic spikes. But it is necessary to scale
the service to meet demand in case of a persistent change in load.

To achieve this, the policy makes use of an
[_Auto Scaler_](/concepts/auto-scale/components/auto-scaler.md) component to
adjust the number of instances allocated to the service. Load-based auto-scaling
is achieved by defining a scale-out _Controller_ that acts on the load
multiplier signal from the service protection policy. This signal measures the
fraction of traffic that the
[_Load Scheduler_](/concepts/flow-control/components/load-scheduler.md) is
throttling into a queue. The _Auto Scaler_ is configured to scale-out using a
_Gradient Controller_ based on this signal and a setpoint of 1.0.

In addition to load-based scaling, the policy includes a scale-in _Controller_
based on CPU utilization. These _Controllers_ adjust the resources allocated to
the service based on changes in CPU usage, ensuring that the service can handle
the workload efficiently.

## Policy Key Concepts

At a high level, this policy consists of:

- Service protection based on response latency trend of the service.
- An _Auto Scaler_ that adjusts the number of replicas of the Kubernetes
  Deployment for the service.
- Load-based scale-out is done based on `OBSERVED_LOAD_MULTIPLIER` signal from
  the blueprint. This signal measures the fraction of traffic that the _Load
  Scheduler_ is throttling into a queue. The _Auto Scaler_ is configured to
  scale-out based on a _Gradient Controller_ using this signal and a setpoint of
  1.0.
- In addition to the load-based scale-out, the policy also includes a scale-in
  _Controller_ based on CPU utilization which adjusts the instances of the
  service based on changes in CPU usage, ensuring that the service is not
  over-provisioned.

Some of the key concepts used in this policy are:

- [Load Scheduler](../../../concepts/flow-control/components/load-scheduler.md):
  The Load Scheduler prevents chaos by managing incoming request traffic
  efficiently. It's tasked with limiting the concurrent requests to a service
  and assigning different priorities and weights to workloads to ensures that
  high-priority requests get served first during heavy traffic.
- [Selector](../../../concepts/flow-control/selector.md): Selectors are the
  traffic signal managers for flow control and observability components in the
  Aperture Agents. They lay down the traffic rules determining how these
  components should select flows for their operations.
- [Control Point](../../../concepts/flow-control/selector.md): Think of Control
  Points as designated checkpoints in your code or data plane. They're the
  strategic points where flow control decisions are applied. Developers define
  these using SDKs or during API Gateways or Service Meshes integration.
- [FluxMeter](../../../concepts/flow-control/resources/flux-meter.md): Flux
  Meter converts a flux of flows matching a Flow Selector into a Prometheus
  histogram. By default, it tracks the workload duration of a flow. However,
  it's flexible enough to track any metric from OpenTelemetry attributes based
  on the method of insertion.

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
