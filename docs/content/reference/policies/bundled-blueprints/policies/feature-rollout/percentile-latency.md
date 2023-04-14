<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/percentile-latency`}>policies/feature-rollout/percentile-latency</a>

### Parameters

#### common {#common}

<a id="common-policy-name"></a> <ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

---

<a id="policy"></a> <ParameterDescription
    name="policy"
    type="
Object (policies/feature-rollout/base:schema:rollout_policy)"
    reference="../../../bundled-blueprints/policies/feature-rollout/base#rollout-policy"
    value="{'components': [], 'drivers': {'percentile_latency_drivers': [{'backward': {'threshold': '__REQUIRED_FIELD__'}, 'flux_meter': {'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'static_buckets': {'buckets': [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}, 'forward': {'threshold': '__REQUIRED_FIELD__'}, 'percentile': 95, 'reset': {'threshold': '__REQUIRED_FIELD__'}}]}, 'evaluation_interval': '1s', 'load_shaper': {'flow_regulator_parameters': {'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'label_key': ''}, 'steps': [{'duration': '__REQUIRED_FIELD__', 'target_accept_percentage': '__REQUIRED_FIELD__'}]}, 'resources': {'flow_control': {'classifiers': []}}}"
    description='Parameters for the Feature Rollout policy.' />

---

<a id="dashboard"></a> <ParameterDescription
    name="dashboard"
    type="
Object (policies/feature-rollout/base:param:dashboard)"
    reference="../../../bundled-blueprints/policies/feature-rollout/base#dashboard"
    value="{'datasource': {'filter_regex': '', 'name': '$datasource'}, 'refresh_interval': '5s', 'time_from': 'now-15m', 'time_to': 'now'}"
    description='Configuration for the Grafana dashboard accompanying this policy.' />

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<a id="load-shaper"></a> <ParameterDescription
    name="load_shaper"
    type="
Object (aperture.spec.v1.FlowRegulatorDynamicConfig)"
    reference="../../../spec#flow-regulator-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Default configuration for flow regulator that can be updated at the runtime without shutting down the policy.' />

---
