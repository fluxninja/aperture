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
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/ema-latency`}>policies/feature-rollout/ema-latency</a>

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

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Parameters for the Feature Rollout policy.'
    type='Object (policies/feature-rollout/base:schema:rollout_policy)'
    reference='../../../bundled-blueprints/policies/feature-rollout/base#rollout-policy'
    value='{"components": [], "drivers": {"ema_latency_drivers": [{"criteria": {"backward": {"latency_tolerance_multiplier": 1.05}, "forward": {"latency_tolerance_multiplier": 1.05}, "reset": {"latency_tolerance_multiplier": 1.25}}, "ema": {"ema_window": "1500s", "warmup_window": "60s"}, "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}]}, "evaluation_interval": "1s", "load_ramp": {"regulator_parameters": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "resources": {"flow_control": {"classifiers": []}}}'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="dashboard"></a>

<ParameterDescription
    name='dashboard'
    description='Configuration for the Grafana dashboard accompanying this policy.'
    type='Object (policies/feature-rollout/base:param:dashboard)'
    reference='../../../bundled-blueprints/policies/feature-rollout/base#dashboard'
    value='{"datasource": {"filter_regex": "", "name": "$datasource"}, "refresh_interval": "5s", "time_from": "now-15m", "time_to": "now"}'
/>

<!-- vale on -->

---
