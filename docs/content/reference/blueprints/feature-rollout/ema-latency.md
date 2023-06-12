---
title: Feature Rollout with Exponential Moving Average Latency Feedback
---

## Introduction

This policy rolls out new features based on exponential moving average of
latency as the rollout criteria. The current average latency is compared with
multiples of exponential moving average latency to determine conditions for
advancing, reversing, or resetting the rollout to its initial state. The rollout
process consists of a series of steps that progress if the feature is considered
healthy.

:::info

This blueprint is the same as [`Feature Rollout blueprint`](base.md), with
relevant fields highlighted in the aperturectl generated values file.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/ema-latency`}>policies/feature-rollout/ema-latency</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Parameters for the Feature Rollout policy.'
    type='Object (policies/feature-rollout/base:schema:rollout_policy)'
    reference='../../../blueprints/feature-rollout/base#rollout-policy'
    value='{"components": [], "drivers": {"ema_latency_drivers": [{"criteria": {"backward": {"latency_tolerance_multiplier": 1.05}, "forward": {"latency_tolerance_multiplier": 1.05}, "reset": {"latency_tolerance_multiplier": 1.25}}, "ema": {"ema_window": "1500s", "warmup_window": "60s"}, "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}]}, "evaluation_interval": "1s", "load_ramp": {"sampler": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "policy_name": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}, "rollout": false}'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="dashboard"></a>

<ParameterDescription
    name='dashboard'
    description='Configuration for the Grafana dashboard accompanying this policy.'
    type='Object (policies/feature-rollout/base:param:dashboard)'
    reference='../../../blueprints/feature-rollout/base#dashboard'
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
