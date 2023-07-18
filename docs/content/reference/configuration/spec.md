---
title: Policy Language Specification
sidebar_position: 3
sidebar_label: Policy Specification
---

<!-- markdownlint-disable -->
<!-- vale off -->

<head>
  <body className="schema-docs" />
</head>

<!-- vale on -->

Reference for all objects used in
[the Policy language](/concepts/advanced/policy.md).

The top-level object representing a policy is [Policy](#policy).

<!---
Generated File Starts
-->

## Objects

---

<!-- vale off -->

### AdaptiveLoadScheduler {#adaptive-load-scheduler}

<!-- vale on -->

The _Adaptive Load Scheduler_ adjusts the accepted token rate based on the
deviation of the input signal from the setpoint.

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

Decides whether to run the load scheduler in dry-run mode. In dry run mode the
scheduler acts as pass through to all flow and does not queue flows. It is
useful for observing the behavior of load scheduler without disrupting any real
traffic.

</dd>
<dt>dry_run_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for setting dry run mode through dynamic configuration.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadSchedulerIns](#adaptive-load-scheduler-ins))

<!-- vale on -->

Collection of input ports for the _Adaptive Load Scheduler_ component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadSchedulerOuts](#adaptive-load-scheduler-outs))

<!-- vale on -->

Collection of output ports for the _Adaptive Load Scheduler_ component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadSchedulerParameters](#adaptive-load-scheduler-parameters))

<!-- vale on -->

Parameters for the _Adaptive Load Scheduler_ component.

</dd>
</dl>

---

<!-- vale off -->

### AdaptiveLoadSchedulerIns {#adaptive-load-scheduler-ins}

<!-- vale on -->

Input ports for the _Adaptive Load Scheduler_ component.

<dl>
<dt>overload_confirmation</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The `overload_confirmation` port provides additional criteria to determine
overload state which results in _Flow_ throttling at the service.

</dd>
<dt>setpoint</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The setpoint input to the controller.

</dd>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The input signal to the controller.

</dd>
</dl>

---

<!-- vale off -->

### AdaptiveLoadSchedulerOuts {#adaptive-load-scheduler-outs}

<!-- vale on -->

Output ports for the _Adaptive Load Scheduler_ component.

<dl>
<dt>desired_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Desired Load multiplier is the ratio of desired token rate to the incoming token
rate.

</dd>
<dt>is_overload</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal that indicates whether the service is in overload state.

</dd>
<dt>observed_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Observed Load multiplier is the ratio of accepted token rate to the incoming
token rate.

</dd>
</dl>

---

<!-- vale off -->

### AdaptiveLoadSchedulerParameters {#adaptive-load-scheduler-parameters}

<!-- vale on -->

Parameters for the _Adaptive Load Scheduler_ component.

<dl>
<dt>alerter</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration parameters for the embedded Alerter.

</dd>
<dt>gradient</dt>
<dd>

<!-- vale off -->

([GradientControllerParameters](#gradient-controller-parameters))

<!-- vale on -->

Parameters for the _Gradient Controller_.

</dd>
<dt>load_multiplier_linear_increment</dt>
<dd>

<!-- vale off -->

(float64, default: `0.0025`)

<!-- vale on -->

Linear increment to load multiplier in each execution tick when the system is
not in the overloaded state, up until the `max_load_multiplier` is reached.

</dd>
<dt>load_scheduler</dt>
<dd>

<!-- vale off -->

([LoadSchedulerParameters](#load-scheduler-parameters))

<!-- vale on -->

Parameters for the _Load Scheduler_.

</dd>
<dt>max_load_multiplier</dt>
<dd>

<!-- vale off -->

(float64, default: `2`)

<!-- vale on -->

The maximum load multiplier that can be reached during recovery from an overload
state.

- Helps protect the service from request bursts while the system is still
  recovering.
- Once this value is reached, the scheduler enters the pass-through mode,
  allowing requests to bypass the scheduler and be sent directly to the service.
- Any future overload state is detected by the control policy, and the load
  multiplier increment cycle is restarted.

</dd>
</dl>

---

<!-- vale off -->

### AddressExtractor {#address-extractor}

<!-- vale on -->

Display an [Address][ext-authz-address] as a single string, for example,
`<ip>:<port>`

IP addresses in attribute context are defined as objects with separate IP and
port fields. This is a helper to display an address as a single string.

:::caution

This might introduce high-cardinality flow label values.

:::

[ext-authz-address]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/address.proto#config-core-v3-address

Example:

```yaml
from: "source.address # or destination.address"
```

<dl>
<dt>from</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Attribute path pointing to some string - for example, `source.address`.

</dd>
</dl>

---

<!-- vale off -->

### Alerter {#alerter}

<!-- vale on -->

Alerter reacts to a signal and generates alert to send to alert manager.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([AlerterIns](#alerter-ins))

<!-- vale on -->

Input ports for the Alerter component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Alerter configuration

</dd>
</dl>

---

<!-- vale off -->

### AlerterIns {#alerter-ins}

<!-- vale on -->

Inputs for the Alerter component.

<dl>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Signal which Alerter is monitoring. If the signal greater than 0, Alerter
generates an alert.

</dd>
</dl>

---

<!-- vale off -->

### AlerterParameters {#alerter-parameters}

<!-- vale on -->

Alerter Parameters configure parameters such as alert name, severity, resolve
timeout, alert channels and labels.

<dl>
<dt>alert_channels</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

A list of alert channel strings.

</dd>
<dt>alert_name</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Name of the alert.

</dd>
<dt>labels</dt>
<dd>

<!-- vale off -->

(map of string)

<!-- vale on -->

Additional labels to add to alert.

</dd>
<dt>resolve_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Duration of alert resolver. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>severity</dt>
<dd>

<!-- vale off -->

(string, one of: `info | warn | crit`, default: `"info"`)

<!-- vale on -->

Severity of the alert, one of 'info', 'warn' or 'crit'.

</dd>
</dl>

---

<!-- vale off -->

### And {#and}

<!-- vale on -->

Logical AND.

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

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([AndIns](#and-ins))

<!-- vale on -->

Input ports for the And component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([AndOuts](#and-outs))

<!-- vale on -->

Output ports for the And component.

</dd>
</dl>

---

<!-- vale off -->

### AndIns {#and-ins}

<!-- vale on -->

Inputs for the And component.

<dl>
<dt>inputs</dt>
<dd>

<!-- vale off -->

([[]InPort](#in-port))

<!-- vale on -->

Array of input signals.

</dd>
</dl>

---

<!-- vale off -->

### AndOuts {#and-outs}

<!-- vale on -->

Output ports for the And component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Result of logical AND of all the input signals.

Will always be 0 (false), 1 (true) or invalid (unknown).

</dd>
</dl>

---

<!-- vale off -->

### ArithmeticCombinator {#arithmetic-combinator}

<!-- vale on -->

Type of Combinator that computes the arithmetic operation on the operand signals

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([ArithmeticCombinatorIns](#arithmetic-combinator-ins))

<!-- vale on -->

Input ports for the Arithmetic Combinator component.

</dd>
<dt>operator</dt>
<dd>

<!-- vale off -->

(string, one of: `add | sub | mul | div | xor | lshift | rshift`)

<!-- vale on -->

Operator of the arithmetic operation.

The arithmetic operation can be addition, subtraction, multiplication, division,
XOR, right bit shift or left bit shift. In case of XOR and bit shifts, value of
signals is cast to integers before performing the operation.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([ArithmeticCombinatorOuts](#arithmetic-combinator-outs))

<!-- vale on -->

Output ports for the Arithmetic Combinator component.

</dd>
</dl>

---

<!-- vale off -->

### ArithmeticCombinatorIns {#arithmetic-combinator-ins}

<!-- vale on -->

Inputs for the Arithmetic Combinator component.

<dl>
<dt>lhs</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Left hand side of the arithmetic operation.

</dd>
<dt>rhs</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Right hand side of the arithmetic operation.

</dd>
</dl>

---

<!-- vale off -->

### ArithmeticCombinatorOuts {#arithmetic-combinator-outs}

<!-- vale on -->

Outputs for the Arithmetic Combinator component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Result of arithmetic operation.

</dd>
</dl>

---

<!-- vale off -->

### AutoScale {#auto-scale}

<!-- vale on -->

AutoScale components are used to scale a service.

<dl>
<dt>auto_scaler</dt>
<dd>

<!-- vale off -->

([AutoScaler](#auto-scaler))

<!-- vale on -->

_AutoScaler_ provides auto-scaling functionality for any scalable resource.

</dd>
<dt>pod_scaler</dt>
<dd>

<!-- vale off -->

([PodScaler](#pod-scaler))

<!-- vale on -->

PodScaler provides pod horizontal scaling functionality for scalable Kubernetes
resources.

</dd>
</dl>

---

<!-- vale off -->

### AutoScaler {#auto-scaler}

<!-- vale on -->

_AutoScaler_ provides auto-scaling functionality for any scalable resource.
Multiple _Controllers_ can be defined on the _AutoScaler_ for performing
scale-out or scale-in. The _AutoScaler_ can interface with infrastructure APIs
such as Kubernetes to perform auto-scale.

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

Dry run mode ensures that no scaling is invoked by this auto scaler. This is
useful for observing the behavior of auto scaler without disrupting any real
deployment. This parameter sets the default value of dry run setting which can
be overridden at runtime using dynamic configuration.

</dd>
<dt>dry_run_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for overriding dry run setting through dynamic configuration.

</dd>
<dt>scale_in_controllers</dt>
<dd>

<!-- vale off -->

([[]ScaleInController](#scale-in-controller))

<!-- vale on -->

List of _Controllers_ for scaling in.

</dd>
<dt>scale_out_controllers</dt>
<dd>

<!-- vale off -->

([[]ScaleOutController](#scale-out-controller))

<!-- vale on -->

List of _Controllers_ for scaling out.

</dd>
<dt>scaling_backend</dt>
<dd>

<!-- vale off -->

([AutoScalerScalingBackend](#auto-scaler-scaling-backend))

<!-- vale on -->

</dd>
<dt>scaling_parameters</dt>
<dd>

<!-- vale off -->

([AutoScalerScalingParameters](#auto-scaler-scaling-parameters))

<!-- vale on -->

Parameters that define the scaling behavior.

</dd>
</dl>

---

<!-- vale off -->

### AutoScalerScalingBackend {#auto-scaler-scaling-backend}

<!-- vale on -->

<dl>
<dt>kubernetes_replicas</dt>
<dd>

<!-- vale off -->

([AutoScalerScalingBackendKubernetesReplicas](#auto-scaler-scaling-backend-kubernetes-replicas))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### AutoScalerScalingBackendKubernetesReplicas {#auto-scaler-scaling-backend-kubernetes-replicas}

<!-- vale on -->

KubernetesReplicas defines a horizontal pod scaler for Kubernetes.

<dl>
<dt>kubernetes_object_selector</dt>
<dd>

<!-- vale off -->

([KubernetesObjectSelector](#kubernetes-object-selector))

<!-- vale on -->

The Kubernetes object on which horizontal scaling is applied.

</dd>
<dt>max_replicas</dt>
<dd>

<!-- vale off -->

(string, default: `"9223372036854775807"`)

<!-- vale on -->

The maximum replicas to which the _AutoScaler_ can scale-out.

</dd>
<dt>min_replicas</dt>
<dd>

<!-- vale off -->

(string, default: `"0"`)

<!-- vale on -->

The minimum replicas to which the _AutoScaler_ can scale-in.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([AutoScalerScalingBackendKubernetesReplicasOuts](#auto-scaler-scaling-backend-kubernetes-replicas-outs))

<!-- vale on -->

Output ports for _Kubernetes Replicas_.

</dd>
</dl>

---

<!-- vale off -->

### AutoScalerScalingBackendKubernetesReplicasOuts {#auto-scaler-scaling-backend-kubernetes-replicas-outs}

<!-- vale on -->

Outputs

<dl>
<dt>actual_replicas</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
<dt>configured_replicas</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
<dt>desired_replicas</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### AutoScalerScalingParameters {#auto-scaler-scaling-parameters}

<!-- vale on -->

<dl>
<dt>cooldown_override_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `50`)

<!-- vale on -->

Cooldown override percentage defines a threshold change in scale-out beyond
which previous cooldown is overridden. For example, if the cooldown is 5 minutes
and the cooldown override percentage is 10%, then if the scale-increases by 10%
or more, the previous cooldown is cancelled. Defaults to 50%.

</dd>
<dt>max_scale_in_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

The maximum decrease of scale (for example, pods) at one time. Defined as
percentage of current scale value. Can never go below one even if percentage
computation is less than one. Defaults to 1% of current scale value.

</dd>
<dt>max_scale_out_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `10`)

<!-- vale on -->

The maximum increase of scale (for example, pods) at one time. Defined as
percentage of current scale value. Can never go below one even if percentage
computation is less than one. Defaults to 10% of current scale value.

</dd>
<dt>scale_in_alerter</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for scale-in Alerter.

</dd>
<dt>scale_in_cooldown</dt>
<dd>

<!-- vale off -->

(string, default: `"120s"`)

<!-- vale on -->

The amount of time to wait after a scale-in operation for another scale-in
operation. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>scale_out_alerter</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for scale-out Alerter.

</dd>
<dt>scale_out_cooldown</dt>
<dd>

<!-- vale off -->

(string, default: `"30s"`)

<!-- vale on -->

The amount of time to wait after a scale-out operation for another scale-out or
scale-in operation. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### BoolVariable {#bool-variable}

<!-- vale on -->

Component that emits a constant Boolean signal which can be changed at runtime
through dynamic configuration.

<dl>
<dt>config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for overriding value setting through dynamic configuration.

</dd>
<dt>constant_output</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

The constant Boolean signal emitted by this component. The value of the constant
Boolean signal can be overridden at runtime through dynamic configuration.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([BoolVariableOuts](#bool-variable-outs))

<!-- vale on -->

Output ports for the BoolVariable component.

</dd>
</dl>

---

<!-- vale off -->

### BoolVariableOuts {#bool-variable-outs}

<!-- vale on -->

Outputs for the BoolVariable component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The value is emitted to the output port.

</dd>
</dl>

---

<!-- vale off -->

### Circuit {#circuit}

<!-- vale on -->

Circuit is graph of inter-connected signal processing components.

:::info

See also [Circuit overview](/concepts/advanced/circuit.md).

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

<dl>
<dt>components</dt>
<dd>

<!-- vale off -->

([[]Component](#component))

<!-- vale on -->

Defines a signal processing graph as a list of components.

</dd>
<dt>evaluation_interval</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Evaluation interval (tick) is the time between consecutive runs of the policy
circuit. This interval is typically aligned with how often the corrective action
(actuation) needs to be taken. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### Classifier {#classifier}

<!-- vale on -->

Set of classification rules sharing a common selector

:::info

See also [Classifier overview](/concepts/classifier.md).

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

<dl>
<dt>rego</dt>
<dd>

<!-- vale off -->

([Rego](#rego))

<!-- vale on -->

Rego based classification

Rego is a policy language used to express complex policies in a concise and
declarative way. It can be used to define flow classification rules by writing
custom queries that extract values from request metadata. For simple cases, such
as directly reading a value from header or a field from JSON body, declarative
extractors are recommended.

</dd>
<dt>rules</dt>
<dd>

<!-- vale off -->

(map of [Rule](#rule))

<!-- vale on -->

A map of {key, value} pairs mapping from [flow label](/concepts/flow-label.md)
keys to rules that define how to extract and propagate flow labels with that
key.

</dd>
<dt>selectors</dt>
<dd>

<!-- vale off -->

([[]Selector](#selector), **required**)

<!-- vale on -->

Selectors for flows that will be classified by this _Classifier_.

</dd>
</dl>

---

<!-- vale off -->

### Component {#component}

<!-- vale on -->

Computational block that forms the circuit

:::info

See also [Components overview](/concepts/advanced/circuit.md#components).

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

<dl>
<dt>alerter</dt>
<dd>

<!-- vale off -->

([Alerter](#alerter))

<!-- vale on -->

Alerter reacts to a signal and generates alert to send to alert manager.

</dd>
<dt>and</dt>
<dd>

<!-- vale off -->

([And](#and))

<!-- vale on -->

Logical AND.

</dd>
<dt>arithmetic_combinator</dt>
<dd>

<!-- vale off -->

([ArithmeticCombinator](#arithmetic-combinator))

<!-- vale on -->

Applies the given operator on input operands (signals) and emits the result.

</dd>
<dt>auto_scale</dt>
<dd>

<!-- vale off -->

([AutoScale](#auto-scale))

<!-- vale on -->

AutoScale components are used to scale the service.

</dd>
<dt>bool_variable</dt>
<dd>

<!-- vale off -->

([BoolVariable](#bool-variable))

<!-- vale on -->

BoolVariable emits a constant Boolean signal which can be changed at runtime
through dynamic configuration.

</dd>
<dt>decider</dt>
<dd>

<!-- vale off -->

([Decider](#decider))

<!-- vale on -->

Decider emits the binary result of comparison operator on two operands.

</dd>
<dt>differentiator</dt>
<dd>

<!-- vale off -->

([Differentiator](#differentiator))

<!-- vale on -->

Differentiator calculates rate of change per tick.

</dd>
<dt>ema</dt>
<dd>

<!-- vale off -->

([EMA](#e-m-a))

<!-- vale on -->

Exponential Moving Average filter.

</dd>
<dt>extrapolator</dt>
<dd>

<!-- vale off -->

([Extrapolator](#extrapolator))

<!-- vale on -->

Takes an input signal and emits the extrapolated value; either mirroring the
input value or repeating the last known value up to the maximum extrapolation
interval.

</dd>
<dt>first_valid</dt>
<dd>

<!-- vale off -->

([FirstValid](#first-valid))

<!-- vale on -->

Picks the first valid input signal and emits it.

</dd>
<dt>flow_control</dt>
<dd>

<!-- vale off -->

([FlowControl](#flow-control))

<!-- vale on -->

FlowControl components are used to regulate requests flow.

</dd>
<dt>gradient_controller</dt>
<dd>

<!-- vale off -->

([GradientController](#gradient-controller))

<!-- vale on -->

Gradient controller calculates the ratio between the signal and the setpoint to
determine the magnitude of the correction that need to be applied. This
controller can be used to build AIMD (Additive Increase, Multiplicative
Decrease) or MIMD style response.

</dd>
<dt>holder</dt>
<dd>

<!-- vale off -->

([Holder](#holder))

<!-- vale on -->

Holds the last valid signal value for the specified duration then waits for next
valid value to hold.

</dd>
<dt>integrator</dt>
<dd>

<!-- vale off -->

([Integrator](#integrator))

<!-- vale on -->

Accumulates sum of signal every tick.

</dd>
<dt>inverter</dt>
<dd>

<!-- vale off -->

([Inverter](#inverter))

<!-- vale on -->

Logical NOT.

</dd>
<dt>max</dt>
<dd>

<!-- vale off -->

([Max](#max))

<!-- vale on -->

Emits the maximum of the input signals.

</dd>
<dt>min</dt>
<dd>

<!-- vale off -->

([Min](#min))

<!-- vale on -->

Emits the minimum of the input signals.

</dd>
<dt>nested_circuit</dt>
<dd>

<!-- vale off -->

([NestedCircuit](#nested-circuit))

<!-- vale on -->

Nested circuit defines a sub-circuit as a high-level component. It consists of a
list of components and a map of input and output ports.

</dd>
<dt>nested_signal_egress</dt>
<dd>

<!-- vale off -->

([NestedSignalEgress](#nested-signal-egress))

<!-- vale on -->

Nested signal egress is a special type of component that allows to extract a
signal from a nested circuit.

</dd>
<dt>nested_signal_ingress</dt>
<dd>

<!-- vale off -->

([NestedSignalIngress](#nested-signal-ingress))

<!-- vale on -->

Nested signal ingress is a special type of component that allows to inject a
signal into a nested circuit.

</dd>
<dt>or</dt>
<dd>

<!-- vale off -->

([Or](#or))

<!-- vale on -->

Logical OR.

</dd>
<dt>pulse_generator</dt>
<dd>

<!-- vale off -->

([PulseGenerator](#pulse-generator))

<!-- vale on -->

Generates 0 and 1 in turns.

</dd>
<dt>query</dt>
<dd>

<!-- vale off -->

([Query](#query))

<!-- vale on -->

Query components that are query databases such as Prometheus.

</dd>
<dt>signal_generator</dt>
<dd>

<!-- vale off -->

([SignalGenerator](#signal-generator))

<!-- vale on -->

Generates the specified signal.

</dd>
<dt>sma</dt>
<dd>

<!-- vale off -->

([SMA](#s-m-a))

<!-- vale on -->

Simple Moving Average filter.

</dd>
<dt>switcher</dt>
<dd>

<!-- vale off -->

([Switcher](#switcher))

<!-- vale on -->

Switcher acts as a switch that emits one of the two signals based on third
signal.

</dd>
<dt>unary_operator</dt>
<dd>

<!-- vale off -->

([UnaryOperator](#unary-operator))

<!-- vale on -->

Takes an input signal and emits the square root of the input signal.

</dd>
<dt>variable</dt>
<dd>

<!-- vale off -->

([Variable](#variable))

<!-- vale on -->

Emits a variable signal which can be changed at runtime through dynamic
configuration.

</dd>
</dl>

---

<!-- vale off -->

### ConstantSignal {#constant-signal}

<!-- vale on -->

Special constant input for ports and Variable component. Can provide either a
constant value or special Nan/+-Inf value.

<dl>
<dt>special_value</dt>
<dd>

<!-- vale off -->

(string, one of: `NaN | +Inf | -Inf`)

<!-- vale on -->

A special value such as NaN, +Inf, -Inf.

</dd>
<dt>value</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

A constant value.

</dd>
</dl>

---

<!-- vale off -->

### Decider {#decider}

<!-- vale on -->

Type of Combinator that computes the comparison operation on LHS and RHS signals

The comparison operator can be greater-than, less-than, greater-than-or-equal,
less-than-or-equal, equal, or not-equal.

This component also supports time-based response (the output) transitions
between 1.0 or 0.0 signal if the decider condition is true or false for at least
`true_for` or `false_for` duration. If `true_for` and `false_for` durations are
zero then the transitions are instantaneous.

<dl>
<dt>false_for</dt>
<dd>

<!-- vale off -->

(string, default: `"0s"`)

<!-- vale on -->

Duration of time to wait before changing to false state. If the duration is
zero, the change will happen instantaneously. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([DeciderIns](#decider-ins))

<!-- vale on -->

Input ports for the Decider component.

</dd>
<dt>operator</dt>
<dd>

<!-- vale off -->

(string, one of: `gt | lt | gte | lte | eq | neq`)

<!-- vale on -->

Comparison operator that computes operation on LHS and RHS input signals.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([DeciderOuts](#decider-outs))

<!-- vale on -->

Output ports for the Decider component.

</dd>
<dt>true_for</dt>
<dd>

<!-- vale off -->

(string, default: `"0s"`)

<!-- vale on -->

Duration of time to wait before changing to true state. If the duration is zero,
the change will happen instantaneously.``` This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### DeciderIns {#decider-ins}

<!-- vale on -->

Inputs for the Decider component.

<dl>
<dt>lhs</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Left hand side input signal for the comparison operation.

</dd>
<dt>rhs</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Right hand side input signal for the comparison operation.

</dd>
</dl>

---

<!-- vale off -->

### DeciderOuts {#decider-outs}

<!-- vale on -->

Outputs for the Decider component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Selected signal (1.0 or 0.0).

</dd>
</dl>

---

<!-- vale off -->

### DecreasingGradient {#decreasing-gradient}

<!-- vale on -->

Decreasing Gradient defines a controller for scaling in based on Gradient
Controller.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([DecreasingGradientIns](#decreasing-gradient-ins))

<!-- vale on -->

Input ports for the Gradient.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([DecreasingGradientParameters](#decreasing-gradient-parameters))

<!-- vale on -->

Gradient parameters for the controller. Defaults and constraints:

- `slope` = 1
- `min_gradient` = -Inf (must be less than 1)
- `max_gradient` = 1 (cannot be changed)

</dd>
</dl>

---

<!-- vale off -->

### DecreasingGradientIns {#decreasing-gradient-ins}

<!-- vale on -->

Inputs for Gradient.

<dl>
<dt>setpoint</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The setpoint to use for scale-in.

</dd>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The signal to use for scale-in.

</dd>
</dl>

---

<!-- vale off -->

### DecreasingGradientParameters {#decreasing-gradient-parameters}

<!-- vale on -->

This allows subset of parameters with constrained values compared to a regular
gradient controller. For full documentation of these parameters, refer to the
[GradientControllerParameters](#gradient-controller-parameters).

<dl>
<dt>min_gradient</dt>
<dd>

<!-- vale off -->

(float64, default: `-1.7976931348623157e+308`)

<!-- vale on -->

</dd>
<dt>slope</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### Differentiator {#differentiator}

<!-- vale on -->

Differentiator calculates rate of change per tick.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([DifferentiatorIns](#differentiator-ins))

<!-- vale on -->

Input ports for the Differentiator component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([DifferentiatorOuts](#differentiator-outs))

<!-- vale on -->

Output ports for the Differentiator component.

</dd>
<dt>window</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

The window of time over which differentiator operates. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### DifferentiatorIns {#differentiator-ins}

<!-- vale on -->

Inputs for the Differentiator component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### DifferentiatorOuts {#differentiator-outs}

<!-- vale on -->

Outputs for the Differentiator component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### EMA {#e-m-a}

<!-- vale on -->

Exponential Moving Average (EMA) is a type of moving average that applies
exponentially more weight to recent signal readings

At any time EMA component operates in one of the following states:

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

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([EMAIns](#e-m-a-ins))

<!-- vale on -->

Input ports for the EMA component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([EMAOuts](#e-m-a-outs))

<!-- vale on -->

Output ports for the EMA component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([EMAParameters](#e-m-a-parameters))

<!-- vale on -->

Parameters for the EMA component.

</dd>
</dl>

---

<!-- vale off -->

### EMAIns {#e-m-a-ins}

<!-- vale on -->

Inputs for the EMA component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Input signal to be used for the EMA computation.

</dd>
<dt>max_envelope</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Upper bound of the moving average.

When the signal exceeds `max_envelope` it is multiplied by
`correction_factor_on_max_envelope_violation` **once per tick**.

:::note

If the signal deviates from `max_envelope` faster than the correction faster, it
might end up exceeding the envelope.

:::

</dd>
<dt>min_envelope</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Lower bound of the moving average.

Behavior is similar to `max_envelope`.

</dd>
</dl>

---

<!-- vale off -->

### EMAOuts {#e-m-a-outs}

<!-- vale on -->

Outputs for the EMA component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Exponential moving average of the series of reading as an output signal.

</dd>
</dl>

---

<!-- vale off -->

### EMAParameters {#e-m-a-parameters}

<!-- vale on -->

Parameters for the EMA component.

<dl>
<dt>correction_factor_on_max_envelope_violation</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, default: `1`)

<!-- vale on -->

Correction factor to apply on the output value if its in violation of the max
envelope.

</dd>
<dt>correction_factor_on_min_envelope_violation</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

Correction factor to apply on the output value if its in violation of the min
envelope.

</dd>
<dt>ema_window</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Duration of EMA sampling window. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>valid_during_warmup</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Whether the output is valid during the warm-up stage.

</dd>
<dt>warmup_window</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Duration of EMA warming up window.

The initial value of the EMA is the average of signal readings received during
the warm up window. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### EqualsMatchExpression {#equals-match-expression}

<!-- vale on -->

Label selector expression of the equal form `label == value`.

<dl>
<dt>label</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Name of the label to equal match the value.

</dd>
<dt>value</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Exact value that the label should be equal to.

</dd>
</dl>

---

<!-- vale off -->

### Extractor {#extractor}

<!-- vale on -->

Defines a high-level way to specify how to extract a flow label value given HTTP
request metadata, without a need to write Rego code

There are multiple variants of extractor, specify exactly one.

<dl>
<dt>address</dt>
<dd>

<!-- vale off -->

([AddressExtractor](#address-extractor))

<!-- vale on -->

Display an address as a single string - `<ip>:<port>`.

</dd>
<dt>from</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Use an attribute with no conversion

Attribute path is a dot-separated path to attribute.

Should be either:

- one of the fields of
  [Attribute Context](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto),
  or
- a special `request.http.bearer` pseudo-attribute. For example,
  `request.http.method` or `request.http.header.user-agent`

Note: The same attribute path syntax is shared by other extractor variants,
wherever attribute path is needed in their "from" syntax.

Example:

```yaml
from: request.http.headers.user-agent
```

</dd>
<dt>json</dt>
<dd>

<!-- vale off -->

([JSONExtractor](#json-extractor))

<!-- vale on -->

Parse JSON, and extract one of the fields.

</dd>
<dt>jwt</dt>
<dd>

<!-- vale off -->

([JWTExtractor](#j-w-t-extractor))

<!-- vale on -->

Parse the attribute as JWT and read the payload.

</dd>
<dt>path_templates</dt>
<dd>

<!-- vale off -->

([PathTemplateMatcher](#path-template-matcher))

<!-- vale on -->

Match HTTP Path to given path templates.

</dd>
</dl>

---

<!-- vale off -->

### Extrapolator {#extrapolator}

<!-- vale on -->

Extrapolates the input signal by repeating the last valid value during the
period in which it is invalid

It does so until `maximum_extrapolation_interval` is reached, beyond which it
emits invalid signal unless input signal becomes valid again.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([ExtrapolatorIns](#extrapolator-ins))

<!-- vale on -->

Input ports for the Extrapolator component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([ExtrapolatorOuts](#extrapolator-outs))

<!-- vale on -->

Output ports for the Extrapolator component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([ExtrapolatorParameters](#extrapolator-parameters))

<!-- vale on -->

Parameters for the Extrapolator component.

</dd>
</dl>

---

<!-- vale off -->

### ExtrapolatorIns {#extrapolator-ins}

<!-- vale on -->

Inputs for the Extrapolator component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Input signal for the Extrapolator component.

</dd>
</dl>

---

<!-- vale off -->

### ExtrapolatorOuts {#extrapolator-outs}

<!-- vale on -->

Outputs for the Extrapolator component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Extrapolated signal.

</dd>
</dl>

---

<!-- vale off -->

### ExtrapolatorParameters {#extrapolator-parameters}

<!-- vale on -->

Parameters for the Extrapolator component.

<dl>
<dt>max_extrapolation_interval</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Maximum time interval to repeat the last valid value of input signal. This field
employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### FirstValid {#first-valid}

<!-- vale on -->

Picks the first valid input signal from the array of input signals and emits it
as an output signal

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([FirstValidIns](#first-valid-ins))

<!-- vale on -->

Input ports for the FirstValid component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([FirstValidOuts](#first-valid-outs))

<!-- vale on -->

Output ports for the FirstValid component.

</dd>
</dl>

---

<!-- vale off -->

### FirstValidIns {#first-valid-ins}

<!-- vale on -->

Inputs for the FirstValid component.

<dl>
<dt>inputs</dt>
<dd>

<!-- vale off -->

([[]InPort](#in-port))

<!-- vale on -->

Array of input signals.

</dd>
</dl>

---

<!-- vale off -->

### FirstValidOuts {#first-valid-outs}

<!-- vale on -->

Outputs for the FirstValid component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

First valid input signal as an output signal.

</dd>
</dl>

---

<!-- vale off -->

### FlowControl {#flow-control}

<!-- vale on -->

_Flow Control_ encompasses components that manage the flow of requests or access
to features within a service.

<dl>
<dt>adaptive_load_scheduler</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadScheduler](#adaptive-load-scheduler))

<!-- vale on -->

_Adaptive Load Scheduler_ component is based on additive increase and
multiplicative decrease of token rate. It takes a signal and setpoint as inputs
and reduces token rate proportionally (or any arbitrary power) based on
deviation of the signal from setpoint.

</dd>
<dt>load_ramp</dt>
<dd>

<!-- vale off -->

([LoadRamp](#load-ramp))

<!-- vale on -->

_Load Ramp_ smoothly regulates the flow of requests over specified steps.

</dd>
<dt>load_scheduler</dt>
<dd>

<!-- vale off -->

([LoadScheduler](#load-scheduler))

<!-- vale on -->

_Load Scheduler_ provides service protection by creating a prioritized workload
queue in front of the service using Weighted Fair Queuing.

</dd>
<dt>quota_scheduler</dt>
<dd>

<!-- vale off -->

([QuotaScheduler](#quota-scheduler))

<!-- vale on -->

</dd>
<dt>rate_limiter</dt>
<dd>

<!-- vale off -->

([RateLimiter](#rate-limiter))

<!-- vale on -->

_Rate Limiter_ provides service protection by applying rate limits using the
token bucket algorithm.

</dd>
<dt>sampler</dt>
<dd>

<!-- vale off -->

([Sampler](#sampler))

<!-- vale on -->

Sampler is a component that regulates the flow of requests to the service by
allowing only the specified percentage of requests or sticky sessions.

</dd>
</dl>

---

<!-- vale off -->

### FlowControlResources {#flow-control-resources}

<!-- vale on -->

FlowControl Resources

<dl>
<dt>classifiers</dt>
<dd>

<!-- vale off -->

([[]Classifier](#classifier))

<!-- vale on -->

Classifiers are installed in the data-plane and are used to label the requests
based on payload content.

The flow labels created by Classifiers can be matched by Flux Meters to create
metrics for control purposes.

</dd>
<dt>flux_meters</dt>
<dd>

<!-- vale off -->

(map of [FluxMeter](#flux-meter))

<!-- vale on -->

Flux Meters are installed in the data-plane and form the observability leg of
the feedback loop.

Flux Meter created metrics can be consumed as input to the circuit through the
PromQL component.

</dd>
</dl>

---

<!-- vale off -->

### FluxMeter {#flux-meter}

<!-- vale on -->

Flux Meter gathers metrics for the traffic that matches its selector. The
histogram created by Flux Meter measures the workload latency by default.

:::info

See also [Flux Meter overview](/concepts/flux-meter.md).

::: Example:

```yaml
static_buckets:
  buckets:
    [
      5.0,
      10.0,
      25.0,
      50.0,
      100.0,
      250.0,
      500.0,
      1000.0,
      2500.0,
      5000.0,
      10000.0,
    ]
selectors:
  - agent_group: demoapp
    service: service1-demo-app.demoapp.svc.cluster.local
    control_point: ingress
attribute_key: response_duration_ms
```

<dl>
<dt>attribute_key</dt>
<dd>

<!-- vale off -->

(string, default: `"workload_duration_ms"`)

<!-- vale on -->

Key of the attribute in access log or span from which the metric for this flux
meter is read.

:::info

For list of available attributes in Envoy access logs, refer
[Envoy Filter](/integrations/istio/istio.md#envoy-filter)

:::

</dd>
<dt>exponential_buckets</dt>
<dd>

<!-- vale off -->

([FluxMeterExponentialBuckets](#flux-meter-exponential-buckets))

<!-- vale on -->

</dd>
<dt>exponential_buckets_range</dt>
<dd>

<!-- vale off -->

([FluxMeterExponentialBucketsRange](#flux-meter-exponential-buckets-range))

<!-- vale on -->

</dd>
<dt>linear_buckets</dt>
<dd>

<!-- vale off -->

([FluxMeterLinearBuckets](#flux-meter-linear-buckets))

<!-- vale on -->

</dd>
<dt>selectors</dt>
<dd>

<!-- vale off -->

([[]Selector](#selector), **required**)

<!-- vale on -->

Selectors for flows that will be metered by this _Flux Meter_.

</dd>
<dt>static_buckets</dt>
<dd>

<!-- vale off -->

([FluxMeterStaticBuckets](#flux-meter-static-buckets))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### FluxMeterExponentialBuckets {#flux-meter-exponential-buckets}

<!-- vale on -->

ExponentialBuckets creates `count` number of buckets where the lowest bucket has
an upper bound of `start` and each following bucket's upper bound is `factor`
times the previous bucket's upper bound. The final +inf bucket is not counted.

<dl>
<dt>count</dt>
<dd>

<!-- vale off -->

(int32, minimum: `0`)

<!-- vale on -->

Number of buckets.

</dd>
<dt>factor</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

Factor to be multiplied to the previous bucket's upper bound to calculate the
following bucket's upper bound.

</dd>
<dt>start</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

Upper bound of the lowest bucket.

</dd>
</dl>

---

<!-- vale off -->

### FluxMeterExponentialBucketsRange {#flux-meter-exponential-buckets-range}

<!-- vale on -->

ExponentialBucketsRange creates `count` number of buckets where the lowest
bucket is `min` and the highest bucket is `max`. The final +inf bucket is not
counted.

<dl>
<dt>count</dt>
<dd>

<!-- vale off -->

(int32, minimum: `0`)

<!-- vale on -->

Number of buckets.

</dd>
<dt>max</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

Highest bucket.

</dd>
<dt>min</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

Lowest bucket.

</dd>
</dl>

---

<!-- vale off -->

### FluxMeterLinearBuckets {#flux-meter-linear-buckets}

<!-- vale on -->

LinearBuckets creates `count` number of buckets, each `width` wide, where the
lowest bucket has an upper bound of `start`. The final +inf bucket is not
counted.

<dl>
<dt>count</dt>
<dd>

<!-- vale off -->

(int32, minimum: `0`)

<!-- vale on -->

Number of buckets.

</dd>
<dt>start</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

Upper bound of the lowest bucket.

</dd>
<dt>width</dt>
<dd>

<!-- vale off -->

(float64)

<!-- vale on -->

Width of each bucket.

</dd>
</dl>

---

<!-- vale off -->

### FluxMeterStaticBuckets {#flux-meter-static-buckets}

<!-- vale on -->

StaticBuckets holds the static value of the buckets where latency histogram will
be stored.

<dl>
<dt>buckets</dt>
<dd>

<!-- vale off -->

([]float64, default: `[5,10,25,50,100,250,500,1000,2500,5000,10000]`)

<!-- vale on -->

The buckets in which latency histogram will be stored.

</dd>
</dl>

---

<!-- vale off -->

### GradientController {#gradient-controller}

<!-- vale on -->

Gradient controller is a type of controller which tries to adjust the control
variable proportionally to the relative difference between setpoint and actual
value of the signal

The `gradient` describes a corrective factor that should be applied to the
control variable to get the signal closer to the setpoint. It's computed as
follows:

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

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([GradientControllerIns](#gradient-controller-ins))

<!-- vale on -->

Input ports of the Gradient Controller.

</dd>
<dt>manual_mode</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

In manual mode, the controller does not adjust the control variable. It emits
the same output as the control variable input. This setting can be adjusted at
runtime through dynamic configuration without restarting the policy.

</dd>
<dt>manual_mode_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for overriding `manual_mode` setting through dynamic
configuration.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([GradientControllerOuts](#gradient-controller-outs))

<!-- vale on -->

Output ports of the Gradient Controller.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([GradientControllerParameters](#gradient-controller-parameters))

<!-- vale on -->

Gradient Parameters.

</dd>
</dl>

---

<!-- vale off -->

### GradientControllerIns {#gradient-controller-ins}

<!-- vale on -->

Inputs for the Gradient Controller component.

<dl>
<dt>control_variable</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Actual current value of the control variable.

This signal is multiplied by the gradient to produce the output.

</dd>
<dt>max</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Maximum value to limit the output signal.

</dd>
<dt>min</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Minimum value to limit the output signal.

</dd>
<dt>setpoint</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Setpoint to be used for the gradient computation.

</dd>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Signal to be used for the gradient computation.

</dd>
</dl>

---

<!-- vale off -->

### GradientControllerOuts {#gradient-controller-outs}

<!-- vale on -->

Outputs for the Gradient Controller component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Computed desired value of the control variable.

</dd>
</dl>

---

<!-- vale off -->

### GradientControllerParameters {#gradient-controller-parameters}

<!-- vale on -->

Gradient Parameters.

<dl>
<dt>max_gradient</dt>
<dd>

<!-- vale off -->

(float64, default: `1.7976931348623157e+308`)

<!-- vale on -->

Maximum gradient which clamps the computed gradient value to the range,
`[min_gradient, max_gradient]`.

</dd>
<dt>min_gradient</dt>
<dd>

<!-- vale off -->

(float64, default: `-1.7976931348623157e+308`)

<!-- vale on -->

Minimum gradient which clamps the computed gradient value to the range,
`[min_gradient, max_gradient]`.

</dd>
<dt>slope</dt>
<dd>

<!-- vale off -->

(float64, **required**)

<!-- vale on -->

Slope controls the aggressiveness and direction of the Gradient Controller.

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

:::

</dd>
</dl>

---

<!-- vale off -->

### Holder {#holder}

<!-- vale on -->

Holds the last valid signal value for the specified duration then waits for next
valid value to hold. If it is holding a value that means it ignores both valid
and invalid new signals until the `hold_for` duration is finished.

<dl>
<dt>hold_for</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Holding the last valid signal value for the `hold_for` duration. This field
employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([HolderIns](#holder-ins))

<!-- vale on -->

Input ports for the Holder component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([HolderOuts](#holder-outs))

<!-- vale on -->

Output ports for the Holder component.

</dd>
</dl>

---

<!-- vale off -->

### HolderIns {#holder-ins}

<!-- vale on -->

Inputs for the Holder component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The input signal.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Resets the holder output to the current input signal when reset signal is valid
and non-zero.

</dd>
</dl>

---

<!-- vale off -->

### HolderOuts {#holder-outs}

<!-- vale on -->

Outputs for the Holder component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The output signal.

</dd>
</dl>

---

<!-- vale off -->

### InPort {#in-port}

<!-- vale on -->

Components receive input from other components through InPorts

<dl>
<dt>constant_signal</dt>
<dd>

<!-- vale off -->

([ConstantSignal](#constant-signal))

<!-- vale on -->

Constant value to be used for this InPort instead of a signal.

</dd>
<dt>signal_name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the incoming Signal on the InPort.

</dd>
</dl>

---

<!-- vale off -->

### IncreasingGradient {#increasing-gradient}

<!-- vale on -->

Increasing Gradient defines a controller for scaling out based on _Gradient
Controller_.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([IncreasingGradientIns](#increasing-gradient-ins))

<!-- vale on -->

Input ports for the Gradient.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([IncreasingGradientParameters](#increasing-gradient-parameters))

<!-- vale on -->

Gradient parameters for the controller. Defaults and constraints:

- `slope` = 1
- `min_gradient` = 1 (cannot be changed)
- `max_gradient` = +Inf (must be greater than 1)

</dd>
</dl>

---

<!-- vale off -->

### IncreasingGradientIns {#increasing-gradient-ins}

<!-- vale on -->

Inputs for Gradient.

<dl>
<dt>setpoint</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The setpoint to use for scale-out.

</dd>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The signal to use for scale-out.

</dd>
</dl>

---

<!-- vale off -->

### IncreasingGradientParameters {#increasing-gradient-parameters}

<!-- vale on -->

This allows subset of parameters with constrained values compared to a regular
gradient controller. For full documentation of these parameters, refer to the
[GradientControllerParameters](#gradient-controller-parameters).

<dl>
<dt>max_gradient</dt>
<dd>

<!-- vale off -->

(float64, default: `1.7976931348623157e+308`)

<!-- vale on -->

</dd>
<dt>slope</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### InfraMeter {#infra-meter}

<!-- vale on -->

InfraMeter is a resource that sets up OpenTelemetry pipelines. It defines
receivers, processors, and a single metrics pipeline which will be exported to
the configured Prometheus instance. Environment variables can be used in the
configuration using format `${ENV_VAR_NAME}`.

:::info

See also
[Get Started / Setup Integrations / Metrics](/integrations/metrics/metrics.md).

:::

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

AgentGroup is the agent group to sync this InfraMeter with.

</dd>
<dt>per_agent_group</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

PerAgentGroup marks the pipeline to be instantiated only once per agent group.
This is helpful for receivers that scrape for example, some cluster-wide
metrics. When not set, pipeline will be instantiated on every Agent.

</dd>
<dt>pipeline</dt>
<dd>

<!-- vale off -->

([InfraMeterMetricsPipeline](#infra-meter-metrics-pipeline))

<!-- vale on -->

Pipeline is an OTel metrics pipeline definition, which **only** uses receivers
and processors defined above. Exporter would be added automatically.

If there are no processors defined or only one processor is defined, the
pipeline definition can be omitted. In such cases, the pipeline will
automatically use all given receivers and the defined processor (if any).
However, if there are more than one processor, the pipeline must be defined
explicitly.

</dd>
<dt>processors</dt>
<dd>

<!-- vale off -->

(map of any )

<!-- vale on -->

Processors define processors to be used in custom metrics pipelines. This should
be in
[OTel format](https://opentelemetry.io/docs/collector/configuration/#processors).

</dd>
<dt>receivers</dt>
<dd>

<!-- vale off -->

(map of any )

<!-- vale on -->

Receivers define receivers to be used in custom metrics pipelines. This should
be in
[OTel format](https://opentelemetry.io/docs/collector/configuration/#receivers).

</dd>
</dl>

---

<!-- vale off -->

### InfraMeterMetricsPipeline {#infra-meter-metrics-pipeline}

<!-- vale on -->

MetricsPipelineConfig defines a custom metrics pipeline.

<dl>
<dt>processors</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

</dd>
<dt>receivers</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### Integrator {#integrator}

<!-- vale on -->

Accumulates sum of signal every tick.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([IntegratorIns](#integrator-ins))

<!-- vale on -->

Input ports for the Integrator component.

</dd>
<dt>initial_value</dt>
<dd>

<!-- vale off -->

(float64, default: `0`)

<!-- vale on -->

Initial value of the integrator.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([IntegratorOuts](#integrator-outs))

<!-- vale on -->

Output ports for the Integrator component.

</dd>
</dl>

---

<!-- vale off -->

### IntegratorIns {#integrator-ins}

<!-- vale on -->

Inputs for the Integrator component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The input signal.

</dd>
<dt>max</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The maximum output.

</dd>
<dt>min</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The minimum output.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Resets the integrator output to zero when reset signal is valid and non-zero.
Reset also resets the max and min constraints.

</dd>
</dl>

---

<!-- vale off -->

### IntegratorOuts {#integrator-outs}

<!-- vale on -->

Outputs for the Integrator component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### Inverter {#inverter}

<!-- vale on -->

Logical NOT.

See [And component](#and) on how signals are mapped onto Boolean values.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([InverterIns](#inverter-ins))

<!-- vale on -->

Input ports for the Inverter component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([InverterOuts](#inverter-outs))

<!-- vale on -->

Output ports for the Inverter component.

</dd>
</dl>

---

<!-- vale off -->

### InverterIns {#inverter-ins}

<!-- vale on -->

Inputs for the Inverter component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Signal to be negated.

</dd>
</dl>

---

<!-- vale off -->

### InverterOuts {#inverter-outs}

<!-- vale on -->

Output ports for the Inverter component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Logical negation of the input signal.

Will always be 0 (false), 1 (true) or invalid (unknown).

</dd>
</dl>

---

<!-- vale off -->

### JSONExtractor {#json-extractor}

<!-- vale on -->

Parse JSON, and extract one of the fields

Example:

```yaml
from: request.http.body
pointer: /user/name
```

<dl>
<dt>from</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Attribute path pointing to some strings - for example, `request.http.body`.

</dd>
<dt>pointer</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

JSON pointer represents a parsed JSON pointer which allows to select a specified
field from the payload.

Note: Uses [JSON pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
for example, `/foo/bar`. If the pointer points into an object, it'd be converted
to a string.

</dd>
</dl>

---

<!-- vale off -->

### JWTExtractor {#j-w-t-extractor}

<!-- vale on -->

Parse the attribute as JWT and read the payload

Specify a field to be extracted from payload using `json_pointer`.

Note: The signature is not verified against the secret (assuming there's some
other part of the system that handles such verification).

Example:

```yaml
from: request.http.bearer
json_pointer: /user/email
```

<dl>
<dt>from</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

JWT (JSON Web Token) can be extracted from any input attribute, but most likely
you'd want to use `request.http.bearer`.

</dd>
<dt>json_pointer</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

JSON pointer allowing to select a specified field from the payload.

Note: Uses [JSON pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
for example, `/foo/bar`. If the pointer points into an object, it'd be converted
to a string.

</dd>
</dl>

---

<!-- vale off -->

### K8sLabelMatcherRequirement {#k8s-label-matcher-requirement}

<!-- vale on -->

Label selector requirement which is a selector that contains values, a key, and
an operator that relates the key and values.

<dl>
<dt>key</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Label key that the selector applies to.

</dd>
<dt>operator</dt>
<dd>

<!-- vale off -->

(string, one of: `In | NotIn | Exists | DoesNotExists`)

<!-- vale on -->

Logical operator which represents a key's relationship to a set of values. Valid
operators are In, NotIn, Exists and DoesNotExist.

</dd>
<dt>values</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

An array of string values that relates to the key by an operator. If the
operator is In or NotIn, the values array must be non-empty. If the operator is
Exists or DoesNotExist, the values array must be empty.

</dd>
</dl>

---

<!-- vale off -->

### KubernetesObjectSelector {#kubernetes-object-selector}

<!-- vale on -->

Describes which pods a control or observability component should apply to.

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

Which [agent-group](/concepts/selector.md#agent-group) this selector applies to.

</dd>
<dt>api_version</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

API version of Kubernetes resource

</dd>
<dt>kind</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Kubernetes resource type.

</dd>
<dt>name</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Kubernetes resource name.

</dd>
<dt>namespace</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Kubernetes namespace that the resource belongs to.

</dd>
</dl>

---

<!-- vale off -->

### LabelMatcher {#label-matcher}

<!-- vale on -->

Allows to define rules whether a map of [labels](/concepts/flow-label.md) should
be considered a match or not

It provides three ways to define requirements:

- match labels
- match expressions
- arbitrary expression

If multiple requirements are set, they're all combined using the logical AND
operator. An empty label matcher always matches.

<dl>
<dt>expression</dt>
<dd>

<!-- vale off -->

([MatchExpression](#match-expression))

<!-- vale on -->

An arbitrary expression to be evaluated on the labels.

</dd>
<dt>match_expressions</dt>
<dd>

<!-- vale off -->

([[]K8sLabelMatcherRequirement](#k8s-label-matcher-requirement))

<!-- vale on -->

List of Kubernetes-style label matcher requirements.

Note: The requirements are combined using the logical AND operator.

</dd>
<dt>match_labels</dt>
<dd>

<!-- vale off -->

(map of string)

<!-- vale on -->

A map of {key,value} pairs representing labels to be matched. A single
{key,value} in the `match_labels` requires that the label `key` is present and
equal to `value`.

Note: The requirements are combined using the logical AND operator.

</dd>
</dl>

---

<!-- vale off -->

### LoadRamp {#load-ramp}

<!-- vale on -->

The _Load Ramp_ produces a smooth and continuous traffic load that changes
progressively over time, based on the specified steps.

Each step is defined by two parameters:

- The `target_accept_percentage`.
- The `duration` for the signal to change from the previous step's
  `target_accept_percentage` to the current step's `target_accept_percentage`.

The percentage of requests accepted starts at the `target_accept_percentage`
defined in the first step and gradually ramps up or down linearly from the
previous step's `target_accept_percentage` to the next
`target_accept_percentage`, over the `duration` specified for each step.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadRampIns](#load-ramp-ins))

<!-- vale on -->

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([LoadRampOuts](#load-ramp-outs))

<!-- vale on -->

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadRampParameters](#load-ramp-parameters))

<!-- vale on -->

</dd>
<dt>pass_through_label_values</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

Specify certain label values to be always accepted by the _Sampler_ regardless
of accept percentage.

</dd>
<dt>pass_through_label_values_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for setting pass through label values through dynamic
configuration.

</dd>
</dl>

---

<!-- vale off -->

### LoadRampIns {#load-ramp-ins}

<!-- vale on -->

Inputs for the _Load Ramp_ component.

<dl>
<dt>backward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the _Load Ramp_ towards the previous step.

</dd>
<dt>forward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the _Load Ramp_ towards the next step.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to reset the _Load Ramp_ to the first step.

</dd>
</dl>

---

<!-- vale off -->

### LoadRampOuts {#load-ramp-outs}

<!-- vale on -->

Outputs for the _Load Ramp_ component.

<dl>
<dt>accept_percentage</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The percentage of flows being accepted by the _Load Ramp_.

</dd>
<dt>at_end</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Load Ramp_ is at the end of signal
generation.

</dd>
<dt>at_start</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Load Ramp_ is at the start of signal
generation.

</dd>
</dl>

---

<!-- vale off -->

### LoadRampParameters {#load-ramp-parameters}

<!-- vale on -->

Parameters for the _Load Ramp_ component.

<dl>
<dt>sampler</dt>
<dd>

<!-- vale off -->

([SamplerParameters](#sampler-parameters))

<!-- vale on -->

Parameters for the _Sampler_.

</dd>
<dt>steps</dt>
<dd>

<!-- vale off -->

([[]LoadRampParametersStep](#load-ramp-parameters-step), **required**)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadRampParametersStep {#load-ramp-parameters-step}

<!-- vale on -->

<dl>
<dt>duration</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Duration for which the step is active. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>target_accept_percentage</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, maximum: `100`)

<!-- vale on -->

The value of the step.

</dd>
</dl>

---

<!-- vale off -->

### LoadScheduler {#load-scheduler}

<!-- vale on -->

_Load Scheduler_ creates a queue for flows in front of a service to provide
active service protection

:::info

See also [_Load Scheduler_ overview](/concepts/scheduler/load-scheduler.md).

:::

To make scheduling decisions the Flows are mapped into Workloads by providing
match rules. A workload determines the priority and cost of admitting each Flow
that belongs to it. Scheduling of Flows is based on Weighted Fair Queuing
principles. _Load Scheduler_ measures and controls the incoming tokens per
second, which can translate to (avg. latency \* in-flight requests) (Little's
Law) in concurrency limiting use-case.

The signal at port `load_multiplier` determines the fraction of incoming tokens
that get admitted.

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

Decides whether to run the load scheduler in dry-run mode. In dry run mode the
scheduler acts as pass through to all flow and does not queue flows. It is
useful for observing the behavior of load scheduler without disrupting any real
traffic.

</dd>
<dt>dry_run_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for setting dry run mode through dynamic configuration.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadSchedulerIns](#load-scheduler-ins))

<!-- vale on -->

Input ports for the LoadScheduler component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([LoadSchedulerOuts](#load-scheduler-outs))

<!-- vale on -->

Output ports for the LoadScheduler component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadSchedulerParameters](#load-scheduler-parameters))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerIns {#load-scheduler-ins}

<!-- vale on -->

Input for the LoadScheduler component.

<dl>
<dt>load_multiplier</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Load multiplier is proportion of incoming token rate that needs to be accepted.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerOuts {#load-scheduler-outs}

<!-- vale on -->

Output for the LoadScheduler component.

<dl>
<dt>observed_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Observed load multiplier is the proportion of incoming token rate that is being
accepted.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerParameters {#load-scheduler-parameters}

<!-- vale on -->

Parameters for _Load Scheduler_ component

<dl>
<dt>scheduler</dt>
<dd>

<!-- vale off -->

([Scheduler](#scheduler))

<!-- vale on -->

Configuration of Weighted Fair Queuing-based workload scheduler.

Contains configuration of per-agent scheduler

</dd>
<dt>selectors</dt>
<dd>

<!-- vale off -->

([[]Selector](#selector), **required**)

<!-- vale on -->

Selectors for the component.

</dd>
<dt>workload_latency_based_tokens</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Automatically estimate the size flows within each workload, based on historical
latency. Each workload's `tokens` will be set to average latency of flows in
that workload during the last few seconds (exact duration of this average can
change). This setting is useful in concurrency limiting use-case, where the
concurrency is calculated as ``(avg. latency \* in-flight flows)`.

The value of tokens estimated takes a lower precedence than the value of
`tokens` specified in the workload definition and `tokens` explicitly specified
in the flow labels.

</dd>
</dl>

---

<!-- vale off -->

### MatchExpression {#match-expression}

<!-- vale on -->

MatchExpression has multiple variants, exactly one should be set.

Example:

```yaml
all:
  of:
    - label_exists: foo
    - label_equals:
        label: app
        value: frobnicator
```

<dl>
<dt>all</dt>
<dd>

<!-- vale off -->

([MatchExpressionList](#match-expression-list))

<!-- vale on -->

The expression is true when all sub expressions are true.

</dd>
<dt>any</dt>
<dd>

<!-- vale off -->

([MatchExpressionList](#match-expression-list))

<!-- vale on -->

The expression is true when any sub expression is true.

</dd>
<dt>label_equals</dt>
<dd>

<!-- vale off -->

([EqualsMatchExpression](#equals-match-expression))

<!-- vale on -->

The expression is true when label value equals given value.

</dd>
<dt>label_exists</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

The expression is true when label with given name exists.

</dd>
<dt>label_matches</dt>
<dd>

<!-- vale off -->

([MatchesMatchExpression](#matches-match-expression))

<!-- vale on -->

The expression is true when label matches given regular expression.

</dd>
<dt>not</dt>
<dd>

<!-- vale off -->

([MatchExpression](#match-expression))

<!-- vale on -->

The expression negates the result of sub expression.

</dd>
</dl>

---

<!-- vale off -->

### MatchExpressionList {#match-expression-list}

<!-- vale on -->

List of MatchExpressions that is used for all or any matching

for example, `{any: {of: [expr1, expr2]}}`.

<dl>
<dt>of</dt>
<dd>

<!-- vale off -->

([[]MatchExpression](#match-expression))

<!-- vale on -->

List of sub expressions of the match expression.

</dd>
</dl>

---

<!-- vale off -->

### MatchesMatchExpression {#matches-match-expression}

<!-- vale on -->

Label selector expression of the form `label matches regex`.

<dl>
<dt>label</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Name of the label to match the regular expression.

</dd>
<dt>regex</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Regular expression that should match the label value. It uses
[Go's regular expression syntax](https://github.com/google/re2/wiki/Syntax).

</dd>
</dl>

---

<!-- vale off -->

### Max {#max}

<!-- vale on -->

Takes a list of input signals and emits the signal with the maximum value

Max: output = max([]inputs).

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([MaxIns](#max-ins))

<!-- vale on -->

Input ports for the Max component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([MaxOuts](#max-outs))

<!-- vale on -->

Output ports for the Max component.

</dd>
</dl>

---

<!-- vale off -->

### MaxIns {#max-ins}

<!-- vale on -->

Inputs for the Max component.

<dl>
<dt>inputs</dt>
<dd>

<!-- vale off -->

([[]InPort](#in-port))

<!-- vale on -->

Array of input signals.

</dd>
</dl>

---

<!-- vale off -->

### MaxOuts {#max-outs}

<!-- vale on -->

Output for the Max component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Signal with maximum value as an output signal.

</dd>
</dl>

---

<!-- vale off -->

### Min {#min}

<!-- vale on -->

Takes an array of input signals and emits the signal with the minimum value Min:
output = min([]inputs).

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([MinIns](#min-ins))

<!-- vale on -->

Input ports for the Min component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([MinOuts](#min-outs))

<!-- vale on -->

Output ports for the Min component.

</dd>
</dl>

---

<!-- vale off -->

### MinIns {#min-ins}

<!-- vale on -->

Inputs for the Min component.

<dl>
<dt>inputs</dt>
<dd>

<!-- vale off -->

([[]InPort](#in-port))

<!-- vale on -->

Array of input signals.

</dd>
</dl>

---

<!-- vale off -->

### MinOuts {#min-outs}

<!-- vale on -->

Output ports for the Min component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Signal with minimum value as an output signal.

</dd>
</dl>

---

<!-- vale off -->

### NestedCircuit {#nested-circuit}

<!-- vale on -->

Nested circuit defines a sub-circuit as a high-level component. It consists of a
list of components and a map of input and output ports.

<dl>
<dt>components</dt>
<dd>

<!-- vale off -->

([[]Component](#component))

<!-- vale on -->

List of components in the nested circuit.

</dd>
<dt>in_ports_map</dt>
<dd>

<!-- vale off -->

(map of [InPort](#in-port))

<!-- vale on -->

Maps input port names to input ports.

</dd>
<dt>name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the nested circuit component. This name is displayed by graph
visualization tools.

</dd>
<dt>out_ports_map</dt>
<dd>

<!-- vale off -->

(map of [OutPort](#out-port))

<!-- vale on -->

Maps output port names to output ports.

</dd>
<dt>short_description</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Short description of the nested circuit component. This description is displayed
by graph visualization tools.

</dd>
</dl>

---

<!-- vale off -->

### NestedSignalEgress {#nested-signal-egress}

<!-- vale on -->

Nested signal egress is a special type of component that allows to extract a
signal from a nested circuit.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([NestedSignalEgressIns](#nested-signal-egress-ins))

<!-- vale on -->

Input ports for the NestedSignalEgress component.

</dd>
<dt>port_name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the port.

</dd>
</dl>

---

<!-- vale off -->

### NestedSignalEgressIns {#nested-signal-egress-ins}

<!-- vale on -->

Inputs for the NestedSignalEgress component.

<dl>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Egress signal.

</dd>
</dl>

---

<!-- vale off -->

### NestedSignalIngress {#nested-signal-ingress}

<!-- vale on -->

Nested signal ingress is a special type of component that allows to inject a
signal into a nested circuit.

<dl>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([NestedSignalIngressOuts](#nested-signal-ingress-outs))

<!-- vale on -->

Output ports for the NestedSignalIngress component.

</dd>
<dt>port_name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the port.

</dd>
</dl>

---

<!-- vale off -->

### NestedSignalIngressOuts {#nested-signal-ingress-outs}

<!-- vale on -->

Outputs for the NestedSignalIngress component.

<dl>
<dt>signal</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Ingress signal.

</dd>
</dl>

---

<!-- vale off -->

### Or {#or}

<!-- vale on -->

Logical OR.

See [And component](#and) on how signals are mapped onto Boolean values.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([OrIns](#or-ins))

<!-- vale on -->

Input ports for the Or component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([OrOuts](#or-outs))

<!-- vale on -->

Output ports for the Or component.

</dd>
</dl>

---

<!-- vale off -->

### OrIns {#or-ins}

<!-- vale on -->

Inputs for the Or component.

<dl>
<dt>inputs</dt>
<dd>

<!-- vale off -->

([[]InPort](#in-port))

<!-- vale on -->

Array of input signals.

</dd>
</dl>

---

<!-- vale off -->

### OrOuts {#or-outs}

<!-- vale on -->

Output ports for the Or component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Result of logical OR of all the input signals.

Will always be 0 (false), 1 (true) or invalid (unknown).

</dd>
</dl>

---

<!-- vale off -->

### OutPort {#out-port}

<!-- vale on -->

Components produce output for other components through OutPorts

<dl>
<dt>signal_name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the outgoing Signal on the OutPort.

</dd>
</dl>

---

<!-- vale off -->

### PathTemplateMatcher {#path-template-matcher}

<!-- vale on -->

Matches HTTP Path to given path templates

HTTP path will be matched against given path templates. If a match occurs, the
value associated with the path template will be treated as a result. In case of
multiple path templates matching, the most specific one will be chosen.

<dl>
<dt>template_values</dt>
<dd>

<!-- vale off -->

(map of string)

<!-- vale on -->

Template value keys are OpenAPI-inspired path templates.

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

```yaml
/register: register
"/user/{userId}": user
/static/*: other
```

</dd>
</dl>

---

<!-- vale off -->

### PeriodicDecrease {#periodic-decrease}

<!-- vale on -->

PeriodicDecrease defines a controller for scaling in based on a periodic timer.

<dl>
<dt>period</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

The period of the timer. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>scale_in_percentage</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, maximum: `100`, **required**)

<!-- vale on -->

The percentage of scale to reduce.

</dd>
</dl>

---

<!-- vale off -->

### PodScaler {#pod-scaler}

<!-- vale on -->

Component for scaling pods based on a signal.

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

Dry run mode ensures that no scaling is invoked by this pod scaler. This is
useful for observing the behavior of pod scaler without disrupting any real
deployment. This parameter sets the default value of dry run setting which can
be overridden at runtime using dynamic configuration.

</dd>
<dt>dry_run_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for overriding dry run setting through dynamic configuration.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([PodScalerIns](#pod-scaler-ins))

<!-- vale on -->

Input ports for the PodScaler component.

</dd>
<dt>kubernetes_object_selector</dt>
<dd>

<!-- vale off -->

([KubernetesObjectSelector](#kubernetes-object-selector))

<!-- vale on -->

The Kubernetes object to which this pod scaler applies.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([PodScalerOuts](#pod-scaler-outs))

<!-- vale on -->

Output ports for the PodScaler component.

</dd>
</dl>

---

<!-- vale off -->

### PodScalerIns {#pod-scaler-ins}

<!-- vale on -->

Inputs for the PodScaler component.

<dl>
<dt>replicas</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The number of replicas to configure for the Kubernetes resource

</dd>
</dl>

---

<!-- vale off -->

### PodScalerOuts {#pod-scaler-outs}

<!-- vale on -->

Outputs for the PodScaler component.

<dl>
<dt>actual_replicas</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The number of replicas that are currently running.

</dd>
<dt>configured_replicas</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The number of replicas that are desired.

</dd>
</dl>

---

<!-- vale off -->

### Policy {#policy}

<!-- vale on -->

Policy expresses observability-driven control logic.

:::info

See also [Policy overview](/concepts/advanced/policy.md).

:::

Policy specification contains a circuit that defines the controller logic and
resources that need to be setup.

<dl>
<dt>circuit</dt>
<dd>

<!-- vale off -->

([Circuit](#circuit))

<!-- vale on -->

Defines the control-loop logic of the policy.

</dd>
<dt>resources</dt>
<dd>

<!-- vale off -->

([Resources](#resources))

<!-- vale on -->

Resources (such as Flux Meters, Classifiers) to setup.

</dd>
</dl>

---

<!-- vale off -->

### PromQL {#prom-q-l}

<!-- vale on -->

Component that runs a Prometheus query periodically and returns the result as an
output signal

<dl>
<dt>evaluation_interval</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Describes the interval between successive evaluations of the Prometheus query.
This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([PromQLOuts](#prom-q-l-outs))

<!-- vale on -->

Output ports for the PromQL component.

</dd>
<dt>query_string</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Describes the
[PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) query to
be run.

:::note

The query must return a single value either as a scalar or as a vector with a
single element.

:::

:::info Usage with Flux Meter

[Flux Meter](/concepts/flux-meter.md) metrics can be queried using PromQL. Flux
Meter defines histogram type of metrics in Prometheus. Therefore, one can refer
to `flux_meter_sum`, `flux_meter_count` and `flux_meter_bucket`. The particular
Flux Meter can be identified with the `flux_meter_name` label. There are
additional labels available on a Flux Meter such as `valid`, `flow_status`,
`http_status_code` and `decision_type`.

:::

:::info Usage with OpenTelemetry Metrics

Aperture supports OpenTelemetry metrics. See
[reference](/integrations/metrics/metrics.md) for more details.

:::

</dd>
</dl>

---

<!-- vale off -->

### PromQLOuts {#prom-q-l-outs}

<!-- vale on -->

Output for the PromQL component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The result of the Prometheus query as an output signal.

</dd>
</dl>

---

<!-- vale off -->

### PulseGenerator {#pulse-generator}

<!-- vale on -->

Generates 0 and 1 in turns.

<dl>
<dt>false_for</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Emitting 0 for the `false_for` duration. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([PulseGeneratorOuts](#pulse-generator-outs))

<!-- vale on -->

Output ports for the PulseGenerator component.

</dd>
<dt>true_for</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Emitting 1 for the `true_for` duration. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
</dl>

---

<!-- vale off -->

### PulseGeneratorOuts {#pulse-generator-outs}

<!-- vale on -->

Outputs for the PulseGenerator component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### Query {#query}

<!-- vale on -->

Query components that are query databases such as Prometheus.

<dl>
<dt>promql</dt>
<dd>

<!-- vale off -->

([PromQL](#prom-q-l))

<!-- vale on -->

Periodically runs a Prometheus query in the background and emits the result.

</dd>
</dl>

---

<!-- vale off -->

### QuotaScheduler {#quota-scheduler}

<!-- vale on -->

Schedules the traffic based on token-bucket based quotas.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([RateLimiterIns](#rate-limiter-ins))

<!-- vale on -->

</dd>
<dt>rate_limiter</dt>
<dd>

<!-- vale off -->

([RateLimiterParameters](#rate-limiter-parameters))

<!-- vale on -->

</dd>
<dt>scheduler</dt>
<dd>

<!-- vale off -->

([Scheduler](#scheduler))

<!-- vale on -->

</dd>
<dt>selectors</dt>
<dd>

<!-- vale off -->

([[]Selector](#selector), **required**)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### RateLimiter {#rate-limiter}

<!-- vale on -->

Limits the traffic on a control point to specified rate

:::info

See also [_Rate Limiter_ overview](/concepts/rate-limiter.md).

:::

RateLimiting is done on per-label-value (`label_key`) basis and it uses the
_Token Bucket Algorithm_.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([RateLimiterIns](#rate-limiter-ins))

<!-- vale on -->

Input ports for the RateLimiter component

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([RateLimiterParameters](#rate-limiter-parameters))

<!-- vale on -->

Parameters for the RateLimiter component

</dd>
<dt>selectors</dt>
<dd>

<!-- vale off -->

([[]Selector](#selector), **required**)

<!-- vale on -->

Selectors for the component.

</dd>
</dl>

---

<!-- vale off -->

### RateLimiterIns {#rate-limiter-ins}

<!-- vale on -->

Inputs for the RateLimiter component

<dl>
<dt>bucket_capacity</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Capacity of the bucket.

</dd>
<dt>fill_amount</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Number of tokens to fill within an `interval`.

</dd>
<dt>pass_through</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

PassThrough port determines whether all requests

</dd>
</dl>

---

<!-- vale off -->

### RateLimiterParameters {#rate-limiter-parameters}

<!-- vale on -->

<dl>
<dt>continuous_fill</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Continuous fill determines whether the token bucket should be filled
continuously or only on discrete intervals.

</dd>
<dt>interval</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Interval defines the time interval in which the token bucket will fill tokens
specified by `fill_amount` signal. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>label_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Specifies which label the rate limiter should be keyed by.

Rate limiting is done independently for each value of the
[label](/concepts/flow-label.md) with given key. For example, to give each user
a separate limit, assuming you have a _user_ flow label set up, set
`label_key: "user"`. If no label key is specified, then all requests matching
the selectors will be rate limited based on the global bucket.

</dd>
<dt>lazy_sync</dt>
<dd>

<!-- vale off -->

([RateLimiterParametersLazySync](#rate-limiter-parameters-lazy-sync))

<!-- vale on -->

Configuration of lazy-syncing behavior of rate limiter

</dd>
<dt>max_idle_time</dt>
<dd>

<!-- vale off -->

(string, default: `"7200s"`)

<!-- vale on -->

Max idle time before token bucket state for a label is removed. If set to 0, the
state is never removed. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>tokens_label_key</dt>
<dd>

<!-- vale off -->

(string, default: `"tokens"`)

<!-- vale on -->

Flow label key that will be used to override the number of tokens for this
request. This is an optional parameter and takes highest precedence when
assigning tokens to a request. The label value must be a valid uint64 number.

</dd>
</dl>

---

<!-- vale off -->

### RateLimiterParametersLazySync {#rate-limiter-parameters-lazy-sync}

<!-- vale on -->

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Enables lazy sync

</dd>
<dt>num_sync</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, default: `4`)

<!-- vale on -->

Number of times to lazy sync within the `interval`.

</dd>
</dl>

---

<!-- vale off -->

### Rego {#rego}

<!-- vale on -->

Rego define a set of labels that are extracted after evaluating a Rego module.

:::info

You can use the [live-preview](/concepts/classifier.md#live-previewing-requests)
feature to first preview the input to the classifier before writing the labeling
logic.

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
```

<dl>
<dt>labels</dt>
<dd>

<!-- vale off -->

(map of [RegoLabelProperties](#rego-label-properties), **required**)

<!-- vale on -->

A map of {key, value} pairs mapping from [flow label](/concepts/flow-label.md)
keys to queries that define how to extract and propagate flow labels with that
key. The name of the label maps to a variable in the Rego module. It maps to
`data.<package>.<label>` variable.

</dd>
<dt>module</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Source code of the Rego module.

:::note

Must include a "package" declaration.

:::

</dd>
</dl>

---

<!-- vale off -->

### RegoLabelProperties {#rego-label-properties}

<!-- vale on -->

<dl>
<dt>telemetry</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Decides if the created flow label should be available as an attribute in OLAP
telemetry and propagated in [baggage](/concepts/flow-label.md#baggage)

:::note

The flow label is always accessible in Aperture Policies regardless of this
setting.

:::

:::caution

When using [FluxNinja ARC extension](/arc/extension.md), telemetry enabled
labels are sent to FluxNinja ARC for observability. Telemetry should be disabled
for sensitive labels.

:::

</dd>
</dl>

---

<!-- vale off -->

### Resources {#resources}

<!-- vale on -->

Resources that need to be setup for the policy to function

:::info

See also [Resources overview](/concepts/advanced/policy.md).

:::

<dl>
<dt>flow_control</dt>
<dd>

<!-- vale off -->

([FlowControlResources](#flow-control-resources))

<!-- vale on -->

FlowControlResources are resources that are provided by flow control
integration.

</dd>
<dt>infra_meters</dt>
<dd>

<!-- vale off -->

(map of [InfraMeter](#infra-meter))

<!-- vale on -->

_Infra Meters_ configure custom metrics OpenTelemetry collector pipelines, which
will receive and process telemetry at the agents and send metrics to the
configured Prometheus. Key in this map refers to OTel pipeline name. Prefixing
pipeline name with `metrics/` is optional, as all the components and pipeline
names would be normalized.

Example:

```yaml
infra_meters:
  rabbitmq:
    agent_group: default
    per_agent_group: true
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

```

:::caution

Validate the OTel configuration before applying it to the production cluster.
Incorrect configuration will get rejected at the agents and might cause shutdown
of the agent(s).

:::

</dd>
<dt>telemetry_collectors</dt>
<dd>

<!-- vale off -->

([[]TelemetryCollector](#telemetry-collector))

<!-- vale on -->

TelemetryCollector configures OpenTelemetry collector integration. Deprecated:
v3.0.0. Use `infra_meters` instead.

</dd>
</dl>

---

<!-- vale off -->

### Rule {#rule}

<!-- vale on -->

Rule describes a single classification Rule

Example of a JSON extractor:

```yaml
extractor:
  json:
    from: request.http.body
    pointer: /user/name
```

<dl>
<dt>extractor</dt>
<dd>

<!-- vale off -->

([Extractor](#extractor))

<!-- vale on -->

High-level declarative extractor.

</dd>
<dt>telemetry</dt>
<dd>

<!-- vale off -->

(bool, default: `true`)

<!-- vale on -->

Decides if the created flow label should be available as an attribute in OLAP
telemetry and propagated in [baggage](/concepts/flow-label.md#baggage)

:::note

The flow label is always accessible in Aperture Policies regardless of this
setting.

:::

:::caution

When using [FluxNinja ARC extension](/arc/extension.md), telemetry enabled
labels are sent to FluxNinja ARC for observability. Telemetry should be disabled
for sensitive labels.

:::

</dd>
</dl>

---

<!-- vale off -->

### SMA {#s-m-a}

<!-- vale on -->

Simple Moving Average (SMA) is a type of moving average that computes the
average of a fixed number of signal readings.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([SMAIns](#s-m-a-ins))

<!-- vale on -->

Input ports for the SMA component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([SMAOuts](#s-m-a-outs))

<!-- vale on -->

Output ports for the SMA component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([SMAParameters](#s-m-a-parameters))

<!-- vale on -->

Parameters for the SMA component.

</dd>
</dl>

---

<!-- vale off -->

### SMAIns {#s-m-a-ins}

<!-- vale on -->

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Signal to be used for the moving average computation.

</dd>
</dl>

---

<!-- vale off -->

### SMAOuts {#s-m-a-outs}

<!-- vale on -->

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Computed moving average.

</dd>
</dl>

---

<!-- vale off -->

### SMAParameters {#s-m-a-parameters}

<!-- vale on -->

<dl>
<dt>sma_window</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Window of time over which the moving average is computed. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>valid_during_warmup</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Whether the output is valid during the warm-up stage.

</dd>
</dl>

---

<!-- vale off -->

### Sampler {#sampler}

<!-- vale on -->

_Sampler_ is a component that regulates the load at a
[_Control Point_](/concepts/selector.md/#control-point) by allowing only a
specified percentage of flows at random or by sticky sessions.

:::info

See also [_Sampler_ overview](/concepts/load-ramp.md#sampler).

:::

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([SamplerIns](#sampler-ins))

<!-- vale on -->

Input ports for the _Sampler_.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([SamplerParameters](#sampler-parameters))

<!-- vale on -->

Parameters for the _Sampler_.

</dd>
<dt>pass_through_label_values</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

Specify certain label values to be always accepted by this _Sampler_ regardless
of accept percentage.

</dd>
<dt>pass_through_label_values_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for setting pass through label values through dynamic
configuration.

</dd>
</dl>

---

<!-- vale off -->

### SamplerIns {#sampler-ins}

<!-- vale on -->

<dl>
<dt>accept_percentage</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The percentage of requests to accept.

</dd>
</dl>

---

<!-- vale off -->

### SamplerParameters {#sampler-parameters}

<!-- vale on -->

<dl>
<dt>label_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

The flow label key for identifying sessions.

- When label key is specified, _Sampler_ acts as a sticky filter. The series of
  flows with the same value of label key get the same decision provided that the
  `accept_percentage` is same or higher.
- When label key is not specified, _Sampler_ acts as a stateless filter.
  Percentage of flows are selected randomly for rejection.

</dd>
<dt>selectors</dt>
<dd>

<!-- vale off -->

([[]Selector](#selector), **required**)

<!-- vale on -->

Selectors for the component.

</dd>
</dl>

---

<!-- vale off -->

### ScaleInController {#scale-in-controller}

<!-- vale on -->

<dl>
<dt>alerter</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for embedded Alerter.

</dd>
<dt>controller</dt>
<dd>

<!-- vale off -->

([ScaleInControllerController](#scale-in-controller-controller))

<!-- vale on -->

Controller

</dd>
</dl>

---

<!-- vale off -->

### ScaleInControllerController {#scale-in-controller-controller}

<!-- vale on -->

<dl>
<dt>gradient</dt>
<dd>

<!-- vale off -->

([DecreasingGradient](#decreasing-gradient))

<!-- vale on -->

</dd>
<dt>periodic</dt>
<dd>

<!-- vale off -->

([PeriodicDecrease](#periodic-decrease))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### ScaleOutController {#scale-out-controller}

<!-- vale on -->

<dl>
<dt>alerter</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for embedded Alerter.

</dd>
<dt>controller</dt>
<dd>

<!-- vale off -->

([ScaleOutControllerController](#scale-out-controller-controller))

<!-- vale on -->

Controller

</dd>
</dl>

---

<!-- vale off -->

### ScaleOutControllerController {#scale-out-controller-controller}

<!-- vale on -->

<dl>
<dt>gradient</dt>
<dd>

<!-- vale off -->

([IncreasingGradient](#increasing-gradient))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### Scheduler {#scheduler}

<!-- vale on -->

Weighted Fair Queuing-based workload scheduler

:::note

Each Agent instantiates an independent copy of the scheduler, but output signals
for accepted and incoming token rate are aggregated across all agents.

:::

<dl>
<dt>decision_deadline_margin</dt>
<dd>

<!-- vale off -->

(string, default: `"0.01s"`)

<!-- vale on -->

Decision deadline margin is the amount of time that the scheduler will subtract
from the request deadline to determine the deadline for the decision. This is to
ensure that the scheduler has enough time to make a decision before the request
deadline happens, accounting for processing delays. The request deadline is
based on the [gRPC deadline](https://grpc.io/blog/deadlines) or the
[`grpc-timeout` HTTP header](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests).

Fail-open logic is use for flow control APIs, so if the gRPC deadline reaches,
the flow will end up being unconditionally allowed while it is still waiting on
the scheduler. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>default_workload_parameters</dt>
<dd>

<!-- vale off -->

([SchedulerWorkloadParameters](#scheduler-workload-parameters))

<!-- vale on -->

Parameters to be used if none of workloads specified in `workloads` match.

</dd>
<dt>priorities_label_key</dt>
<dd>

<!-- vale off -->

(string, default: `"priorities"`)

<!-- vale on -->

- Key for a flow label that can be used to override the default priority for
  this flow.
- The value associated with this key must be a valid uint64 number. Higher
  numbers means higher priority.
- If this parameter is not provided, the priority for the flow will be
  determined by the matched workload's priority.

</dd>
<dt>tokens_label_key</dt>
<dd>

<!-- vale off -->

(string, default: `"tokens"`)

<!-- vale on -->

- Key for a flow label that can be used to override the default number of tokens
  for this flow.
- The value associated with this key must be a valid uint64 number.
- If this parameter is not provided, the number of tokens for the flow will be
  determined by the matched workload's token count.

</dd>
<dt>workloads</dt>
<dd>

<!-- vale off -->

([[]SchedulerWorkload](#scheduler-workload))

<!-- vale on -->

List of workloads to be used in scheduler.

Categorizing flows into workloads allows for load throttling to be "intelligent"
instead of queueing flows in an arbitrary order. There are two aspects of this
"intelligence":

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
[workload definition in the concepts section](/concepts/scheduler/scheduler.md#workload).

:::

</dd>
</dl>

---

<!-- vale off -->

### SchedulerWorkload {#scheduler-workload}

<!-- vale on -->

Workload defines a class of flows that preferably have similar properties such
as response latency and desired priority.

<dl>
<dt>label_matcher</dt>
<dd>

<!-- vale off -->

([LabelMatcher](#label-matcher))

<!-- vale on -->

Label Matcher to select a Workload based on
[flow labels](/concepts/flow-label.md).

</dd>
<dt>name</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Name of the workload.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([SchedulerWorkloadParameters](#scheduler-workload-parameters))

<!-- vale on -->

Parameters associated with flows matching the label matcher.

</dd>
</dl>

---

<!-- vale off -->

### SchedulerWorkloadParameters {#scheduler-workload-parameters}

<!-- vale on -->

Parameters such as priority and tokens that are applicable to flows within a
workload.

<dl>
<dt>priority</dt>
<dd>

<!-- vale off -->

(float64, minimum: `0`, default: `1`)

<!-- vale on -->

Describes priority level of the flows within the workload. Priority level is
unbounded and can be any positive number. Higher numbers means higher priority
level. The following formula is used to determine the position of a flow in the
queue based on virtual finish time:

$$
inverted\_priority = {\frac {1} {priority}}
$$

$$
virtual\_finish\_time = virtual\_time + \left(tokens \cdot inverted\_priority\right)
$$

</dd>
<dt>queue_timeout</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Timeout for the flow in the workload. If timeout is provided on the Check call
as well, we pick the minimum of the two. If this override is not provided, the
timeout provided in the check call is used. 0 timeout value implies that the
request will not wait in the queue and will be accepted/dropped immediately.
This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>tokens</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Tokens determines the cost of admitting a single flow in the workload, which is
typically defined as milliseconds of flow latency (time to response or duration
of a feature) or simply equal to 1 if the resource being accessed is constrained
by the number of flows (3rd party rate limiters). This override is applicable
only if tokens for the flow aren't specified in the flow labels.

</dd>
</dl>

---

<!-- vale off -->

### Selector {#selector}

<!-- vale on -->

Selects flows based on control point, flow labels, agent group and the service
that the flow control component will operate on.

:::info

See also [Selector overview](/concepts/selector.md).

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

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

[_Agent Group_](/concepts/selector.md#agent-group) this selector applies to.

:::info

Agent Groups are used to scope policies to a subset of agents connected to the
same controller. The agents within an agent group receive exact same policy
configuration and form a peer to peer cluster to constantly share state.

:::

</dd>
<dt>control_point</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

[Control Point](/concepts/control-point.md) identifies location within services
where policies can act on flows. For an SDK based insertion, a _Control Point_
can represent a particular feature or execution block within a service. In case
of service mesh or middleware insertion, a _Control Point_ can identify ingress
or egress calls or distinct listeners or filter chains.

</dd>
<dt>label_matcher</dt>
<dd>

<!-- vale off -->

([LabelMatcher](#label-matcher))

<!-- vale on -->

[Label Matcher](/concepts/selector.md#label-matcher) can be used to match flows
based on flow labels.

</dd>
<dt>service</dt>
<dd>

<!-- vale off -->

(string, default: `"any"`)

<!-- vale on -->

The Fully Qualified Domain Name of the [service](/concepts/selector.md) to
select.

In Kubernetes, this is the FQDN of the Service object.

:::info

`any` matches all services.

:::

:::info

An entity (for example, Kubernetes pod) might belong to multiple services.

:::

</dd>
</dl>

---

<!-- vale off -->

### SignalGenerator {#signal-generator}

<!-- vale on -->

The _Signal Generator_ component generates a smooth and continuous signal by
following a sequence of specified steps. Each step has two parameters:

- `target_output`: The desired output value at the end of the step.
- `duration`: The time it takes for the signal to change linearly from the
  previous step's `target_output` to the current step's `target_output`.

The output signal starts at the `target_output` of the first step and changes
linearly between steps based on their `duration`. The _Signal Generator_ can be
controlled to move forwards, backwards, or reset to the beginning based on input
signals.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([SignalGeneratorIns](#signal-generator-ins))

<!-- vale on -->

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([SignalGeneratorOuts](#signal-generator-outs))

<!-- vale on -->

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([SignalGeneratorParameters](#signal-generator-parameters))

<!-- vale on -->

Parameters for the _Signal Generator_ component.

</dd>
</dl>

---

<!-- vale off -->

### SignalGeneratorIns {#signal-generator-ins}

<!-- vale on -->

Inputs for the _Signal Generator_ component.

<dl>
<dt>backward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the _Signal Generator_ towards the previous step.

</dd>
<dt>forward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the _Signal Generator_ towards the next step.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to reset the _Signal Generator_ to the first step.

</dd>
</dl>

---

<!-- vale off -->

### SignalGeneratorOuts {#signal-generator-outs}

<!-- vale on -->

Outputs for the _Signal Generator_ component.

<dl>
<dt>at_end</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Signal Generator_ is at the end of
signal generation.

</dd>
<dt>at_start</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Signal Generator_ is at the start of
signal generation.

</dd>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The generated signal.

</dd>
</dl>

---

<!-- vale off -->

### SignalGeneratorParameters {#signal-generator-parameters}

<!-- vale on -->

Parameters for the _Signal Generator_ component.

<dl>
<dt>steps</dt>
<dd>

<!-- vale off -->

([[]SignalGeneratorParametersStep](#signal-generator-parameters-step),
**required**)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### SignalGeneratorParametersStep {#signal-generator-parameters-step}

<!-- vale on -->

<dl>
<dt>duration</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Duration for which the step is active. This field employs the
[Duration](https://developers.google.com/protocol-buffers/docs/proto3#json) JSON
representation from Protocol Buffers. The format accommodates fractional seconds
up to nine digits after the decimal point, offering nanosecond precision. Every
duration value must be suffixed with an "s" to indicate 'seconds.' For example,
a value of "10s" would signify a duration of 10 seconds.

</dd>
<dt>target_output</dt>
<dd>

<!-- vale off -->

([ConstantSignal](#constant-signal))

<!-- vale on -->

The value of the step.

</dd>
</dl>

---

<!-- vale off -->

### Switcher {#switcher}

<!-- vale on -->

Type of Combinator that switches between `on_signal` and `off_signal` signals
based on switch input

`on_signal` will be returned if switch input is valid and not equal to 0.0 ,
otherwise `off_signal` will be returned.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([SwitcherIns](#switcher-ins))

<!-- vale on -->

Input ports for the Switcher component.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([SwitcherOuts](#switcher-outs))

<!-- vale on -->

Output ports for the Switcher component.

</dd>
</dl>

---

<!-- vale off -->

### SwitcherIns {#switcher-ins}

<!-- vale on -->

Inputs for the Switcher component.

<dl>
<dt>off_signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Output signal when switch is invalid or 0.0.

</dd>
<dt>on_signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Output signal when switch is valid and not 0.0.

</dd>
<dt>switch</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Decides whether to return `on_signal` or `off_signal`.

</dd>
</dl>

---

<!-- vale off -->

### SwitcherOuts {#switcher-outs}

<!-- vale on -->

Outputs for the Switcher component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Selected signal (`on_signal` or `off_signal`).

</dd>
</dl>

---

<!-- vale off -->

### TelemetryCollector {#telemetry-collector}

<!-- vale on -->

TelemetryCollector defines the telemetry configuration to be synced with the
agents. Deprecated: v3.0.0. Use `InfraMeter` instead. It consists of two parts:

- Agent Group: Agent group to sync telemetry configuration with
- Infra Meters: OTel compatible metrics pipelines

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

</dd>
<dt>infra_meters</dt>
<dd>

<!-- vale off -->

(map of [InfraMeter](#infra-meter))

<!-- vale on -->

_Infra Meters_ configure custom metrics OpenTelemetry collector pipelines, which
will receive and process telemetry at the agents and send metrics to the
configured Prometheus. Key in this map refers to OTel pipeline name. Prefixing
pipeline name with `metrics/` is optional, as all the components and pipeline
names would be normalized.

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

Validate the OTel configuration before applying it to the production cluster.
Incorrect configuration will get rejected at the agents and might cause shutdown
of the agent(s).

:::

</dd>
</dl>

---

<!-- vale off -->

### UnaryOperator {#unary-operator}

<!-- vale on -->

Takes an input signal and emits the output after applying the specified unary
operator

$$
\text{output} = \unary_operator{\text{input}}
$$

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([UnaryOperatorIns](#unary-operator-ins))

<!-- vale on -->

Input ports for the UnaryOperator component.

</dd>
<dt>operator</dt>
<dd>

<!-- vale off -->

(string, one of:
`abs | acos | acosh | asin | asinh | atan | atanh | cbrt | ceil | cos | cosh | erf | erfc | erfcinv | erfinv | exp | exp2 | expm1 | floor | gamma | j0 | j1 | lgamma | log | log10 | log1p | log2 | round | roundtoeven | sin | sinh | sqrt | tan | tanh | trunc | y0 | y1`)

<!-- vale on -->

Unary Operator to apply.

The unary operator can be one of the following:

- `abs`: Absolute value with the sign removed.
- `acos`: `arccosine`, in radians.
- `acosh`: Inverse hyperbolic cosine.
- `asin`: `arcsine`, in radians.
- `asinh`: Inverse hyperbolic sine.
- `atan`: `arctangent`, in radians.
- `atanh`: Inverse hyperbolic tangent.
- `cbrt`: Cube root.
- `ceil`: Least integer value greater than or equal to input signal.
- `cos`: `cosine`, in radians.
- `cosh`: Hyperbolic cosine.
- `erf`: Error function.
- `erfc`: Complementary error function.
- `erfcinv`: Inverse complementary error function.
- `erfinv`: Inverse error function.
- `exp`: The base-e exponential of input signal.
- `exp2`: The base-2 exponential of input signal.
- `expm1`: The base-e exponential of input signal minus 1.
- `floor`: Greatest integer value less than or equal to input signal.
- `gamma`: Gamma function.
- `j0`: Bessel function of the first kind of order 0.
- `j1`: Bessel function of the first kind of order 1.
- `lgamma`: Natural logarithm of the absolute value of the gamma function.
- `log`: Natural logarithm of input signal.
- `log10`: Base-10 logarithm of input signal.
- `log1p`: Natural logarithm of input signal plus 1.
- `log2`: Base-2 logarithm of input signal.
- `round`: Round to nearest integer.
- `roundtoeven`: Round to nearest integer, with ties going to the nearest even
  integer.
- `sin`: `sine`, in radians.
- `sinh`: Hyperbolic sine.
- `sqrt`: Square root.
- `tan`: `tangent`, in radians.
- `tanh`: Hyperbolic tangent.
- `trunc`: Truncate to integer.
- `y0`: Bessel function of the second kind of order 0.
- `y1`: Bessel function of the second kind of order 1.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([UnaryOperatorOuts](#unary-operator-outs))

<!-- vale on -->

Output ports for the UnaryOperator component.

</dd>
</dl>

---

<!-- vale off -->

### UnaryOperatorIns {#unary-operator-ins}

<!-- vale on -->

Inputs for the UnaryOperator component.

<dl>
<dt>input</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Input signal.

</dd>
</dl>

---

<!-- vale off -->

### UnaryOperatorOuts {#unary-operator-outs}

<!-- vale on -->

Outputs for the UnaryOperator component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Output signal.

</dd>
</dl>

---

<!-- vale off -->

### Variable {#variable}

<!-- vale on -->

Component that emits a constant signal which can be changed at runtime through
dynamic configuration.

<dl>
<dt>config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for overriding value setting through dynamic configuration.

</dd>
<dt>constant_output</dt>
<dd>

<!-- vale off -->

([ConstantSignal](#constant-signal))

<!-- vale on -->

The constant signal emitted by this component. The value of the constant signal
can be overridden at runtime through dynamic configuration.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([VariableOuts](#variable-outs))

<!-- vale on -->

Output ports for the Variable component.

</dd>
</dl>

---

<!-- vale off -->

### VariableOuts {#variable-outs}

<!-- vale on -->

Outputs for the Variable component.

<dl>
<dt>output</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The value is emitted to the output port.

</dd>
</dl>

---

<!---
Generated File Ends
-->
