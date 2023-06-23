---
title:
  Service Protection based on Average Latency Feedback with overload
  confirmation from JMX metrics
keywords:
  - blueprints
sidebar_position: 3
---

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection/jmx`}>policies/service-protection/jmx</a>

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
    type='Array of Object (overload_confirmation)'
    reference='#overload-confirmation'
    value='[]'
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

##### policy.auto_scaling {#policy-auto-scaling}

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-dry-run"></a>

<ParameterDescription
    name='policy.auto_scaling.dry_run'
    description='Dry run mode ensures that no scaling is invoked by the auto scaler escalation. This config can be updated at runtime without restarting the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-promql-scale-out-controllers"></a>

<ParameterDescription
    name='policy.auto_scaling.promql_scale_out_controllers'
    description='List of scale out controllers.'
    type='Array of Object (promql_scale_out_controller)'
    reference='#promql-scale-out-controller'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-promql-scale-in-controllers"></a>

<ParameterDescription
    name='policy.auto_scaling.promql_scale_in_controllers'
    description='List of scale in controllers.'
    type='Array of Object (promql_scale_in_controller)'
    reference='#promql-scale-in-controller'
    value='[]'
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

<a id="policy-auto-scaling-scaling-backend"></a>

<ParameterDescription
    name='policy.auto_scaling.scaling_backend'
    description='Scaling backend for the policy.'
    type='Object (aperture.spec.v1.AutoScalerScalingBackend)'
    reference='../../../spec#auto-scaler-scaling-backend'
    value='{"kubernetes_replicas": {"kubernetes_object_selector": "__REQUIRED_FIELD__", "max_replicas": "__REQUIRED_FIELD__", "min_replicas": "__REQUIRED_FIELD__"}}'
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
    value='"15s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a>

<ParameterDescription
    name='dashboard.time_from'
    description='Time from of dashboard.'
    type='string'
    reference=''
    value='"now-15m"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name='dashboard.time_to'
    description='Time to of dashboard.'
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
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the overload confirmation criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="overload-confirmation-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the overload confirmation criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

---

<!-- vale off -->

#### promql_scale_out_controller {#promql-scale-out-controller}

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-threshold"></a>

<ParameterDescription
    name='threshold'
    description='Threshold for the controller.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-gradient"></a>

<ParameterDescription
    name='gradient'
    description='Gradient parameters for the controller.'
    type='Object (aperture.spec.v1.IncreasingGradientParameters)'
    reference='../../../spec#increasing-gradient-parameters'
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-alerter"></a>

<ParameterDescription
    name='alerter'
    description='Alerter parameters for the controller.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../../spec#alerter-parameters'
    value='null'
/>

<!-- vale on -->

---

<!-- vale off -->

#### promql_scale_in_controller {#promql-scale-in-controller}

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-threshold"></a>

<ParameterDescription
    name='threshold'
    description='Threshold for the controller.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-gradient"></a>

<ParameterDescription
    name='gradient'
    description='Gradient parameters for the controller.'
    type='Object (aperture.spec.v1.DecreasingGradientParameters)'
    reference='../../../spec#decreasing-gradient-parameters'
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-alerter"></a>

<ParameterDescription
    name='alerter'
    description='Alerter parameters for the controller.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../../spec#alerter-parameters'
    value='null'
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
    description='Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. It is useful for observing the behavior of load scheduler without disrupting any real traffic.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
