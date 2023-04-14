---
title: Basic Concurrency Limiting
keywords:
  - policies
  - concurrency
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

The most effective technique to protect services from cascading failures is to
limit the concurrency of the service to match the processing capacity of the
service. However, figuring out the concurrency limit of a service is a hard
problem in the face of continuously changing service infrastructure. Each new
version deployed, horizontal scaling, or a change in access patterns can change
the concurrency limit of a service.

To accurately model the concurrency limit of a service, it's critical to track
its
[golden signals](https://sre.google/sre-book/monitoring-distributed-systems/#xref_monitoring_golden-signals).
For instance, queue buildup can be detected by tracking deviation of current
latency from historically normal values.

## Policy

This policy uses the Latency based AIMD (Additive Increase, Multiplicative
Decrease) Concurrency Limiting
[Blueprint](/reference/policies/bundled-blueprints/policies/latency-aimd-concurrency-limiting.md)
and is instantiated via Jsonnet. The Signal Processing tutorials describe
various building blocks used in the policy separately.

At a high-level, this policy consists of:

- Latency EMA-based overload detection: A Flux Meter is used to gather latency
  metrics from a
  [service control point](/concepts/flow-control/flow-selector.md). The latency
  signal is then fed into an Exponential Moving Average (EMA) component to
  establish a long-term trend that can be compared to the current latency to
  detect overloads. For more information on how this is achieved, see the
  tutorial on
  [Detecting Overload](/tutorials/signal-processing/detecting-overload.md).
- Gradient Controller: Set point latency and current latency signals are fed to
  the gradient controller that calculates the proportional response to adjust
  the Accepted Concurrency (Control Variable).
- Integral Optimizer: When the service is detected to be in the normal state, an
  integral optimizer is used to additively increase the concurrency of the
  service in each execution cycle of the circuit. This design allows warming-up
  a service from a cold start state. This also protects applications from sudden
  spikes in traffic, as it sets an upper bound to the concurrency allowed on a
  service in each execution cycle of the circuit based on the observed incoming
  concurrency.
- Concurrency Limiting Actuator: The concurrency limits are actuated via a
  [weighted-fair queuing scheduler](/concepts/flow-control/components/concurrency-limiter.md).
  The output of the adjustments to accepted concurrency made by gradient
  controller and optimizer logic are translated to a load multiplier that's
  synchronized with Aperture Agents via etcd. The load multiplier adjusts
  (increases or decreases) the token bucket fill rates based on the incoming
  concurrency observed at each agent.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/basic-concurrency-limiting/values.yaml}
```

```mdx-code-block
</TabItem>
<TabItem value="Jsonnet Mixin">
```

```jsonnet
{@include: ./assets/basic-concurrency-limiting/basic-concurrency-limiting.jsonnet}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/basic-concurrency-limiting/basic-concurrency-limiting.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/basic-concurrency-limiting/basic-concurrency-limiting.mmd.svg)
for this policy.

:::

### Playground

When the above policy is loaded in Aperture's
[Playground](/get-started/playground/playground.md), it can be observed that as
traffic spikes above the concurrency limit of
`service1-demo-app.demoapp.svc.cluster.local`, the controller triggers load
shedding for a proportion of requests matching the Selector. This helps to
protect the service from becoming unresponsive and keeps the latency within the
tolerance limit (`1.1`) configured in the circuit.

<Zoom>

![Basic Concurrency Limiting](./assets/basic-concurrency-limiting/basic-concurrency-limiting-playground.png)

</Zoom>

### Dry Run Mode

You can run this policy in the `Dry Run` mode by setting the
`defaultConfig.dry_run` option to `true`. In the `Dry Run` mode, the policy
doesn't actuate (that's traffic is never dropped) while still evaluating the
decision it would take in each cycle. This helps understand how the policy would
behave as the input signals change.

:::note

The `Dry Run` mode can also be toggled dynamically at runtime, without reloading
the policy.

:::

### Demo Video

The below demo video shows the basic concurrency limiter and workload
prioritization policy in action within Aperture Playground.

[![Demo Video](https://img.youtube.com/vi/m070bAvrDHM/0.jpg)](https://www.youtube.com/watch?v=m070bAvrDHM)
