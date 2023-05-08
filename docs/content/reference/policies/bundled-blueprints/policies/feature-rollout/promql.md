---
title: Feature Rollout with PromQL Feedback
---

## Introduction

This policy rolls out new features based on the result of PromQL query as the
rollout criteria. The result of PromQL query is compared with thresholds to
determine conditions for advancing, reversing, or resetting the rollout to its
initial state. The rollout process consists of a series of steps that progress
if the feature is considered healthy.

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
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/promql`}>policies/feature-rollout/promql</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Parameters for the Feature Rollout policy.'
    type='Object (policies/feature-rollout/base:schema:rollout_policy)'
    reference='../../../bundled-blueprints/policies/feature-rollout/base#rollout-policy'
    value='{"components": [], "drivers": {"promql_drivers": [{"criteria": {"backward": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}, "forward": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}, "reset": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}}, "query_string": "__REQUIRED_FIELD__"}]}, "evaluation_interval": "1s", "load_ramp": {"regulator_parameters": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "policy_name": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}}'
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
