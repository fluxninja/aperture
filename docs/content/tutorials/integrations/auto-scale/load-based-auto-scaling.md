---
title: Load-based Auto Scaling
sidebar_position: 2
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
  - HPA
---

# **Load-based Auto Scaling**

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Load-based auto-scaling is a technique used to dynamically adjust the number of
instances or resources allocated to a service based on workload demands. In
Aperture, the
[_AutoScaler_](concepts/integrations/auto-scale/components/auto-scaler.md)
component enables load-based auto-scaling by interfacing with infrastructure
APIs such as Kubernetes to automatically adjust the resources allocated to a
service.

This tutorial builds upon the
[basic concurrency limiting policy](tutorials/integrations/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
to add an escalation for auto-scaling. The concurrency limiter protects the
service from sudden traffic spikes, but to protect from persistent changes in
load, it's necessary to scale the service in response to demand.

Load-based auto-scaling extends the concurrency limiting policy by automatically
scaling the service to meet persistent changes in demand. In Aperture,
load-based auto-scaling is accomplished by configuring the _AutoScaler_
component with Controllers that can adjust the number of instances or resources
allocated to the service.

Observed Load Multiplier The load-based auto-scaling policy makes use of the
OBSERVED*LOAD_MULTIPLIER signal from the
[AIMDConcurrencyController](reference/policies/spec.md#a-i-m-d-concurrency-controller)
component. This signal measures the amount of traffic that is being load shed by
the Concurrency Controller. The \_AutoScaler* is configured to scale out based
on this signal and a setpoint of 1.0.

CPU Utilization In addition to the load shedding signal, the load-based
auto-scaling policy also includes scale in and scale out Controllers based on
CPU utilization. These Controllers adjust the resources allocated to the service
based on changes in CPU usage, ensuring that the service can handle the workload
efficiently.

## Policy

In this policy we will be using the Latency based AIMD (Additive
Increase,Multiplicative Decrease) Concurrency Limiting
[Blueprint](reference/policies/bundled-blueprints/policies/latency-aimd-concurrency-limiting.md).
Policy will do load-based autoscaling for a Kubernetes deployment.

At a high Level, this policy consist:

This policy includes:

- Using a Flux Meter to gather latency metrics from a service control point,
  which are then fed into an EMA component to establish a long-term trend that
  can be compared against current latency to detect overloads
- A Gradient Controller that adjusts the Accepted Concurrency based on Setpoint
  Latency and current Latency signals
- An Integral Optimizer that additively increases the concurrency on the service
  in each execution cycle of the circuit when the service is in the normal
  state, allowing warming-up from a cold start state and protecting applications
  from sudden spikes in traffic
- A Concurrency Limiting Actuator that actuates concurrency limits via a
  weighted-fair queueing scheduler, with adjustments to accepted concurrency
  made by gradient controller and optimizer logic translated to a load
  multiplier that adjusts token bucket fill rates based on incoming concurrency
  observed at each agent
- An Auto Scale feature that adjusts the number of replicas of the configured
  Kubernetes Resource, including a scale-in controller that adjusts the number
  of replicas based on observed CPU utilization and a scale-out controller that
  adjusts the number of replicas based on the observed load multiplier

The Policy also includes a set of resources classifiers and flux meters, which
are used to define the components of the policy that should be applied to the
deployment based on certain criteria (in this case, the "user_type" field in the
request headers).

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
{@include: ./assets/load-based-auto-scale.yaml}
```

</p>
</details>

### Circuit Diagram

<Zoom>

```mermaid
{@include: ./assets/load-based-auto-scale.mmd}
```

</Zoom>

### Playground

When the above policy is loaded in Aperture's
[Playground](/get-started/playground/playground.md), we will see that as the
traffic spikes above the concurrency limit of
`service1-demo-app.demoapp.svc.cluster.local`, controller triggers load-shed for
a proportion of requests matching the Selector and average cpu signal and load
multiplier triggers a signal to do auto-scale.

<Zoom>

![Auto-Scale](./assets/auto-scale-playground.png)

</Zoom>
