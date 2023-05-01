---
title: RabbitMQ Queue Buildup Policy
---

## Introduction

This policy detects RabbitMQ queue buildup by looking at the number of messages
in "ready" state. A gradient controller calculates a proportional response to
limit accepted concurrency, which is increased additively when the overload is
no longer detected.

## Build Instructions

Verify that the Aperture binary is
[built](/reference/aperturectl/build/agent/agent.md) with the `rabbitmqreceiver`
extension enabled.

```yaml
bundled_extensions:
  - integrations/otel/rabbitmqreceiver
```

Use the following agent configuration to allow the agent to collect RabbitMQ
metrics.

```yaml
otel:
  custom_metrics:
    rabbitmq:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - rabbitmq
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        rabbitmq:
          collection_interval: 1s
          endpoint: http://rabbitmq.rabbitmq.svc.cluster.local:15672
          password: secretpassword
          username: admin
```

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/rabbitmq-queue-buildup`}>policies/rabbitmq-queue-buildup</a>

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

<!-- vale off -->

<a id="common-queue-name"></a>

<ParameterDescription
    name='common.queue_name'
    description='Name of the queue to watch for buildup.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-classifiers"></a>

<ParameterDescription
    name='policy.classifiers'
    description='List of classification rules.'
    type='Array of Object (aperture.spec.v1.Classifier)'
    reference='../../spec#classifier'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../spec#component'
    value='[]'
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
    reference='../../spec#e-m-a-parameters'
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

<!-- vale off -->

##### policy.concurrency_controller {#policy-concurrency-controller}

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-selectors"></a>

<ParameterDescription
    name='policy.concurrency_controller.selectors'
    description='Concurrency Limiter flow selectors.'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-scheduler"></a>

<ParameterDescription
    name='policy.concurrency_controller.scheduler'
    description='Scheduler parameters.'
    type='Object (aperture.spec.v1.SchedulerParameters)'
    reference='../../spec#scheduler-parameters'
    value='{"auto_tokens": true}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-gradient"></a>

<ParameterDescription
    name='policy.concurrency_controller.gradient'
    description='Gradient Controller parameters.'
    type='Object (aperture.spec.v1.GradientControllerParameters)'
    reference='../../spec#gradient-controller-parameters'
    value='{"max_gradient": 1, "min_gradient": 0.1, "slope": -1}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-alerter"></a>

<ParameterDescription
    name='policy.concurrency_controller.alerter'
    description='Whether tokens for workloads are computed dynamically or set statically by the user.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../spec#alerter-parameters'
    value='{"alert_name": "Load Shed Event"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-max-load-multiplier"></a>

<ParameterDescription
    name='policy.concurrency_controller.max_load_multiplier'
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.'
    type='Number (double)'
    reference=''
    value='2'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-queue-buildup-setpoint"></a>

<ParameterDescription
    name='policy.concurrency_controller.queue_buildup_setpoint'
    description='Queue buildup setpoint in number of messages.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-load-multiplier-linear-increment"></a>

<ParameterDescription
    name='policy.concurrency_controller.load_multiplier_linear_increment'
    description='Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.'
    type='Number (double)'
    reference=''
    value='0.0025'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-default-config"></a>

<ParameterDescription
    name='policy.concurrency_controller.default_config'
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.'
    type='Object (aperture.spec.v1.LoadActuatorDynamicConfig)'
    reference='../../spec#load-actuator-dynamic-config'
    value='{"dry_run": false}'
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

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="concurrency-controller"></a>

<ParameterDescription
    name='concurrency_controller'
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.'
    type='Object (aperture.spec.v1.LoadActuatorDynamicConfig)'
    reference='../../spec#load-actuator-dynamic-config'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
