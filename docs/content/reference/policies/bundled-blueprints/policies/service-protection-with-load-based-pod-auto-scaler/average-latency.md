---
title:
  Service Protection and Load-based Pod Auto-Scaler with Average Latency
  Feedback
---

## Introduction

This policy detects traffic overloads and cascading failure build-up by
comparing the real-time latency with its exponential moving average. A gradient
controller calculates a proportional response to limit accepted concurrency. The
concurrency is reduced by a multiplicative factor when the service is
overloaded, and increased by an additive factor while the service is no longer
overloaded. An auto-scaler controller is used to dynamically adjust the number
of instances or resources allocated to a service based on workload demands. The
basic service protection policy protects the service from sudden traffic spikes.
It is necessary to scale the service to meet demand in case of a persistent
change in load.

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
- An _Auto Scaler_ that adjusts the number of replicas of the Kubernetes
  Deployment for the service.
- Load-based scale-out is done based on `OBSERVED_LOAD_MULTIPLIER` signal from
  the blueprint. This signal measures the fraction of traffic that the _Load
  Scheduler_ is throttling into a queue. The _Auto Scaler_ is configured to
  scale-out based on a _Gradient Controller_ using this signal and a setpoint of
  1.0.

:::info

Please see reference for the
[`AdaptiveLoadScheduler`](/reference/policies/spec.md#adaptive-load-scheduler)
and [`AutoScale`](/reference/policies/spec.md#auto-scale) components that are
used within this blueprint.

:::

:::info

See tutorials on
[Load-based Auto Scaling](/use-cases/auto-scale/load-based-auto-scaling/load-based-auto-scaling.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection-with-load-based-pod-auto-scaler/average-latency`}>policies/service-protection-with-load-based-pod-auto-scaler/average-latency</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-policy-name"></a>

<ParameterDescription
    name='policy.policy_name'
    description='Name of the policy.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../../spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-resources"></a>

<ParameterDescription
    name='policy.resources'
    description='Additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../../spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-evaluation-interval"></a>

<ParameterDescription
    name='policy.evaluation_interval'
    description='The interval between successive evaluations of the Circuit.'
    type='string'
    reference=''
    value='"1s"'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.service_protection_core {#policy-service-protection-core}

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-overload-confirmations"></a>

<ParameterDescription
    name='policy.service_protection_core.overload_confirmations'
    description='List of overload confirmation criteria. Load scheduler can throttle flows when all of the specified overload confirmation criteria are met.'
    type='Array of Object (policies/service-protection/average-latency:schema:overload_confirmation)'
    reference='../../../bundled-blueprints/policies/service-protection/average-latency#overload-confirmation'
    value='[{"operator": "__REQUIRED_FIELD__", "query_string": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler'
    description='Parameters for Adaptive Load Scheduler.'
    type='Object (aperture.spec.v1.AdaptiveLoadSchedulerParameters)'
    reference='../../../spec#adaptive-load-scheduler-parameters'
    value='{"alerter": {"alert_name": "Load Throttling Event"}, "gradient": {"max_gradient": 1, "min_gradient": 0.1, "slope": -1}, "load_multiplier_linear_increment": 0.0025, "load_scheduler": {"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "max_load_multiplier": 2}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-dry-run"></a>

<ParameterDescription
    name='policy.service_protection_core.dry_run'
    description='Default configuration for setting dry run mode on Load Scheduler. In dry run mode, the Load Scheduler acts as a passthrough and does not throttle flows. This config can be updated at runtime without restarting the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.latency_baseliner {#policy-latency-baseliner}

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-flux-meter"></a>

<ParameterDescription
    name='policy.latency_baseliner.flux_meter'
    description='Flux Meter defines the scope of latency measurements.'
    type='Object (aperture.spec.v1.FluxMeter)'
    reference='../../../spec#flux-meter'
    value='{"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-ema"></a>

<ParameterDescription
    name='policy.latency_baseliner.ema'
    description='EMA parameters.'
    type='Object (aperture.spec.v1.EMAParameters)'
    reference='../../../spec#e-m-a-parameters'
    value='{"correction_factor_on_max_envelope_violation": 0.95, "ema_window": "1500s", "warmup_window": "60s"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='policy.latency_baseliner.latency_tolerance_multiplier'
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.'
    type='Number (double)'
    reference=''
    value='1.1'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-latency-ema-limit-multiplier"></a>

<ParameterDescription
    name='policy.latency_baseliner.latency_ema_limit_multiplier'
    description='Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.'
    type='Number (double)'
    reference=''
    value='2'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.auto_scaling {#policy-auto-scaling}

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-kubernetes-replicas"></a>

<ParameterDescription
    name='policy.auto_scaling.kubernetes_replicas'
    description='Kubernetes replicas scaling backend.'
    type='Object (aperture.spec.v1.AutoScalerScalingBackendKubernetesReplicas)'
    reference='../../../spec#auto-scaler-scaling-backend-kubernetes-replicas'
    value='{"kubernetes_object_selector": "__REQUIRED_FIELD__", "max_replicas": "__REQUIRED_FIELD__", "min_replicas": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-dry-run"></a>

<ParameterDescription
    name='policy.auto_scaling.dry_run'
    description='Dry run mode ensures that no scaling is invoked by this auto scaler. This config can be updated at runtime without restarting the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-scaling-parameters"></a>

<ParameterDescription
    name='policy.auto_scaling.scaling_parameters'
    description='Parameters that define the scaling behavior.'
    type='Object (aperture.spec.v1.AutoScalerScalingParameters)'
    reference='../../../spec#auto-scaler-scaling-parameters'
    value='{"scale_in_alerter": {"alert_name": "Auto-scaler is scaling in"}, "scale_in_cooldown": "40s", "scale_out_alerter": {"alert_name": "Auto-scaler is scaling out"}, "scale_out_cooldown": "30s"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-disable-periodic-scale-in"></a>

<ParameterDescription
    name='policy.auto_scaling.disable_periodic_scale_in'
    description='Disable periodic scale in.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

###### policy.auto_scaling.periodic_decrease {#policy-auto-scaling-periodic-decrease}

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-periodic-decrease-period"></a>

<ParameterDescription
    name='policy.auto_scaling.periodic_decrease.period'
    description='Period for periodic scale in.'
    type='string'
    reference=''
    value='"60s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-periodic-decrease-scale-in-percentage"></a>

<ParameterDescription
    name='policy.auto_scaling.periodic_decrease.scale_in_percentage'
    description='Percentage of replicas to scale in.'
    type='Number (double)'
    reference=''
    value='10'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a>

<ParameterDescription
    name='dashboard.refresh_interval'
    description='Refresh interval for dashboard panels.'
    type='string'
    reference=''
    value='"15s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a>

<ParameterDescription
    name='dashboard.time_from'
    description='From time of dashboard.'
    type='string'
    reference=''
    value='"now-15m"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name='dashboard.time_to'
    description='To time of dashboard.'
    type='string'
    reference=''
    value='"now"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-extra-filters"></a>

<ParameterDescription
    name='dashboard.extra_filters'
    description='Additional filters to pass to each query to Grafana datasource.'
    type='Object (map[string]string)'
    reference='#map-string-string'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-name"></a>

<ParameterDescription
    name='dashboard.datasource.name'
    description='Datasource name.'
    type='string'
    reference=''
    value='"$datasource"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-filter-regex"></a>

<ParameterDescription
    name='dashboard.datasource.filter_regex'
    description='Datasource filter regex.'
    type='string'
    reference=''
    value='""'
/>

<!-- vale on -->

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="dry-run"></a>

<ParameterDescription
    name='dry_run'
    description='Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. The Auto Scaler does not perform any scaling in dry mode. This mode is useful for observing the behavior of load scheduler and auto scaler without disrupting any real deployment or traffic.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
