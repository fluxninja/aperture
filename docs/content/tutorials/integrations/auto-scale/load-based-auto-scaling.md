---
title: Load-Based Auto Scaling
sidebar_position: 2
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
  - HPA
---

# **Load-Based Auto Scaling**

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Overload protection is done by adaptive concurrency limiting it is also
important to scale your services up and down based on the traffic. Aperture
provides a way to scale your services up and down based on the metrics using
Autoscale and Pod Scale components.

Though there is Horizontal Pod Autoscaler (HPA) in Kubernetes but it has
limitations. The HPA can only scale based on a limited set of metrics, such as
CPU or memory utilization. Additionally, it can't handle complex scaling
behavior based on multiple metrics, making it difficult to optimize application
performance and efficiency.

Whereas Aperture Auto Scale has ability to handle complex scaling behavior based
on multiple metrics. For example, Aperture Auto Scale can use a combination of
latency, throughput, and error rates to determine how many replicas to deploy,
providing a more accurate and nuanced approach to scaling. This is essential for
applications where performance is critical, such as e-commerce or financial
applications.

In this tutorial, we will see how to use Aperture's Auto Scale component to do a
load-based autoscaling for a Kubernetes deployment.

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
{@include: ./assets/load-based-autoscale-service1.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/load-based-autoscale-service1-cr.yaml}
```

</p>
</details>

### Circuit Diagram

<Zoom>

```mermaid
{@include: ./assets/load-based-autoscale-service1-cr.mmd}
```

</Zoom>

### Playground

When the above policy is loaded in Aperture's
[Playground](/get-started/playground/playground.md), we will see that as the
traffic spikes above the concurrency limit of
`service1-demo-app.demoapp.svc.cluster.local`, controller triggers load-shed for
a proportion of requests matching the Selector and average cpu signal and load
multiplier triggers a signal to do autoscale.

<Zoom>

![Auto-Scale](./assets/auto-scale-playground.png)

</Zoom>
