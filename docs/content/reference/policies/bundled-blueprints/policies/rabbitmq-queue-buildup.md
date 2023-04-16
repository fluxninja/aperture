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

<a id="common-policy-name"></a> <ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<!-- vale on -->

<!-- vale off -->

<a id="common-queue-name"></a> <ParameterDescription
    name="common.queue_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the queue to watch for buildup.' />

<!-- vale on -->

---

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-classifiers"></a> <ParameterDescription
    name="policy.classifiers"
    type="
Array of
Object (aperture.spec.v1.Classifier)"
    reference="../../spec#classifier"
    value="[]"
    description='List of classification rules.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a> <ParameterDescription
    name="policy.components"
    type="
Array of
Object (aperture.spec.v1.Component)"
    reference="../../spec#component"
    value="[]"
    description='List of additional circuit components.' />

<!-- vale on -->

<!-- vale off -->

##### policy.latency_baseliner {#policy-latency-baseliner}

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-ema"></a> <ParameterDescription
    name="policy.latency_baseliner.ema"
    type="
Object (aperture.spec.v1.EMAParameters)"
    reference="../../spec#e-m-a-parameters"
    value="{'correction_factor_on_max_envelope_violation': 0.95, 'ema_window': '1500s', 'warmup_window': '60s'}"
    description='EMA parameters.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-latency-tolerance-multiplier"></a>
<ParameterDescription
    name="policy.latency_baseliner.latency_tolerance_multiplier"
    type="
Number (double)"
    reference=""
    value="1.1"
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-latency-baseliner-latency-ema-limit-multiplier"></a>
<ParameterDescription
    name="policy.latency_baseliner.latency_ema_limit_multiplier"
    type="
Number (double)"
    reference=""
    value="2"
    description='Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.' />

<!-- vale on -->

<!-- vale off -->

##### policy.concurrency_controller {#policy-concurrency-controller}

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-flow-selector"></a> <ParameterDescription
    name="policy.concurrency_controller.flow_selector"
    type="
Object (aperture.spec.v1.FlowSelector)"
    reference="../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='Concurrency Limiter flow selector.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-scheduler"></a> <ParameterDescription
    name="policy.concurrency_controller.scheduler"
    type="
Object (aperture.spec.v1.SchedulerParameters)"
    reference="../../spec#scheduler-parameters"
    value="{'auto_tokens': True}"
    description='Scheduler parameters.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-gradient"></a> <ParameterDescription
    name="policy.concurrency_controller.gradient"
    type="
Object (aperture.spec.v1.GradientControllerParameters)"
    reference="../../spec#gradient-controller-parameters"
    value="{'max_gradient': 1, 'min_gradient': 0.1, 'slope': -1}"
    description='Gradient Controller parameters.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-alerter"></a> <ParameterDescription
    name="policy.concurrency_controller.alerter"
    type="
Object (aperture.spec.v1.AlerterParameters)"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Load Shed Event'}"
    description='Whether tokens for workloads are computed dynamically or set statically by the user.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-max-load-multiplier"></a>
<ParameterDescription
    name="policy.concurrency_controller.max_load_multiplier"
    type="
Number (double)"
    reference=""
    value="2"
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-queue-buildup-setpoint"></a>
<ParameterDescription
    name="policy.concurrency_controller.queue_buildup_setpoint"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Queue buildup setpoint in number of messages.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-load-multiplier-linear-increment"></a>
<ParameterDescription
    name="policy.concurrency_controller.load_multiplier_linear_increment"
    type="
Number (double)"
    reference=""
    value="0.0025"
    description='Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-controller-default-config"></a> <ParameterDescription
    name="policy.concurrency_controller.default_config"
    type="
Object (aperture.spec.v1.LoadActuatorDynamicConfig)"
    reference="../../spec#load-actuator-dynamic-config"
    value="{'dry_run': False}"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a> <ParameterDescription
    name="dashboard.refresh_interval"
    type="
string"
    reference=""
    value="'5s'"
    description='Refresh interval for dashboard panels.' />

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a> <ParameterDescription
    name="dashboard.time_from"
    type="
string"
    reference=""
    value="'now-15m'"
    description='From time of dashboard.' />

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a> <ParameterDescription
    name="dashboard.time_to"
    type="
string"
    reference=""
    value="'now'"
    description='To time of dashboard.' />

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-name"></a> <ParameterDescription
    name="dashboard.datasource.name"
    type="
string"
    reference=""
    value="'$datasource'"
    description='Datasource name.' />

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-filter-regex"></a> <ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="
string"
    reference=""
    value="''"
    description='Datasource filter regex.' />

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

<a id="concurrency-controller"></a> <ParameterDescription
    name="concurrency_controller"
    type="
Object (aperture.spec.v1.LoadActuatorDynamicConfig)"
    reference="../../spec#load-actuator-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />

<!-- vale on -->

---
