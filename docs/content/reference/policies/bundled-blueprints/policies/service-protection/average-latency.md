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
[`AdaptiveLoadScheduler`](/reference/policies/spec.md#adaptive-load-scheduler)
component that is used within this blueprint.

:::

:::info

See tutorials on
[Basic Service Protection](/applying-policies/service-protection/basic-service-protection.md)
and
[Workload Prioritization](/applying-policies/service-protection/workload-prioritization.md)
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
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection/average-latency`}>policies/service-protection/average-latency</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### common {#common}

<!-- vale on -->

<!-- vale off -->

<a id="common-policy-name"></a>

<ParameterDescription
    name='common.policy_name'
    description='Name of the policy.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### policy {#policy}

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
    description='List of overload confirmation criteria. Load scheduler can shed flows when all of the specified overload confirmation criteria are met.'
    type='Array of Object (overload_confirmation)'
    reference='#overload-confirmation'
    value='[{"operator": "__REQUIRED_FIELD__", "query_string": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

###### policy.service_protection_core.adaptive_load_scheduler {#policy-service-protection-core-adaptive-load-scheduler}

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-selectors"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.selectors'
    description='The selectors determine the flows that are protected by this policy.'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-scheduler"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.scheduler'
    description='Scheduler parameters.'
    type='Object (aperture.spec.v1.SchedulerParameters)'
    reference='../../../spec#scheduler-parameters'
    value='{"auto_tokens": true}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-gradient"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.gradient'
    description='Gradient Controller parameters.'
    type='Object (aperture.spec.v1.GradientControllerParameters)'
    reference='../../../spec#gradient-controller-parameters'
    value='{"max_gradient": 1, "min_gradient": 0.1, "slope": -1}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-alerter"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.alerter'
    description='Parameters for the Alerter that detects load throttling.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../../spec#alerter-parameters'
    value='{"alert_name": "Load Throttling Event"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-max-load-multiplier"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.max_load_multiplier'
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.'
    type='Number (double)'
    reference=''
    value='2'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-load-multiplier-linear-increment"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.load_multiplier_linear_increment'
    description='Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.'
    type='Number (double)'
    reference=''
    value='0.0025'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-service-protection-core-adaptive-load-scheduler-default-config"></a>

<ParameterDescription
    name='policy.service_protection_core.adaptive_load_scheduler.default_config'
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the'
    type='Object (aperture.spec.v1.LoadActuatorDynamicConfig)'
    reference='../../../spec#load-actuator-dynamic-config'
    value='{"dry_run": false}'
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
    value='"5s"'
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

### Schemas

<!-- vale off -->

#### overload_confirmation {#overload-confirmation}

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the overload confirmation criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
