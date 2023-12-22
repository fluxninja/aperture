<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../apertureVersion.js'
import {ParameterDescription} from '../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/concurrency-limiting`}>concurrency-limiting</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../configuration/spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-limiter"></a>

<ParameterDescription
    name='policy.concurrency_limiter'
    description='Concurrency limiter.'
    type='Object (aperture.spec.v1.ConcurrencyLimiter)'
    reference='../configuration/spec#concurrency-limiter'
    value='{"max_concurrency": "__REQUIRED_FIELD__", "parameters": {"limit_by_label_key": "limit_key", "max_idle_time": "__REQUIRED_FIELD__", "max_inflight_duration": "__REQUIRED_FIELD__"}, "request_parameters": {}}'
/>

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

<a id="policy-resources"></a>

<ParameterDescription
    name='policy.resources'
    description='Additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../configuration/spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.concurency_limiter {#policy-concurency-limiter}

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurency-limiter-parameters"></a>

<ParameterDescription
    name='policy.concurency_limiter.parameters'
    description='Parameters.'
    type='Object (aperture.spec.v1.ConcurrencyLimiterParameters)'
    reference='../configuration/spec#concurrency-limiter-parameters'
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurency-limiter-request-parameters"></a>

<ParameterDescription
    name='policy.concurency_limiter.request_parameters'
    description='Request Parameters.'
    type='Object (aperture.spec.v1.ConcurrencyLimiterRequestParameters)'
    reference='../configuration/spec#concurrency-limiter-request-parameters'
    value='null'
/>

<!-- vale on -->

---
