---
title: Service Protection with Average Latency Feedback
---

## Introduction

This policy detects traffic overloads and cascading failure build-up by
comparing the real-time latency with its exponential moving average. A gradient
controller calculates a proportional response to limit accepted concurrency. The
concurrency is reduced by a multiplicative factor when the service is
overloaded, and increased by an additive factor while the service is no longer
overloaded.

At a high level, this policy works as follows:

- Latency EMA-based overload detection: A Flux Meter is used to gather latency
  metrics from a [service control point](/concepts/flow-control/selector.md).
  The latency signal gets fed into an Exponential Moving Average (EMA) component
  to establish a long-term trend that can be compared to the current latency to
  detect overloads.
- Gradient Controller: Set point latency and current latency signals are fed to
  the gradient controller that calculates the proportional response to adjust
  the accepted concurrency (Control Variable).
- Integral Optimizer: When the service is detected to be in the normal state, an
  integral optimizer is used to additively increase the concurrency of the
  service in each execution cycle of the circuit. This design allows warming-up
  a service from an initial inactive state. This also protects applications from
  sudden spikes in traffic, as it sets an upper bound to the concurrency allowed
  on a service in each execution cycle of the circuit based on the observed
  incoming concurrency.
- Load Scheduler and Actuator: The Accepted Concurrency at the service is
  throttled by a
  [weighted-fair queuing scheduler](/concepts/flow-control/components/load-scheduler.md).
  The output of the adjustments to accepted concurrency made by gradient
  controller and optimizer logic are translated to a load multiplier that is
  synchronized with Aperture Agents through etcd. The load multiplier adjusts
  (increases or decreases) the token bucket fill rates based on the incoming
  concurrency observed at each agent.

:::info

Please see reference for the
[`AdaptiveLoadScheduler`](/reference/configuration/spec.md#adaptive-load-scheduler)
component that is used within this blueprint.

:::

:::info

See tutorials on
[Basic Service Protection](/use-cases/service-protection/protection.md) and
[Workload Prioritization](/use-cases/service-protection/prioritization.md) to
see this blueprint in use.

:::

<!-- Configuration Marker -->
```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration
<!-- vale off -->

Blueprint name: <a href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection/average-latency`}>policies/service-protection/average-latency</a>

<!-- vale on -->


### Parameters

## Dynamic Configuration



:::note

The following configuration parameters can be [dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md) at runtime, without reloading the policy.

:::



### Parameters

<!-- vale off -->

<a id="dry-run"></a>

<ParameterDescription
    name='dry_run'
    description='Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. It is useful for observing the behavior of load scheduler without disrupting any real traffic.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---