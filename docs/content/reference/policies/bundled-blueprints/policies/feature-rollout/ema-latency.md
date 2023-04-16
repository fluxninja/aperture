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
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name="policy"
    type="
Object (policies/feature-rollout/base:schema:rollout_policy)"
    reference="../../../bundled-blueprints/policies/feature-rollout/base#rollout-policy"
    value="{'components': [], 'drivers': {'ema_latency_drivers': [{'backward': {'latency_tolerance_multiplier': 1.05}, 'ema': {'ema_window': '1500s', 'warmup_window': '60s'}, 'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'forward': {'latency_tolerance_multiplier': 1.05}, 'reset': {'latency_tolerance_multiplier': 1.25}}]}, 'evaluation_interval': '1s', 'load_shaper': {'flow_regulator_parameters': {'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'label_key': ''}, 'steps': [{'duration': '__REQUIRED_FIELD__', 'target_accept_percentage': '__REQUIRED_FIELD__'}]}, 'resources': {'flow_control': {'classifiers': []}}}"
    description='Parameters for the Feature Rollout policy.'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="dashboard"></a>

<ParameterDescription
    name="dashboard"
    type="
Object (policies/feature-rollout/base:param:dashboard)"
    reference="../../../bundled-blueprints/policies/feature-rollout/base#dashboard"
    value="{'datasource': {'filter_regex': '', 'name': '$datasource'}, 'refresh_interval': '5s', 'time_from': 'now-15m', 'time_to': 'now'}"
    description='Configuration for the Grafana dashboard accompanying this policy.'
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

<a id="load-shaper"></a>

<ParameterDescription
    name="load_shaper"
    type="
Object (aperture.spec.v1.FlowRegulatorDynamicConfig)"
    reference="../../../spec#flow-regulator-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Default configuration for flow regulator that can be updated at the runtime without shutting down the policy.'
/>

<!-- vale on -->

---
