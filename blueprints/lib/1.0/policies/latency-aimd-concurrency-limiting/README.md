# Latency AIMD Concurrency Limiting Policy

This policy detect overloads/cascading failures by comparing the real-time
latency with it's exponential moving average. Gradient controller is then used
to calculate a proportional response that limits the accepted concurrency.
Concurrency is increased additively when the overload is no longer detected.

:::info

See tutorials on
[Basic Concurrency Limiting](/tutorials/integrations/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
and
[Workload Prioritization](/tutorials/integrations/flow-control/concurrency-limiting/workload-prioritization.md)
to see this blueprint in use.

:::

## Configuration

<!-- Configuration Marker -->

export const ParameterHeading = ({children}) => ( <span
style={{fontWeight: "bold"}}>{children}</span> );

export const WrappedDescription = ({children}) => ( <span
style={{wordWrap: "normal"}}>{children}</span> );

export const ParameterDescription = ({name, type, value, description}) => (

  <table class="blueprints-params">
  <tr>
    <td><ParameterHeading>Parameter</ParameterHeading></td>
    <td><code>{name}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Type</ParameterHeading></td>
    <td><code>{type}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Default Value</ParameterHeading></td>
    <td><code>{value != '' ? value : "REQUIRED VALUE"}</code></td>
  </tr>
  <tr>
    <td colspan="2" class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);

### Common

<ParameterDescription
    name="common.policy_name"
    type="string"
    value=''
    description='Name of the policy.' />

### Policy

<ParameterDescription
    name="policy.flux_meter"
    type="aperture.spec.v1.FluxMeter"
    value=''
    description='Flux Meter.' />

<ParameterDescription
    name="policy.classifiers"
    type="[]aperture.spec.v1.Classifier"
    value=''
    description='List of classification rules.' />

<ParameterDescription
    name="policy.components"
    type="[]aperture.spec.v1.Component"
    value=''
    description='List of additional circuit components.' />

#### Latency Baseliner

<ParameterDescription
    name="policy.latency_baseliner.ema"
    type="aperture.spec.v1.EMAParameters"
    value=''
    description='EMA parameters.' />

<ParameterDescription
    name="policy.latency_baseliner.latency_tolerance_multiplier"
    type="float64"
    value=''
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.' />

<ParameterDescription
    name="policy.latency_baseliner.latency_ema_limit_multiplier"
    type="float64"
    value=''
    description='Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.' />

#### Concurrency Controller

<ParameterDescription
    name="policy.concurrency_controller.flow_selector"
    type="aperture.spec.v1.FlowSelector"
    value=''
    description='Concurrency Limiter flow selector.' />

<ParameterDescription
    name="policy.concurrency_controller.scheduler"
    type="aperture.spec.v1.SchedulerParameters"
    value=''
    description='Scheduler parameters.' />

<ParameterDescription
    name="policy.concurrency_controller.gradient"
    type="aperture.spec.v1.GradientParameters"
    value=''
    description='Gradient parameters.' />

<ParameterDescription
    name="policy.concurrency_controller.alerter"
    type="aperture.spec.v1.AlerterParameters"
    value=''
    description='Whether tokens for workloads are computed dynamically or set statically by the user.' />

<ParameterDescription
    name="policy.concurrency_controller.concurrency_limit_multiplier"
    type="float64"
    value=''
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.' />

<ParameterDescription
    name="policy.concurrency_controller.concurrency_linear_increment"
    type="float64"
    value=''
    description='Linear increment to concurrency in each execution tick when the system is not in overloaded state.' />

<ParameterDescription
    name="policy.concurrency_controller.concurrency_sqrt_increment_multiplier"
    type="float64"
    value=''
    description='Scale factor to multiply square root of current accepted concurrrency. This, along with concurrency_linear_increment helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect).' />

<ParameterDescription
    name="policy.concurrency_controller.dynamic_config"
    type="aperture.v1.LoadActuatorDynamicConfig"
    value=''
    description='Dynamic configuration for concurrency controller.' />

### Dashboard

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    value=''
    description='Refresh interval for dashboard panels.' />

#### Datasource

<ParameterDescription
    name="dashboard.datasource.name"
    type="string"
    value=''
    description='Datasource name.' />

<ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="string"
    value=''
    description='Datasource filter regex.' />
