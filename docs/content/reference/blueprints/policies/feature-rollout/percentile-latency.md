---
title: Feature Rollout with Percentile Latency Feedback
---

## Introduction

This policy rolls out new features based on percentile of latency as the rollout
criteria. The percentile of latency is compared with thresholds to determine
conditions for advancing, reversing, or resetting the rollout to its initial
state. The rollout process consists of a series of steps that progress if the
feature is considered healthy.

:::info

This blueprint is the same as [`Feature Rollout blueprint`](base.md), with
relevant fields highlighted in the aperturectl generated values file.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/percentile-latency`}>policies/feature-rollout/percentile-latency</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Parameters for the Feature Rollout policy.'
    type='Object (policies/feature-rollout/base:schema:rollout_policy)'
    reference='../../../blueprints/policies/feature-rollout/base#rollout-policy'
    value='{"components": [], "drivers": {"percentile_latency_drivers": [{"criteria": {"backward": {"threshold": "__REQUIRED_FIELD__"}, "forward": {"threshold": "__REQUIRED_FIELD__"}, "reset": {"threshold": "__REQUIRED_FIELD__"}}, "flux_meter": {"selector": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}], "static_buckets": {"buckets": [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}, "percentile": 95}]}, "evaluation_interval": "1s", "load_ramp": {"sampler": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "policy_name": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}, "rollout": false}'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="dashboard"></a>

<ParameterDescription
    name='dashboard'
    description='Configuration for the Grafana dashboard accompanying this policy.'
    type='Object (policies/feature-rollout/base:param:dashboard)'
    reference='../../../blueprints/policies/feature-rollout/base#dashboard'
    value='{"datasource": {"filter_regex": "", "name": "$datasource"}, "extra_filters": {}, "refresh_interval": "5s", "time_from": "now-15m", "time_to": "now", "title": "Aperture Feature Rollout"}'
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

<a id="pass-through-label-values"></a>

<ParameterDescription
    name='pass_through_label_values'
    description='Specify certain label values to be always accepted by the _Sampler_ regardless of accept percentage. This configuration can be updated at the runtime without shutting down the policy.'
    type='Array of string'
    reference=''
    value='["__REQUIRED_FIELD__"]'
/>

<!-- vale on -->

---
