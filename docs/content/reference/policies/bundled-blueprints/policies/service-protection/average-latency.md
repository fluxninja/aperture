---
title: Service Protection with Average Latency Feedback
---

## Introduction

This policy detects traffic overloads and cascading failure build-up by
comparing the real-time latency with its exponential moving average. A gradient
controller calculates a proportional response to limit accepted concurrency. The
concurrency is reduced by a multiplicative factor when the service is
overloaded, and increased by an additive factor while the service is no longer
overloaded.

:::info

Please see reference for the
[`AdaptiveLoadScheduler`](/reference/policies/spec.md#adaptive-load-scheduler)
component that is used within this blueprint.

:::

:::info

See tutorials on
[Basic Concurrency Limiting](/tutorials/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
and
[Workload Prioritization](/tutorials/flow-control/concurrency-limiting/workload-prioritization.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/service-protection/average-latency`}>policies/service-protection/average-latency</a>

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

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a>

<ParameterDescription
    name='dashboard.refresh_interval'
    description='Refresh interval for dashboard panels.'
    type='string'
    reference=''
    value='"5s"'
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

<!-- vale off -->

#### policy {#policy}

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
    description='List of additional resources.'
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

<a id="policy-service-protection-core"></a>

<ParameterDescription
    name='policy.service_protection_core'
    description='Core parameters for Service Protection policy.'
    type='Object (policies/service-protection/base:schema:service_protection_core)'
    reference='../../../bundled-blueprints/policies/service-protection/base#service-protection-core'
    value='{"adaptive_load_scheduler": {"alerter": {"alert_name": "Load Shed Event"}, "default_config": {"dry_run": false}, "flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "gradient": {"max_gradient": 1, "min_gradient": 0.1, "slope": -1}, "load_multiplier_linear_increment": 0.0025, "max_load_multiplier": 2, "scheduler": {"auto_tokens": true}}, "overload_confirmations": [{"operator": "__REQUIRED_FIELD__", "query_string": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}]}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-flux-meter"></a>

<ParameterDescription
    name='policy.flux_meter'
    description='Flux Meter.'
    type='Object (aperture.spec.v1.FluxMeter)'
    reference='../../../spec#flux-meter'
    value='{"flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}}'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.latency_baseliner {#policy-latency-baseliner}

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
