# Policy

Policy

## Content negotiation

### URI Schemes

- http

### Consumes

- application/json

### Produces

- application/json

## All endpoints

### policy_configuration

| Method | URI                  | Name              | Summary |
| ------ | -------------------- | ----------------- | ------- |
| POST   | /aperture-controller | [policy](#policy) |         |

## Paths

### <span id="policy"></span> policy (_Policy_)

```
POST /aperture-controller
```

#### Parameters

| Name | Source | Type              | Go type         | Separator | Required | Default | Description |
| ---- | ------ | ----------------- | --------------- | --------- | :------: | ------- | ----------- |
| body | `body` | [Policy](#policy) | `models.Policy` |           |    ✓     |         |             |

#### All responses

| Code               | Status | Description            | Has headers | Schema                       |
| ------------------ | ------ | ---------------------- | :---------: | ---------------------------- |
| [200](#policy-200) | OK     | A successful response. |             | [schema](#policy-200-schema) |

#### Responses

##### <span id="policy-200"></span> 200 - A successful response.

Status: OK

###### <span id="policy-200-schema"></span> Schema

## Models

### <span id="adaptive-load-scheduler"></span> AdaptiveLoadScheduler

> The _Adaptive Load Scheduler_ adjusts the accepted token rate based on the
> deviation of the input signal from the setpoint.

**Properties**

| Name                                                                                           | Type                                                                   | Go type                           | Required | Default | Description                                                                                                                                         | Example |
| ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------- | --------------------------------- | :------: | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| dry_run                                                                                        | boolean                                                                | `bool`                            |          |         | Decides whether to run the load scheduler in dry-run mode. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. |
| It is useful for observing the behavior of load scheduler without disrupting any real traffic. |                                                                        |
| dry_run_config_key                                                                             | string                                                                 | `string`                          |          |         | Configuration key for setting dry run mode through dynamic configuration.                                                                           |         |
| in_ports                                                                                       | [AdaptiveLoadSchedulerIns](#adaptive-load-scheduler-ins)               | `AdaptiveLoadSchedulerIns`        |          |         | Collection of input ports for the _Adaptive Load Scheduler_ component.                                                                              |         |
| out_ports                                                                                      | [AdaptiveLoadSchedulerOuts](#adaptive-load-scheduler-outs)             | `AdaptiveLoadSchedulerOuts`       |          |         | Collection of output ports for the _Adaptive Load Scheduler_ component.                                                                             |         |
| parameters                                                                                     | [AdaptiveLoadSchedulerParameters](#adaptive-load-scheduler-parameters) | `AdaptiveLoadSchedulerParameters` |          |         | Parameters for the _Adaptive Load Scheduler_ component.                                                                                             |         |

### <span id="adaptive-load-scheduler-ins"></span> AdaptiveLoadSchedulerIns

> Input ports for the _Adaptive Load Scheduler_ component.

**Properties**

| Name                                         | Type               | Go type  | Required | Default | Description                                                                                     | Example |
| -------------------------------------------- | ------------------ | -------- | :------: | ------- | ----------------------------------------------------------------------------------------------- | ------- |
| overload_confirmation                        | [InPort](#in-port) | `InPort` |          |         | The `overload_confirmation` port provides additional criteria to determine overload state which |
| results in _Flow_ throttling at the service. |                    |
| setpoint                                     | [InPort](#in-port) | `InPort` |          |         | The setpoint input to the controller.                                                           |         |
| signal                                       | [InPort](#in-port) | `InPort` |          |         | The input signal to the controller.                                                             |         |

### <span id="adaptive-load-scheduler-outs"></span> AdaptiveLoadSchedulerOuts

> Output ports for the _Adaptive Load Scheduler_ component.

**Properties**

| Name                     | Type                 | Go type   | Required | Default | Description                                                                              | Example |
| ------------------------ | -------------------- | --------- | :------: | ------- | ---------------------------------------------------------------------------------------- | ------- |
| desired_load_multiplier  | [OutPort](#out-port) | `OutPort` |          |         | Desired Load multiplier is the ratio of desired token rate to the incoming token rate.   |         |
| is_overload              | [OutPort](#out-port) | `OutPort` |          |         | A Boolean signal that indicates whether the service is in overload state.                |         |
| observed_load_multiplier | [OutPort](#out-port) | `OutPort` |          |         | Observed Load multiplier is the ratio of accepted token rate to the incoming token rate. |         |

### <span id="adaptive-load-scheduler-parameters"></span> AdaptiveLoadSchedulerParameters

> Parameters for the _Adaptive Load Scheduler_ component.

**Properties**

| Name                                                                        | Type                                                            | Go type                        | Required | Default  | Description                                                                             | Example |
| --------------------------------------------------------------------------- | --------------------------------------------------------------- | ------------------------------ | :------: | -------- | --------------------------------------------------------------------------------------- | ------- |
| alerter                                                                     | [AlerterParameters](#alerter-parameters)                        | `AlerterParameters`            |          |          | Configuration parameters for the embedded Alerter.                                      |         |
| gradient                                                                    | [GradientControllerParameters](#gradient-controller-parameters) | `GradientControllerParameters` |          |          | Parameters for the _Gradient Controller_.                                               |         |
| load_multiplier_linear_increment                                            | double (formatted number)                                       | `float64`                      |          | `0.0025` | Linear increment to load multiplier in each execution tick when the system is           |
| not in the overloaded state, up until the `max_load_multiplier` is reached. |                                                                 |
| load_scheduler                                                              | [LoadSchedulerParameters](#load-scheduler-parameters)           | `LoadSchedulerParameters`      |          |          | Parameters for the _Load Scheduler_.                                                    |         |
| max_load_multiplier                                                         | double (formatted number)                                       | `float64`                      |          | `2`      | The maximum load multiplier that can be reached during recovery from an overload state. |

- Helps protect the service from request bursts while the system is still
  recovering.
- Once this value is reached, the scheduler enters the pass-through mode,
  allowing requests to bypass the scheduler and be sent directly to the service.
- Any future overload state is detected by the control policy, and the load
  multiplier increment cycle is restarted. | |

### <span id="address-extractor"></span> AddressExtractor

> IP addresses in attribute context are defined as objects with separate IP and
> port fields. This is a helper to display an address as a single string.

:::caution

This might introduce high-cardinality flow label values.

:::

[ext-authz-address]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/address.proto#config-core-v3-address

Example:

```yaml
from: "source.address # or destination.address"
```

**Properties**

| Name | Type   | Go type  | Required | Default | Description                                                             | Example |
| ---- | ------ | -------- | :------: | ------- | ----------------------------------------------------------------------- | ------- |
| from | string | `string` |    ✓     |         | Attribute path pointing to some string - for example, `source.address`. |         |

### <span id="alerter"></span> Alerter

> Alerter reacts to a signal and generates alert to send to alert manager.

**Properties**

| Name       | Type                                     | Go type             | Required | Default | Description                            | Example |
| ---------- | ---------------------------------------- | ------------------- | :------: | ------- | -------------------------------------- | ------- |
| in_ports   | [AlerterIns](#alerter-ins)               | `AlerterIns`        |          |         | Input ports for the Alerter component. |         |
| parameters | [AlerterParameters](#alerter-parameters) | `AlerterParameters` |          |         |                                        |         |

### <span id="alerter-ins"></span> AlerterIns

> Inputs for the Alerter component.

**Properties**

| Name   | Type               | Go type  | Required | Default | Description                                                                                   | Example |
| ------ | ------------------ | -------- | :------: | ------- | --------------------------------------------------------------------------------------------- | ------- |
| signal | [InPort](#in-port) | `InPort` |          |         | Signal which Alerter is monitoring. If the signal greater than 0, Alerter generates an alert. |         |

### <span id="alerter-parameters"></span> AlerterParameters

> Alerter Parameters configure parameters such as alert name, severity, resolve
> timeout, alert channels and labels.

**Properties**

| Name            | Type          | Go type             | Required | Default  | Description                                             | Example |
| --------------- | ------------- | ------------------- | :------: | -------- | ------------------------------------------------------- | ------- |
| alert_channels  | []string      | `[]string`          |          |          | A list of alert channel strings.                        |         |
| alert_name      | string        | `string`            |    ✓     |          | Name of the alert.                                      |         |
| labels          | map of string | `map[string]string` |          |          | Additional labels to add to alert.                      |         |
| resolve_timeout | string        | `string`            |          | `"5s"`   | Duration of alert resolver.                             |         |
| severity        | string        | `string`            |          | `"info"` | Severity of the alert, one of 'info', 'warn' or 'crit'. |         |

### <span id="and"></span> And

> Logical AND.

Signals are mapped to Boolean values as follows:

- Zero is treated as false.
- Any non-zero is treated as true.
- Invalid inputs are considered unknown.

  :::note

  Treating invalid inputs as "unknowns" has a consequence that the result might
  end up being valid even when some inputs are invalid. For example,
  `unknown && false == false`, because the result would end up false no matter
  if first signal was true or false. Conversely, `unknown && true == unknown`.

  :::

**Properties**

| Name      | Type                 | Go type   | Required | Default | Description                         | Example |
| --------- | -------------------- | --------- | :------: | ------- | ----------------------------------- | ------- |
| in_ports  | [AndIns](#and-ins)   | `AndIns`  |          |         | Input ports for the And component.  |         |
| out_ports | [AndOuts](#and-outs) | `AndOuts` |          |         | Output ports for the And component. |         |

### <span id="and-ins"></span> AndIns

> Inputs for the And component.

**Properties**

| Name   | Type                 | Go type     | Required | Default | Description             | Example |
| ------ | -------------------- | ----------- | :------: | ------- | ----------------------- | ------- |
| inputs | [][InPort](#in-port) | `[]*InPort` |          |         | Array of input signals. |         |

### <span id="and-outs"></span> AndOuts

> Output ports for the And component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                                     | Example |
| ------ | -------------------- | --------- | :------: | ------- | ----------------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Result of logical AND of all the input signals. |

Will always be 0 (false), 1 (true) or invalid (unknown). | |

### <span id="arithmetic-combinator"></span> ArithmeticCombinator

**Properties**

| Name     | Type                                                  | Go type                   | Required | Default | Description                                          | Example |
| -------- | ----------------------------------------------------- | ------------------------- | :------: | ------- | ---------------------------------------------------- | ------- |
| in_ports | [ArithmeticCombinatorIns](#arithmetic-combinator-ins) | `ArithmeticCombinatorIns` |          |         | Input ports for the Arithmetic Combinator component. |         |
| operator | string                                                | `string`                  |          |         | Operator of the arithmetic operation.                |

The arithmetic operation can be addition, subtraction, multiplication, division,
XOR, right bit shift or left bit shift. In case of XOR and bit shifts, value of
signals is cast to integers before performing the operation. | | | out_ports |
[ArithmeticCombinatorOuts](#arithmetic-combinator-outs)|
`ArithmeticCombinatorOuts` | | | Output ports for the Arithmetic Combinator
component. | |

### <span id="arithmetic-combinator-ins"></span> ArithmeticCombinatorIns

> Inputs for the Arithmetic Combinator component.

**Properties**

| Name | Type               | Go type  | Required | Default | Description                                  | Example |
| ---- | ------------------ | -------- | :------: | ------- | -------------------------------------------- | ------- |
| lhs  | [InPort](#in-port) | `InPort` |          |         | Left hand side of the arithmetic operation.  |         |
| rhs  | [InPort](#in-port) | `InPort` |          |         | Right hand side of the arithmetic operation. |         |

### <span id="arithmetic-combinator-outs"></span> ArithmeticCombinatorOuts

> Outputs for the Arithmetic Combinator component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                     | Example |
| ------ | -------------------- | --------- | :------: | ------- | ------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Result of arithmetic operation. |         |

### <span id="auto-scale"></span> AutoScale

> AutoScale components are used to scale a service.

**Properties**

| Name        | Type                       | Go type      | Required | Default | Description                                                                                | Example |
| ----------- | -------------------------- | ------------ | :------: | ------- | ------------------------------------------------------------------------------------------ | ------- |
| auto_scaler | [AutoScaler](#auto-scaler) | `AutoScaler` |          |         | _AutoScaler_ provides auto-scaling functionality for any scalable resource.                |         |
| pod_scaler  | [PodScaler](#pod-scaler)   | `PodScaler`  |          |         | PodScaler provides pod horizontal scaling functionality for scalable Kubernetes resources. |         |

### <span id="auto-scale-kubernetes-control-point"></span> AutoScaleKubernetesControlPoint

**Properties**

| Name        | Type   | Go type  | Required | Default | Description | Example |
| ----------- | ------ | -------- | :------: | ------- | ----------- | ------- |
| api_version | string | `string` |          |         |             |         |
| kind        | string | `string` |          |         |             |         |
| name        | string | `string` |          |         |             |         |
| namespace   | string | `string` |          |         |             |         |

### <span id="auto-scaler"></span> AutoScaler

> _AutoScaler_ provides auto-scaling functionality for any scalable resource.
> Multiple _Controllers_ can be defined on the _AutoScaler_ for performing
> scale-out or scale-in. The _AutoScaler_ can interface with infrastructure APIs
> such as Kubernetes to perform auto-scale.

**Properties**

| Name    | Type    | Go type | Required | Default | Description                                                          | Example |
| ------- | ------- | ------- | :------: | ------- | -------------------------------------------------------------------- | ------- |
| dry_run | boolean | `bool`  |          |         | Dry run mode ensures that no scaling is invoked by this auto scaler. |

This is useful for observing the behavior of auto scaler without disrupting any
real deployment. This parameter sets the default value of dry run setting which
can be overridden at runtime using dynamic configuration. | | |
dry*run_config_key | string| `string` | | | Configuration key for overriding dry
run setting through dynamic configuration. | | | scale_in_controllers |
[][ScaleInController](#scale-in-controller)|
`[]*ScaleInController`| | | List of \_Controllers* for scaling in. | | | scale*out_controllers | [][ScaleOutController](#scale-out-controller)|`[]_ScaleOutController`
| | | List of \_Controllers_ for scaling out. | | | scaling_backend |
[AutoScalerScalingBackend](#auto-scaler-scaling-backend)|
`AutoScalerScalingBackend` | | | | | | scaling_parameters |
[AutoScalerScalingParameters](#auto-scaler-scaling-parameters)|
`AutoScalerScalingParameters` | | | Parameters that define the scaling behavior.
| |

### <span id="auto-scaler-scaling-backend"></span> AutoScalerScalingBackend

**Properties**

| Name                | Type                                                                                           | Go type                                      | Required | Default | Description | Example |
| ------------------- | ---------------------------------------------------------------------------------------------- | -------------------------------------------- | :------: | ------- | ----------- | ------- |
| kubernetes_replicas | [AutoScalerScalingBackendKubernetesReplicas](#auto-scaler-scaling-backend-kubernetes-replicas) | `AutoScalerScalingBackendKubernetesReplicas` |          |         |             |         |

### <span id="auto-scaler-scaling-backend-kubernetes-replicas"></span> AutoScalerScalingBackendKubernetesReplicas

> KubernetesReplicas defines a horizontal pod scaler for Kubernetes.

**Properties**

| Name                       | Type                                                                                                    | Go type                                          | Required | Default                 | Description                                                   | Example |
| -------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------ | :------: | ----------------------- | ------------------------------------------------------------- | ------- |
| kubernetes_object_selector | [KubernetesObjectSelector](#kubernetes-object-selector)                                                 | `KubernetesObjectSelector`                       |          |                         | The Kubernetes object on which horizontal scaling is applied. |         |
| max_replicas               | int64 (formatted string)                                                                                | `string`                                         |          | `"9223372036854775807"` | The maximum replicas to which the _AutoScaler_ can scale-out. |         |
| min_replicas               | int64 (formatted string)                                                                                | `string`                                         |          | `"0"`                   | The minimum replicas to which the _AutoScaler_ can scale-in.  |         |
| out_ports                  | [AutoScalerScalingBackendKubernetesReplicasOuts](#auto-scaler-scaling-backend-kubernetes-replicas-outs) | `AutoScalerScalingBackendKubernetesReplicasOuts` |          |                         | Output ports for _Kubernetes Replicas_.                       |         |

### <span id="auto-scaler-scaling-backend-kubernetes-replicas-outs"></span> AutoScalerScalingBackendKubernetesReplicasOuts

**Properties**

| Name                | Type                 | Go type   | Required | Default | Description | Example |
| ------------------- | -------------------- | --------- | :------: | ------- | ----------- | ------- |
| actual_replicas     | [OutPort](#out-port) | `OutPort` |          |         |             |         |
| configured_replicas | [OutPort](#out-port) | `OutPort` |          |         |             |         |
| desired_replicas    | [OutPort](#out-port) | `OutPort` |          |         |             |         |

### <span id="auto-scaler-scaling-parameters"></span> AutoScalerScalingParameters

**Properties**

| Name                         | Type                      | Go type   | Required | Default | Description                                                                                                        | Example |
| ---------------------------- | ------------------------- | --------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------ | ------- |
| cooldown_override_percentage | double (formatted number) | `float64` |          | `50`    | Cooldown override percentage defines a threshold change in scale-out beyond which previous cooldown is overridden. |

For example, if the cooldown is 5 minutes and the cooldown override percentage
is 10%, then if the scale-increases by 10% or more, the previous cooldown is
cancelled. Defaults to 50%. | | | max_scale_in_percentage | double (formatted
number)| `float64` | | `1`| The maximum decrease of scale (for example, pods) at
one time. Defined as percentage of current scale value. Can never go below one
even if percentage computation is less than one. Defaults to 1% of current scale
value. | | | max_scale_out_percentage | double (formatted number)| `float64` | |
`10`| The maximum increase of scale (for example, pods) at one time. Defined as
percentage of current scale value. Can never go below one even if percentage
computation is less than one. Defaults to 10% of current scale value. | | |
scale_in_alerter | [AlerterParameters](#alerter-parameters)| `AlerterParameters`
| | | Configuration for scale-in Alerter. | | | scale_in_cooldown | string|
`string` | | `"120s"`| The amount of time to wait after a scale-in operation for
another scale-in operation. | | | scale_out_alerter |
[AlerterParameters](#alerter-parameters)| `AlerterParameters` | | |
Configuration for scale-out Alerter. | | | scale_out_cooldown | string| `string`
| | `"30s"`| The amount of time to wait after a scale-out operation for another
scale-out or scale-in operation. | |

### <span id="bool-variable"></span> BoolVariable

> Component that emits a constant Boolean signal which can be changed at runtime
> through dynamic configuration.

**Properties**

| Name            | Type                                    | Go type            | Required | Default | Description                                                                                                                                                 | Example |
| --------------- | --------------------------------------- | ------------------ | :------: | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| config_key      | string                                  | `string`           |          |         | Configuration key for overriding value setting through dynamic configuration.                                                                               |         |
| constant_output | boolean                                 | `bool`             |          |         | The constant Boolean signal emitted by this component. The value of the constant Boolean signal can be overridden at runtime through dynamic configuration. |         |
| out_ports       | [BoolVariableOuts](#bool-variable-outs) | `BoolVariableOuts` |          |         | Output ports for the BoolVariable component.                                                                                                                |         |

### <span id="bool-variable-outs"></span> BoolVariableOuts

> Outputs for the BoolVariable component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                              | Example |
| ------ | -------------------- | --------- | :------: | ------- | ---------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | The value is emitted to the output port. |         |

### <span id="circuit"></span> Circuit

> Circuit is graph of inter-connected signal processing components.

:::info

See also [Circuit overview](/concepts/policy/circuit.md).

:::

Signals flow between components through ports. As signals traverse the circuit,
they get processed, stored within components or get acted upon (for example,
load-shed, rate-limit, auto-scale and so on). Circuit is evaluated periodically
to respond to changes in signal readings.

:::info Signals

Signals are floating point values.

A signal can also have a special **Invalid** value. It's usually used to
communicate that signal does not have a meaningful value at the moment, for
example, [PromQL](#prom-q-l) emits such a value if it cannot execute a query.
Components know when their input signals are invalid and can act accordingly.
They can either propagate the invalid signal, by making their output itself
invalid (for example, [ArithmeticCombinator](#arithmetic-combinator)) or use
some different logic, for example, [Extrapolator](#extrapolator). Refer to a
component's docs on how exactly it handles invalid inputs.

:::

**Properties**

| Name                                                                                                   | Type                      | Go type        | Required | Default  | Description                                                                            | Example |
| ------------------------------------------------------------------------------------------------------ | ------------------------- | -------------- | :------: | -------- | -------------------------------------------------------------------------------------- | ------- |
| components                                                                                             | [][Component](#component) | `[]*Component` |          |          | Defines a signal processing graph as a list of components.                             |         |
| evaluation_interval                                                                                    | string                    | `string`       |          | `"0.5s"` | Evaluation interval (tick) is the time between consecutive runs of the policy circuit. |
| This interval is typically aligned with how often the corrective action (actuation) needs to be taken. |                           |

### <span id="classifier"></span> Classifier

> :::info

See also [Classifier overview](/concepts/flow-control/resources/classifier.md).

::: Example

```yaml
selectors:
  - agent_group: demoapp
    service: service1-demo-app.demoapp.svc.cluster.local
    control_point: ingress
    label_matcher:
      match_labels:
        user_tier: gold
      match_expressions:
        - key: user_type
          operator: In
rules:
  user:
    extractor:
      from: request.http.headers.user-agent
  telemetry: false
```

**Properties**

| Name | Type          | Go type | Required | Default | Description                                                                                  | Example |
| ---- | ------------- | ------- | :------: | ------- | -------------------------------------------------------------------------------------------- | ------- |
| rego | [Rego](#rego) | `Rego`  |          |         | Rego is a policy language used to express complex policies in a concise and declarative way. |

It can be used to define flow classification rules by writing custom queries
that extract values from request metadata. For simple cases, such as directly
reading a value from header or a field from JSON body, declarative extractors
are recommended. | | | rules | map of [Rule](#rule)| `map[string]Rule` | | | A
map of {key, value} pairs mapping from
[flow label](/concepts/flow-control/flow-label.md) keys to rules that define how
to extract and propagate flow labels with that key. | | | selectors |
[][Selector](#selector)| `[]*Selector` | ✓ | | Selectors for flows that will be
classified by this _Classifier_. | |

### <span id="classifier-info-error"></span> ClassifierInfoError

> Error information.

| Name                | Type   | Go type | Default        | Description        | Example |
| ------------------- | ------ | ------- | -------------- | ------------------ | ------- |
| ClassifierInfoError | string | string  | `"ERROR_NONE"` | Error information. |         |

### <span id="common-attributes"></span> CommonAttributes

**Properties**

| Name         | Type   | Go type  | Required | Default | Description                             | Example |
| ------------ | ------ | -------- | :------: | ------- | --------------------------------------- | ------- |
| component_id | string | `string` |          |         | The id of Component within the circuit. |         |
| policy_hash  | string | `string` |          |         | Hash of the entire Policy spec.         |         |
| policy_name  | string | `string` |          |         | Name of the Policy.                     |         |

### <span id="component"></span> Component

> :::info

See also [Components overview](/concepts/policy/circuit.md#components).

:::

Signals flow into the components from input ports and results are emitted on
output ports. Components are wired to each other based on signal names forming
an execution graph of the circuit.

:::note

Loops are broken by the runtime at the earliest component index that is part of
the loop. The looped signals are saved in the tick they're generated and served
in the subsequent tick.

:::

See also [Policy](#policy) for a higher-level explanation of circuits.

**Properties**

| Name                                                                                                           | Type                                           | Go type                | Required | Default | Description                                                                                                                                                          | Example |
| -------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- | ---------------------- | :------: | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| alerter                                                                                                        | [Alerter](#alerter)                            | `Alerter`              |          |         | Alerter reacts to a signal and generates alert to send to alert manager.                                                                                             |         |
| and                                                                                                            | [And](#and)                                    | `And`                  |          |         | Logical AND.                                                                                                                                                         |         |
| arithmetic_combinator                                                                                          | [ArithmeticCombinator](#arithmetic-combinator) | `ArithmeticCombinator` |          |         | Applies the given operator on input operands (signals) and emits the result.                                                                                         |         |
| auto_scale                                                                                                     | [AutoScale](#auto-scale)                       | `AutoScale`            |          |         | AutoScale components are used to scale the service.                                                                                                                  |         |
| bool_variable                                                                                                  | [BoolVariable](#bool-variable)                 | `BoolVariable`         |          |         | BoolVariable emits a constant Boolean signal which can be changed at runtime via dynamic configuration.                                                              |         |
| decider                                                                                                        | [Decider](#decider)                            | `Decider`              |          |         | Decider emits the binary result of comparison operator on two operands.                                                                                              |         |
| differentiator                                                                                                 | [Differentiator](#differentiator)              | `Differentiator`       |          |         | Differentiator calculates rate of change per tick.                                                                                                                   |         |
| ema                                                                                                            | [EMA](#e-m-a)                                  | `EMA`                  |          |         | Exponential Moving Average filter.                                                                                                                                   |         |
| extrapolator                                                                                                   | [Extrapolator](#extrapolator)                  | `Extrapolator`         |          |         | Takes an input signal and emits the extrapolated value; either mirroring the input value or repeating the last known value up to the maximum extrapolation interval. |         |
| first_valid                                                                                                    | [FirstValid](#first-valid)                     | `FirstValid`           |          |         | Picks the first valid input signal and emits it.                                                                                                                     |         |
| flow_control                                                                                                   | [FlowControl](#flow-control)                   | `FlowControl`          |          |         | FlowControl components are used to regulate requests flow.                                                                                                           |         |
| gradient_controller                                                                                            | [GradientController](#gradient-controller)     | `GradientController`   |          |         | Gradient controller calculates the ratio between the signal and the setpoint to determine the magnitude of the correction that need to be applied.                   |
| This controller can be used to build AIMD (Additive Increase, Multiplicative Decrease) or MIMD style response. |                                                |
| holder                                                                                                         | [Holder](#holder)                              | `Holder`               |          |         | Holds the last valid signal value for the specified duration then waits for next valid value to hold.                                                                |         |
| integrator                                                                                                     | [Integrator](#integrator)                      | `Integrator`           |          |         | Accumulates sum of signal every tick.                                                                                                                                |         |
| inverter                                                                                                       | [Inverter](#inverter)                          | `Inverter`             |          |         | Logical NOT.                                                                                                                                                         |         |
| max                                                                                                            | [Max](#max)                                    | `Max`                  |          |         | Emits the maximum of the input signals.                                                                                                                              |         |
| min                                                                                                            | [Min](#min)                                    | `Min`                  |          |         | Emits the minimum of the input signals.                                                                                                                              |         |
| nested_circuit                                                                                                 | [NestedCircuit](#nested-circuit)               | `NestedCircuit`        |          |         | Nested circuit defines a sub-circuit as a high-level component. It consists of a list of components and a map of input and output ports.                             |         |
| nested_signal_egress                                                                                           | [NestedSignalEgress](#nested-signal-egress)    | `NestedSignalEgress`   |          |         | Nested signal egress is a special type of component that allows to extract a signal from a nested circuit.                                                           |         |
| nested_signal_ingress                                                                                          | [NestedSignalIngress](#nested-signal-ingress)  | `NestedSignalIngress`  |          |         | Nested signal ingress is a special type of component that allows to inject a signal into a nested circuit.                                                           |         |
| or                                                                                                             | [Or](#or)                                      | `Or`                   |          |         | Logical OR.                                                                                                                                                          |         |
| pulse_generator                                                                                                | [PulseGenerator](#pulse-generator)             | `PulseGenerator`       |          |         | Generates 0 and 1 in turns.                                                                                                                                          |         |
| query                                                                                                          | [Query](#query)                                | `Query`                |          |         | Query components that are query databases such as Prometheus.                                                                                                        |         |
| signal_generator                                                                                               | [SignalGenerator](#signal-generator)           | `SignalGenerator`      |          |         | Generates the specified signal.                                                                                                                                      |         |
| sma                                                                                                            | [SMA](#s-m-a)                                  | `SMA`                  |          |         | Simple Moving Average filter.                                                                                                                                        |         |
| switcher                                                                                                       | [Switcher](#switcher)                          | `Switcher`             |          |         | Switcher acts as a switch that emits one of the two signals based on third signal.                                                                                   |         |
| unary_operator                                                                                                 | [UnaryOperator](#unary-operator)               | `UnaryOperator`        |          |         | Takes an input signal and emits the square root of the input signal.                                                                                                 |         |
| variable                                                                                                       | [Variable](#variable)                          | `Variable`             |          |         | Emits a variable signal which can be changed at runtime via dynamic configuration.                                                                                   |         |

### <span id="constant-signal"></span> ConstantSignal

> Special constant input for ports and Variable component. Can provide either a
> constant value or special Nan/+-Inf value.

**Properties**

| Name          | Type                      | Go type   | Required | Default | Description                              | Example |
| ------------- | ------------------------- | --------- | :------: | ------- | ---------------------------------------- | ------- |
| special_value | string                    | `string`  |          |         | A special value such as NaN, +Inf, -Inf. |         |
| value         | double (formatted number) | `float64` |          |         | A constant value.                        |         |

### <span id="d-map"></span> DMap

**Properties**

| Name       | Type                     | Go type    | Required | Default | Description | Example |
| ---------- | ------------------------ | ---------- | :------: | ------- | ----------- | ------- |
| length     | int64 (formatted string) | `string`   |          |         |             |         |
| num_tables | int64 (formatted string) | `string`   |          |         |             |         |
| slab_info  | [SlabInfo](#slab-info)   | `SlabInfo` |          |         |             |         |

### <span id="decider"></span> Decider

> The comparison operator can be greater-than, less-than, greater-than-or-equal,
> less-than-or-equal, equal, or not-equal.

This component also supports time-based response (the output) transitions
between 1.0 or 0.0 signal if the decider condition is true or false for at least
`true_for` or `false_for` duration. If `true_for` and `false_for` durations are
zero then the transitions are instantaneous.

**Properties**

| Name                                                             | Type                         | Go type       | Required | Default | Description                                                               | Example |
| ---------------------------------------------------------------- | ---------------------------- | ------------- | :------: | ------- | ------------------------------------------------------------------------- | ------- |
| false_for                                                        | string                       | `string`      |          | `"0s"`  | Duration of time to wait before changing to false state.                  |
| If the duration is zero, the change will happen instantaneously. |                              |
| in_ports                                                         | [DeciderIns](#decider-ins)   | `DeciderIns`  |          |         | Input ports for the Decider component.                                    |         |
| operator                                                         | string                       | `string`      |          |         | Comparison operator that computes operation on LHS and RHS input signals. |         |
| out_ports                                                        | [DeciderOuts](#decider-outs) | `DeciderOuts` |          |         | Output ports for the Decider component.                                   |         |
| true_for                                                         | string                       | `string`      |          | `"0s"`  |                                                                           |         |

### <span id="decider-ins"></span> DeciderIns

> Inputs for the Decider component.

**Properties**

| Name | Type               | Go type  | Required | Default | Description                                                | Example |
| ---- | ------------------ | -------- | :------: | ------- | ---------------------------------------------------------- | ------- |
| lhs  | [InPort](#in-port) | `InPort` |          |         | Left hand side input signal for the comparison operation.  |         |
| rhs  | [InPort](#in-port) | `InPort` |          |         | Right hand side input signal for the comparison operation. |         |

### <span id="decider-outs"></span> DeciderOuts

> Outputs for the Decider component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                   | Example |
| ------ | -------------------- | --------- | :------: | ------- | ----------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Selected signal (1.0 or 0.0). |         |

### <span id="decreasing-gradient"></span> DecreasingGradient

> Decreasing Gradient defines a controller for scaling in based on Gradient
> Controller.

**Properties**

| Name       | Type                                                            | Go type                        | Required | Default | Description                   | Example |
| ---------- | --------------------------------------------------------------- | ------------------------------ | :------: | ------- | ----------------------------- | ------- |
| in_ports   | [DecreasingGradientIns](#decreasing-gradient-ins)               | `DecreasingGradientIns`        |          |         | Input ports for the Gradient. |         |
| parameters | [DecreasingGradientParameters](#decreasing-gradient-parameters) | `DecreasingGradientParameters` |          |         |                               |         |

### <span id="decreasing-gradient-ins"></span> DecreasingGradientIns

> Inputs for Gradient.

**Properties**

| Name     | Type               | Go type  | Required | Default | Description                       | Example |
| -------- | ------------------ | -------- | :------: | ------- | --------------------------------- | ------- |
| setpoint | [InPort](#in-port) | `InPort` |          |         | The setpoint to use for scale-in. |         |
| signal   | [InPort](#in-port) | `InPort` |          |         | The signal to use for scale-in.   |         |

### <span id="decreasing-gradient-parameters"></span> DecreasingGradientParameters

> This allows subset of parameters with constrained values compared to a regular
> gradient controller. For full documentation of these parameters, refer to the
> [GradientControllerParameters](#gradient-controller-parameters).

**Properties**

| Name         | Type                      | Go type   | Required | Default                    | Description | Example |
| ------------ | ------------------------- | --------- | :------: | -------------------------- | ----------- | ------- |
| min_gradient | double (formatted number) | `float64` |          | `-1.7976931348623157e+308` |             |         |
| slope        | double (formatted number) | `float64` |          | `1`                        |             |         |

### <span id="differentiator"></span> Differentiator

> Differentiator calculates rate of change per tick.

**Properties**

| Name      | Type                                       | Go type              | Required | Default | Description                                            | Example |
| --------- | ------------------------------------------ | -------------------- | :------: | ------- | ------------------------------------------------------ | ------- |
| in_ports  | [DifferentiatorIns](#differentiator-ins)   | `DifferentiatorIns`  |          |         | Input ports for the Differentiator component.          |         |
| out_ports | [DifferentiatorOuts](#differentiator-outs) | `DifferentiatorOuts` |          |         | Output ports for the Differentiator component.         |         |
| window    | string                                     | `string`             |          | `"5s"`  | The window of time over which differentiator operates. |         |

### <span id="differentiator-ins"></span> DifferentiatorIns

> Inputs for the Differentiator component.

**Properties**

| Name  | Type               | Go type  | Required | Default | Description | Example |
| ----- | ------------------ | -------- | :------: | ------- | ----------- | ------- |
| input | [InPort](#in-port) | `InPort` |          |         |             |         |

### <span id="differentiator-outs"></span> DifferentiatorOuts

> Outputs for the Differentiator component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description | Example |
| ------ | -------------------- | --------- | :------: | ------- | ----------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         |             |         |

### <span id="e-m-a"></span> EMA

> At any time EMA component operates in one of the following states:

1. Warm up state: The first `warmup_window` samples are used to compute the
   initial EMA. If an invalid reading is received during the `warmup_window`,
   the last good average is emitted and the state gets reset back to beginning
   of warm up state.
2. Normal state: The EMA is computed using following formula.

The EMA for a series $Y$ is calculated recursively as:

<!-- vale off -->

$$
\text{EMA} _t =
\begin{cases}
  Y_0, &\text{for } t = 0 \\
  \alpha Y_t + (1 - \alpha) \text{EMA}_{t-1}, &\text{for }t > 0
\end{cases}
$$

The coefficient $\alpha$ represents the degree of weighting decrease, a constant
smoothing factor between 0 and 1. A higher $\alpha$ discounts older observations
faster. The $\alpha$ is computed using ema_window:

$$
\alpha = \frac{2}{N + 1} \quad\text{where } N = \frac{\text{ema\_window}}{\text{evaluation\_period}}
$$

<!-- vale on -->

**Properties**

| Name       | Type                               | Go type         | Required | Default | Description                         | Example |
| ---------- | ---------------------------------- | --------------- | :------: | ------- | ----------------------------------- | ------- |
| in_ports   | [EMAIns](#e-m-a-ins)               | `EMAIns`        |          |         | Input ports for the EMA component.  |         |
| out_ports  | [EMAOuts](#e-m-a-outs)             | `EMAOuts`       |          |         | Output ports for the EMA component. |         |
| parameters | [EMAParameters](#e-m-a-parameters) | `EMAParameters` |          |         | Parameters for the EMA component.   |         |

### <span id="e-m-a-ins"></span> EMAIns

> Inputs for the EMA component.

**Properties**

| Name         | Type               | Go type  | Required | Default | Description                                      | Example |
| ------------ | ------------------ | -------- | :------: | ------- | ------------------------------------------------ | ------- |
| input        | [InPort](#in-port) | `InPort` |          |         | Input signal to be used for the EMA computation. |         |
| max_envelope | [InPort](#in-port) | `InPort` |          |         | Upper bound of the moving average.               |

When the signal exceeds `max_envelope` it is multiplied by
`correction_factor_on_max_envelope_violation` **once per tick**.

:::note

If the signal deviates from `max_envelope` faster than the correction faster, it
might end up exceeding the envelope.

::: | | | min_envelope | [InPort](#in-port)| `InPort` | | | Lower bound of the
moving average.

Behavior is similar to `max_envelope`. | |

### <span id="e-m-a-outs"></span> EMAOuts

> Outputs for the EMA component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                                                              | Example |
| ------ | -------------------- | --------- | :------: | ------- | ------------------------------------------------------------------------ | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Exponential moving average of the series of reading as an output signal. |         |

### <span id="e-m-a-parameters"></span> EMAParameters

> Parameters for the EMA component.

**Properties**

| Name                                        | Type                      | Go type   | Required | Default | Description                                                                             | Example |
| ------------------------------------------- | ------------------------- | --------- | :------: | ------- | --------------------------------------------------------------------------------------- | ------- |
| correction_factor_on_max_envelope_violation | double (formatted number) | `float64` |          | `1`     | Correction factor to apply on the output value if its in violation of the max envelope. |         |
| correction_factor_on_min_envelope_violation | double (formatted number) | `float64` |          | `1`     | Correction factor to apply on the output value if its in violation of the min envelope. |         |
| ema_window                                  | string                    | `string`  |    ✓     |         | Duration of EMA sampling window.                                                        |         |
| valid_during_warmup                         | boolean                   | `bool`    |          |         | Whether the output is valid during the warm-up stage.                                   |         |
| warmup_window                               | string                    | `string`  |    ✓     |         | Duration of EMA warming up window.                                                      |

The initial value of the EMA is the average of signal readings received during
the warm up window. | |

### <span id="entity"></span> Entity

> Entity represents a pod, VM, and so on.

**Properties**

| Name       | Type     | Go type    | Required | Default | Description                                          | Example |
| ---------- | -------- | ---------- | :------: | ------- | ---------------------------------------------------- | ------- |
| ip_address | string   | `string`   |    ✓     |         | IP address of the entity.                            |         |
| name       | string   | `string`   |          |         | Name of the entity. For example, pod name.           |         |
| namespace  | string   | `string`   |          |         | Namespace of the entity. For example, pod namespace. |         |
| node_name  | string   | `string`   |          |         | Node name of the entity. For example, hostname.      |         |
| services   | []string | `[]string` |          |         | Services of the entity.                              |         |
| uid        | string   | `string`   |    ✓     |         | Unique identifier of the entity.                     |         |

### <span id="equals-match-expression"></span> EqualsMatchExpression

> Label selector expression of the equal form `label == value`.

**Properties**

| Name  | Type   | Go type  | Required | Default | Description                                    | Example |
| ----- | ------ | -------- | :------: | ------- | ---------------------------------------------- | ------- |
| label | string | `string` |    ✓     |         | Name of the label to equal match the value.    |         |
| value | string | `string` |          |         | Exact value that the label should be equal to. |         |

### <span id="extractor"></span> Extractor

> There are multiple variants of extractor, specify exactly one.

**Properties**

| Name    | Type                                   | Go type            | Required | Default | Description                                            | Example |
| ------- | -------------------------------------- | ------------------ | :------: | ------- | ------------------------------------------------------ | ------- |
| address | [AddressExtractor](#address-extractor) | `AddressExtractor` |          |         | Display an address as a single string - `<ip>:<port>`. |         |
| from    | string                                 | `string`           |          |         | Attribute path is a dot-separated path to attribute.   |

Should be either:

- one of the fields of
  [Attribute Context](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto),
  or
- a special `request.http.bearer` pseudo-attribute. For example,
  `request.http.method` or `request.http.header.user-agent`

Note: The same attribute path syntax is shared by other extractor variants,
wherever attribute path is needed in their "from" syntax.

Example:

````yaml
from: request.http.headers.user-agent
``` |  |
| json | [JSONExtractor](#json-extractor)| `JSONExtractor` |  | | Parse JSON, and extract one of the fields. |  |
| jwt | [JWTExtractor](#j-w-t-extractor)| `JWTExtractor` |  | | Parse the attribute as JWT and read the payload. |  |
| path_templates | [PathTemplateMatcher](#path-template-matcher)| `PathTemplateMatcher` |  | | Match HTTP Path to given path templates. |  |



### <span id="extrapolator"></span> Extrapolator


> It does so until `maximum_extrapolation_interval` is reached, beyond which it emits invalid signal unless input signal becomes valid again.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [ExtrapolatorIns](#extrapolator-ins)| `ExtrapolatorIns` |  | | Input ports for the Extrapolator component. |  |
| out_ports | [ExtrapolatorOuts](#extrapolator-outs)| `ExtrapolatorOuts` |  | | Output ports for the Extrapolator component. |  |
| parameters | [ExtrapolatorParameters](#extrapolator-parameters)| `ExtrapolatorParameters` |  | | Parameters for the Extrapolator component. |  |



### <span id="extrapolator-ins"></span> ExtrapolatorIns


> Inputs for the Extrapolator component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| input | [InPort](#in-port)| `InPort` |  | | Input signal for the Extrapolator component. |  |



### <span id="extrapolator-outs"></span> ExtrapolatorOuts


> Outputs for the Extrapolator component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | | Extrapolated signal. |  |



### <span id="extrapolator-parameters"></span> ExtrapolatorParameters


> Parameters for the Extrapolator component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| max_extrapolation_interval | string| `string` | ✓ | | Maximum time interval to repeat the last valid value of input signal. |  |



### <span id="first-valid"></span> FirstValid






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [FirstValidIns](#first-valid-ins)| `FirstValidIns` |  | | Input ports for the FirstValid component. |  |
| out_ports | [FirstValidOuts](#first-valid-outs)| `FirstValidOuts` |  | | Output ports for the FirstValid component. |  |



### <span id="first-valid-ins"></span> FirstValidIns


> Inputs for the FirstValid component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| inputs | [][InPort](#in-port)| `[]*InPort` |  | | Array of input signals. |  |



### <span id="first-valid-outs"></span> FirstValidOuts


> Outputs for the FirstValid component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | | First valid input signal as an output signal. |  |



### <span id="flow-control"></span> FlowControl


> _Flow Control_ encompasses components that manage the flow of requests or access to features within a service.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| adaptive_load_scheduler | [AdaptiveLoadScheduler](#adaptive-load-scheduler)| `AdaptiveLoadScheduler` |  | | _Adaptive Load Scheduler_ component is based on additive increase and multiplicative decrease of token rate. It takes a signal and setpoint as inputs and reduces token rate proportionally (or any arbitrary power) based on deviation of the signal from setpoint. |  |
| load_ramp | [LoadRamp](#load-ramp)| `LoadRamp` |  | | _Load Ramp_ smoothly regulates the flow of requests over specified steps. |  |
| load_ramp_series | [LoadRampSeries](#load-ramp-series)| `LoadRampSeries` |  | | _Load Ramp Series_ is a series of _Load Ramp_ components that can shape load one after another at same or different _Control Points_. |  |
| load_scheduler | [LoadScheduler](#load-scheduler)| `LoadScheduler` |  | | _Load Scheduler_ provides service protection by creating a prioritized workload queue in front of the service using Weighted Fair Queuing. |  |
| quota_scheduler | [QuotaScheduler](#quota-scheduler)| `QuotaScheduler` |  | |  |  |
| rate_limiter | [RateLimiter](#rate-limiter)| `RateLimiter` |  | | _Rate Limiter_ provides service protection by applying rate limits using the token bucket algorithm. |  |
| regulator | [Regulator](#regulator)| `Regulator` |  | | Regulator is a component that regulates the flow of requests to the service by allowing only the specified percentage of requests or sticky sessions. |  |



### <span id="flow-control-point"></span> FlowControlPoint






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| control_point | string| `string` |  | |  |  |
| service | string| `string` |  | |  |  |
| type | string| `string` |  | |  |  |



### <span id="flow-control-resources"></span> FlowControlResources






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| classifiers | [][Classifier](#classifier)| `[]*Classifier` |  | | Classifiers are installed in the data-plane and are used to label the requests based on payload content.

The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes. |  |
| flux_meters | map of [FluxMeter](#flux-meter)| `map[string]FluxMeter` |  | | Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.

Flux Meter created metrics can be consumed as input to the circuit through the PromQL component. |  |



### <span id="flux-meter"></span> FluxMeter


> Flux Meter gathers metrics for the traffic that matches its selector.
The histogram created by Flux Meter measures the workload latency by default.

:::info

See also [Flux Meter overview](/concepts/flow-control/resources/flux-meter.md).

:::
Example:
```yaml
static_buckets:
   buckets: [5.0,10.0,25.0,50.0,100.0,250.0,500.0,1000.0,2500.0,5000.0,10000.0]
selectors:
   - agent_group: demoapp
     service: service1-demo-app.demoapp.svc.cluster.local
     control_point: ingress
attribute_key: response_duration_ms
````

**Properties**

| Name          | Type   | Go type  | Required | Default                  | Description                                                                                   | Example |
| ------------- | ------ | -------- | :------: | ------------------------ | --------------------------------------------------------------------------------------------- | ------- |
| attribute_key | string | `string` |          | `"workload_duration_ms"` | Key of the attribute in access log or span from which the metric for this flux meter is read. |

:::info

For list of available attributes in Envoy access logs, refer
[Envoy Filter](/integrations/istio/istio.md#envoy-filter)

::: | | | exponential*buckets |
[FluxMeterExponentialBuckets](#flux-meter-exponential-buckets)|
`FluxMeterExponentialBuckets` | | | | | | exponential_buckets_range |
[FluxMeterExponentialBucketsRange](#flux-meter-exponential-buckets-range)|
`FluxMeterExponentialBucketsRange` | | | | | | linear_buckets |
[FluxMeterLinearBuckets](#flux-meter-linear-buckets)| `FluxMeterLinearBuckets` |
| | | | | selectors | [][Selector](#selector)|
`[]*Selector`| ✓ | | Selectors for flows that will be metered by this \_Flux Meter*. | | | static_buckets | [FluxMeterStaticBuckets](#flux-meter-static-buckets)|`FluxMeterStaticBuckets`
| | | | |

### <span id="flux-meter-exponential-buckets"></span> FluxMeterExponentialBuckets

> ExponentialBuckets creates `count` number of buckets where the lowest bucket
> has an upper bound of `start` and each following bucket's upper bound is
> `factor` times the previous bucket's upper bound. The final +inf bucket is not
> counted.

**Properties**

| Name   | Type                      | Go type   | Required | Default | Description                                                                                                   | Example |
| ------ | ------------------------- | --------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------- | ------- |
| count  | int32 (formatted integer) | `int32`   |          |         | Number of buckets.                                                                                            |         |
| factor | double (formatted number) | `float64` |          |         | Factor to be multiplied to the previous bucket's upper bound to calculate the following bucket's upper bound. |         |
| start  | double (formatted number) | `float64` |          |         | Upper bound of the lowest bucket.                                                                             |         |

### <span id="flux-meter-exponential-buckets-range"></span> FluxMeterExponentialBucketsRange

> ExponentialBucketsRange creates `count` number of buckets where the lowest
> bucket is `min` and the highest bucket is `max`. The final +inf bucket is not
> counted.

**Properties**

| Name  | Type                      | Go type   | Required | Default | Description        | Example |
| ----- | ------------------------- | --------- | :------: | ------- | ------------------ | ------- |
| count | int32 (formatted integer) | `int32`   |          |         | Number of buckets. |         |
| max   | double (formatted number) | `float64` |          |         | Highest bucket.    |         |
| min   | double (formatted number) | `float64` |          |         | Lowest bucket.     |         |

### <span id="flux-meter-linear-buckets"></span> FluxMeterLinearBuckets

> LinearBuckets creates `count` number of buckets, each `width` wide, where the
> lowest bucket has an upper bound of `start`. The final +inf bucket is not
> counted.

**Properties**

| Name  | Type                      | Go type   | Required | Default | Description                       | Example |
| ----- | ------------------------- | --------- | :------: | ------- | --------------------------------- | ------- |
| count | int32 (formatted integer) | `int32`   |          |         | Number of buckets.                |         |
| start | double (formatted number) | `float64` |          |         | Upper bound of the lowest bucket. |         |
| width | double (formatted number) | `float64` |          |         | Width of each bucket.             |         |

### <span id="flux-meter-static-buckets"></span> FluxMeterStaticBuckets

> StaticBuckets holds the static value of the buckets where latency histogram
> will be stored.

**Properties**

| Name    | Type                        | Go type     | Required | Default                                         | Description                                            | Example |
| ------- | --------------------------- | ----------- | :------: | ----------------------------------------------- | ------------------------------------------------------ | ------- |
| buckets | []double (formatted number) | `[]float64` |          | `[5,10,25,50,100,250,500,1000,2500,5000,10000]` | The buckets in which latency histogram will be stored. |         |

### <span id="gradient-controller"></span> GradientController

> The `gradient` describes a corrective factor that should be applied to the
> control variable to get the signal closer to the setpoint. It's computed as
> follows:

$$
\text{gradient} = \left(\frac{\text{signal}}{\text{setpoint}}\right)^{\text{slope}}
$$

`gradient` is then clamped to `[min_gradient, max_gradient]` range.

The output of gradient controller is computed as follows:

$$
\text{output} = \text{gradient}_{\text{clamped}} \cdot \text{control\_variable} + \text{optimize}.
$$

Note the additional `optimize` signal, that can be used to "nudge" the
controller into desired idle state.

The output can be _optionally_ clamped to desired range using `max` and `min`
input.

**Properties**

| Name                   | Type                                                            | Go type                        | Required | Default | Description                                                                                                                                                                                                                       | Example |
| ---------------------- | --------------------------------------------------------------- | ------------------------------ | :------: | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| in_ports               | [GradientControllerIns](#gradient-controller-ins)               | `GradientControllerIns`        |          |         | Input ports of the Gradient Controller.                                                                                                                                                                                           |         |
| manual_mode            | boolean                                                         | `bool`                         |          |         | In manual mode, the controller does not adjust the control variable. It emits the same output as the control variable input. This setting can be adjusted at runtime through dynamic configuration without restarting the policy. |         |
| manual_mode_config_key | string                                                          | `string`                       |          |         | Configuration key for overriding `manual_mode` setting through dynamic configuration.                                                                                                                                             |         |
| out_ports              | [GradientControllerOuts](#gradient-controller-outs)             | `GradientControllerOuts`       |          |         | Output ports of the Gradient Controller.                                                                                                                                                                                          |         |
| parameters             | [GradientControllerParameters](#gradient-controller-parameters) | `GradientControllerParameters` |          |         | Gradient Parameters.                                                                                                                                                                                                              |         |

### <span id="gradient-controller-ins"></span> GradientControllerIns

> Inputs for the Gradient Controller component.

**Properties**

| Name             | Type               | Go type  | Required | Default | Description                                   | Example |
| ---------------- | ------------------ | -------- | :------: | ------- | --------------------------------------------- | ------- |
| control_variable | [InPort](#in-port) | `InPort` |          |         | Actual current value of the control variable. |

This signal is multiplied by the gradient to produce the output. | | | max |
[InPort](#in-port)| `InPort` | | | Maximum value to limit the output signal. | |
| min | [InPort](#in-port)| `InPort` | | | Minimum value to limit the output
signal. | | | setpoint | [InPort](#in-port)| `InPort` | | | Setpoint to be used
for the gradient computation. | | | signal | [InPort](#in-port)| `InPort` | | |
Signal to be used for the gradient computation. | |

### <span id="gradient-controller-outs"></span> GradientControllerOuts

> Outputs for the Gradient Controller component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                                     | Example |
| ------ | -------------------- | --------- | :------: | ------- | ----------------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Computed desired value of the control variable. |         |

### <span id="gradient-controller-parameters"></span> GradientControllerParameters

> Gradient Parameters.

**Properties**

| Name         | Type                      | Go type   | Required | Default                    | Description                                                                                             | Example |
| ------------ | ------------------------- | --------- | :------: | -------------------------- | ------------------------------------------------------------------------------------------------------- | ------- |
| max_gradient | double (formatted number) | `float64` |          | `1.7976931348623157e+308`  | Maximum gradient which clamps the computed gradient value to the range, `[min_gradient, max_gradient]`. |         |
| min_gradient | double (formatted number) | `float64` |          | `-1.7976931348623157e+308` | Minimum gradient which clamps the computed gradient value to the range, `[min_gradient, max_gradient]`. |         |
| slope        | double (formatted number) | `float64` |    ✓     |                            | Slope controls the aggressiveness and direction of the Gradient Controller.                             |

Slope is used as exponent on the signal to setpoint ratio in computation of the
gradient (see the [main description](#gradient-controller) for exact equation).
This parameter decides how aggressive the controller responds to the deviation
of signal from the setpoint. for example:

- $\text{slope} = 1$: when signal is too high, increase control variable,
- $\text{slope} = -1$: when signal is too high, decrease control variable,
- $\text{slope} = -0.5$: when signal is too high, decrease control variable
  gradually.

The sign of slope depends on correlation between the signal and control
variable:

- Use $\text{slope} < 0$ if there is a _positive_ correlation between the signal
  and the control variable (for example, Per-pod CPU usage and total
  concurrency).
- Use $\text{slope} > 0$ if there is a _negative_ correlation between the signal
  and the control variable (for example, Per-pod CPU usage and number of pods).

:::note

You need to set _negative_ slope for a _positive_ correlation, as you're
describing the _action_ which controller should make when the signal increases.

:::

The magnitude of slope describes how aggressively should the controller react to
a deviation of signal. With $|\text{slope}| = 1$, the controller will aim to
bring the signal to the setpoint in one tick (assuming linear correlation with
signal and setpoint). Smaller magnitudes of slope will make the controller
adjust the control variable gradually.

Setting $|\text{slope}| < 1$ (for example, $\pm0.8$) is recommended. If you
experience overshooting, consider lowering the magnitude even more. Values of
$|\text{slope}| > 1$ aren't recommended.

:::note

Remember that the gradient and output signal can be (optionally) clamped, so the
_slope_ might not fully describe aggressiveness of the controller.

::: | |

### <span id="group-status"></span> GroupStatus

> Groups is nested structure that holds status information about the node and a
> pointer to the next node.

**Properties**

| Name   | Type                                | Go type                  | Required | Default | Description | Example |
| ------ | ----------------------------------- | ------------------------ | :------: | ------- | ----------- | ------- |
| groups | map of [GroupStatus](#group-status) | `map[string]GroupStatus` |          |         |             |         |
| status | [Status](#status)                   | `Status`                 |          |         |             |         |

### <span id="holder"></span> Holder

> Holds the last valid signal value for the specified duration then waits for
> next valid value to hold. If it is holding a value that means it ignores both
> valid and invalid new signals until the `hold_for` duration is finished.

**Properties**

| Name      | Type                       | Go type      | Required | Default | Description                                                      | Example |
| --------- | -------------------------- | ------------ | :------: | ------- | ---------------------------------------------------------------- | ------- |
| hold_for  | string                     | `string`     |          | `"5s"`  | Holding the last valid signal value for the `hold_for` duration. |         |
| in_ports  | [HolderIns](#holder-ins)   | `HolderIns`  |          |         | Input ports for the Holder component.                            |         |
| out_ports | [HolderOuts](#holder-outs) | `HolderOuts` |          |         | Output ports for the Holder component.                           |         |

### <span id="holder-ins"></span> HolderIns

> Inputs for the Holder component.

**Properties**

| Name  | Type               | Go type  | Required | Default | Description                                                                                   | Example |
| ----- | ------------------ | -------- | :------: | ------- | --------------------------------------------------------------------------------------------- | ------- |
| input | [InPort](#in-port) | `InPort` |          |         | The input signal.                                                                             |         |
| reset | [InPort](#in-port) | `InPort` |          |         | Resets the holder output to the current input signal when reset signal is valid and non-zero. |         |

### <span id="holder-outs"></span> HolderOuts

> Outputs for the Holder component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description        | Example |
| ------ | -------------------- | --------- | :------: | ------- | ------------------ | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | The output signal. |         |

### <span id="in-port"></span> InPort

**Properties**

| Name            | Type                               | Go type          | Required | Default | Description                                                    | Example |
| --------------- | ---------------------------------- | ---------------- | :------: | ------- | -------------------------------------------------------------- | ------- |
| constant_signal | [ConstantSignal](#constant-signal) | `ConstantSignal` |          |         | Constant value to be used for this InPort instead of a signal. |         |
| signal_name     | string                             | `string`         |          |         | Name of the incoming Signal on the InPort.                     |         |

### <span id="increasing-gradient"></span> IncreasingGradient

> Increasing Gradient defines a controller for scaling out based on _Gradient
> Controller_.

**Properties**

| Name       | Type                                                            | Go type                        | Required | Default | Description                   | Example |
| ---------- | --------------------------------------------------------------- | ------------------------------ | :------: | ------- | ----------------------------- | ------- |
| in_ports   | [IncreasingGradientIns](#increasing-gradient-ins)               | `IncreasingGradientIns`        |          |         | Input ports for the Gradient. |         |
| parameters | [IncreasingGradientParameters](#increasing-gradient-parameters) | `IncreasingGradientParameters` |          |         |                               |         |

### <span id="increasing-gradient-ins"></span> IncreasingGradientIns

> Inputs for Gradient.

**Properties**

| Name     | Type               | Go type  | Required | Default | Description                        | Example |
| -------- | ------------------ | -------- | :------: | ------- | ---------------------------------- | ------- |
| setpoint | [InPort](#in-port) | `InPort` |          |         | The setpoint to use for scale-out. |         |
| signal   | [InPort](#in-port) | `InPort` |          |         | The signal to use for scale-out.   |         |

### <span id="increasing-gradient-parameters"></span> IncreasingGradientParameters

> This allows subset of parameters with constrained values compared to a regular
> gradient controller. For full documentation of these parameters, refer to the
> [GradientControllerParameters](#gradient-controller-parameters).

**Properties**

| Name         | Type                      | Go type   | Required | Default                   | Description | Example |
| ------------ | ------------------------- | --------- | :------: | ------------------------- | ----------- | ------- |
| max_gradient | double (formatted number) | `float64` |          | `1.7976931348623157e+308` |             |         |
| slope        | double (formatted number) | `float64` |          | `1`                       |             |         |

### <span id="infra-meter"></span> InfraMeter

> InfraMeter is a resource that sets up OpenTelemetry pipelines. It defines
> receivers, processors, and a single metrics pipeline which will be exported to
> the configured Prometheus instance. Environment variables can be used in the
> configuration using format `${ENV_VAR_NAME}`.

:::info

See also
[Get Started / Setup Integrations / Metrics](/integrations/metrics/metrics.md).

:::

**Properties**

| Name            | Type    | Go type | Required | Default | Description                                                             | Example |
| --------------- | ------- | ------- | :------: | ------- | ----------------------------------------------------------------------- | ------- |
| per_agent_group | boolean | `bool`  |          |         | PerAgentGroup marks the pipeline to be instantiated only once per agent |

group. This is helpful for receivers that scrape for example, some cluster-wide
metrics. When not set, pipeline will be instantiated on every Agent. | | |
pipeline | [InfraMeterMetricsPipeline](#infra-meter-metrics-pipeline)|
`InfraMeterMetricsPipeline` | | | Pipeline is an OTel metrics pipeline
definition, which **only** uses receivers and processors defined above. Exporter
would be added automatically.

If there are no processors defined or only one processor is defined, the
pipeline definition can be omitted. In such cases, the pipeline will
automatically use all given receivers and the defined processor (if any).
However, if there are more than one processor, the pipeline must be defined
explicitly. | | | processors | map of any | `map[string]interface{}` | | |
Processors define processors to be used in custom metrics pipelines. This should
be in
[OTel format](https://opentelemetry.io/docs/collector/configuration/#processors).
| | | receivers | map of any | `map[string]interface{}` | | | Receivers define
receivers to be used in custom metrics pipelines. This should be in
[OTel format](https://opentelemetry.io/docs/collector/configuration/#receivers).
| |

### <span id="infra-meter-metrics-pipeline"></span> InfraMeterMetricsPipeline

> MetricsPipelineConfig defines a custom metrics pipeline.

**Properties**

| Name       | Type     | Go type    | Required | Default | Description | Example |
| ---------- | -------- | ---------- | :------: | ------- | ----------- | ------- |
| processors | []string | `[]string` |          |         |             |         |
| receivers  | []string | `[]string` |          |         |             |         |

### <span id="integrator"></span> Integrator

> Accumulates sum of signal every tick.

**Properties**

| Name          | Type                               | Go type          | Required | Default | Description                                | Example |
| ------------- | ---------------------------------- | ---------------- | :------: | ------- | ------------------------------------------ | ------- |
| in_ports      | [IntegratorIns](#integrator-ins)   | `IntegratorIns`  |          |         | Input ports for the Integrator component.  |         |
| initial_value | double (formatted number)          | `float64`        |          |         | Initial value of the integrator.           |         |
| out_ports     | [IntegratorOuts](#integrator-outs) | `IntegratorOuts` |          |         | Output ports for the Integrator component. |         |

### <span id="integrator-ins"></span> IntegratorIns

> Inputs for the Integrator component.

**Properties**

| Name  | Type               | Go type  | Required | Default | Description                                                                                                                  | Example |
| ----- | ------------------ | -------- | :------: | ------- | ---------------------------------------------------------------------------------------------------------------------------- | ------- |
| input | [InPort](#in-port) | `InPort` |          |         | The input signal.                                                                                                            |         |
| max   | [InPort](#in-port) | `InPort` |          |         | The maximum output.                                                                                                          |         |
| min   | [InPort](#in-port) | `InPort` |          |         | The minimum output.                                                                                                          |         |
| reset | [InPort](#in-port) | `InPort` |          |         | Resets the integrator output to zero when reset signal is valid and non-zero. Reset also resets the max and min constraints. |         |

### <span id="integrator-outs"></span> IntegratorOuts

> Outputs for the Integrator component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description | Example |
| ------ | -------------------- | --------- | :------: | ------- | ----------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         |             |         |

### <span id="inverter"></span> Inverter

> Logical NOT.

See [And component](#and) on how signals are mapped onto Boolean values.

**Properties**

| Name      | Type                           | Go type        | Required | Default | Description                              | Example |
| --------- | ------------------------------ | -------------- | :------: | ------- | ---------------------------------------- | ------- |
| in_ports  | [InverterIns](#inverter-ins)   | `InverterIns`  |          |         | Input ports for the Inverter component.  |         |
| out_ports | [InverterOuts](#inverter-outs) | `InverterOuts` |          |         | Output ports for the Inverter component. |         |

### <span id="inverter-ins"></span> InverterIns

> Inputs for the Inverter component.

**Properties**

| Name  | Type               | Go type  | Required | Default | Description           | Example |
| ----- | ------------------ | -------- | :------: | ------- | --------------------- | ------- |
| input | [InPort](#in-port) | `InPort` |          |         | Signal to be negated. |         |

### <span id="inverter-outs"></span> InverterOuts

> Output ports for the Inverter component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                           | Example |
| ------ | -------------------- | --------- | :------: | ------- | ------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Logical negation of the input signal. |

Will always be 0 (false), 1 (true) or invalid (unknown). | |

### <span id="json-extractor"></span> JSONExtractor

> Example:

```yaml
from: request.http.body
pointer: /user/name
```

**Properties**

| Name    | Type   | Go type  | Required | Default | Description                                                                                              | Example |
| ------- | ------ | -------- | :------: | ------- | -------------------------------------------------------------------------------------------------------- | ------- |
| from    | string | `string` |    ✓     |         | Attribute path pointing to some strings - for example, `request.http.body`.                              |         |
| pointer | string | `string` |          |         | JSON pointer represents a parsed JSON pointer which allows to select a specified field from the payload. |

Note: Uses [JSON pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
for example, `/foo/bar`. If the pointer points into an object, it'd be converted
to a string. | |

### <span id="j-w-t-extractor"></span> JWTExtractor

> Specify a field to be extracted from payload using `json_pointer`.

Note: The signature is not verified against the secret (assuming there's some
other part of the system that handles such verification).

Example:

```yaml
from: request.http.bearer
json_pointer: /user/email
```

**Properties**

| Name         | Type   | Go type  | Required | Default | Description                                                                                                              | Example |
| ------------ | ------ | -------- | :------: | ------- | ------------------------------------------------------------------------------------------------------------------------ | ------- |
| from         | string | `string` |    ✓     |         | JWT (JSON Web Token) can be extracted from any input attribute, but most likely you'd want to use `request.http.bearer`. |         |
| json_pointer | string | `string` |          |         | JSON pointer allowing to select a specified field from the payload.                                                      |

Note: Uses [JSON pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
for example, `/foo/bar`. If the pointer points into an object, it'd be converted
to a string. | |

### <span id="k8s-label-matcher-requirement"></span> K8sLabelMatcherRequirement

> Label selector requirement which is a selector that contains values, a key,
> and an operator that relates the key and values.

**Properties**

| Name                                                    | Type     | Go type    | Required | Default | Description                                                                | Example |
| ------------------------------------------------------- | -------- | ---------- | :------: | ------- | -------------------------------------------------------------------------- | ------- |
| key                                                     | string   | `string`   |    ✓     |         | Label key that the selector applies to.                                    |         |
| operator                                                | string   | `string`   |          |         | Logical operator which represents a key's relationship to a set of values. |
| Valid operators are In, NotIn, Exists and DoesNotExist. |          |
| values                                                  | []string | `[]string` |          |         | An array of string values that relates to the key by an operator.          |

If the operator is In or NotIn, the values array must be non-empty. If the
operator is Exists or DoesNotExist, the values array must be empty. | |

### <span id="kubernetes-object-selector"></span> KubernetesObjectSelector

> Describes which pods a control or observability component should apply to.

**Properties**

| Name                 | Type   | Go type  | Required | Default     | Description                                                              | Example |
| -------------------- | ------ | -------- | :------: | ----------- | ------------------------------------------------------------------------ | ------- |
| agent_group          | string | `string` |          | `"default"` | Which [agent-group](/concepts/flow-control/selector.md#agent-group) this |
| selector applies to. |        |
| api_version          | string | `string` |    ✓     |             |                                                                          |         |
| kind                 | string | `string` |    ✓     |             | Kubernetes resource type.                                                |         |
| name                 | string | `string` |    ✓     |             | Kubernetes resource name.                                                |         |
| namespace            | string | `string` |    ✓     |             | Kubernetes namespace that the resource belongs to.                       |         |

### <span id="label-matcher"></span> LabelMatcher

> It provides three ways to define requirements:

- match labels
- match expressions
- arbitrary expression

If multiple requirements are set, they're all combined using the logical AND
operator. An empty label matcher always matches.

**Properties**

| Name              | Type                                                           | Go type                         | Required | Default | Description                                            | Example |
| ----------------- | -------------------------------------------------------------- | ------------------------------- | :------: | ------- | ------------------------------------------------------ | ------- |
| expression        | [MatchExpression](#match-expression)                           | `MatchExpression`               |          |         | An arbitrary expression to be evaluated on the labels. |         |
| match_expressions | [][K8sLabelMatcherRequirement](#k8s-label-matcher-requirement) | `[]*K8sLabelMatcherRequirement` |          |         | List of Kubernetes-style label matcher requirements.   |

Note: The requirements are combined using the logical AND operator. | | |
match_labels | map of string| `map[string]string` | | | A map of {key,value}
pairs representing labels to be matched. A single {key,value} in the
`match_labels` requires that the label `key` is present and equal to `value`.

Note: The requirements are combined using the logical AND operator. | |

### <span id="limiter-decision-limiter-reason"></span> LimiterDecisionLimiterReason

| Name                         | Type   | Go type | Default                        | Description | Example |
| ---------------------------- | ------ | ------- | ------------------------------ | ----------- | ------- |
| LimiterDecisionLimiterReason | string | string  | `"LIMITER_REASON_UNSPECIFIED"` |             |         |

### <span id="limiter-decision-quota-scheduler-info"></span> LimiterDecisionQuotaSchedulerInfo

**Properties**

| Name           | Type                                                             | Go type                        | Required | Default | Description | Example |
| -------------- | ---------------------------------------------------------------- | ------------------------------ | :------: | ------- | ----------- | ------- |
| label          | string                                                           | `string`                       |          |         |             |         |
| scheduler_info | [LimiterDecisionSchedulerInfo](#limiter-decision-scheduler-info) | `LimiterDecisionSchedulerInfo` |          |         |             |         |

### <span id="limiter-decision-rate-limiter-info"></span> LimiterDecisionRateLimiterInfo

**Properties**

| Name            | Type                      | Go type   | Required | Default | Description | Example |
| --------------- | ------------------------- | --------- | :------: | ------- | ----------- | ------- |
| current         | double (formatted number) | `float64` |          |         |             |         |
| label           | string                    | `string`  |          |         |             |         |
| remaining       | double (formatted number) | `float64` |          |         |             |         |
| tokens_consumed | double (formatted number) | `float64` |          |         |             |         |

### <span id="limiter-decision-regulator-info"></span> LimiterDecisionRegulatorInfo

**Properties**

| Name  | Type   | Go type  | Required | Default | Description | Example |
| ----- | ------ | -------- | :------: | ------- | ----------- | ------- |
| label | string | `string` |          |         |             |         |

### <span id="limiter-decision-scheduler-info"></span> LimiterDecisionSchedulerInfo

**Properties**

| Name            | Type                      | Go type  | Required | Default | Description | Example |
| --------------- | ------------------------- | -------- | :------: | ------- | ----------- | ------- |
| tokens_consumed | uint64 (formatted string) | `string` |          |         |             |         |
| workload_index  | string                    | `string` |          |         |             |         |

### <span id="load-ramp"></span> LoadRamp

> The _Load Ramp_ produces a smooth and continuous traffic load that changes
> progressively over time, based on the specified steps.

Each step is defined by two parameters:

- The `target_accept_percentage`.
- The `duration` for the signal to change from the previous step's
  `target_accept_percentage` to the current step's `target_accept_percentage`.

The percentage of requests accepted starts at the `target_accept_percentage`
defined in the first step and gradually ramps up or down linearly from the
previous step's `target_accept_percentage` to the next
`target_accept_percentage`, over the `duration` specified for each step.

**Properties**

| Name                                 | Type                                        | Go type              | Required | Default | Description                                                                                            | Example |
| ------------------------------------ | ------------------------------------------- | -------------------- | :------: | ------- | ------------------------------------------------------------------------------------------------------ | ------- |
| in_ports                             | [LoadRampIns](#load-ramp-ins)               | `LoadRampIns`        |          |         |                                                                                                        |         |
| out_ports                            | [LoadRampOuts](#load-ramp-outs)             | `LoadRampOuts`       |          |         |                                                                                                        |         |
| parameters                           | [LoadRampParameters](#load-ramp-parameters) | `LoadRampParameters` |          |         |                                                                                                        |         |
| pass_through_label_values            | []string                                    | `[]string`           |          |         | Specify certain label values to be always accepted by the _Regulator_ regardless of accept percentage. |         |
| pass_through_label_values_config_key | string                                      | `string`             |          |         | Configuration key for setting pass through label values through dynamic configuration.                 |         |

### <span id="load-ramp-ins"></span> LoadRampIns

> Inputs for the _Load Ramp_ component.

**Properties**

| Name     | Type               | Go type  | Required | Default | Description                                                    | Example |
| -------- | ------------------ | -------- | :------: | ------- | -------------------------------------------------------------- | ------- |
| backward | [InPort](#in-port) | `InPort` |          |         | Whether to progress the _Load Ramp_ towards the previous step. |         |
| forward  | [InPort](#in-port) | `InPort` |          |         | Whether to progress the _Load Ramp_ towards the next step.     |         |
| reset    | [InPort](#in-port) | `InPort` |          |         | Whether to reset the _Load Ramp_ to the first step.            |         |

### <span id="load-ramp-outs"></span> LoadRampOuts

> Outputs for the _Load Ramp_ component.

**Properties**

| Name              | Type                 | Go type   | Required | Default | Description                                                                               | Example |
| ----------------- | -------------------- | --------- | :------: | ------- | ----------------------------------------------------------------------------------------- | ------- |
| accept_percentage | [OutPort](#out-port) | `OutPort` |          |         | The percentage of flows being accepted by the _Load Ramp_.                                |         |
| at_end            | [OutPort](#out-port) | `OutPort` |          |         | A Boolean signal indicating whether the _Load Ramp_ is at the end of signal generation.   |         |
| at_start          | [OutPort](#out-port) | `OutPort` |          |         | A Boolean signal indicating whether the _Load Ramp_ is at the start of signal generation. |         |

### <span id="load-ramp-parameters"></span> LoadRampParameters

> Parameters for the _Load Ramp_ component.

**Properties**

| Name      | Type                                                   | Go type                     | Required | Default | Description                     | Example |
| --------- | ------------------------------------------------------ | --------------------------- | :------: | ------- | ------------------------------- | ------- |
| regulator | [RegulatorParameters](#regulator-parameters)           | `RegulatorParameters`       |          |         | Parameters for the _Regulator_. |         |
| steps     | [][LoadRampParametersStep](#load-ramp-parameters-step) | `[]*LoadRampParametersStep` |    ✓     |         |                                 |         |

### <span id="load-ramp-parameters-step"></span> LoadRampParametersStep

**Properties**

| Name                     | Type                      | Go type   | Required | Default | Description                            | Example |
| ------------------------ | ------------------------- | --------- | :------: | ------- | -------------------------------------- | ------- |
| duration                 | string                    | `string`  |    ✓     |         | Duration for which the step is active. |         |
| target_accept_percentage | double (formatted number) | `float64` |          |         | The value of the step.                 |         |

### <span id="load-ramp-series"></span> LoadRampSeries

> _LoadRampSeries_ is a component that applies a series of _Load Ramps_ in
> order.

**Properties**

| Name       | Type                                                     | Go type                    | Required | Default | Description | Example |
| ---------- | -------------------------------------------------------- | -------------------------- | :------: | ------- | ----------- | ------- |
| in_ports   | [LoadRampSeriesIns](#load-ramp-series-ins)               | `LoadRampSeriesIns`        |          |         |             |         |
| parameters | [LoadRampSeriesParameters](#load-ramp-series-parameters) | `LoadRampSeriesParameters` |          |         |             |         |

### <span id="load-ramp-series-ins"></span> LoadRampSeriesIns

> Inputs for the _LoadRampSeries_ component.

**Properties**

| Name     | Type               | Go type  | Required | Default | Description                                                         | Example |
| -------- | ------------------ | -------- | :------: | ------- | ------------------------------------------------------------------- | ------- |
| backward | [InPort](#in-port) | `InPort` |          |         | Whether to progress the load ramp series towards the previous step. |         |
| forward  | [InPort](#in-port) | `InPort` |          |         | Whether to progress the load ramp series towards the next step.     |         |
| reset    | [InPort](#in-port) | `InPort` |          |         | Whether to reset the load ramp series to the first step.            |         |

### <span id="load-ramp-series-load-ramp-instance"></span> LoadRampSeriesLoadRampInstance

**Properties**

| Name      | Type                                        | Go type              | Required | Default | Description    | Example |
| --------- | ------------------------------------------- | -------------------- | :------: | ------- | -------------- | ------- |
| load_ramp | [LoadRampParameters](#load-ramp-parameters) | `LoadRampParameters` |          |         | The load ramp. |         |
| out_ports | [LoadRampOuts](#load-ramp-outs)             | `LoadRampOuts`       |          |         |                |         |

### <span id="load-ramp-series-parameters"></span> LoadRampSeriesParameters

> Parameters for the _LoadRampSeries_ component.

**Properties**

| Name       | Type                                                                     | Go type                             | Required | Default | Description                                              | Example |
| ---------- | ------------------------------------------------------------------------ | ----------------------------------- | :------: | ------- | -------------------------------------------------------- | ------- |
| load_ramps | [][LoadRampSeriesLoadRampInstance](#load-ramp-series-load-ramp-instance) | `[]*LoadRampSeriesLoadRampInstance` |    ✓     |         | An ordered list of load ramps that get applied in order. |         |

### <span id="load-scheduler"></span> LoadScheduler

> :::info

See also
[_Load Scheduler_ overview](/concepts/flow-control/components/load-scheduler.md).

:::

To make scheduling decisions the Flows are mapped into Workloads by providing
match rules. A workload determines the priority and cost of admitting each Flow
that belongs to it. Scheduling of Flows is based on Weighted Fair Queuing
principles. _Load Scheduler_ measures and controls the incoming tokens per
second, which can translate to (avg. latency \* in-flight requests) (Little's
Law) in concurrency limiting use-case.

The signal at port `load_multiplier` determines the fraction of incoming tokens
that get admitted.

**Properties**

| Name                                                                                           | Type                                                  | Go type                   | Required | Default | Description                                                                                                                                         | Example |
| ---------------------------------------------------------------------------------------------- | ----------------------------------------------------- | ------------------------- | :------: | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| dry_run                                                                                        | boolean                                               | `bool`                    |          |         | Decides whether to run the load scheduler in dry-run mode. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. |
| It is useful for observing the behavior of load scheduler without disrupting any real traffic. |                                                       |
| dry_run_config_key                                                                             | string                                                | `string`                  |          |         | Configuration key for setting dry run mode through dynamic configuration.                                                                           |         |
| in_ports                                                                                       | [LoadSchedulerIns](#load-scheduler-ins)               | `LoadSchedulerIns`        |          |         | Input ports for the LoadScheduler component.                                                                                                        |         |
| out_ports                                                                                      | [LoadSchedulerOuts](#load-scheduler-outs)             | `LoadSchedulerOuts`       |          |         | Output ports for the LoadScheduler component.                                                                                                       |         |
| parameters                                                                                     | [LoadSchedulerParameters](#load-scheduler-parameters) | `LoadSchedulerParameters` |          |         |                                                                                                                                                     |         |

### <span id="load-scheduler-ins"></span> LoadSchedulerIns

> Input for the LoadScheduler component.

**Properties**

| Name                                  | Type               | Go type  | Required | Default | Description                               | Example |
| ------------------------------------- | ------------------ | -------- | :------: | ------- | ----------------------------------------- | ------- |
| load_multiplier                       | [InPort](#in-port) | `InPort` |          |         | Load multiplier is proportion of incoming |
| token rate that needs to be accepted. |                    |

### <span id="load-scheduler-outs"></span> LoadSchedulerOuts

> Output for the LoadScheduler component.

**Properties**

| Name                     | Type                 | Go type   | Required | Default | Description                                                                               | Example |
| ------------------------ | -------------------- | --------- | :------: | ------- | ----------------------------------------------------------------------------------------- | ------- |
| observed_load_multiplier | [OutPort](#out-port) | `OutPort` |          |         | Observed load multiplier is the proportion of incoming token rate that is being accepted. |         |

### <span id="load-scheduler-parameters"></span> LoadSchedulerParameters

**Properties**

| Name      | Type                    | Go type     | Required | Default | Description                                                      | Example |
| --------- | ----------------------- | ----------- | :------: | ------- | ---------------------------------------------------------------- | ------- |
| scheduler | [Scheduler](#scheduler) | `Scheduler` |          |         | Configuration of Weighted Fair Queuing-based workload scheduler. |

Contains configuration of per-agent scheduler | | | selectors |
[][Selector](#selector)| `[]*Selector` | ✓ | | Selectors for the component. | |
| workload_latency_based_tokens | boolean| `bool` | | `true`| Automatically
estimate the size flows within each workload, based on historical latency. Each
workload's `tokens` will be set to average latency of flows in that workload
during the last few seconds (exact duration of this average can change). This
setting is useful in concurrency limiting use-case, where the concurrency is
calculated as ``(avg. latency \* in-flight flows)`.

The value of tokens estimated takes a lower precedence than the value of
`tokens` specified in the workload definition and `tokens` explicitly specified
in the flow labels. | |

### <span id="match-expression"></span> MatchExpression

> MatchExpression has multiple variants, exactly one should be set.

Example:

```yaml
all:
  of:
    - label_exists: foo
    - label_equals:
        label: app
        value: frobnicator
```

**Properties**

| Name          | Type                                                | Go type                  | Required | Default | Description                                                         | Example |
| ------------- | --------------------------------------------------- | ------------------------ | :------: | ------- | ------------------------------------------------------------------- | ------- |
| all           | [MatchExpressionList](#match-expression-list)       | `MatchExpressionList`    |          |         | The expression is true when all sub expressions are true.           |         |
| any           | [MatchExpressionList](#match-expression-list)       | `MatchExpressionList`    |          |         | The expression is true when any sub expression is true.             |         |
| label_equals  | [EqualsMatchExpression](#equals-match-expression)   | `EqualsMatchExpression`  |          |         | The expression is true when label value equals given value.         |         |
| label_exists  | string                                              | `string`                 |          |         | The expression is true when label with given name exists.           |         |
| label_matches | [MatchesMatchExpression](#matches-match-expression) | `MatchesMatchExpression` |          |         | The expression is true when label matches given regular expression. |         |
| not           | [MatchExpression](#match-expression)                | `MatchExpression`        |          |         | The expression negates the result of sub expression.                |         |

### <span id="match-expression-list"></span> MatchExpressionList

> for example, `{any: {of: [expr1, expr2]}}`.

**Properties**

| Name | Type                                   | Go type              | Required | Default | Description                                      | Example |
| ---- | -------------------------------------- | -------------------- | :------: | ------- | ------------------------------------------------ | ------- |
| of   | [][MatchExpression](#match-expression) | `[]*MatchExpression` |          |         | List of sub expressions of the match expression. |         |

### <span id="matches-match-expression"></span> MatchesMatchExpression

> Label selector expression of the form `label matches regex`.

**Properties**

| Name                                                                                 | Type   | Go type  | Required | Default | Description                                           | Example |
| ------------------------------------------------------------------------------------ | ------ | -------- | :------: | ------- | ----------------------------------------------------- | ------- |
| label                                                                                | string | `string` |    ✓     |         | Name of the label to match the regular expression.    |         |
| regex                                                                                | string | `string` |    ✓     |         | Regular expression that should match the label value. |
| It uses [Go's regular expression syntax](https://github.com/google/re2/wiki/Syntax). |        |

### <span id="max"></span> Max

> Max: output = max([]inputs).

**Properties**

| Name      | Type                 | Go type   | Required | Default | Description                         | Example |
| --------- | -------------------- | --------- | :------: | ------- | ----------------------------------- | ------- |
| in_ports  | [MaxIns](#max-ins)   | `MaxIns`  |          |         | Input ports for the Max component.  |         |
| out_ports | [MaxOuts](#max-outs) | `MaxOuts` |          |         | Output ports for the Max component. |         |

### <span id="max-ins"></span> MaxIns

> Inputs for the Max component.

**Properties**

| Name   | Type                 | Go type     | Required | Default | Description             | Example |
| ------ | -------------------- | ----------- | :------: | ------- | ----------------------- | ------- |
| inputs | [][InPort](#in-port) | `[]*InPort` |          |         | Array of input signals. |         |

### <span id="max-outs"></span> MaxOuts

> Output for the Max component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                                    | Example |
| ------ | -------------------- | --------- | :------: | ------- | ---------------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Signal with maximum value as an output signal. |         |

### <span id="member"></span> Member

**Properties**

| Name      | Type                      | Go type  | Required | Default | Description | Example |
| --------- | ------------------------- | -------- | :------: | ------- | ----------- | ------- |
| birthdate | int64 (formatted string)  | `string` |          |         |             |         |
| id        | uint64 (formatted string) | `string` |          |         |             |         |
| name      | string                    | `string` |          |         |             |         |

### <span id="min"></span> Min

> Takes an array of input signals and emits the signal with the minimum value
> Min: output = min([]inputs).

**Properties**

| Name      | Type                 | Go type   | Required | Default | Description                         | Example |
| --------- | -------------------- | --------- | :------: | ------- | ----------------------------------- | ------- |
| in_ports  | [MinIns](#min-ins)   | `MinIns`  |          |         | Input ports for the Min component.  |         |
| out_ports | [MinOuts](#min-outs) | `MinOuts` |          |         | Output ports for the Min component. |         |

### <span id="min-ins"></span> MinIns

> Inputs for the Min component.

**Properties**

| Name   | Type                 | Go type     | Required | Default | Description             | Example |
| ------ | -------------------- | ----------- | :------: | ------- | ----------------------- | ------- |
| inputs | [][InPort](#in-port) | `[]*InPort` |          |         | Array of input signals. |         |

### <span id="min-outs"></span> MinOuts

> Output ports for the Min component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                                    | Example |
| ------ | -------------------- | --------- | :------: | ------- | ---------------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Signal with minimum value as an output signal. |         |

### <span id="nested-circuit"></span> NestedCircuit

> Nested circuit defines a sub-circuit as a high-level component. It consists of
> a list of components and a map of input and output ports.

**Properties**

| Name              | Type                        | Go type              | Required | Default | Description                                                                                                    | Example |
| ----------------- | --------------------------- | -------------------- | :------: | ------- | -------------------------------------------------------------------------------------------------------------- | ------- |
| components        | [][Component](#component)   | `[]*Component`       |          |         | List of components in the nested circuit.                                                                      |         |
| in_ports_map      | map of [InPort](#in-port)   | `map[string]InPort`  |          |         | Maps input port names to input ports.                                                                          |         |
| name              | string                      | `string`             |          |         | Name of the nested circuit component. This name is displayed by graph visualization tools.                     |         |
| out_ports_map     | map of [OutPort](#out-port) | `map[string]OutPort` |          |         | Maps output port names to output ports.                                                                        |         |
| short_description | string                      | `string`             |          |         | Short description of the nested circuit component. This description is displayed by graph visualization tools. |         |

### <span id="nested-signal-egress"></span> NestedSignalEgress

> Nested signal egress is a special type of component that allows to extract a
> signal from a nested circuit.

**Properties**

| Name      | Type                                               | Go type                 | Required | Default | Description                                       | Example |
| --------- | -------------------------------------------------- | ----------------------- | :------: | ------- | ------------------------------------------------- | ------- |
| in_ports  | [NestedSignalEgressIns](#nested-signal-egress-ins) | `NestedSignalEgressIns` |          |         | Input ports for the NestedSignalEgress component. |         |
| port_name | string                                             | `string`                |          |         | Name of the port.                                 |         |

### <span id="nested-signal-egress-ins"></span> NestedSignalEgressIns

> Inputs for the NestedSignalEgress component.

**Properties**

| Name   | Type               | Go type  | Required | Default | Description    | Example |
| ------ | ------------------ | -------- | :------: | ------- | -------------- | ------- |
| signal | [InPort](#in-port) | `InPort` |          |         | Egress signal. |         |

### <span id="nested-signal-ingress"></span> NestedSignalIngress

> Nested signal ingress is a special type of component that allows to inject a
> signal into a nested circuit.

**Properties**

| Name      | Type                                                   | Go type                   | Required | Default | Description                                         | Example |
| --------- | ------------------------------------------------------ | ------------------------- | :------: | ------- | --------------------------------------------------- | ------- |
| out_ports | [NestedSignalIngressOuts](#nested-signal-ingress-outs) | `NestedSignalIngressOuts` |          |         | Output ports for the NestedSignalIngress component. |         |
| port_name | string                                                 | `string`                  |          |         | Name of the port.                                   |         |

### <span id="nested-signal-ingress-outs"></span> NestedSignalIngressOuts

> Outputs for the NestedSignalIngress component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description     | Example |
| ------ | -------------------- | --------- | :------: | ------- | --------------- | ------- |
| signal | [OutPort](#out-port) | `OutPort` |          |         | Ingress signal. |         |

### <span id="or"></span> Or

> Logical OR.

See [And component](#and) on how signals are mapped onto Boolean values.

**Properties**

| Name      | Type               | Go type  | Required | Default | Description                        | Example |
| --------- | ------------------ | -------- | :------: | ------- | ---------------------------------- | ------- |
| in_ports  | [OrIns](#or-ins)   | `OrIns`  |          |         | Input ports for the Or component.  |         |
| out_ports | [OrOuts](#or-outs) | `OrOuts` |          |         | Output ports for the Or component. |         |

### <span id="or-ins"></span> OrIns

> Inputs for the Or component.

**Properties**

| Name   | Type                 | Go type     | Required | Default | Description             | Example |
| ------ | -------------------- | ----------- | :------: | ------- | ----------------------- | ------- |
| inputs | [][InPort](#in-port) | `[]*InPort` |          |         | Array of input signals. |         |

### <span id="or-outs"></span> OrOuts

> Output ports for the Or component.

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description                                    | Example |
| ------ | -------------------- | --------- | :------: | ------- | ---------------------------------------------- | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Result of logical OR of all the input signals. |

Will always be 0 (false), 1 (true) or invalid (unknown). | |

### <span id="out-port"></span> OutPort

**Properties**

| Name        | Type   | Go type  | Required | Default | Description                                 | Example |
| ----------- | ------ | -------- | :------: | ------- | ------------------------------------------- | ------- |
| signal_name | string | `string` |          |         | Name of the outgoing Signal on the OutPort. |         |

### <span id="overlapping-service"></span> OverlappingService

> OverlappingService contains info about a service that overlaps with another
> one.

**Properties**

| Name           | Type                      | Go type  | Required | Default | Description | Example |
| -------------- | ------------------------- | -------- | :------: | ------- | ----------- | ------- |
| entities_count | int32 (formatted integer) | `int32`  |          |         |             |         |
| service1       | string                    | `string` |          |         |             |         |
| service2       | string                    | `string` |          |         |             |         |

### <span id="path-template-matcher"></span> PathTemplateMatcher

> HTTP path will be matched against given path templates. If a match occurs, the
> value associated with the path template will be treated as a result. In case
> of multiple path templates matching, the most specific one will be chosen.

**Properties**

| Name            | Type          | Go type             | Required | Default | Description                                              | Example |
| --------------- | ------------- | ------------------- | :------: | ------- | -------------------------------------------------------- | ------- |
| template_values | map of string | `map[string]string` |          |         | Template value keys are OpenAPI-inspired path templates. |

- Static path segment `/foo` matches a path segment exactly
- `/{param}` matches arbitrary path segment. (The parameter name is ignored and
  can be omitted (`{}`))
- The parameter must cover whole segment.
- Additionally, path template can end with `/*` wildcard to match arbitrary
  number of trailing segments (0 or more).
- Multiple consecutive `/` are ignored, as well as trailing `/`.
- Parametrized path segments must come after static segments.
- `*`, if present, must come last.
- Most specific template "wins" (`/foo` over `/{}` and `/{}` over `/*`).

See also <https://swagger.io/specification/#path-templating-matching>

Example:

````yaml
/register: register
"/user/{userId}": user
/static/*: other
``` |  |



### <span id="peer"></span> Peer


> Peer holds peer info and services.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| address | string| `string` |  | |  |  |
| hostname | string| `string` |  | |  |  |
| services | map of string| `map[string]string` |  | |  |  |



### <span id="periodic-decrease"></span> PeriodicDecrease


> PeriodicDecrease defines a controller for scaling in based on a periodic timer.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| period | string| `string` | ✓ | | The period of the timer. |  |
| scale_in_percentage | double (formatted number)| `float64` | ✓ | | The percentage of scale to reduce. |  |



### <span id="pod-scaler"></span> PodScaler


> Component for scaling pods based on a signal.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| dry_run | boolean| `bool` |  | | Dry run mode ensures that no scaling is invoked by this pod scaler.
This is useful for observing the behavior of pod scaler without disrupting any real deployment.
This parameter sets the default value of dry run setting which can be overridden at runtime using dynamic configuration. |  |
| dry_run_config_key | string| `string` |  | | Configuration key for overriding dry run setting through dynamic configuration. |  |
| in_ports | [PodScalerIns](#pod-scaler-ins)| `PodScalerIns` |  | | Input ports for the PodScaler component. |  |
| kubernetes_object_selector | [KubernetesObjectSelector](#kubernetes-object-selector)| `KubernetesObjectSelector` |  | | The Kubernetes object to which this pod scaler applies. |  |
| out_ports | [PodScalerOuts](#pod-scaler-outs)| `PodScalerOuts` |  | | Output ports for the PodScaler component. |  |



### <span id="pod-scaler-ins"></span> PodScalerIns


> Inputs for the PodScaler component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| replicas | [InPort](#in-port)| `InPort` |  | |  |  |



### <span id="pod-scaler-outs"></span> PodScalerOuts


> Outputs for the PodScaler component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| actual_replicas | [OutPort](#out-port)| `OutPort` |  | | The number of replicas that are currently running. |  |
| configured_replicas | [OutPort](#out-port)| `OutPort` |  | | The number of replicas that are desired. |  |



### <span id="policy"></span> Policy


> Policy expresses observability-driven control logic.

:::info

See also [Policy overview](/concepts/policy/policy.md).

:::

Policy specification contains a circuit that defines the controller logic and resources that need to be setup.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| circuit | [Circuit](#circuit)| `Circuit` |  | | Defines the control-loop logic of the policy. |  |
| resources | [Resources](#resources)| `Resources` |  | | Resources (such as Flux Meters, Classifiers) to setup. |  |



### <span id="policy-wrapper"></span> PolicyWrapper






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| common_attributes | [CommonAttributes](#common-attributes)| `CommonAttributes` |  | |  |  |
| policy | [Policy](#policy)| `Policy` |  | |  |  |



### <span id="preview-flow-labels-response-flow-labels"></span> PreviewFlowLabelsResponseFlowLabels






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| labels | map of string| `map[string]string` |  | |  |  |



### <span id="prom-q-l"></span> PromQL






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| evaluation_interval | string| `string` |  | `"10s"`| Describes the interval between successive evaluations of the Prometheus query. |  |
| out_ports | [PromQLOuts](#prom-q-l-outs)| `PromQLOuts` |  | | Output ports for the PromQL component. |  |
| query_string | string| `string` |  | | Describes the [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) query to be run.

:::note

The query must return a single value either as a scalar or as a vector with a single element.

:::

:::info Usage with Flux Meter

[Flux Meter](/concepts/flow-control/resources/flux-meter.md) metrics can be queries using PromQL. Flux Meter defines histogram type of metrics in Prometheus.
Therefore, one can refer to `flux_meter_sum`, `flux_meter_count` and `flux_meter_bucket`.
The particular Flux Meter can be identified with the `flux_meter_name` label.
There are additional labels available on a Flux Meter such as `valid`, `flow_status`, `http_status_code` and `decision_type`.

:::

:::info Usage with OpenTelemetry Metrics

Aperture supports OpenTelemetry metrics. See [reference](/integrations/metrics/metrics.md) for more details.

::: |  |



### <span id="prom-q-l-outs"></span> PromQLOuts


> Output for the PromQL component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | | The result of the Prometheus query as an output signal. |  |



### <span id="pulse-generator"></span> PulseGenerator


> Generates 0 and 1 in turns.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| false_for | string| `string` |  | `"5s"`| Emitting 0 for the `false_for` duration. |  |
| out_ports | [PulseGeneratorOuts](#pulse-generator-outs)| `PulseGeneratorOuts` |  | | Output ports for the PulseGenerator component. |  |
| true_for | string| `string` |  | `"5s"`| Emitting 1 for the `true_for` duration. |  |



### <span id="pulse-generator-outs"></span> PulseGeneratorOuts


> Outputs for the PulseGenerator component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | |  |  |



### <span id="query"></span> Query


> Query components that are query databases such as Prometheus.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| promql | [PromQL](#prom-q-l)| `PromQL` |  | | Periodically runs a Prometheus query in the background and emits the result. |  |



### <span id="quota-scheduler"></span> QuotaScheduler


> Schedules the traffic based on token-bucket based quotas.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [RateLimiterIns](#rate-limiter-ins)| `RateLimiterIns` |  | |  |  |
| rate_limiter | [RateLimiterParameters](#rate-limiter-parameters)| `RateLimiterParameters` |  | |  |  |
| scheduler | [Scheduler](#scheduler)| `Scheduler` |  | |  |  |
| selectors | [][Selector](#selector)| `[]*Selector` | ✓ | |  |  |



### <span id="rate-limiter"></span> RateLimiter


> :::info

See also [_Rate Limiter_ overview](/concepts/flow-control/components/rate-limiter.md).

:::

RateLimiting is done on per-label-value (`label_key`) basis and it uses the _Token Bucket Algorithm_.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [RateLimiterIns](#rate-limiter-ins)| `RateLimiterIns` |  | |  |  |
| parameters | [RateLimiterParameters](#rate-limiter-parameters)| `RateLimiterParameters` |  | |  |  |
| selectors | [][Selector](#selector)| `[]*Selector` | ✓ | | Selectors for the component. |  |



### <span id="rate-limiter-ins"></span> RateLimiterIns






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| bucket_capacity | [InPort](#in-port)| `InPort` |  | | Capacity of the bucket. |  |
| fill_amount | [InPort](#in-port)| `InPort` |  | | Number of tokens to fill within an `interval`. |  |
| pass_through | [InPort](#in-port)| `InPort` |  | |  |  |



### <span id="rate-limiter-parameters"></span> RateLimiterParameters






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| continuous_fill | boolean| `bool` |  | `true`| Continuous fill determines whether the token bucket should be filled
continuously or only on discrete intervals. |  |
| interval | string| `string` | ✓ | | Interval defines the time interval in which the token bucket
will fill tokens specified by `fill_amount` signal. |  |
| label_key | string| `string` |  | | Specifies which label the rate limiter should be keyed by.

Rate limiting is done independently for each value of the
[label](/concepts/flow-control/flow-label.md) with given key.
For example, to give each user a separate limit, assuming you
have a _user_ flow
label set up, set `label_key: "user"`.
If no label key is specified, then all requests matching the
selectors will be rate limited based on the global bucket. |  |
| lazy_sync | [RateLimiterParametersLazySync](#rate-limiter-parameters-lazy-sync)| `RateLimiterParametersLazySync` |  | |  |  |
| max_idle_time | string| `string` |  | `"7200s"`| Max idle time before token bucket state for a label is removed.
If set to 0, the state is never removed. |  |
| tokens_label_key | string| `string` |  | `"tokens"`| Flow label key that will be used to override the number of tokens
for this request.
This is an optional parameter and takes highest precedence
when assigning tokens to a request.
The label value must be a valid uint64 number. |  |



### <span id="rate-limiter-parameters-lazy-sync"></span> RateLimiterParametersLazySync






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| enabled | boolean| `bool` |  | |  |  |
| num_sync | int64 (formatted integer)| `int64` |  | `4`| Number of times to lazy sync within the `interval`. |  |



### <span id="rego"></span> Rego


> Rego define a set of labels that are extracted after evaluating a Rego module.

:::info

You can use the [live-preview](/concepts/flow-control/resources/classifier.md#live-previewing-requests) feature to first preview the input to the classifier before writing the labeling logic.

:::

Example of Rego module which also disables telemetry visibility of label:
```yaml
rego:
  labels:
    user:
      telemetry: false
  module: |
    package user_from_cookie
    cookies := split(input.attributes.request.http.headers.cookie, "; ")
    user := user {
        cookie := cookies[_]
        startswith(cookie, "session=")
        session := substring(cookie, count("session="), -1)
        parts := split(session, ".")
        object := json.unmarshal(base64url.decode(parts[0]))
        user := object.user
    }
````

**Properties**

| Name   | Type                                                 | Go type                          | Required | Default | Description                              | Example |
| ------ | ---------------------------------------------------- | -------------------------------- | :------: | ------- | ---------------------------------------- | ------- |
| labels | map of [RegoLabelProperties](#rego-label-properties) | `map[string]RegoLabelProperties` |    ✓     |         | A map of {key, value} pairs mapping from |

[flow label](/concepts/flow-control/flow-label.md) keys to queries that define
how to extract and propagate flow labels with that key. The name of the label
maps to a variable in the Rego module. It maps to `data.<package>.<label>`
variable. | | | module | string| `string` | ✓ | | Source code of the Rego
module.

:::note

Must include a "package" declaration.

::: | |

### <span id="rego-label-properties"></span> RegoLabelProperties

**Properties**

| Name      | Type    | Go type | Required | Default | Description | Example |
| --------- | ------- | ------- | :------: | ------- | ----------- | ------- |
| telemetry | boolean | `bool`  |          | `true`  | :::note     |

The flow label is always accessible in Aperture Policies regardless of this
setting.

:::

:::caution

When using [FluxNinja ARC extension](/arc/extension.md), telemetry enabled
labels are sent to FluxNinja ARC for observability. Telemetry should be disabled
for sensitive labels.

::: | |

### <span id="regulator"></span> Regulator

> _Regulator_ is a component that regulates the load at a
> [_Control Point_](/concepts/flow-control/selector.md/#control-point) by
> allowing only a specified percentage of flows at random or by sticky sessions.

:::info

See also
[\_Load Regulator overview](/concepts/flow-control/components/regulator.md).

:::

**Properties**

| Name                                 | Type                                         | Go type               | Required | Default | Description                                                                                             | Example |
| ------------------------------------ | -------------------------------------------- | --------------------- | :------: | ------- | ------------------------------------------------------------------------------------------------------- | ------- |
| in_ports                             | [RegulatorIns](#regulator-ins)               | `RegulatorIns`        |          |         | Input ports for the _Regulator_.                                                                        |         |
| parameters                           | [RegulatorParameters](#regulator-parameters) | `RegulatorParameters` |          |         | Parameters for the _Regulator_.                                                                         |         |
| pass_through_label_values            | []string                                     | `[]string`            |          |         | Specify certain label values to be always accepted by this _Regulator_ regardless of accept percentage. |         |
| pass_through_label_values_config_key | string                                       | `string`              |          |         | Configuration key for setting pass through label values through dynamic configuration.                  |         |

### <span id="regulator-ins"></span> RegulatorIns

**Properties**

| Name              | Type               | Go type  | Required | Default | Description                           | Example |
| ----------------- | ------------------ | -------- | :------: | ------- | ------------------------------------- | ------- |
| accept_percentage | [InPort](#in-port) | `InPort` |          |         | The percentage of requests to accept. |         |

### <span id="regulator-parameters"></span> RegulatorParameters

**Properties**

| Name      | Type   | Go type  | Required | Default | Description                                  | Example |
| --------- | ------ | -------- | :------: | ------- | -------------------------------------------- | ------- |
| label_key | string | `string` |          |         | The flow label key for identifying sessions. |

- When label key is specified, _Regulator_ acts as a sticky filter. The series
  of flows with the same value of label key get the same decision provided that
  the `accept_percentage` is same or higher.
- When label key is not specified, _Regulator_ acts as a stateless filter.
  Percentage of flows are selected randomly for rejection. | | | selectors |
  [][Selector](#selector)| `[]*Selector` | ✓ | | Selectors for the component. |
  |

### <span id="resources"></span> Resources

> :::info

See also [Resources overview](/concepts/policy/resources.md).

:::

**Properties**

| Name                 | Type                                            | Go type                 | Required | Default | Description                                                                       | Example |
| -------------------- | ----------------------------------------------- | ----------------------- | :------: | ------- | --------------------------------------------------------------------------------- | ------- |
| flow_control         | [FlowControlResources](#flow-control-resources) | `FlowControlResources`  |          |         | FlowControlResources are resources that are provided by flow control integration. |         |
| telemetry_collectors | [][TelemetryCollector](#telemetry-collector)    | `[]*TelemetryCollector` |          |         | TelemetryCollector configures OpenTelemetry collector integration.                |         |

### <span id="rule"></span> Rule

> Example of a JSON extractor:

```yaml
extractor:
  json:
    from: request.http.body
    pointer: /user/name
```

**Properties**

| Name      | Type                    | Go type     | Required | Default | Description                       | Example |
| --------- | ----------------------- | ----------- | :------: | ------- | --------------------------------- | ------- |
| extractor | [Extractor](#extractor) | `Extractor` |          |         | High-level declarative extractor. |         |
| telemetry | boolean                 | `bool`      |          | `true`  | :::note                           |

The flow label is always accessible in Aperture Policies regardless of this
setting.

:::

:::caution

When using [FluxNinja ARC extension](/arc/extension.md), telemetry enabled
labels are sent to FluxNinja ARC for observability. Telemetry should be disabled
for sensitive labels.

::: | |

### <span id="s-m-a"></span> SMA

> Simple Moving Average (SMA) is a type of moving average that computes the
> average of a fixed number of signal readings.

**Properties**

| Name       | Type                               | Go type         | Required | Default | Description                         | Example |
| ---------- | ---------------------------------- | --------------- | :------: | ------- | ----------------------------------- | ------- |
| in_ports   | [SMAIns](#s-m-a-ins)               | `SMAIns`        |          |         | Input ports for the SMA component.  |         |
| out_ports  | [SMAOuts](#s-m-a-outs)             | `SMAOuts`       |          |         | Output ports for the SMA component. |         |
| parameters | [SMAParameters](#s-m-a-parameters) | `SMAParameters` |          |         | Parameters for the SMA component.   |         |

### <span id="s-m-a-ins"></span> SMAIns

**Properties**

| Name  | Type               | Go type  | Required | Default | Description                                           | Example |
| ----- | ------------------ | -------- | :------: | ------- | ----------------------------------------------------- | ------- |
| input | [InPort](#in-port) | `InPort` |          |         | Signal to be used for the moving average computation. |         |

### <span id="s-m-a-outs"></span> SMAOuts

**Properties**

| Name   | Type                 | Go type   | Required | Default | Description              | Example |
| ------ | -------------------- | --------- | :------: | ------- | ------------------------ | ------- |
| output | [OutPort](#out-port) | `OutPort` |          |         | Computed moving average. |         |

### <span id="s-m-a-parameters"></span> SMAParameters

**Properties**

| Name                | Type    | Go type  | Required | Default | Description                                               | Example |
| ------------------- | ------- | -------- | :------: | ------- | --------------------------------------------------------- | ------- |
| sma_window          | string  | `string` |    ✓     |         | Window of time over which the moving average is computed. |         |
| valid_during_warmup | boolean | `bool`   |          |         | Whether the output is valid during the warm-up stage.     |         |

### <span id="scale-in-controller"></span> ScaleInController

**Properties**

| Name       | Type                                                           | Go type                       | Required | Default | Description                         | Example |
| ---------- | -------------------------------------------------------------- | ----------------------------- | :------: | ------- | ----------------------------------- | ------- |
| alerter    | [AlerterParameters](#alerter-parameters)                       | `AlerterParameters`           |          |         | Configuration for embedded Alerter. |         |
| controller | [ScaleInControllerController](#scale-in-controller-controller) | `ScaleInControllerController` |          |         |                                     |         |

### <span id="scale-in-controller-controller"></span> ScaleInControllerController

**Properties**

| Name     | Type                                       | Go type              | Required | Default | Description | Example |
| -------- | ------------------------------------------ | -------------------- | :------: | ------- | ----------- | ------- |
| gradient | [DecreasingGradient](#decreasing-gradient) | `DecreasingGradient` |          |         |             |         |
| periodic | [PeriodicDecrease](#periodic-decrease)     | `PeriodicDecrease`   |          |         |             |         |

### <span id="scale-out-controller"></span> ScaleOutController

**Properties**

| Name       | Type                                                             | Go type                        | Required | Default | Description                         | Example |
| ---------- | ---------------------------------------------------------------- | ------------------------------ | :------: | ------- | ----------------------------------- | ------- |
| alerter    | [AlerterParameters](#alerter-parameters)                         | `AlerterParameters`            |          |         | Configuration for embedded Alerter. |         |
| controller | [ScaleOutControllerController](#scale-out-controller-controller) | `ScaleOutControllerController` |          |         |                                     |         |

### <span id="scale-out-controller-controller"></span> ScaleOutControllerController

**Properties**

| Name     | Type                                       | Go type              | Required | Default | Description | Example |
| -------- | ------------------------------------------ | -------------------- | :------: | ------- | ----------- | ------- |
| gradient | [IncreasingGradient](#increasing-gradient) | `IncreasingGradient` |          |         |             |         |

### <span id="scheduler"></span> Scheduler

> :::note

Each Agent instantiates an independent copy of the scheduler, but output signals
for accepted and incoming token rate are aggregated across all agents.

:::

**Properties**

| Name                     | Type   | Go type  | Required | Default   | Description                                                            | Example |
| ------------------------ | ------ | -------- | :------: | --------- | ---------------------------------------------------------------------- | ------- |
| decision_deadline_margin | string | `string` |          | `"0.01s"` | Decision deadline margin is the amount of time that the scheduler will |

subtract from the request deadline to determine the deadline for the decision.
This is to ensure that the scheduler has enough time to make a decision before
the request deadline happens, accounting for processing delays. The request
deadline is based on the [gRPC deadline](https://grpc.io/blog/deadlines) or the
[`grpc-timeout` HTTP header](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests).

Fail-open logic is use for flow control APIs, so if the gRPC deadline reaches,
the flow will end up being unconditionally allowed while it is still waiting on
the scheduler. | | | default_workload_parameters |
[SchedulerWorkloadParameters](#scheduler-workload-parameters)|
`SchedulerWorkloadParameters` | | | Parameters to be used if none of workloads
specified in `workloads` match. | | | tokens_label_key | string| `string` | |
`"tokens"`| \* Key for a flow label that can be used to override the default
number of tokens for this flow.

- The value associated with this key must be a valid uint64 number.
- If this parameter is not provided, the number of tokens for the flow will be
  determined by the matched workload's token count. | | | workloads |
  [][SchedulerWorkload](#scheduler-workload)| `[]*SchedulerWorkload` | | | List
  of workloads to be used in scheduler.

Categorizing [flows](/concepts/flow-control/flow-control.md#flow) into workloads
allows for load throttling to be "intelligent" instead of queueing flows in an
arbitrary order. There are two aspects of this "intelligence":

- Scheduler can more precisely calculate concurrency if it understands that
  flows belonging to different classes have different weights (for example,
  insert queries compared to select queries).
- Setting different priorities to different workloads lets the scheduler avoid
  dropping important traffic during overload.

Each workload in this list specifies also a matcher that is used to determine
which flow will be categorized into which workload. In case of multiple matching
workloads, the first matching one will be used. If none of workloads match,
`default_workload` will be used.

:::info

See also
[workload definition in the concepts section](/concepts/flow-control/components/load-scheduler.md#workload).

::: | |

### <span id="scheduler-workload"></span> SchedulerWorkload

> Workload defines a class of flows that preferably have similar properties such
> as response latency and desired priority.

**Properties**

| Name                                                 | Type                                                          | Go type                       | Required | Default | Description                                                  | Example |
| ---------------------------------------------------- | ------------------------------------------------------------- | ----------------------------- | :------: | ------- | ------------------------------------------------------------ | ------- |
| label_matcher                                        | [LabelMatcher](#label-matcher)                                | `LabelMatcher`                |          |         | Label Matcher to select a Workload based on                  |
| [flow labels](/concepts/flow-control/flow-label.md). |                                                               |
| parameters                                           | [SchedulerWorkloadParameters](#scheduler-workload-parameters) | `SchedulerWorkloadParameters` |          |         | Parameters associated with flows matching the label matcher. |         |

### <span id="scheduler-workload-parameters"></span> SchedulerWorkloadParameters

> Parameters such as priority, tokens and fairness key that are applicable to
> flows within a workload.

**Properties**

| Name         | Type   | Go type  | Required | Default | Description                                                                         | Example |
| ------------ | ------ | -------- | :------: | ------- | ----------------------------------------------------------------------------------- | ------- |
| fairness_key | string | `string` |          |         | Fairness key is a label key that can be used to provide fairness within a workload. |

Any [flow label](/concepts/flow-control/flow-label.md) can be used here. For
example, if you have a classifier that sets `user` flow label, you might want to
set `fairness_key = "user"`. | | | priority | int64 (formatted integer)| `int64`
| | | $$ \text{virtual_finish_time} = \text{virtual_time} + \left(\text{tokens}
\cdot \left(\text{256} - \text{priority}\right)\right)

$$
|  |
| tokens | uint64 (formatted string)| `string` |  | | Tokens determines the cost of admitting a single flow in the workload,
which is typically defined as milliseconds of flow latency (time to response or duration of a feature) or
simply equal to 1 if the resource being accessed is constrained by the
number of flows (3rd party rate limiters).
This override is applicable only if tokens for the flow aren't specified
in the flow labels. |  |



### <span id="selector"></span> Selector


> Selects flows based on control point, flow labels, agent group and the service
that the [flow control component](/concepts/flow-control/flow-control.md#components)
will operate on.

:::info

See also [Selector overview](/concepts/flow-control/selector.md).

:::

Example:
```yaml
control_point: ingress
label_matcher:
  match_labels:
    user_tier: gold
  match_expressions:
    - key: query
      operator: In
      values:
        - insert
        - delete
  expression:
    label_matches:
        - label: user_agent
          regex: ^(?!.*Chrome).*Safari
```






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| agent_group | string| `string` |  | `"default"`| [_Agent Group_](/concepts/flow-control/selector.md#agent-group) this
selector applies to.

:::info

Agent Groups are used to scope policies to a subset of agents connected to the same controller.
The agents within an agent group receive exact same policy configuration and
form a peer to peer cluster to constantly share state.

::: |  |
| control_point | string| `string` | ✓ | | [Control Point](/concepts/flow-control/selector.md#control-point)
identifies location within services where policies can act on flows.
For an SDK based insertion,
a _Control Point_ can represent a particular feature or execution
block within a service. In case of service mesh or middleware insertion, a
_Control Point_ can identify ingress or egress calls or distinct listeners
or filter chains. |  |
| label_matcher | [LabelMatcher](#label-matcher)| `LabelMatcher` |  | | [Label Matcher](/concepts/flow-control/selector.md#label-matcher)
can be used to match flows based on flow labels. |  |
| service | string| `string` |  | `"any"`| The Fully Qualified Domain Name of the
[service](/concepts/flow-control/selector.md) to select.

In Kubernetes, this is the FQDN of the Service object.

:::info

`any` matches all services.

:::

:::info

An entity (for example, Kubernetes pod) might belong to multiple services.

::: |  |



### <span id="service"></span> Service


> Service contains information about single service discovered in agent group by a
particular agent.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| entities_count | int32 (formatted integer)| `int32` |  | |  |  |
| name | string| `string` |  | |  |  |



### <span id="signal-generator"></span> SignalGenerator


> The _Signal Generator_ component generates a smooth and continuous signal
by following a sequence of specified steps. Each step has two parameters:
- `target_output`: The desired output value at the end of the step.
- `duration`: The time it takes for the signal to change linearly from the
  previous step's `target_output` to the current step's `target_output`.

The output signal starts at the `target_output` of the first step and
changes linearly between steps based on their `duration`. The _Signal
Generator_ can be controlled to move forwards, backwards, or reset to the
beginning based on input signals.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [SignalGeneratorIns](#signal-generator-ins)| `SignalGeneratorIns` |  | |  |  |
| out_ports | [SignalGeneratorOuts](#signal-generator-outs)| `SignalGeneratorOuts` |  | |  |  |
| parameters | [SignalGeneratorParameters](#signal-generator-parameters)| `SignalGeneratorParameters` |  | | Parameters for the _Signal Generator_ component. |  |



### <span id="signal-generator-ins"></span> SignalGeneratorIns


> Inputs for the _Signal Generator_ component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| backward | [InPort](#in-port)| `InPort` |  | | Whether to progress the _Signal Generator_ towards the previous step. |  |
| forward | [InPort](#in-port)| `InPort` |  | | Whether to progress the _Signal Generator_ towards the next step. |  |
| reset | [InPort](#in-port)| `InPort` |  | | Whether to reset the _Signal Generator_ to the first step. |  |



### <span id="signal-generator-outs"></span> SignalGeneratorOuts


> Outputs for the _Signal Generator_ component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| at_end | [OutPort](#out-port)| `OutPort` |  | | A Boolean signal indicating whether the _Signal Generator_ is at the end of signal generation. |  |
| at_start | [OutPort](#out-port)| `OutPort` |  | | A Boolean signal indicating whether the _Signal Generator_ is at the start of signal generation. |  |
| output | [OutPort](#out-port)| `OutPort` |  | | The generated signal. |  |



### <span id="signal-generator-parameters"></span> SignalGeneratorParameters


> Parameters for the _Signal Generator_ component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| steps | [][SignalGeneratorParametersStep](#signal-generator-parameters-step)| `[]*SignalGeneratorParametersStep` | ✓ | |  |  |



### <span id="signal-generator-parameters-step"></span> SignalGeneratorParametersStep






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| duration | string| `string` | ✓ | | Duration for which the step is active. |  |
| target_output | [ConstantSignal](#constant-signal)| `ConstantSignal` |  | | The value of the step. |  |



### <span id="slab-info"></span> SlabInfo






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| allocated | int64 (formatted string)| `string` |  | |  |  |
| garbage | int64 (formatted string)| `string` |  | |  |  |
| inuse | int64 (formatted string)| `string` |  | |  |  |



### <span id="socket-address-protocol"></span> SocketAddressProtocol




| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| SocketAddressProtocol | string| string | `"TCP"`|  |  |



### <span id="status"></span> Status


> Status holds details about a status that can be reported to the registry.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| error | [StatusError](#status-error)| `StatusError` |  | |  |  |
| message | [GoogleprotobufAny](#googleprotobuf-any)| `GoogleprotobufAny` |  | |  |  |
| timestamp | date-time (formatted string)| `strfmt.DateTime` |  | |  |  |



### <span id="status-error"></span> StatusError


> Error holds raw error message and its cause in a nested field.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| message | string| `string` |  | |  |  |



### <span id="switcher"></span> Switcher


> `on_signal` will be returned if switch input is valid and not equal to 0.0 ,
 otherwise `off_signal` will be returned.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [SwitcherIns](#switcher-ins)| `SwitcherIns` |  | | Input ports for the Switcher component. |  |
| out_ports | [SwitcherOuts](#switcher-outs)| `SwitcherOuts` |  | | Output ports for the Switcher component. |  |



### <span id="switcher-ins"></span> SwitcherIns


> Inputs for the Switcher component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| off_signal | [InPort](#in-port)| `InPort` |  | | Output signal when switch is invalid or 0.0. |  |
| on_signal | [InPort](#in-port)| `InPort` |  | | Output signal when switch is valid and not 0.0. |  |
| switch | [InPort](#in-port)| `InPort` |  | | Decides whether to return `on_signal` or `off_signal`. |  |



### <span id="switcher-outs"></span> SwitcherOuts


> Outputs for the Switcher component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | | Selected signal (`on_signal` or `off_signal`). |  |



### <span id="telemetry-collector"></span> TelemetryCollector






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| agent_group | string| `string` |  | `"default"`|  |  |
| infra_meters | map of [InfraMeter](#infra-meter)| `map[string]InfraMeter` |  | | _Infra Meters_ configure custom metrics OpenTelemetry collector pipelines, which will
receive and process telemetry at the agents and send metrics to the configured Prometheus.
Key in this map refers to OTel pipeline name. Prefixing pipeline name with `metrics/`
is optional, as all the components and pipeline names would be normalized.

Example:

```yaml
 telemetry_collectors:
   - agent_group: default
     infra_meters:
	      rabbitmq:
	        processors:
	          batch:
	            send_batch_size: 10
	            timeout: 10s
	        receivers:
	          rabbitmq:
	            collection_interval: 10s
	            endpoint: http://<rabbitmq-svc-fqdn>:15672
	            password: secretpassword
	            username: admin
	        per_agent_group: true

```

:::caution

Validate the OTel configuration before applying it to the
production cluster.
Incorrect configuration will get rejected at the agents and might cause
shutdown of the agent(s).

::: |  |



### <span id="unary-operator"></span> UnaryOperator


>
$$

\text{output} = \unary_operator{\text{input}}

$$






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_ports | [UnaryOperatorIns](#unary-operator-ins)| `UnaryOperatorIns` |  | | Input ports for the UnaryOperator component. |  |
| operator | string| `string` |  | | Unary Operator to apply.

The unary operator can be one of the following:
* `abs`: Absolute value with the sign removed.
* `acos`: `arccosine`, in radians.
* `acosh`: Inverse hyperbolic cosine.
* `asin`: `arcsine`, in radians.
* `asinh`: Inverse hyperbolic sine.
* `atan`: `arctangent`, in radians.
* `atanh`: Inverse hyperbolic tangent.
* `cbrt`: Cube root.
* `ceil`: Least integer value greater than or equal to input signal.
* `cos`: `cosine`, in radians.
* `cosh`: Hyperbolic cosine.
* `erf`: Error function.
* `erfc`: Complementary error function.
* `erfcinv`: Inverse complementary error function.
* `erfinv`: Inverse error function.
* `exp`: The base-e exponential of input signal.
* `exp2`: The base-2 exponential of input signal.
* `expm1`: The base-e exponential of input signal minus 1.
* `floor`: Greatest integer value less than or equal to input signal.
* `gamma`: Gamma function.
* `j0`: Bessel function of the first kind of order 0.
* `j1`: Bessel function of the first kind of order 1.
* `lgamma`: Natural logarithm of the absolute value of the gamma function.
* `log`: Natural logarithm of input signal.
* `log10`: Base-10 logarithm of input signal.
* `log1p`: Natural logarithm of input signal plus 1.
* `log2`: Base-2 logarithm of input signal.
* `round`: Round to nearest integer.
* `roundtoeven`: Round to nearest integer, with ties going to the nearest even integer.
* `sin`: `sine`, in radians.
* `sinh`: Hyperbolic sine.
* `sqrt`: Square root.
* `tan`: `tangent`, in radians.
* `tanh`: Hyperbolic tangent.
* `trunc`: Truncate to integer.
* `y0`: Bessel function of the second kind of order 0.
* `y1`: Bessel function of the second kind of order 1. |  |
| out_ports | [UnaryOperatorOuts](#unary-operator-outs)| `UnaryOperatorOuts` |  | | Output ports for the UnaryOperator component. |  |



### <span id="unary-operator-ins"></span> UnaryOperatorIns


> Inputs for the UnaryOperator component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| input | [InPort](#in-port)| `InPort` |  | | Input signal. |  |



### <span id="unary-operator-outs"></span> UnaryOperatorOuts


> Outputs for the UnaryOperator component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | | Output signal. |  |



### <span id="variable"></span> Variable


> Component that emits a constant signal which can be changed at runtime through dynamic configuration.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| config_key | string| `string` |  | | Configuration key for overriding value setting through dynamic configuration. |  |
| constant_output | [ConstantSignal](#constant-signal)| `ConstantSignal` |  | | The constant signal emitted by this component. The value of the constant signal can be overridden at runtime through dynamic configuration. |  |
| out_ports | [VariableOuts](#variable-outs)| `VariableOuts` |  | | Output ports for the Variable component. |  |



### <span id="variable-outs"></span> VariableOuts


> Outputs for the Variable component.






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| output | [OutPort](#out-port)| `OutPort` |  | | The value is emitted to the output port. |  |



### <span id="googleprotobuf-any"></span> googleprotobufAny


> `Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C++.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }
    // or ...
    if (any.isSameTypeAs(Foo.getDefaultInstance())) {
      foo = any.unpack(Foo.getDefaultInstance());
    }

Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

Example 4: Pack and unpack a message in Go

     foo := &pb.Foo{...}
     any, err := anypb.New(foo)
     if err != nil {
       ...
     }
     ...
     foo := &pb.Foo{}
     if err := any.UnmarshalTo(foo); err != nil {
       ...
     }

The pack methods provided by protobuf library will by default use
'type.googleapis.com/full.type.name' as the type URL and the unpack
methods only use the fully qualified type name after the last '/'
in the type URL, for example "foo.bar.com/x/y.z" will yield type
name "y.z".

JSON

The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      "@type": "type.googleapis.com/google.profile.Person",
      "firstName": <string>,
      "lastName": <string>
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      "@type": "type.googleapis.com/google.protobuf.Duration",
      "value": "1.212s"
    }






**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| @type | string| `string` |  | | A URL/resource name that uniquely identifies the type of the serialized
protocol buffer message. This string must contain at least
one "/" character. The last segment of the URL's path must represent
the fully qualified name of the type (as in
`path/google.protobuf.Duration`). The name should be in a canonical form
(e.g., leading "." is not accepted).

In practice, teams usually precompile into the binary all types that they
expect it to use in the context of Any. However, for URLs which use the
scheme `http`, `https`, or no scheme, one can optionally set up a type
server that maps type URLs to message definitions as follows:

* If no scheme is provided, `https` is assumed.
* An HTTP GET on the URL must yield a [google.protobuf.Type][]
  value in binary format, or produce an error.
* Applications are allowed to cache lookup results based on the
  URL, or have them precompiled into a binary to avoid any
  lookup. Therefore, binary compatibility needs to be preserved
  on changes to types. (Use versioned type names to manage
  breaking changes.)

Note: this functionality is not currently available in the official
protobuf release, and it is not used for type URLs beginning with
type.googleapis.com.

Schemes other than `http`, `https` (or the empty scheme) might be
used with implementation specific semantics. |  |


$$
