---
title: RabbitMQ Queue Buildup Policy
---

## Introduction

This policy detects RabbitMQ queue buildup by looking at number of messages in
"ready" state. Gradient controller is then used to calculate a proportional
response that limits the accepted concurrency. Concurrency is increased
additively when the overload is no longer detected.

## Build Instructions

Make sure that the aperture binary is
[built](reference/aperturectl/build/agent/) with `rabbitmqreceiver` extension
enabled.

```yaml
bundled_extensions:
  - integrations/otel/rabbitmqreceiver
```

Use the following agent configuration to allow agent to collect RabbitMQ
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

## Configuration

<!-- Configuration Marker -->

```mdx-code-block

export const ParameterHeading = ({children}) => (
  <span style={{fontWeight: "bold"}}>{children}</span>
);

export const WrappedDescription = ({children}) => (
  <span style={{wordWrap: "normal"}}>{children}</span>
);

export const RefType = ({type, reference}) => (
  <a href={reference}>{type}</a>
);

export const ParameterDescription = ({name, type, reference, value, description}) => (
  <table class="blueprints-params">
  <tr>
    <td><ParameterHeading>Parameter</ParameterHeading></td>
    <td><code>{name}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Type</ParameterHeading></td>
    <td><em>{reference == "" ? type : <RefType type={type} reference={reference} />}</em></td>
  </tr>
  <tr>
    <td class="blueprints-default-heading"><ParameterHeading>Default Value</ParameterHeading></td>
    <td><code>{value}</code></td>
  </tr>
  <tr>
    <td class="blueprints-description"><ParameterHeading>Description</ParameterHeading></td>
    <td class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
```

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
```

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/rabbitmq-queue-buildup`}>policies/rabbitmq-queue-buildup</a>

<h3 class="blueprints-h3">Common</h3>

<ParameterDescription
    name="common.policy_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<ParameterDescription
    name="common.queue_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the queue to watch for buildup.' />

<h3 class="blueprints-h3">Policy</h3>

<ParameterDescription
    name="policy.classifiers"
    type="[]aperture.spec.v1.Classifier"
    reference="../../spec#classifier"
    value="[]"
    description='List of classification rules.' />

<ParameterDescription
    name="policy.components"
    type="[]aperture.spec.v1.Component"
    reference="../../spec#component"
    value="[]"
    description='List of additional circuit components.' />

<h4 class="blueprints-h4">Latency Baseliner</h4>

<ParameterDescription
    name="policy.latency_baseliner.ema"
    type="aperture.spec.v1.EMAParameters"
    reference="../../spec#e-m-a-parameters"
    value="{'correction_factor_on_max_envelope_violation': 0.95, 'ema_window': '1500s', 'warmup_window': '60s'}"
    description='EMA parameters.' />

<ParameterDescription
    name="policy.latency_baseliner.latency_tolerance_multiplier"
    type="float64"
    reference=""
    value="1.1"
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.' />

<ParameterDescription
    name="policy.latency_baseliner.latency_ema_limit_multiplier"
    type="float64"
    reference=""
    value="2"
    description='Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.' />

<h4 class="blueprints-h4">Concurrency Controller</h4>

<ParameterDescription
    name="policy.concurrency_controller.flow_selector"
    type="aperture.spec.v1.FlowSelector"
    reference="../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='Concurrency Limiter flow selector.' />

<ParameterDescription
    name="policy.concurrency_controller.flow_selector.service_selector.service"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Service Name.' />

<ParameterDescription
    name="policy.concurrency_controller.flow_selector.flow_matcher.control_point"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Control Point Name.' />

<ParameterDescription
    name="policy.concurrency_controller.scheduler"
    type="aperture.spec.v1.SchedulerParameters"
    reference="../../spec#scheduler-parameters"
    value="{'auto_tokens': True}"
    description='Scheduler parameters.' />

<ParameterDescription
    name="policy.concurrency_controller.scheduler.auto_tokens"
    type="bool"
    reference=""
    value="true"
    description='Automatically estimate cost (tokens) for workload requests.' />

<ParameterDescription
    name="policy.concurrency_controller.gradient"
    type="aperture.spec.v1.GradientControllerParameters"
    reference="../../spec#gradient-controller-parameters"
    value="{'max_gradient': 1, 'min_gradient': 0.1, 'slope': -1}"
    description='Gradient Controller parameters.' />

<ParameterDescription
    name="policy.concurrency_controller.alerter"
    type="aperture.spec.v1.AlerterParameters"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Load Shed Event'}"
    description='Whether tokens for workloads are computed dynamically or set statically by the user.' />

<ParameterDescription
    name="policy.concurrency_controller.max_load_multiplier"
    type="float64"
    reference=""
    value="2"
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.' />

<ParameterDescription
    name="policy.concurrency_controller.queue_buildup_setpoint"
    type="float64"
    reference=""
    value="1000"
    description='Queue buildup setpoint in number of messages.' />

<ParameterDescription
    name="policy.concurrency_controller.load_multiplier_linear_increment"
    type="float64"
    reference=""
    value="0.0025"
    description='Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.' />

<ParameterDescription
    name="policy.concurrency_controller.default_config"
    type="aperture.spec.v1.LoadActuatorDynamicConfig"
    reference="../../spec#load-actuator-dynamic-config"
    value="{'dry_run': False}"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />

<h3 class="blueprints-h3">Dashboard</h3>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    reference=""
    value="'5s'"
    description='Refresh interval for dashboard panels.' />

<ParameterDescription
    name="dashboard.time_from"
    type="string"
    reference=""
    value="'now-15m'"
    description='From time of dashboard.' />

<ParameterDescription
    name="dashboard.time_to"
    type="string"
    reference=""
    value="'now'"
    description='To time of dashboard.' />

<h4 class="blueprints-h4">Datasource</h4>

<ParameterDescription
    name="dashboard.datasource.name"
    type="string"
    reference=""
    value="'$datasource'"
    description='Datasource name.' />

<ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="string"
    reference=""
    value="''"
    description='Datasource filter regex.' />

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

<h3 class="blueprints-h3">Dynamic Configuration</h3>

<ParameterDescription
    name="concurrency_controller"
    type="aperture.spec.v1.LoadActuatorDynamicConfig"
    reference="../../spec#load-actuator-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />
