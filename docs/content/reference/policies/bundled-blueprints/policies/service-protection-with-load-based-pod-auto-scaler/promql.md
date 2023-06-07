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

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Configuration for the Service Protection policy.'
    type='Object (policies/service-protection/promql:param:policy)'
    reference='../../../bundled-blueprints/policies/service-protection/promql#policy'
    value='{"auto_scaling": {"dry_run": false, "periodic_decrease": {"period": "60s", "scale_in_percentage": 10}, "promql_scale_in_controllers": [], "promql_scale_out_controllers": [], "scaling_backend": {"kubernetes_replicas": {"kubernetes_object_selector": "__REQUIRED_FIELD__", "max_replicas": "__REQUIRED_FIELD__", "min_replicas": "__REQUIRED_FIELD__"}}, "scaling_parameters": {"scale_in_alerter": {"alert_name": "Auto-scaler is scaling in"}, "scale_in_cooldown": "40s", "scale_out_alerter": {"alert_name": "Auto-scaler is scaling out"}, "scale_out_cooldown": "30s"}}, "components": [], "evaluation_interval": "1s", "policy_name": "__REQUIRED_FIELD__", "promql_query": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}, "service_protection_core": {"adaptive_load_scheduler": {"alerter": {"alert_name": "Load Throttling Event"}, "gradient": {"max_gradient": 1, "min_gradient": 0.1, "slope": -1}, "load_multiplier_linear_increment": 0.0025, "load_scheduler": {"selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "max_load_multiplier": 2}, "dry_run": false, "overload_confirmations": []}, "setpoint": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="dashboard"></a>

<ParameterDescription
    name='dashboard'
    description='Configuration for the Grafana dashboard accompanying this policy.'
    type='Object (policies/service-protection/promql:param:dashboard)'
    reference='../../../bundled-blueprints/policies/service-protection/promql#dashboard'
    value='{"datasource": {"filter_regex": "", "name": "$datasource"}, "extra_filters": {}, "refresh_interval": "15s", "time_from": "now-15m", "time_to": "now", "variant_name": "PromQL Output"}'
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
