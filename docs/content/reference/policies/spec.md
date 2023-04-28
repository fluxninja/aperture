---
title: Policy Language Specification
sidebar_position: 2
sidebar_label: Specification
---

<!-- vale off -->

<head>
  <body className="schema-docs" />
</head>

<!-- vale on -->

Reference for all objects used in [the Policy language](/concepts/policy/policy.md).

The top-level object representing a policy is [Policy](#policy).

<!---
Generated File Starts
-->

## Objects

---

<!-- vale off -->

### AIMDConcurrencyController {#a-i-m-d-concurrency-controller}

<!-- vale on -->

High level concurrency control component. Baselines a signal using exponential moving average and applies concurrency limits based on deviation of signal from the baseline. Internally implemented as a nested circuit.

<dl>
<dt>alerter_parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for embedded Alerter.

</dd>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([LoadActuatorDynamicConfig](#load-actuator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Dynamic configuration key for load actuation.

</dd>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

Flow Selector decides the service and flows at which the concurrency limiter is applied.

</dd>
<dt>gradient_parameters</dt>
<dd>

<!-- vale off -->

([GradientControllerParameters](#gradient-controller-parameters))

<!-- vale on -->

Gradient parameters for the controller.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([AIMDConcurrencyControllerIns](#a-i-m-d-concurrency-controller-ins))

<!-- vale on -->

Input ports for the AIMDConcurrencyController component.

</dd>
<dt>load_multiplier_linear_increment</dt>
<dd>

<!-- vale off -->

(float64, default: `0.0025`)

<!-- vale on -->

Linear increment to load multiplier in each execution tick when the system is not in overloaded state.

</dd>
<dt>max_load_multiplier</dt>
<dd>

<!-- vale off -->

(float64, default: `2`)

<!-- vale on -->

Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([AIMDConcurrencyControllerOuts](#a-i-m-d-concurrency-controller-outs))

<!-- vale on -->

Output ports for the AIMDConcurrencyController component.

</dd>
<dt>scheduler_parameters</dt>
<dd>

<!-- vale off -->

([SchedulerParameters](#scheduler-parameters))

<!-- vale on -->

Scheduler parameters.

</dd>
</dl>

---

<!-- vale off -->

### AIMDConcurrencyControllerIns {#a-i-m-d-concurrency-controller-ins}

<!-- vale on -->

Inputs for the AIMDConcurrencyController component.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The enabled port controls whether the _Adaptive Load Scheduler_ can load shed _Flows_. By default, the _Adaptive Load Scheduler_ is enabled.

</dd>
<dt>setpoint</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The setpoint to the controller.

</dd>
<dt>signal</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The signal to the controller.

</dd>
</dl>

---

<!-- vale off -->

### AIMDConcurrencyControllerOuts {#a-i-m-d-concurrency-controller-outs}

<!-- vale on -->

Outputs for the AIMDConcurrencyController component.

<dl>
<dt>accepted_concurrency</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Accepted concurrency is the number of concurrent requests that are accepted by the service.

</dd>
<dt>desired_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Desired Load multiplier is the ratio of desired concurrency to the incoming concurrency.

</dd>
<dt>incoming_concurrency</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

IncomingConcurrency is the number of concurrent requests that are received by the service.

</dd>
<dt>is_overload</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Is overload is a Boolean signal that indicates whether the service is overloaded based on the deviation of the signal from the setpoint taking into account some tolerance.
Deprecated: 1.6.0

</dd>
<dt>observed_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Observed Load multiplier is the ratio of accepted concurrency to the incoming concurrency.

</dd>
</dl>

---

<!-- vale off -->

### AdaptiveLoadScheduler {#adaptive-load-scheduler}

<!-- vale on -->

The _Adaptive Load Scheduler_ adjusts the accepted token rate based on the deviation of the input signal from the setpoint..

<dl>
<dt>alerter_parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration parameters for the embedded Alerter.

</dd>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([LoadSchedulerActuatorDynamicConfig](#load-scheduler-actuator-dynamic-config))

<!-- vale on -->

Default dynamic configuration for load actuation.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Dynamic configuration key for load actuation.

</dd>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

_Flow Selector_ is responsible for choosing the _Flows_ to which the _Load Scheduler_ is applied.

</dd>
<dt>gradient_parameters</dt>
<dd>

<!-- vale off -->

([GradientControllerParameters](#gradient-controller-parameters))

<!-- vale on -->

Parameters for the _Gradient Controller_.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadSchedulerIns](#adaptive-load-scheduler-ins))

<!-- vale on -->

Collection of input ports for the _Adaptive Load Scheduler_ component.

</dd>
<dt>load_multiplier_linear_increment</dt>
<dd>

<!-- vale off -->

(float64, default: `0.0025`)

<!-- vale on -->

Linear increment to load multiplier in each execution tick when the system is not in overloaded state.

</dd>
<dt>max_load_multiplier</dt>
<dd>

<!-- vale off -->

(float64, default: `2`)

<!-- vale on -->

The accepted token rate is multiplied by this value to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) states, helping to protect the Service from sudden spikes in incoming token rate.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadSchedulerOuts](#adaptive-load-scheduler-outs))

<!-- vale on -->

Collection of output ports for the _Adaptive Load Scheduler_ component.

</dd>
<dt>scheduler_parameters</dt>
<dd>

<!-- vale off -->

([LoadSchedulerSchedulerParameters](#load-scheduler-scheduler-parameters))

<!-- vale on -->

Parameters for the _Load Scheduler_.

</dd>
</dl>

---

<!-- vale off -->

### AdaptiveLoadSchedulerIns {#adaptive-load-scheduler-ins}

<!-- vale on -->

Input ports for the _Adaptive Load Scheduler_ component.

<dl>
<dt>enabled</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

The enabled port determines if the _Load Scheduler_ can shed _Flows_. By default, the Load Scheduler is enabled.

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
<dt>accepted_token_rate</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Accepted token rate is the number of tokens per second accepted by the service.

</dd>
<dt>desired_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Desired Load multiplier is the ratio of desired token rate to the incoming token rate.

</dd>
<dt>incoming_token_rate</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Incoming token rate is the number of tokens per second incoming to the service (including rejected ones).

</dd>
<dt>is_overload</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal that indicates whether the service is overloaded based on the deviation of the signal from the setpoint, considering a certain tolerance.
Deprecated: 1.6.0

</dd>
<dt>observed_load_multiplier</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Observed Load multiplier is the ratio of accepted token rate to the incoming token rate.

</dd>
</dl>

---

<!-- vale off -->

### AddressExtractor {#address-extractor}

<!-- vale on -->

Display an [Address][ext-authz-address] as a single string, for example, `<ip>:<port>`

IP addresses in attribute context are defined as objects with separate IP and port fields.
This is a helper to display an address as a single string.

:::caution

This might introduce high-cardinality flow label values.

:::

[ext-authz-address]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/address.proto#config-core-v3-address

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

Signal which Alerter is monitoring. If the signal greater than 0, Alerter generates an alert.

</dd>
</dl>

---

<!-- vale off -->

### AlerterParameters {#alerter-parameters}

<!-- vale on -->

Alerter Parameters configure parameters such as alert name, severity, resolve timeout, alert channels and labels.

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

Duration of alert resolver.

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

  Treating invalid inputs as "unknowns" has a consequence that the result
  might end up being valid even when some inputs are invalid. For example, `unknown && false == false`,
  because the result would end up false no matter if
  first signal was true or false. Conversely, `unknown && true == unknown`.

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

The arithmetic operation can be addition, subtraction, multiplication, division, XOR, right bit shift or left bit shift.
In case of XOR and bit shifts, value of signals is cast to integers before performing the operation.

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
<dt>pod_auto_scaler</dt>
<dd>

<!-- vale off -->

([PodAutoScaler](#pod-auto-scaler))

<!-- vale on -->

_PodAutoScaler_ provides auto-scaling functionality for scalable Kubernetes resource.

</dd>
<dt>pod_scaler</dt>
<dd>

<!-- vale off -->

([PodScaler](#pod-scaler))

<!-- vale on -->

PodScaler provides pod horizontal scaling functionality for scalable Kubernetes resources.

</dd>
</dl>

---

<!-- vale off -->

### AutoScaler {#auto-scaler}

<!-- vale on -->

_AutoScaler_ provides auto-scaling functionality for any scalable resource. Multiple _Controllers_ can be defined on the _AutoScaler_ for performing scale-out or scale-in. The _AutoScaler_ can interface with infrastructure APIs such as Kubernetes to perform auto-scale.

<dl>
<dt>cooldown_override_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `50`)

<!-- vale on -->

Cooldown override percentage defines a threshold change in scale-out beyond which previous cooldown is overridden.
For example, if the cooldown is 5 minutes and the cooldown override percentage is 10%, then if the
scale-increases by 10% or more, the previous cooldown is cancelled. Defaults to 50%.

</dd>
<dt>max_scale</dt>
<dd>

<!-- vale off -->

(string, default: `"9223372036854775807"`)

<!-- vale on -->

The maximum scale to which the _AutoScaler_ can scale-out. For example, in case of KubernetesReplicas Scaler, this is the maximum number of replicas.

</dd>
<dt>max_scale_in_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

The maximum decrease of scale (for example, pods) at one time. Defined as percentage of current scale value. Can never go below one even if percentage computation is less than one. Defaults to 1% of current scale value.

</dd>
<dt>max_scale_out_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `10`)

<!-- vale on -->

The maximum increase of scale (for example, pods) at one time. Defined as percentage of current scale value. Can never go below one even if percentage computation is less than one. Defaults to 10% of current scale value.

</dd>
<dt>min_scale</dt>
<dd>

<!-- vale off -->

(string, default: `"0"`)

<!-- vale on -->

The minimum scale to which the _AutoScaler_ can scale-in. For example, in case of KubernetesReplicas Scaler, this is the minimum number of replicas.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([AutoScalerOuts](#auto-scaler-outs))

<!-- vale on -->

Output ports for the _AutoScaler_.

</dd>
<dt>scale_in_alerter_parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for scale-in Alerter.

</dd>
<dt>scale_in_controllers</dt>
<dd>

<!-- vale off -->

([[]ScaleInController](#scale-in-controller))

<!-- vale on -->

List of _Controllers_ for scaling in.

</dd>
<dt>scale_in_cooldown</dt>
<dd>

<!-- vale off -->

(string, default: `"120s"`)

<!-- vale on -->

The amount of time to wait after a scale-in operation for another scale-in operation.

</dd>
<dt>scale_out_alerter_parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for scale-out Alerter.

</dd>
<dt>scale_out_controllers</dt>
<dd>

<!-- vale off -->

([[]ScaleOutController](#scale-out-controller))

<!-- vale on -->

List of _Controllers_ for scaling out.

</dd>
<dt>scale_out_cooldown</dt>
<dd>

<!-- vale off -->

(string, default: `"30s"`)

<!-- vale on -->

The amount of time to wait after a scale-out operation for another scale-out or scale-in operation.

</dd>
<dt>scaler</dt>
<dd>

<!-- vale off -->

([AutoScalerScaler](#auto-scaler-scaler))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### AutoScalerOuts {#auto-scaler-outs}

<!-- vale on -->

Outputs for _AutoScaler_.

<dl>
<dt>actual_scale</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
<dt>configured_scale</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
<dt>desired_scale</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### AutoScalerScaler {#auto-scaler-scaler}

<!-- vale on -->

<dl>
<dt>kubernetes_replicas</dt>
<dd>

<!-- vale off -->

([KubernetesReplicas](#kubernetes-replicas))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### Circuit {#circuit}

<!-- vale on -->

Circuit is graph of inter-connected signal processing components.

:::info

See also [Circuit overview](/concepts/policy/circuit.md).

:::

Signals flow between components through ports.
As signals traverse the circuit, they get processed, stored within components or get acted upon (for example, load-shed, rate-limit, auto-scale and so on).
Circuit is evaluated periodically to respond to changes in signal readings.

:::info Signals

Signals are floating point values.

A signal can also have a special **Invalid** value. It's usually used to
communicate that signal does not have a meaningful value at the moment, for example,
[PromQL](#prom-q-l) emits such a value if it cannot execute a query.
Components know when their input signals are invalid and can act
accordingly. They can either propagate the invalid signal, by making their
output itself invalid (for example,
[ArithmeticCombinator](#arithmetic-combinator)) or use some different
logic, for example, [Extrapolator](#extrapolator). Refer to a component's
docs on how exactly it handles invalid inputs.

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

(string, default: `"0.5s"`)

<!-- vale on -->

Evaluation interval (tick) is the time between consecutive runs of the policy circuit.
This interval is typically aligned with how often the corrective action (actuation) needs to be taken.

</dd>
</dl>

---

<!-- vale off -->

### Classifier {#classifier}

<!-- vale on -->

Set of classification rules sharing a common selector

:::info

See also [Classifier overview](/concepts/flow-control/resources/classifier.md).

:::
Example

```yaml
flow_selector:
  service_selector:
    agent_group: demoapp
    service: service1-demo-app.demoapp.svc.cluster.local
  flow_matcher:
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
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

Defines where to apply the flow classification rule.

</dd>
<dt>rego</dt>
<dd>

<!-- vale off -->

([Rego](#rego))

<!-- vale on -->

Rego based classification

Rego is a policy language used to express complex policies in a concise and declarative way.
It can be used to define flow classification rules by writing custom queries that extract values from request metadata.
For simple cases, such as directly reading a value from header or a field from JSON body, declarative extractors are recommended.

</dd>
<dt>rules</dt>
<dd>

<!-- vale off -->

(map of [Rule](#rule))

<!-- vale on -->

A map of {key, value} pairs mapping from
[flow label](/concepts/flow-control/flow-label.md) keys to rules that define
how to extract and propagate flow labels with that key.

</dd>
</dl>

---

<!-- vale off -->

### Component {#component}

<!-- vale on -->

Computational block that forms the circuit

:::info

See also [Components overview](/concepts/policy/circuit.md#components).

:::

Signals flow into the components from input ports and results are emitted on output ports.
Components are wired to each other based on signal names forming an execution graph of the circuit.

:::note

Loops are broken by the runtime at the earliest component index that is part of the loop.
The looped signals are saved in the tick they're generated and served in the subsequent tick.

:::

There are three categories of components:

- "source" components: they take some sort of input from "the real world" and output
  a signal based on this input. Example: [PromQL](#prom-q-l). In the UI
  they're represented by green color.
- signal processor components: processing components that do not interact with the external systems.
  Examples: [GradientController](#gradient-controller), [Max](#max).

  :::note

  Signal processor components' output can depend on their internal state, in addition to the inputs.
  Eg. see the [Exponential Moving Average filter](#e-m-a).

  :::

- "sink" components: they affect the real world.
  [_Concurrency Limiter_](#concurrency-limiter) and [_Rate Limiter_](#rate-limiter).
  In the UI, represented by orange color. Sink components usually come in pairs with a
  "sources" component which emits a feedback signal, like
  `accepted_concurrency` emitted by _Concurrency Limiter_.

:::tip

Sometimes you might want to use a constant value as one of component's inputs.
You can create an input port containing the constant value instead of being connected to a signal.
To do so, use the [InPort](#in_port)'s .withConstantSignal(constant_signal) method.
You can also use it to provide special math values such as NaN and +- Inf.
If You need to provide the same constant signal to multiple components,
You can use the [Variable](#variable) component.

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

Takes an input signal and emits the extrapolated value; either mirroring the input value or repeating the last known value up to the maximum extrapolation interval.

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

Gradient controller calculates the ratio between the signal and the setpoint to determine the magnitude of the correction that need to be applied.
This controller can be used to build AIMD (Additive Increase, Multiplicative Decrease) or MIMD style response.

</dd>
<dt>holder</dt>
<dd>

<!-- vale off -->

([Holder](#holder))

<!-- vale on -->

Holds the last valid signal value for the specified duration then waits for next valid value to hold.

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

Nested circuit defines a sub-circuit as a high-level component. It consists of a list of components and a map of input and output ports.

</dd>
<dt>nested_signal_egress</dt>
<dd>

<!-- vale off -->

([NestedSignalEgress](#nested-signal-egress))

<!-- vale on -->

Nested signal egress is a special type of component that allows to extract a signal from a nested circuit.

</dd>
<dt>nested_signal_ingress</dt>
<dd>

<!-- vale off -->

([NestedSignalIngress](#nested-signal-ingress))

<!-- vale on -->

Nested signal ingress is a special type of component that allows to inject a signal into a nested circuit.

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

Switcher acts as a switch that emits one of the two signals based on third signal.

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

Emits a variable signal which can be set to invalid.

</dd>
</dl>

---

<!-- vale off -->

### ConcurrencyLimiter {#concurrency-limiter}

<!-- vale on -->

_Concurrency Limiter_ is an actuator component that regulates flows to provide active service protection

It's based on the actuation strategy (for example, load actuator) and workload scheduling
which is based on Weighted Fair Queuing principles.
Concurrency is calculated in terms of total tokens per second, which can translate
to (avg. latency \* in-flight requests) (Little's Law) in concurrency limiting use-case.

ConcurrencyLimiter configuration is split into two parts: An actuation
strategy and a scheduler. At this time, only `load_actuator` strategy is available.

<dl>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

Flow Selector decides the service and flows at which the concurrency limiter is applied.

</dd>
<dt>load_actuator</dt>
<dd>

<!-- vale off -->

([LoadActuator](#load-actuator))

<!-- vale on -->

Actuator based on limiting the accepted concurrency under incoming concurrency \* load multiplier.

Actuation strategy defines the input signal that will drive the scheduler.

</dd>
<dt>scheduler</dt>
<dd>

<!-- vale off -->

([Scheduler](#scheduler))

<!-- vale on -->

Configuration of Weighted Fair Queuing-based workload scheduler.

Contains configuration of per-agent scheduler, and also defines some
output signals.

</dd>
</dl>

---

<!-- vale off -->

### ConstantSignal {#constant-signal}

<!-- vale on -->

Special constant input for ports and Variable component. Can provide either a constant value or special Nan/+-Inf value.

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

The comparison operator can be greater-than, less-than, greater-than-or-equal, less-than-or-equal, equal, or not-equal.

This component also supports time-based response (the output)
transitions between 1.0 or 0.0 signal if the decider condition is
true or false for at least `true_for` or `false_for` duration. If
`true_for` and `false_for` durations are zero then the transitions are
instantaneous.

<dl>
<dt>false_for</dt>
<dd>

<!-- vale off -->

(string, default: `"0s"`)

<!-- vale on -->

Duration of time to wait before changing to false state.
If the duration is zero, the change will happen instantaneously.

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

Duration of time to wait before changing to true state.
If the duration is zero, the change will happen instantaneously.```

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

Decreasing Gradient defines a controller for scaling in based on Gradient Controller.

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

This allows subset of parameters with constrained values compared to a regular gradient controller. For full documentation of these parameters, refer to the [GradientControllerParameters](#gradient-controller-parameters).

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

The window of time over which differentiator operates.

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

Exponential Moving Average (EMA) is a type of moving average that applies exponentially more weight to recent signal readings

At any time EMA component operates in one of the following states:

1. Warm up state: The first `warmup_window` samples are used to compute the initial EMA.
   If an invalid reading is received during the `warmup_window`, the last good average is emitted and the state gets reset back to beginning of warm up state.
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

The coefficient $\alpha$ represents the degree of weighting decrease, a constant smoothing factor between 0 and 1.
A higher $\alpha$ discounts older observations faster.
The $\alpha$ is computed using ema_window:

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

If the signal deviates from `max_envelope` faster than the correction
faster, it might end up exceeding the envelope.

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

Correction factor to apply on the output value if its in violation of the max envelope.

</dd>
<dt>correction_factor_on_min_envelope_violation</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

Correction factor to apply on the output value if its in violation of the min envelope.

</dd>
<dt>ema_window</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Duration of EMA sampling window.

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

The initial value of the EMA is the average of signal readings received during the warm up window.

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

Defines a high-level way to specify how to extract a flow label value given HTTP request metadata, without a need to write Rego code

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

- one of the fields of [Attribute Context](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto), or
- a special `request.http.bearer` pseudo-attribute.
  For example, `request.http.method` or `request.http.header.user-agent`

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

Extrapolates the input signal by repeating the last valid value during the period in which it is invalid

It does so until `maximum_extrapolation_interval` is reached, beyond which it emits invalid signal unless input signal becomes valid again.

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

Maximum time interval to repeat the last valid value of input signal.

</dd>
</dl>

---

<!-- vale off -->

### FirstValid {#first-valid}

<!-- vale on -->

Picks the first valid input signal from the array of input signals and emits it as an output signal

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

_Flow Control_ encompasses components that manage the flow of requests or access to features within a service.

<dl>
<dt>adaptive_load_scheduler</dt>
<dd>

<!-- vale off -->

([AdaptiveLoadScheduler](#adaptive-load-scheduler))

<!-- vale on -->

_Adaptive Load Scheduler_ component is based on additive increase and multiplicative decrease of token rate. It takes a signal and setpoint as inputs and reduces token rate proportionally (or any arbitrary power) based on deviation of the signal from setpoint.

</dd>
<dt>aimd_concurrency_controller</dt>
<dd>

<!-- vale off -->

([AIMDConcurrencyController](#a-i-m-d-concurrency-controller))

<!-- vale on -->

AIMD Concurrency control component is based on Additive Increase and Multiplicative Decrease of Concurrency. It takes a signal and setpoint as inputs and reduces concurrency limits proportionally (or any arbitrary power) based on deviation of the signal from setpoint. Internally implemented as a nested circuit.
Deprecated: 1.6.0

</dd>
<dt>concurrency_limiter</dt>
<dd>

<!-- vale off -->

([ConcurrencyLimiter](#concurrency-limiter))

<!-- vale on -->

_Concurrency Limiter_ provides service protection by applying prioritized load shedding of flows using a network scheduler (for example, Weighted Fair Queuing).
Deprecated: 1.6.0

</dd>
<dt>flow_regulator</dt>
<dd>

<!-- vale off -->

([FlowRegulator](#flow-regulator))

<!-- vale on -->

Flow Regulator is a component that regulates the flow of requests to the service by allowing only the specified percentage of requests or sticky sessions.
Deprecated: 1.6.0

</dd>
<dt>load_ramp</dt>
<dd>

<!-- vale off -->

([LoadRamp](#load-ramp))

<!-- vale on -->

_Load Ramp_ smoothly regulates the flow of requests over specified steps.

</dd>
<dt>load_ramp_series</dt>
<dd>

<!-- vale off -->

([LoadRampSeries](#load-ramp-series))

<!-- vale on -->

_Load Ramp Series_ is a series of _Load Ramp_ components that can shape load one after another at same or different _Control Points_.

</dd>
<dt>load_scheduler</dt>
<dd>

<!-- vale off -->

([LoadScheduler](#load-scheduler))

<!-- vale on -->

_Load Scheduler_ provides service protection by applying prioritized load shedding of flows using a network scheduler (for example, Weighted Fair Queuing).

</dd>
<dt>load_shaper</dt>
<dd>

<!-- vale off -->

([LoadShaper](#load-shaper))

<!-- vale on -->

_Load Shaper_ is a component that shapes the load at a _Control Point_.
Deprecated: 1.6.0

</dd>
<dt>load_shaper_series</dt>
<dd>

<!-- vale off -->

([LoadShaperSeries](#load-shaper-series))

<!-- vale on -->

_Load Shaper Series_ is a series of _Load Shaper_ components that can shape load one after another at same or different _Control Points_.
Deprecated: 1.6.0

</dd>
<dt>rate_limiter</dt>
<dd>

<!-- vale off -->

([RateLimiter](#rate-limiter))

<!-- vale on -->

_Rate Limiter_ provides service protection by applying rate limits.

</dd>
<dt>regulator</dt>
<dd>

<!-- vale off -->

([Regulator](#regulator))

<!-- vale on -->

Regulator is a component that regulates the flow of requests to the service by allowing only the specified percentage of requests or sticky sessions.

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

Classifiers are installed in the data-plane and are used to label the requests based on payload content.

The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.

</dd>
<dt>flux_meters</dt>
<dd>

<!-- vale off -->

(map of [FluxMeter](#flux-meter))

<!-- vale on -->

Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.

Flux Meter created metrics can be consumed as input to the circuit through the PromQL component.

</dd>
</dl>

---

<!-- vale off -->

### FlowMatcher {#flow-matcher}

<!-- vale on -->

Describes which flows a [flow control
component](/concepts/flow-control/flow-control.md#components) should apply
to

:::info

See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).

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
<dt>control_point</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

[Control Point](/concepts/flow-control/flow-selector.md#control-point)
identifies the location of a Flow within a Service. For an SDK based insertion, a Control Point can represent a particular feature or execution
block within a Service. In case of Service Mesh or Middleware insertion, a Control Point can identify ingress or egress calls or distinct listeners
or filter chains.

</dd>
<dt>label_matcher</dt>
<dd>

<!-- vale off -->

([LabelMatcher](#label-matcher))

<!-- vale on -->

Label matcher allows to add _additional_ condition on
[flow labels](/concepts/flow-control/flow-label.md)
must also be satisfied (in addition to service+control point matching)

:::info

See also [Label Matcher overview](/concepts/flow-control/flow-selector.md#label-matcher).

:::

:::note

[Classifiers](#classifier) _can_ use flow labels created by some other
classifier, but only if they were created at some previous control point
(and propagated in baggage).

This limitation does not apply to selectors of other entities, like
Flux Meters or Actuators. It's valid to create a flow label on a control
point using classifier, and immediately use it for matching on the same
control point.

:::

</dd>
</dl>

---

<!-- vale off -->

### FlowRegulator {#flow-regulator}

<!-- vale on -->

_Flow Regulator_ is a component that regulates the flow of requests to the service by allowing only the specified percentage of requests or sticky sessions.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([FlowRegulatorDynamicConfig](#flow-regulator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([FlowRegulatorIns](#flow-regulator-ins))

<!-- vale on -->

Input ports for the _Flow Regulator_.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([FlowRegulatorParameters](#flow-regulator-parameters))

<!-- vale on -->

Parameters for the _Flow Regulator_.

</dd>
</dl>

---

<!-- vale off -->

### FlowRegulatorDynamicConfig {#flow-regulator-dynamic-config}

<!-- vale on -->

Dynamic Configuration for _Flow Regulator_

<dl>
<dt>enable_label_values</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

Specify certain label values to be accepted by this flow filter regardless of accept percentage.

</dd>
</dl>

---

<!-- vale off -->

### FlowRegulatorIns {#flow-regulator-ins}

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

### FlowRegulatorParameters {#flow-regulator-parameters}

<!-- vale on -->

<dl>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

_Flow Selector_ selects the _Flows_ at which the _Flow Regulator_ is applied.

</dd>
<dt>label_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

The flow label key for identifying sessions.

- When label key is specified, _Flow Regulator_ acts as a sticky filter.
  The series of flows with the same value of label key get the same
  decision provided that the `accept_percentage` is same or higher.
- When label key is not specified, _Flow Regulator_ acts as a stateless filter.
  Percentage of flows are selected randomly for rejection.

</dd>
</dl>

---

<!-- vale off -->

### FlowSelector {#flow-selector}

<!-- vale on -->

Selects flows based on _Control Point_, flow labels, agent group and service that the [flow control
component](/concepts/flow-control/flow-control.md#components) operates on.

:::info

See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).

:::

<dl>
<dt>flow_matcher</dt>
<dd>

<!-- vale off -->

([FlowMatcher](#flow-matcher))

<!-- vale on -->

Match control points and labels

</dd>
<dt>service_selector</dt>
<dd>

<!-- vale off -->

([ServiceSelector](#service-selector))

<!-- vale on -->

Match agent group and service

</dd>
</dl>

---

<!-- vale off -->

### FluxMeter {#flux-meter}

<!-- vale on -->

Flux Meter gathers metrics for the traffic that matches its selector.
The histogram created by Flux Meter measures the workload latency by default.

:::info

See also [Flux Meter overview](/concepts/flow-control/resources/flux-meter.md).

:::
Example:

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
flow_selector:
  service_selector:
    agent_group: demoapp
    service: service1-demo-app.demoapp.svc.cluster.local
  flow_matcher:
    control_point: ingress
attribute_key: response_duration_ms
```

<dl>
<dt>attribute_key</dt>
<dd>

<!-- vale off -->

(string, default: `"workload_duration_ms"`)

<!-- vale on -->

Key of the attribute in access log or span from which the metric for this flux meter is read.

:::info

For list of available attributes in Envoy access logs, refer
[Envoy Filter](/get-started/integrations/flow-control/envoy/istio.md#envoy-filter)

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
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

The selection criteria for the traffic that will be measured.

</dd>
<dt>linear_buckets</dt>
<dd>

<!-- vale off -->

([FluxMeterLinearBuckets](#flux-meter-linear-buckets))

<!-- vale on -->

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

ExponentialBuckets creates `count` number of buckets where the lowest bucket has an upper bound of `start`
and each following bucket's upper bound is `factor` times the previous bucket's upper bound. The final +inf
bucket is not counted.

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

Factor to be multiplied to the previous bucket's upper bound to calculate the following bucket's upper bound.

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

ExponentialBucketsRange creates `count` number of buckets where the lowest bucket is `min` and the highest
bucket is `max`. The final +inf bucket is not counted.

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

LinearBuckets creates `count` number of buckets, each `width` wide, where the lowest bucket has an
upper bound of `start`. The final +inf bucket is not counted.

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

StaticBuckets holds the static value of the buckets where latency histogram will be stored.

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

Gradient controller is a type of controller which tries to adjust the
control variable proportionally to the relative difference between setpoint
and actual value of the signal

The `gradient` describes a corrective factor that should be applied to the
control variable to get the signal closer to the setpoint. It's computed as follows:

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

The output can be _optionally_ clamped to desired range using `max` and
`min` input.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([GradientControllerDynamicConfig](#gradient-controller-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([GradientControllerIns](#gradient-controller-ins))

<!-- vale on -->

Input ports of the Gradient Controller.

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

### GradientControllerDynamicConfig {#gradient-controller-dynamic-config}

<!-- vale on -->

Dynamic Configuration for a Controller

<dl>
<dt>manual_mode</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Decides whether the controller runs in `manual_mode`.
In manual mode, the controller does not adjust the control variable It emits the same output as the control variable input.

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
<dt>optimize</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Optimize signal is added to the output of the gradient calculation.

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

Maximum gradient which clamps the computed gradient value to the range, `[min_gradient, max_gradient]`.

</dd>
<dt>min_gradient</dt>
<dd>

<!-- vale off -->

(float64, default: `-1.7976931348623157e+308`)

<!-- vale on -->

Minimum gradient which clamps the computed gradient value to the range, `[min_gradient, max_gradient]`.

</dd>
<dt>slope</dt>
<dd>

<!-- vale off -->

(float64, **required**)

<!-- vale on -->

Slope controls the aggressiveness and direction of the Gradient Controller.

Slope is used as exponent on the signal to setpoint ratio in computation
of the gradient (see the [main description](#gradient-controller) for
exact equation). This parameter decides how aggressive the controller
responds to the deviation of signal from the setpoint.
for example:

- $\text{slope} = 1$: when signal is too high, increase control variable,
- $\text{slope} = -1$: when signal is too high, decrease control variable,
- $\text{slope} = -0.5$: when signal is too high, decrease control variable gradually.

The sign of slope depends on correlation between the signal and control variable:

- Use $\text{slope} < 0$ if there is a _positive_ correlation between the signal and
  the control variable (for example, Per-pod CPU usage and total concurrency).
- Use $\text{slope} > 0$ if there is a _negative_ correlation between the signal and
  the control variable (for example, Per-pod CPU usage and number of pods).

:::note

You need to set _negative_ slope for a _positive_ correlation, as you're
describing the _action_ which controller should make when the signal
increases.

:::

The magnitude of slope describes how aggressively should the controller
react to a deviation of signal.
With $|\text{slope}| = 1$, the controller will aim to bring the signal to
the setpoint in one tick (assuming linear correlation with signal and setpoint).
Smaller magnitudes of slope will make the controller adjust the control
variable gradually.

Setting $|\text{slope}| < 1$ (for example, $\pm0.8$) is recommended.
If you experience overshooting, consider lowering the magnitude even more.
Values of $|\text{slope}| > 1$ aren't recommended.

:::note

Remember that the gradient and output signal can be (optionally) clamped,
so the _slope_ might not fully describe aggressiveness of the controller.

:::

</dd>
</dl>

---

<!-- vale off -->

### Holder {#holder}

<!-- vale on -->

Holds the last valid signal value for the specified duration then waits for next valid value to hold.
If it is holding a value that means it ignores both valid and invalid new signals until the `hold_for` duration is finished.

<dl>
<dt>hold_for</dt>
<dd>

<!-- vale off -->

(string, default: `"5s"`)

<!-- vale on -->

Holding the last valid signal value for the `hold_for` duration.

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

Resets the holder output to the current input signal when reset signal is valid and non-zero.

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

Increasing Gradient defines a controller for scaling out based on Gradient Controller.

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

This allows subset of parameters with constrained values compared to a regular gradient controller. For full documentation of these parameters, refer to the [GradientControllerParameters](#gradient-controller-parameters).

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

Resets the integrator output to zero when reset signal is valid and non-zero. Reset also resets the max and min constraints.

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

JSON pointer represents a parsed JSON pointer which allows to select a specified field from the payload.

Note: Uses [JSON pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
for example, `/foo/bar`. If the pointer points into an object, it'd be converted to a string.

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

JWT (JSON Web Token) can be extracted from any input attribute, but most likely you'd want to use `request.http.bearer`.

</dd>
<dt>json_pointer</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

JSON pointer allowing to select a specified field from the payload.

Note: Uses [JSON pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
for example, `/foo/bar`. If the pointer points into an object, it'd be converted to a string.

</dd>
</dl>

---

<!-- vale off -->

### K8sLabelMatcherRequirement {#k8s-label-matcher-requirement}

<!-- vale on -->

Label selector requirement which is a selector that contains values, a key, and an operator that relates the key and values.

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

Logical operator which represents a key's relationship to a set of values.
Valid operators are In, NotIn, Exists and DoesNotExist.

</dd>
<dt>values</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

An array of string values that relates to the key by an operator.
If the operator is In or NotIn, the values array must be non-empty.
If the operator is Exists or DoesNotExist, the values array must be empty.

</dd>
</dl>

---

<!-- vale off -->

### KubernetesObjectSelector {#kubernetes-object-selector}

<!-- vale on -->

Describes which pods a control or observability
component should apply to.

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

Which [agent-group](/concepts/flow-control/flow-selector.md#agent-group) this
selector applies to.

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

### KubernetesReplicas {#kubernetes-replicas}

<!-- vale on -->

KubernetesReplicas defines a horizontal pod scaler for Kubernetes.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([PodScalerScaleActuatorDynamicConfig](#pod-scaler-scale-actuator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig

</dd>
<dt>kubernetes_object_selector</dt>
<dd>

<!-- vale off -->

([KubernetesObjectSelector](#kubernetes-object-selector))

<!-- vale on -->

The Kubernetes object on which horizontal scaling is applied.

</dd>
</dl>

---

<!-- vale off -->

### LabelMatcher {#label-matcher}

<!-- vale on -->

Allows to define rules whether a map of
[labels](/concepts/flow-control/flow-label.md)
should be considered a match or not

It provides three ways to define requirements:

- match labels
- match expressions
- arbitrary expression

If multiple requirements are set, they're all combined using the logical AND operator.
An empty label matcher always matches.

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

A map of {key,value} pairs representing labels to be matched.
A single {key,value} in the `match_labels` requires that the label `key` is present and equal to `value`.

Note: The requirements are combined using the logical AND operator.

</dd>
</dl>

---

<!-- vale off -->

### LoadActuator {#load-actuator}

<!-- vale on -->

Takes the load multiplier input signal and publishes it to the schedulers in the data-plane

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([LoadActuatorDynamicConfig](#load-actuator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadActuatorIns](#load-actuator-ins))

<!-- vale on -->

Input ports for the Load Actuator component.

</dd>
</dl>

---

<!-- vale off -->

### LoadActuatorDynamicConfig {#load-actuator-dynamic-config}

<!-- vale on -->

Dynamic Configuration for LoadActuator

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

Decides whether to run the load actuator in dry-run mode. Dry run mode ensures that no traffic gets dropped by this load actuator.
Useful for observing the behavior of Load Actuator without disrupting any real traffic.

</dd>
</dl>

---

<!-- vale off -->

### LoadActuatorIns {#load-actuator-ins}

<!-- vale on -->

Input for the Load Actuator component.

<dl>
<dt>load_multiplier</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Load multiplier is proportion of [incoming
token rate](#scheduler-outs) that needs to be accepted.

</dd>
</dl>

---

<!-- vale off -->

### LoadRamp {#load-ramp}

<!-- vale on -->

The _Load Ramp_ produces a smooth and continuous traffic load
that changes progressively over time, based on the specified steps.

Each step is defined by two parameters:

- The `target_accept_percentage`.
- The `duration` for the signal to change from the
  previous step's `target_accept_percentage` to the current step's
  `target_accept_percentage`.

The percentage of requests accepted starts at the `target_accept_percentage`
defined in the first step and gradually ramps up or down linearly from
the previous step's `target_accept_percentage` to the next
`target_accept_percentage`, over the `duration` specified for each step.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([RegulatorDynamicConfig](#regulator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Dynamic configuration key for flow regulator.

</dd>
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

A Boolean signal indicating whether the _Load Ramp_ is at the end of signal generation.

</dd>
<dt>at_start</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Load Ramp_ is at the start of signal generation.

</dd>
</dl>

---

<!-- vale off -->

### LoadRampParameters {#load-ramp-parameters}

<!-- vale on -->

Parameters for the _Load Ramp_ component.

<dl>
<dt>regulator_parameters</dt>
<dd>

<!-- vale off -->

([RegulatorParameters](#regulator-parameters))

<!-- vale on -->

Parameters for the _Regulator_.

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

Duration for which the step is active.

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

### LoadRampSeries {#load-ramp-series}

<!-- vale on -->

_LoadRampSeries_ is a component that applies a series of _Load Ramps_ in order.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadRampSeriesIns](#load-ramp-series-ins))

<!-- vale on -->

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadRampSeriesParameters](#load-ramp-series-parameters))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadRampSeriesIns {#load-ramp-series-ins}

<!-- vale on -->

Inputs for the _LoadRampSeries_ component.

<dl>
<dt>backward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the load ramp series towards the previous step.

</dd>
<dt>forward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the load ramp series towards the next step.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to reset the load ramp series to the first step.

</dd>
</dl>

---

<!-- vale off -->

### LoadRampSeriesLoadRampInstance {#load-ramp-series-load-ramp-instance}

<!-- vale on -->

<dl>
<dt>load_ramp</dt>
<dd>

<!-- vale off -->

([LoadRampParameters](#load-ramp-parameters))

<!-- vale on -->

The load ramp.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([LoadRampOuts](#load-ramp-outs))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadRampSeriesParameters {#load-ramp-series-parameters}

<!-- vale on -->

Parameters for the _LoadRampSeries_ component.

<dl>
<dt>load_ramps</dt>
<dd>

<!-- vale off -->

([[]LoadRampSeriesLoadRampInstance](#load-ramp-series-load-ramp-instance), **required**)

<!-- vale on -->

An ordered list of load ramps that get applied in order.

</dd>
</dl>

---

<!-- vale off -->

### LoadScheduler {#load-scheduler}

<!-- vale on -->

_Load Scheduler_ is an actuator component that regulates flows to provide active service protection

:::info

See also [_Load Scheduler_ overview](/concepts/flow-control/components/load-scheduler.md).

:::

It's based on the actuation strategy (for example, load actuator) and workload scheduling
which is based on Weighted Fair Queuing principles.
It measures and controls the incoming tokens per second, which can translate
to (avg. latency \* in-flight requests) (Little's Law) in concurrency limiting use-case.

LoadScheduler configuration is split into two parts: An actuation
strategy and a scheduler. At this time, only `load_actuator` strategy is available.

<dl>
<dt>actuator</dt>
<dd>

<!-- vale off -->

([LoadSchedulerActuator](#load-scheduler-actuator))

<!-- vale on -->

Actuator based on limiting the accepted token rate under incoming token rate \* load multiplier.

</dd>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

Flow Selector decides the service and flows at which the _Load Scheduler_ is applied.

</dd>
<dt>scheduler</dt>
<dd>

<!-- vale off -->

([LoadSchedulerScheduler](#load-scheduler-scheduler))

<!-- vale on -->

Configuration of Weighted Fair Queuing-based workload scheduler.

Contains configuration of per-agent scheduler, and also defines some
output signals.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerActuator {#load-scheduler-actuator}

<!-- vale on -->

Takes the load multiplier input signal and publishes it to the schedulers in the data-plane

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([LoadSchedulerActuatorDynamicConfig](#load-scheduler-actuator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadSchedulerActuatorIns](#load-scheduler-actuator-ins))

<!-- vale on -->

Input ports for the Actuator component.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerActuatorDynamicConfig {#load-scheduler-actuator-dynamic-config}

<!-- vale on -->

Dynamic Configuration for Actuator

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool)

<!-- vale on -->

Decides whether to run the actuator in dry-run mode. Dry run mode ensures that no traffic gets dropped by this actuator.
Useful for observing the behavior of actuator without disrupting any real traffic.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerActuatorIns {#load-scheduler-actuator-ins}

<!-- vale on -->

Input for the Actuator component.

<dl>
<dt>load_multiplier</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Load multiplier is proportion of [incoming
token rate](#scheduler-outs) that needs to be accepted.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerScheduler {#load-scheduler-scheduler}

<!-- vale on -->

<dl>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([LoadSchedulerSchedulerOuts](#load-scheduler-scheduler-outs))

<!-- vale on -->

Output ports for the Scheduler component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadSchedulerSchedulerParameters](#load-scheduler-scheduler-parameters))

<!-- vale on -->

Scheduler parameters.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerSchedulerOuts {#load-scheduler-scheduler-outs}

<!-- vale on -->

Output for the Scheduler component.

<dl>
<dt>accepted_token_rate</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Accepted token rate is the tokens admitted per second by the scheduler.
Value of this signal is aggregated from all the relevant schedulers.

</dd>
<dt>incoming_token_rate</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Incoming token rate is the incoming tokens per second for all the
flows entering the scheduler including the rejected ones.

This is computed similar to `accepted_token_rate`,
by summing up tokens from all the flows entering scheduler.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerSchedulerParameters {#load-scheduler-scheduler-parameters}

<!-- vale on -->

Scheduler parameters

<dl>
<dt>auto_tokens</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Automatically estimate the size of a flow in each workload, based on
historical latency. Each workload's `tokens` will be set to average
latency of flows in that workload during last few seconds (exact duration
of this average can change).
This setting is useful in concurrency limiting use-case, where the
concurrency is calculated as (avg. latency \* in-flight flows).

The value of tokens estimated by `auto_tokens` takes lower precedence
than the value of `tokens` specified in the workload definition
and `tokens` explicitly specified in the flow labels.

</dd>
<dt>decision_deadline_margin</dt>
<dd>

<!-- vale off -->

(string, default: `"0.01s"`)

<!-- vale on -->

Decision deadline margin is the amount of time that the scheduler will
subtract from the request deadline to determine the deadline for the
decision. This is to ensure that the scheduler has enough time to
make a decision before the request deadline happens, accounting for
processing delays.
The request deadline is based on the
[gRPC deadline](https://grpc.io/blog/deadlines) or the
[`grpc-timeout` HTTP header](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests).

Fail-open logic is use for flow control APIs, so if the gRPC deadline
reaches, the flow will end up being unconditionally allowed while
it is still waiting on the scheduler.

</dd>
<dt>default_workload_parameters</dt>
<dd>

<!-- vale off -->

([LoadSchedulerSchedulerWorkloadParameters](#load-scheduler-scheduler-workload-parameters))

<!-- vale on -->

Parameters to be used if none of workloads specified in `workloads` match.

</dd>
<dt>max_timeout</dt>
<dd>

<!-- vale off -->

(string, default: `"0s"`)

<!-- vale on -->

Deprecated: 1.5.0. Use `decision_deadline_margin` instead. This value is ignored.

</dd>
<dt>timeout_factor</dt>
<dd>

<!-- vale off -->

(float64, default: `0`)

<!-- vale on -->

Deprecated: 1.5.0. Use `decision_deadline_margin` instead. This value is ignored.

</dd>
<dt>tokens_label_key</dt>
<dd>

<!-- vale off -->

(string, default: `"tokens"`)

<!-- vale on -->

- Key for a flow label that can be used to override the default number of tokens for this flow.
- The value associated with this key must be a valid uint64 number.
- If this parameter is not provided, the number of tokens for the flow will be determined by the matched workload's token count.

</dd>
<dt>workloads</dt>
<dd>

<!-- vale off -->

([[]LoadSchedulerSchedulerWorkload](#load-scheduler-scheduler-workload))

<!-- vale on -->

List of workloads to be used in scheduler.

Categorizing [flows](/concepts/flow-control/flow-control.md#flow) into workloads
allows for load-shedding to be "intelligent" compared to random rejections.
There are two aspects of this "intelligence":

- Scheduler can more precisely calculate concurrency if it understands
  that flows belonging to different classes have different weights (for example, insert queries compared to select queries).
- Setting different priorities to different workloads lets the scheduler
  avoid dropping important traffic during overload.

Each workload in this list specifies also a matcher that is used to
determine which flow will be categorized into which workload.
In case of multiple matching workloads, the first matching one will be used.
If none of workloads match, `default_workload` will be used.

:::info

See also [workload definition in the concepts
section](/concepts/flow-control/components/load-scheduler.md#workload).

:::

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerSchedulerWorkload {#load-scheduler-scheduler-workload}

<!-- vale on -->

Workload defines a class of flows that preferably have similar properties such as response latency and desired priority.

<dl>
<dt>label_matcher</dt>
<dd>

<!-- vale off -->

([LabelMatcher](#label-matcher))

<!-- vale on -->

Label Matcher to select a Workload based on
[flow labels](/concepts/flow-control/flow-label.md).

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadSchedulerSchedulerWorkloadParameters](#load-scheduler-scheduler-workload-parameters))

<!-- vale on -->

Parameters associated with flows matching the label matcher.

</dd>
</dl>

---

<!-- vale off -->

### LoadSchedulerSchedulerWorkloadParameters {#load-scheduler-scheduler-workload-parameters}

<!-- vale on -->

Parameters such as priority, tokens and fairness key that
are applicable to flows within a workload.

<dl>
<dt>fairness_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Fairness key is a label key that can be used to provide fairness within a workload.
Any [flow label](/concepts/flow-control/flow-label.md) can be used here. For example, if
you have a classifier that sets `user` flow label, you might want to set
`fairness_key = "user"`.

</dd>
<dt>priority</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, maximum: `255`, default: `0`)

<!-- vale on -->

Describes priority level of the flows within the workload.
Priority level ranges from 0 to 255.
Higher numbers means higher priority level.
Priority levels have non-linear effect on the workload scheduling. The following formula is used to determine the position of a flow in the queue based on virtual finish time:

$$
\text{virtual\_finish\_time} = \text{virtual\_time} + \left(\text{tokens} \cdot \left(\text{256} - \text{priority}\right)\right)
$$

</dd>
<dt>tokens</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Tokens determines the cost of admitting a single flow in the workload,
which is typically defined as milliseconds of flow latency (time to response or duration of a feature) or
simply equal to 1 if the resource being accessed is constrained by the
number of flows (3rd party rate limiters).
This override is applicable only if tokens for the flow aren't specified
in the flow labels.

</dd>
</dl>

---

<!-- vale off -->

### LoadShaper {#load-shaper}

<!-- vale on -->

The _Load Shaper_ produces a smooth and continuous traffic load
that changes progressively over time, based on the specified steps.

Each step is defined by two parameters:

- The `target_accept_percentage`.
- The `duration` for the signal to change from the
  previous step's `target_accept_percentage` to the current step's
  `target_accept_percentage`.

The percentage of requests accepted starts at the `target_accept_percentage`
defined in the first step and gradually ramps up or down linearly from
the previous step's `target_accept_percentage` to the next
`target_accept_percentage`, over the `duration` specified for each step.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([RegulatorDynamicConfig](#regulator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Dynamic configuration key for flow regulator.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadShaperIns](#load-shaper-ins))

<!-- vale on -->

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([LoadShaperOuts](#load-shaper-outs))

<!-- vale on -->

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadShaperParameters](#load-shaper-parameters))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperIns {#load-shaper-ins}

<!-- vale on -->

Inputs for the _Load Shaper_ component.

<dl>
<dt>backward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the _Load Shaper_ towards the previous step.

</dd>
<dt>forward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the _Load Shaper_ towards the next step.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to reset the _Load Shaper_ to the first step.

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperOuts {#load-shaper-outs}

<!-- vale on -->

Outputs for the _Load Shaper_ component.

<dl>
<dt>accept_percentage</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

The percentage of flows being accepted by the _Load Shaper_.

</dd>
<dt>at_end</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Load Shaper_ is at the end of signal generation.

</dd>
<dt>at_start</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Load Shaper_ is at the start of signal generation.

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperParameters {#load-shaper-parameters}

<!-- vale on -->

Parameters for the _Load Shaper_ component.

<dl>
<dt>flow_regulator_parameters</dt>
<dd>

<!-- vale off -->

([FlowRegulatorParameters](#flow-regulator-parameters))

<!-- vale on -->

Parameters for the _Flow Regulator_.

</dd>
<dt>steps</dt>
<dd>

<!-- vale off -->

([[]LoadShaperParametersStep](#load-shaper-parameters-step), **required**)

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperParametersStep {#load-shaper-parameters-step}

<!-- vale on -->

<dl>
<dt>duration</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Duration for which the step is active.

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

### LoadShaperSeries {#load-shaper-series}

<!-- vale on -->

_LoadShaperSeries_ is a component that applies a series of _Load Shapers_ in order.

<dl>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([LoadShaperSeriesIns](#load-shaper-series-ins))

<!-- vale on -->

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([LoadShaperSeriesParameters](#load-shaper-series-parameters))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperSeriesIns {#load-shaper-series-ins}

<!-- vale on -->

Inputs for the _LoadShaperSeries_ component.

<dl>
<dt>backward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the load shaper series towards the previous step.

</dd>
<dt>forward</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to progress the load shaper series towards the next step.

</dd>
<dt>reset</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Whether to reset the load shaper series to the first step.

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperSeriesLoadShaperInstance {#load-shaper-series-load-shaper-instance}

<!-- vale on -->

<dl>
<dt>load_shaper</dt>
<dd>

<!-- vale off -->

([LoadShaperParameters](#load-shaper-parameters))

<!-- vale on -->

The load shaper.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([LoadShaperOuts](#load-shaper-outs))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### LoadShaperSeriesParameters {#load-shaper-series-parameters}

<!-- vale on -->

Parameters for the _LoadShaperSeries_ component.

<dl>
<dt>load_shapers</dt>
<dd>

<!-- vale off -->

([[]LoadShaperSeriesLoadShaperInstance](#load-shaper-series-load-shaper-instance), **required**)

<!-- vale on -->

An ordered list of load shapers that get applied in order.

</dd>
</dl>

---

<!-- vale off -->

### MatchExpression {#match-expression}

<!-- vale on -->

Defines a `[map<string, string> → bool]` expression to be evaluated on labels

MatchExpression has multiple variants, exactly one should be set.

Example:

```yaml
all:
  of:
    - label_exists: foo
    - label_equals: { label = app, value = frobnicator }
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

(string, **required**)

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

Regular expression that should match the label value.
It uses [Go's regular expression syntax](https://github.com/google/re2/wiki/Syntax).

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

Takes an array of input signals and emits the signal with the minimum value
Min: output = min([]inputs).

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

Nested circuit defines a sub-circuit as a high-level component. It consists of a list of components and a map of input and output ports.

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

Name of the nested circuit component. This name is displayed by graph visualization tools.

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

Short description of the nested circuit component. This description is displayed by graph visualization tools.

</dd>
</dl>

---

<!-- vale off -->

### NestedSignalEgress {#nested-signal-egress}

<!-- vale on -->

Nested signal egress is a special type of component that allows to extract a signal from a nested circuit.

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

Nested signal ingress is a special type of component that allows to inject a signal into a nested circuit.

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

HTTP path will be matched against given path templates.
If a match occurs, the value associated with the path template will be treated as a result.
In case of multiple path templates matching, the most specific one will be chosen.

<dl>
<dt>template_values</dt>
<dd>

<!-- vale off -->

(map of string)

<!-- vale on -->

Template value keys are OpenAPI-inspired path templates.

- Static path segment `/foo` matches a path segment exactly
- `/{param}` matches arbitrary path segment.
  (The parameter name is ignored and can be omitted (`{}`))
- The parameter must cover whole segment.
- Additionally, path template can end with `/*` wildcard to match
  arbitrary number of trailing segments (0 or more).
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

### PodAutoScaler {#pod-auto-scaler}

<!-- vale on -->

_PodAutoScaler_ provides auto-scaling functionality for scalable Kubernetes resource. Multiple _Controllers_ can be defined on the _PodAutoScaler_ for performing scale-out or scale-in. The _PodAutoScaler_ interfaces with Kubernetes infrastructure APIs to perform auto-scale.

<dl>
<dt>cooldown_override_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `50`)

<!-- vale on -->

Cooldown override percentage defines a threshold change in scale-out beyond which previous cooldown is overridden.
For example, if the cooldown is 5 minutes and the cooldown override percentage is 10%, then if the
scale-increases by 10% or more, the previous cooldown is cancelled. Defaults to 50%.

</dd>
<dt>max_replicas</dt>
<dd>

<!-- vale off -->

(string, default: `"9223372036854775807"`)

<!-- vale on -->

The maximum scale to which the _PodAutoScaler_ can scale-out.

</dd>
<dt>max_scale_in_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

The maximum decrease of replicas (for example, pods) at one time. Defined as percentage of current scale value. Can never go below one even if percentage computation is less than one. Defaults to 1% of current scale value.

</dd>
<dt>max_scale_out_percentage</dt>
<dd>

<!-- vale off -->

(float64, default: `10`)

<!-- vale on -->

The maximum increase of replicas (for example, pods) at one time. Defined as percentage of current scale value. Can never go below one even if percentage computation is less than one. Defaults to 10% of current scale value.

</dd>
<dt>min_replicas</dt>
<dd>

<!-- vale off -->

(string, default: `"0"`)

<!-- vale on -->

The minimum replicas to which the _PodAutoScaler_ can scale-in.

</dd>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([PodAutoScalerOuts](#pod-auto-scaler-outs))

<!-- vale on -->

Output ports for the _PodAutoScaler_.

</dd>
<dt>pod_scaler</dt>
<dd>

<!-- vale off -->

([KubernetesReplicas](#kubernetes-replicas))

<!-- vale on -->

</dd>
<dt>scale_in_alerter_parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for scale-in Alerter.

</dd>
<dt>scale_in_controllers</dt>
<dd>

<!-- vale off -->

([[]ScaleInController](#scale-in-controller))

<!-- vale on -->

List of _Controllers_ for scaling in.

</dd>
<dt>scale_in_cooldown</dt>
<dd>

<!-- vale off -->

(string, default: `"120s"`)

<!-- vale on -->

The amount of time to wait after a scale-in operation for another scale-in operation.

</dd>
<dt>scale_out_alerter_parameters</dt>
<dd>

<!-- vale off -->

([AlerterParameters](#alerter-parameters))

<!-- vale on -->

Configuration for scale-out Alerter.

</dd>
<dt>scale_out_controllers</dt>
<dd>

<!-- vale off -->

([[]ScaleOutController](#scale-out-controller))

<!-- vale on -->

List of _Controllers_ for scaling out.

</dd>
<dt>scale_out_cooldown</dt>
<dd>

<!-- vale off -->

(string, default: `"30s"`)

<!-- vale on -->

The amount of time to wait after a scale-out operation for another scale-out or scale-in operation.

</dd>
</dl>

---

<!-- vale off -->

### PodAutoScalerOuts {#pod-auto-scaler-outs}

<!-- vale on -->

Outputs for _PodAutoScaler_.

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

### PodScaler {#pod-scaler}

<!-- vale on -->

Component for scaling pods based on a signal.

<dl>
<dt>kubernetes_object_selector</dt>
<dd>

<!-- vale off -->

([KubernetesObjectSelector](#kubernetes-object-selector))

<!-- vale on -->

The Kubernetes object on which horizontal scaling is applied.

</dd>
<dt>scale_actuator</dt>
<dd>

<!-- vale off -->

([PodScalerScaleActuator](#pod-scaler-scale-actuator))

<!-- vale on -->

Actuates scaling of pods based on a signal.

</dd>
<dt>scale_reporter</dt>
<dd>

<!-- vale off -->

([PodScalerScaleReporter](#pod-scaler-scale-reporter))

<!-- vale on -->

Reports actual and configured number of replicas.

</dd>
</dl>

---

<!-- vale off -->

### PodScalerScaleActuator {#pod-scaler-scale-actuator}

<!-- vale on -->

Actuates scaling of pods based on a signal.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([PodScalerScaleActuatorDynamicConfig](#pod-scaler-scale-actuator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([PodScalerScaleActuatorIns](#pod-scaler-scale-actuator-ins))

<!-- vale on -->

Input ports for the PodScaler component.

</dd>
</dl>

---

<!-- vale off -->

### PodScalerScaleActuatorDynamicConfig {#pod-scaler-scale-actuator-dynamic-config}

<!-- vale on -->

Dynamic Configuration for ScaleActuator

<dl>
<dt>dry_run</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Decides whether to run the pod scaler in dry-run mode. Dry run mode ensures that no scaling is invoked by this pod scaler.
Useful for observing the behavior of Scaler without disrupting any real traffic.

</dd>
</dl>

---

<!-- vale off -->

### PodScalerScaleActuatorIns {#pod-scaler-scale-actuator-ins}

<!-- vale on -->

Inputs for the PodScaler component.

<dl>
<dt>desired_replicas</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

</dd>
</dl>

---

<!-- vale off -->

### PodScalerScaleReporter {#pod-scaler-scale-reporter}

<!-- vale on -->

Reports actual and configured number of replicas.

<dl>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([PodScalerScaleReporterOuts](#pod-scaler-scale-reporter-outs))

<!-- vale on -->

Output ports for the PodScaler component.

</dd>
</dl>

---

<!-- vale off -->

### PodScalerScaleReporterOuts {#pod-scaler-scale-reporter-outs}

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

See also [Policy overview](/concepts/policy/policy.md).

:::

Policy specification contains a circuit that defines the controller logic and resources that need to be setup.

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

Component that runs a Prometheus query periodically and returns the result as an output signal

<dl>
<dt>evaluation_interval</dt>
<dd>

<!-- vale off -->

(string, default: `"10s"`)

<!-- vale on -->

Describes the interval between successive evaluations of the Prometheus query.

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

Describes the [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) query to be run.

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

Aperture supports OpenTelemetry metrics. See [reference](/get-started/integrations/metrics/metrics.md) for more details.

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

Emitting 0 for the `false_for` duration.

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

Emitting 1 for the `true_for` duration.

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

### RateLimiter {#rate-limiter}

<!-- vale on -->

Limits the traffic on a control point to specified rate

:::info

See also [_Rate Limiter_ overview](/concepts/flow-control/components/rate-limiter.md).

:::

RateLimiting is done on per-label-value basis. Use `label_key`
to select which label should be used as key.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([RateLimiterDynamicConfig](#rate-limiter-dynamic-config))

<!-- vale on -->

Default configuration

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig

</dd>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

Which control point to apply this rate limiter to.

</dd>
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
</dl>

---

<!-- vale off -->

### RateLimiterDynamicConfig {#rate-limiter-dynamic-config}

<!-- vale on -->

Dynamic Configuration for the rate limiter

<dl>
<dt>overrides</dt>
<dd>

<!-- vale off -->

([[]RateLimiterOverride](#rate-limiter-override))

<!-- vale on -->

Allows to specify different limits for particular label values.

</dd>
</dl>

---

<!-- vale off -->

### RateLimiterIns {#rate-limiter-ins}

<!-- vale on -->

Inputs for the RateLimiter component

<dl>
<dt>limit</dt>
<dd>

<!-- vale off -->

([InPort](#in-port))

<!-- vale on -->

Number of flows allowed per `limit_reset_interval` per each label.
Negative values disable the rate limiter.

:::tip

Negative limit can be useful to _conditionally_ enable the rate limiter
under certain circumstances. [Decider](#decider) might be helpful.

:::

</dd>
</dl>

---

<!-- vale off -->

### RateLimiterOverride {#rate-limiter-override}

<!-- vale on -->

<dl>
<dt>label_value</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Value of the label for which the override should be applied.

</dd>
<dt>limit_scale_factor</dt>
<dd>

<!-- vale off -->

(float64, default: `1`)

<!-- vale on -->

Amount by which the `in_ports.limit` should be multiplied for
this label value.

</dd>
</dl>

---

<!-- vale off -->

### RateLimiterParameters {#rate-limiter-parameters}

<!-- vale on -->

<dl>
<dt>label_key</dt>
<dd>

<!-- vale off -->

(string, **required**)

<!-- vale on -->

Specifies which label the rate limiter should be keyed by.

Rate limiting is done independently for each value of the
[label](/concepts/flow-control/flow-label.md) with given key.
For example, to give each user a separate limit, assuming you
have a _user_ flow
label set up, set `label_key: "user"`.

</dd>
<dt>lazy_sync</dt>
<dd>

<!-- vale off -->

([RateLimiterParametersLazySync](#rate-limiter-parameters-lazy-sync))

<!-- vale on -->

Configuration of lazy-syncing behaviour of rate limiter

</dd>
<dt>limit_reset_interval</dt>
<dd>

<!-- vale off -->

(string, default: `"60s"`)

<!-- vale on -->

Time after which the limit for a given label value will be reset.

</dd>
<dt>tokens_label_key</dt>
<dd>

<!-- vale off -->

(string, default: `"tokens"`)

<!-- vale on -->

Flow label key that will be used to override the number of tokens
for this request.
This is an optional parameter and takes highest precedence
when assigning tokens to a request.
The label value must be a valid uint64 number.

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

(int64, minimum: `0`, default: `5`)

<!-- vale on -->

Number of times to lazy sync within the `limit_reset_interval`.

</dd>
</dl>

---

<!-- vale off -->

### Rego {#rego}

<!-- vale on -->

Rego define a set of labels that are extracted after evaluating a Rego module.

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
```

<dl>
<dt>labels</dt>
<dd>

<!-- vale off -->

(map of [RegoLabelProperties](#rego-label-properties), **required**)

<!-- vale on -->

A map of {key, value} pairs mapping from
[flow label](/concepts/flow-control/flow-label.md) keys to queries that define
how to extract and propagate flow labels with that key.
The name of the label maps to a variable in the Rego module. It maps to `data.<package>.<label>` variable.

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

Decides if the created flow label should be available as an attribute in OLAP telemetry and
propagated in [baggage](/concepts/flow-control/flow-label.md#baggage)

:::note

The flow label is always accessible in Aperture Policies regardless of this setting.

:::

:::caution

When using [FluxNinja ARC extension](arc/extension.md), telemetry enabled
labels are sent to FluxNinja ARC for observability. Telemetry should be disabled for
sensitive labels.

:::

</dd>
</dl>

---

<!-- vale off -->

### Regulator {#regulator}

<!-- vale on -->

_Regulator_ is a component that regulates the load at a
[_Control Point_][/concepts/flow-control/flow-selector.md/#control-point] by allowing only a specified percentage of
flows at random or by sticky sessions.

:::info

See also [\_Load Regulator overview](/concepts/flow-control/components/regulator.md).

:::

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([RegulatorDynamicConfig](#regulator-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig.

</dd>
<dt>in_ports</dt>
<dd>

<!-- vale off -->

([RegulatorIns](#regulator-ins))

<!-- vale on -->

Input ports for the _Regulator_.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([RegulatorParameters](#regulator-parameters))

<!-- vale on -->

Parameters for the _Regulator_.

</dd>
</dl>

---

<!-- vale off -->

### RegulatorDynamicConfig {#regulator-dynamic-config}

<!-- vale on -->

Dynamic Configuration for _Regulator_

<dl>
<dt>enable_label_values</dt>
<dd>

<!-- vale off -->

([]string)

<!-- vale on -->

Specify certain label values to be accepted by this flow filter regardless of accept percentage.

</dd>
</dl>

---

<!-- vale off -->

### RegulatorIns {#regulator-ins}

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

### RegulatorParameters {#regulator-parameters}

<!-- vale on -->

<dl>
<dt>flow_selector</dt>
<dd>

<!-- vale off -->

([FlowSelector](#flow-selector))

<!-- vale on -->

_Flow Selector_ selects the _Flows_ at which the _Regulator_ is applied.

</dd>
<dt>label_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

The flow label key for identifying sessions.

- When label key is specified, _Regulator_ acts as a sticky filter.
  The series of flows with the same value of label key get the same
  decision provided that the `accept_percentage` is same or higher.
- When label key is not specified, _Regulator_ acts as a stateless filter.
  Percentage of flows are selected randomly for rejection.

</dd>
</dl>

---

<!-- vale off -->

### Resources {#resources}

<!-- vale on -->

Resources that need to be setup for the policy to function

:::info

See also [Resources overview](/concepts/policy/resources.md).

:::

<dl>
<dt>flow_control</dt>
<dd>

<!-- vale off -->

([FlowControlResources](#flow-control-resources))

<!-- vale on -->

FlowControlResources are resources that are provided by flow control integration.

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

Decides if the created flow label should be available as an attribute in OLAP telemetry and
propagated in [baggage](/concepts/flow-control/flow-label.md#baggage)

:::note

The flow label is always accessible in Aperture Policies regardless of this setting.

:::

:::caution

When using [FluxNinja ARC extension](arc/extension.md), telemetry enabled
labels are sent to FluxNinja ARC for observability. Telemetry should be disabled for
sensitive labels.

:::

</dd>
</dl>

---

<!-- vale off -->

### SMA {#s-m-a}

<!-- vale on -->

Simple Moving Average (SMA) is a type of moving average that computes the average of a fixed number of signal readings.

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

Window of time over which the moving average is computed.

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

### ScaleInController {#scale-in-controller}

<!-- vale on -->

<dl>
<dt>alerter_parameters</dt>
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
</dl>

---

<!-- vale off -->

### ScaleOutController {#scale-out-controller}

<!-- vale on -->

<dl>
<dt>alerter_parameters</dt>
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

Each Agent instantiates an independent copy of the scheduler, but output
signals for accepted and incoming token rate are aggregated across all agents.

:::

<dl>
<dt>out_ports</dt>
<dd>

<!-- vale off -->

([SchedulerOuts](#scheduler-outs))

<!-- vale on -->

Output ports for the Scheduler component.

</dd>
<dt>parameters</dt>
<dd>

<!-- vale off -->

([SchedulerParameters](#scheduler-parameters))

<!-- vale on -->

Scheduler parameters.

</dd>
</dl>

---

<!-- vale off -->

### SchedulerOuts {#scheduler-outs}

<!-- vale on -->

Output for the Scheduler component.

<dl>
<dt>accepted_concurrency</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Accepted concurrency is actual concurrency on a control point that this
scheduler is applied on.
Value of this signal is aggregated from all the relevant schedulers.

</dd>
<dt>incoming_concurrency</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

Incoming concurrency is concurrency that'd be needed to accept all the
flows entering the scheduler.

This is computed in the same way as `accepted_concurrency`,
by summing up tokens from all the flows entering scheduler,
including rejected ones.

</dd>
</dl>

---

<!-- vale off -->

### SchedulerParameters {#scheduler-parameters}

<!-- vale on -->

Scheduler parameters

<dl>
<dt>auto_tokens</dt>
<dd>

<!-- vale off -->

(bool, default: `false`)

<!-- vale on -->

Automatically estimate the size of a flow in each workload, based on
historical latency. Each workload's `tokens` will be set to average
latency of flows in that workload during last few seconds (exact duration
of this average can change).
This setting is useful in concurrency limiting use-case, where the
concurrency is calculated as (avg. latency \* in-flight flows).

The value of tokens estimated by `auto_tokens` takes lower precedence
than the value of `tokens` specified in the workload definition
and `tokens` explicitly specified in the flow labels.

</dd>
<dt>decision_deadline_margin</dt>
<dd>

<!-- vale off -->

(string, default: `"0.01s"`)

<!-- vale on -->

Decision deadline margin is the amount of time that the scheduler will
subtract from the request deadline to determine the deadline for the
decision. This is to ensure that the scheduler has enough time to
make a decision before the request deadline happens, accounting for
processing delays.
The request deadline is based on the
[gRPC deadline](https://grpc.io/blog/deadlines) or the
[`grpc-timeout` HTTP header](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests).

Fail-open logic is use for flow control APIs, so if the gRPC deadline
reaches, the flow will end up being unconditionally allowed while
it is still waiting on the scheduler.

</dd>
<dt>default_workload_parameters</dt>
<dd>

<!-- vale off -->

([SchedulerWorkloadParameters](#scheduler-workload-parameters))

<!-- vale on -->

Parameters to be used if none of workloads specified in `workloads` match.

</dd>
<dt>tokens_label_key</dt>
<dd>

<!-- vale off -->

(string, default: `"tokens"`)

<!-- vale on -->

- Key for a flow label that can be used to override the default number of tokens for this flow.
- The value associated with this key must be a valid uint64 number.
- If this parameter is not provided, the number of tokens for the flow will be determined by the matched workload's token count.

</dd>
<dt>workloads</dt>
<dd>

<!-- vale off -->

([[]SchedulerWorkload](#scheduler-workload))

<!-- vale on -->

List of workloads to be used in scheduler.

Categorizing [flows](/concepts/flow-control/flow-control.md#flow) into workloads
allows for load-shedding to be "intelligent" compared to random rejections.
There are two aspects of this "intelligence":

- Scheduler can more precisely calculate concurrency if it understands
  that flows belonging to different classes have different weights (for example, insert queries compared to select queries).
- Setting different priorities to different workloads lets the scheduler
  avoid dropping important traffic during overload.

Each workload in this list specifies also a matcher that is used to
determine which flow will be categorized into which workload.
In case of multiple matching workloads, the first matching one will be used.
If none of workloads match, `default_workload` will be used.

:::info

See also [workload definition in the concepts
section](/concepts/flow-control/components/load-scheduler.md#workload).

:::

</dd>
</dl>

---

<!-- vale off -->

### SchedulerWorkload {#scheduler-workload}

<!-- vale on -->

Workload defines a class of flows that preferably have similar properties such as response latency and desired priority.

<dl>
<dt>label_matcher</dt>
<dd>

<!-- vale off -->

([LabelMatcher](#label-matcher))

<!-- vale on -->

Label Matcher to select a Workload based on
[flow labels](/concepts/flow-control/flow-label.md).

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

Parameters such as priority, tokens and fairness key that
are applicable to flows within a workload.

<dl>
<dt>fairness_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Fairness key is a label key that can be used to provide fairness within a workload.
Any [flow label](/concepts/flow-control/flow-label.md) can be used here. For example, if
you have a classifier that sets `user` flow label, you might want to set
`fairness_key = "user"`.

</dd>
<dt>priority</dt>
<dd>

<!-- vale off -->

(int64, minimum: `0`, maximum: `255`, default: `0`)

<!-- vale on -->

Describes priority level of the flows within the workload.
Priority level ranges from 0 to 255.
Higher numbers means higher priority level.
Priority levels have non-linear effect on the workload scheduling. The following formula is used to determine the position of a flow in the queue based on virtual finish time:

$$
\text{virtual\_finish\_time} = \text{virtual\_time} + \left(\text{tokens} \cdot \left(\text{256} - \text{priority}\right)\right)
$$

</dd>
<dt>tokens</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Tokens determines the cost of admitting a single flow in the workload,
which is typically defined as milliseconds of flow latency (time to response or duration of a feature) or
simply equal to 1 if the resource being accessed is constrained by the
number of flows (3rd party rate limiters).
This override is applicable only if tokens for the flow aren't specified
in the flow labels.

</dd>
</dl>

---

<!-- vale off -->

### ServiceSelector {#service-selector}

<!-- vale on -->

Describes which service a [flow control or observability
component](/concepts/flow-control/flow-control.md#components) should apply
to

:::info

See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).

:::

<dl>
<dt>agent_group</dt>
<dd>

<!-- vale off -->

(string, default: `"default"`)

<!-- vale on -->

Which [agent-group](/concepts/flow-control/flow-selector.md#agent-group) this
selector applies to.

:::info

Agent Groups are used to scope policies to a subset of agents connected to the same controller.
This is especially useful in the Kubernetes sidecar installation because service discovery is switched off in that mode.
The agents within an agent group form a peer to peer cluster and constantly share state.

:::

</dd>
<dt>service</dt>
<dd>

<!-- vale off -->

(string, default: `"any"`)

<!-- vale on -->

The Fully Qualified Domain Name of the
[service](/concepts/flow-control/flow-selector.md) to select.

In Kubernetes, this is the FQDN of the Service object.

:::info

`any` matches all services.

:::

:::info

In the Kubernetes sidecar installation mode, service discovery is switched off by default.
To scope policies to services, the `service` should be set to `any` and instead, `agent_group` name should be used.

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

The _Signal Generator_ component generates a smooth and continuous signal
by following a sequence of specified steps. Each step has two parameters:

- `target_output`: The desired output value at the end of the step.
- `duration`: The time it takes for the signal to change linearly from the
  previous step's `target_output` to the current step's `target_output`.

The output signal starts at the `target_output` of the first step and
changes linearly between steps based on their `duration`. The _Signal
Generator_ can be controlled to move forwards, backwards, or reset to the
beginning based on input signals.

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

A Boolean signal indicating whether the _Signal Generator_ is at the end of signal generation.

</dd>
<dt>at_start</dt>
<dd>

<!-- vale off -->

([OutPort](#out-port))

<!-- vale on -->

A Boolean signal indicating whether the _Signal Generator_ is at the start of signal generation.

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

([[]SignalGeneratorParametersStep](#signal-generator-parameters-step), **required**)

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

Duration for which the step is active.

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

Type of Combinator that switches between `on_signal` and `off_signal` signals based on switch input

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

### UnaryOperator {#unary-operator}

<!-- vale on -->

Takes an input signal and emits the output after applying the specified unary operator

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

(string, one of: `abs | acos | acosh | asin | asinh | atan | atanh | cbrt | ceil | cos | cosh | erf | erfc | erfcinv | erfinv | exp | exp2 | expm1 | floor | gamma | j0 | j1 | lgamma | log | log10 | log1p | log2 | round | roundtoeven | sin | sinh | sqrt | tan | tanh | trunc | y0 | y1`)

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
- `roundtoeven`: Round to nearest integer, with ties going to the nearest even integer.
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

Component that emits a variable value as an output signal, can be defined in dynamic configuration.

<dl>
<dt>default_config</dt>
<dd>

<!-- vale off -->

([VariableDynamicConfig](#variable-dynamic-config))

<!-- vale on -->

Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

<!-- vale off -->

(string)

<!-- vale on -->

Configuration key for DynamicConfig.

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

### VariableDynamicConfig {#variable-dynamic-config}

<!-- vale on -->

<dl>
<dt>constant_signal</dt>
<dd>

<!-- vale off -->

([ConstantSignal](#constant-signal))

<!-- vale on -->

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
