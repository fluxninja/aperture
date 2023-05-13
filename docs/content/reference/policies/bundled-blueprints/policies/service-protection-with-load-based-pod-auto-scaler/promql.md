---
title: Service Protection and Load-based Pod Auto-Scaler Based on PromQL Query
keywords:
  - blueprints
sidebar_position: 3
sidebar_label:
  Service Protection and Load-based Pod Auto Scaler Based on PromQL Query
---

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection-with-load-based-pod-auto-scaler/promql`}>policies/service-protection-with-load-based-pod-auto-scaler/promql</a>

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

<a id="policy-promql-query"></a>

<ParameterDescription
    name='policy.promql_query'
    description='PromQL query.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-setpoint"></a>

<ParameterDescription
    name='policy.setpoint'
    description='Setpoint.'
    type='Number (double)'
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
    type='Array of Object (policies/service-protection/promql:schema:overload_confirmation)'
    reference='../../../bundled-blueprints/policies/service-protection/promql#overload-confirmation'
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

<a id="policy-auto-scaling-promql-scale-in-controllers"></a>

<ParameterDescription
    name='policy.auto_scaling.promql_scale_in_controllers'
    description='List of scale in controllers.'
    type='Array of Object (policies/auto-scaling/pod-auto-scaler:schema:promql_scale_in_controller)'
    reference='../../../bundled-blueprints/policies/auto-scaling/pod-auto-scaler#promql-scale-in-controller'
    value='[{"gradient": {"slope": 1}, "query_string": "__REQUIRED_FIELD__", "threshold": 0.5}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-dry-run"></a>

<ParameterDescription
    name='policy.auto_scaling.dry_run'
    description='Dry run mode ensures that no scaling is invoked by this auto scaler.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-auto-scaling-dry-run-config-key"></a>

<ParameterDescription
    name='policy.auto_scaling.dry_run_config_key'
    description='Configuration key for overriding dry run setting through dynamic configuration.'
    type='string'
    reference=''
    value='"auto_scaling"'
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

<a id="load-scheduler"></a>

<ParameterDescription
    name='load_scheduler'
    description='Default configuration for load scheduler that can be updated at the runtime without shutting down the policy.'
    type='Object (aperture.spec.v1.LoadSchedulerDynamicConfig)'
    reference='../../../spec#load-scheduler-dynamic-config'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="auto-scaling"></a>

<ParameterDescription
    name='auto_scaling'
    description='Dry run mode ensures that no scaling is invoked by this auto scaler.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
