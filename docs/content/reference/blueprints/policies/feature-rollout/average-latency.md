---
title: Feature Rollout with Average Latency Feedback
---

## Introduction

This policy rolls out new features based on average of latency as the rollout
criteria. The average of latency is compared with thresholds to determine
conditions for advancing, reversing, or resetting the rollout to its initial
state. The rollout process consists of a series of steps that progress if the
feature is considered healthy.

:::info

This blueprint is the same as [`Feature Rollout blueprint`](base.md), with
relevant fields highlighted in the aperturectl generated values file.

:::

:::info

See the tutorial on
[Feature Rollout with Average Latency Feedback](/use-cases/feature-rollout/average-latency-feedback.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/average-latency`}>policies/feature-rollout/average-latency</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Configuration for the Feature Rollout policy.'
    type='Object (policies/feature-rollout/base:schema:rollout_policy)'
    reference='../../../blueprints/policies/feature-rollout/base#rollout-policy'
    value='{"components": [], "drivers": {"average_latency_drivers": [{"criteria": {"backward": {"threshold": "__REQUIRED_FIELD__"}, "forward": {"threshold": "__REQUIRED_FIELD__"}, "reset": {"threshold": "__REQUIRED_FIELD__"}}, "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}]}, "evaluation_interval": "1s", "load_ramp": {"sampler": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "policy_name": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}, "rollout": false}'
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

<!-- vale off -->

<a id="rollout"></a>

<ParameterDescription
    name='rollout'
    description='Start feature rollout. This setting can be updated at runtime without shutting down the policy. The feature rollout gets paused if this flag is set to false in the middle of a feature rollout.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="reset"></a>

<ParameterDescription
    name='reset'
    description='Reset feature rollout to the first step. This setting can be updated at the runtime without shutting down the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

---
