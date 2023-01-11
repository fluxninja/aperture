---
title: Policy Configuration Reference
sidebar_position: 1
sidebar_label: Policy
---

Reference for all objects used in [the Policy language](/concepts/policy/policy.md).

The top-level object representing a policy is [v1Policy](#v1-policy).

<!---
Generated File Starts
-->

## Table of contents

### POLICY CONFIGURATION

| Key | Reference         |
| --- | ----------------- |
|     | [Policy](#policy) |

## Reference

### _Policy_ {#policy}

#### Members

<dl>

<dt>body</dt>
<dd>

Type: [V1Policy](#v1-policy)

</dd>

</dl>

## Objects

### FluxMeterExponentialBuckets {#flux-meter-exponential-buckets}

ExponentialBuckets creates `count` number of buckets where the lowest bucket has an upper bound of `start`
and each following bucket's upper bound is `factor` times the previous bucket's upper bound. The final +inf
bucket is not counted.

#### Properties

<dl>
<dt>count</dt>
<dd>

(int32, `gt=0`) Number of buckets.

</dd>
<dt>factor</dt>
<dd>

(float64, `gt=1.0`) Factor to be multiplied to the previous bucket's upper bound to calculate the following bucket's upper bound.

</dd>
<dt>start</dt>
<dd>

(float64, `gt=0`) Upper bound of the lowest bucket.

</dd>
</dl>

### FluxMeterExponentialBucketsRange {#flux-meter-exponential-buckets-range}

ExponentialBucketsRange creates `count` number of buckets where the lowest bucket is `min` and the highest
bucket is `max`. The final +inf bucket is not counted.

#### Properties

<dl>
<dt>count</dt>
<dd>

(int32, `gt=0`) Number of buckets.

</dd>
<dt>max</dt>
<dd>

(float64) Highest bucket.

</dd>
<dt>min</dt>
<dd>

(float64, `gt=0`) Lowest bucket.

</dd>
</dl>

### FluxMeterLinearBuckets {#flux-meter-linear-buckets}

LinearBuckets creates `count` number of buckets, each `width` wide, where the lowest bucket has an
upper bound of `start`. The final +inf bucket is not counted.

#### Properties

<dl>
<dt>count</dt>
<dd>

(int32, `gt=0`) Number of buckets.

</dd>
<dt>start</dt>
<dd>

(float64) Upper bound of the lowest bucket.

</dd>
<dt>width</dt>
<dd>

(float64) Width of each bucket.

</dd>
</dl>

### FluxMeterStaticBuckets {#flux-meter-static-buckets}

StaticBuckets holds the static value of the buckets where latency histogram will be stored.

#### Properties

<dl>
<dt>buckets</dt>
<dd>

([]float64, default: `[5.0,10.0,25.0,50.0,100.0,250.0,500.0,1000.0,2500.0,5000.0,10000.0]`)

</dd>
</dl>

### MatchExpressionList {#match-expression-list}

List of MatchExpressions that is used for all/any matching

eg. {any: {of: [expr1, expr2]}}.

#### Properties

<dl>
<dt>of</dt>
<dd>

([[]V1MatchExpression](#v1-match-expression)) List of subexpressions of the match expression.

</dd>
</dl>

### RateLimiterLazySync {#rate-limiter-lazy-sync}

#### Properties

<dl>
<dt>enabled</dt>
<dd>

(bool) Enables lazy sync

</dd>
<dt>num_sync</dt>
<dd>

(int64, `gt=0`, default: `5`) Number of times to lazy sync within the _limit_reset_interval_.

</dd>
</dl>

### RateLimiterOverride {#rate-limiter-override}

#### Properties

<dl>
<dt>label_value</dt>
<dd>

(string, `required`) Value of the label for which the override should be applied.

</dd>
<dt>limit_scale_factor</dt>
<dd>

(float64, default: `1`) Amount by which the _in_ports.limit_ should be multiplied for this label value.

</dd>
</dl>

### RuleRego {#rule-rego}

Raw rego rules are compiled 1:1 to rego queries

High-level extractor-based rules are compiled into a single rego query.

#### Properties

<dl>
<dt>query</dt>
<dd>

(string, `required`) Query string to extract a value (eg. `data.<mymodulename>.<variablename>`).

Note: The module name must match the package name from the "source".

</dd>
<dt>source</dt>
<dd>

(string, `required`) Source code of the rego module.

Note: Must include a "package" declaration.

</dd>
</dl>

### SchedulerWorkload {#scheduler-workload}

Workload defines a class of requests that preferably have similar properties such as response latency or desired priority.

#### Properties

<dl>
<dt>label_matcher</dt>
<dd>

([V1LabelMatcher](#v1-label-matcher), `required`) Label Matcher to select a Workload based on
[flow labels](/concepts/flow-control/flow-label.md).

</dd>
<dt>workload_parameters</dt>
<dd>

([SchedulerWorkloadParameters](#scheduler-workload-parameters), `required`) WorkloadParameters associated with flows matching the label matcher.

</dd>
</dl>

### SchedulerWorkloadParameters {#scheduler-workload-parameters}

WorkloadParameters defines parameters such as priority, tokens and fairness key that are applicable to flows within a workload.

#### Properties

<dl>
<dt>fairness_key</dt>
<dd>

(string) Fairness key is a label key that can be used to provide fairness within a workload.
Any [flow label](/concepts/flow-control/flow-label.md) can be used here. Eg. if
you have a classifier that sets `user` flow label, you might want to set
`fairness_key = "user"`.

</dd>
<dt>priority</dt>
<dd>

(int64, `gte=0,lte=255`) Describes priority level of the requests within the workload.
Priority level ranges from 0 to 255.
Higher numbers means higher priority level.

</dd>
<dt>tokens</dt>
<dd>

(string, default: `1`) Tokens determines the cost of admitting a single request the workload, which is typically defined as milliseconds of response latency.
This override is applicable only if `auto_tokens` is set to false.

</dd>
</dl>

### v1AddressExtractor {#v1-address-extractor}

Display an [Address][ext-authz-address] as a single string, eg. `<ip>:<port>`

IP addresses in attribute context are defined as objects with separate ip and port fields.
This is a helper to display an address as a single string.

Note: Use with care, as it might accidentally introduce a high-cardinality flow label values.

[ext-authz-address]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/address.proto#config-core-v3-address

Example:

```yaml
from: "source.address # or destination.address"
```

#### Properties

<dl>
<dt>from</dt>
<dd>

(string, `required`) Attribute path pointing to some string - eg. "source.address".

</dd>
</dl>

### v1Alerter {#v1-alerter}

Alerter reacts to a signal and generates alert to send to alert manager.

#### Properties

<dl>
<dt>alerter_config</dt>
<dd>

([V1AlerterConfig](#v1-alerter-config)) Alerter configuration

</dd>
<dt>in_ports</dt>
<dd>

([V1AlerterIns](#v1-alerter-ins)) Input ports for the Alerter component.

</dd>
</dl>

### v1AlerterConfig {#v1-alerter-config}

AlerterConfig is a common config for separate alerter components and alerters embedded in other components.

#### Properties

<dl>
<dt>alert_channels</dt>
<dd>

([]string) A list of alert channel strings.

</dd>
<dt>alert_name</dt>
<dd>

(string, default: `alert`) Name of the alert.

</dd>
<dt>resolve_timeout</dt>
<dd>

(string, default: `300s`) Duration of alert resolver.

</dd>
<dt>severity</dt>
<dd>

(string, `oneof=info warn crit`, default: `info`) Severity of the alert, one of 'info', 'warn' or 'crit'.

</dd>
</dl>

### v1AlerterIns {#v1-alerter-ins}

Inputs for the Alerter component.

#### Properties

<dl>
<dt>alert</dt>
<dd>

([V1InPort](#v1-in-port)) Signal which Alerter is monitoring.

</dd>
</dl>

### v1ArithmeticCombinator {#v1-arithmetic-combinator}

Type of combinator that computes the arithmetic operation on the operand signals

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1ArithmeticCombinatorIns](#v1-arithmetic-combinator-ins)) Input ports for the Arithmetic Combinator component.

</dd>
<dt>operator</dt>
<dd>

(string, `oneof=add sub mul div xor lshift rshift`) Operator of the arithmetic operation.

The arithmetic operation can be addition, subtraction, multiplication, division, XOR, right bit shift or left bit shift.
In case of XOR and bitshifts, value of signals is cast to integers before performing the operation.

</dd>
<dt>out_ports</dt>
<dd>

([V1ArithmeticCombinatorOuts](#v1-arithmetic-combinator-outs)) Output ports for the Arithmetic Combinator component.

</dd>
</dl>

### v1ArithmeticCombinatorIns {#v1-arithmetic-combinator-ins}

Inputs for the Arithmetic Combinator component.

#### Properties

<dl>
<dt>lhs</dt>
<dd>

([V1InPort](#v1-in-port)) Left hand side of the arithmetic operation.

</dd>
<dt>rhs</dt>
<dd>

([V1InPort](#v1-in-port)) Right hand side of the arithmetic operation.

</dd>
</dl>

### v1ArithmeticCombinatorOuts {#v1-arithmetic-combinator-outs}

Outputs for the Arithmetic Combinator component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Result of arithmetic operation.

</dd>
</dl>

### v1Circuit {#v1-circuit}

Circuit is defined as a dataflow graph of inter-connected components

:::info
See also [Circuit overview](/concepts/policy/circuit.md).
:::

Signals flow between components via ports.
As signals traverse the circuit, they get processed, stored within components or get acted upon (e.g. load-shed, rate-limit, auto-scale etc.).
Circuit is evaluated periodically in order to respond to changes in signal readings.

:::info
**Signal**

Signals are floating-point values.

A signal can also have a special **Invalid** value. It's usually used to
communicate that signal doesn't have a meaningful value at the moment, eg.
[PromQL](#v1-prom-q-l) emits such a value if it cannot execute a query.
Components know when their input signals are invalid and can act
accordingly. They can either propagate the invalidness, by making their
output itself invalid (like eg.
[ArithmeticCombinator](#v1-arithmetic-combinator)) or use some different
logic, like eg. [Extrapolator](#v1-extrapolator). Refer to a component's
docs on how exactly it handles invalid inputs.
:::

#### Properties

<dl>
<dt>components</dt>
<dd>

([[]V1Component](#v1-component)) Defines a signal processing graph as a list of components.

</dd>
<dt>evaluation_interval</dt>
<dd>

(string, default: `0.5s`) Evaluation interval (tick) is the time period between consecutive runs of the policy circuit.
This interval is typically aligned with how often the corrective action (actuation) needs to be taken.

</dd>
</dl>

### v1Classifier {#v1-classifier}

Set of classification rules sharing a common selector

:::info
See also [Classifier overview](/concepts/flow-control/flow-classifier.md).
:::

Example:

```yaml
selector:
  service_selector:
    service: service1.default.svc.cluster.local
  flow_selector:
    control_point:
      traffic: ingress
rules:
  user:
    extractor:
      from: request.http.headers.user
```

#### Properties

<dl>
<dt>flow_selector</dt>
<dd>

([V1FlowSelector](#v1-flow-selector), `required`) Defines where to apply the flow classification rule.

</dd>
<dt>rules</dt>
<dd>

(map of [V1Rule](#v1-rule), `required,gt=0,dive,keys,required,endkeys,required`) A map of {key, value} pairs mapping from
[flow label](/concepts/flow-control/flow-label.md) keys to rules that define
how to extract and propagate flow labels with that key.

</dd>
</dl>

### v1Component {#v1-component}

Computational block that form the circuit

:::info
See also [Components overview](/concepts/policy/circuit.md#components).
:::

Signals flow into the components via input ports and results are emitted on output ports.
Components are wired to each other based on signal names forming an execution graph of the circuit.

:::note
Loops are broken by the runtime at the earliest component index that is part of the loop.
The looped signals are saved in the tick they are generated and served in the subsequent tick.
:::

There are three categories of components:

- "source" components – they take some sort of input from "the real world" and output
  a signal based on this input. Example: [PromQL](#v1-prom-q-l). In the UI
  they're represented by green color.
- signal processor components – "pure" components that don't interact with the "real world".
  Examples: [GradientController](#v1-gradient-controller), [Max](#v1-max).
  :::note
  Signal processor components's output can depend on their internal state, in addition to the inputs.
  Eg. see the [Exponential Moving Average filter](#v1-e-m-a).
  :::
- "sink" components – they affect the real world.
  [ConcurrencyLimiter.LoadActuator](#v1-concurrency-limiter) and [RateLimiter](#v1-rate-limiter).
  In the UI, represented by orange color. Sink components usually come in pairs with a
  "sources" component which emits a feedback signal, like
  `accepted_concurrency` emitted by ConcurrencyLimiter.Scheduler.

:::tip
Sometimes you may want to use a constant value as one of component's inputs.
You can create an input port containing the constant value instead of being connected to a signal.
To do so, use the [InPort](#v1-in_port)'s .withConstantValue(constant_value) method.
If You need to provide the same constant signal to multiple components,
You can use the [Constant](#v1-constant) component.
:::

See also [Policy](#v1-policy) for a higher-level explanation of circuits.

#### Properties

<dl>
<dt>alerter</dt>
<dd>

([V1Alerter](#v1-alerter)) Alerter reacts to a signal and generates alert to send to alert manager.

</dd>
<dt>arithmetic_combinator</dt>
<dd>

([V1ArithmeticCombinator](#v1-arithmetic-combinator)) Applies the given operator on input operands (signals) and emits the result.

</dd>
<dt>concurrency_limiter</dt>
<dd>

([V1ConcurrencyLimiter](#v1-concurrency-limiter)) Concurrency Limiter provides service protection by applying prioritized load shedding of flows using a network scheduler (e.g. Weighted Fair Queuing).

</dd>
<dt>constant</dt>
<dd>

([V1Constant](#v1-constant)) Emits a constant signal.

</dd>
<dt>decider</dt>
<dd>

([V1Decider](#v1-decider)) Decider emits the binary result of comparison operator on two operands.

</dd>
<dt>differentiator</dt>
<dd>

([V1Differentiator](#v1-differentiator)) Differentiator calculates rate of change per tick.

</dd>
<dt>ema</dt>
<dd>

([V1EMA](#v1-e-m-a)) Exponential Moving Average filter.

</dd>
<dt>extrapolator</dt>
<dd>

([V1Extrapolator](#v1-extrapolator)) Takes an input signal and emits the extrapolated value; either mirroring the input value or repeating the last known value up to the maximum extrapolation interval.

</dd>
<dt>first_valid</dt>
<dd>

([V1FirstValid](#v1-first-valid)) Picks the first valid input signal and emits it.

</dd>
<dt>gradient_controller</dt>
<dd>

([V1GradientController](#v1-gradient-controller)) Gradient controller basically calculates the ratio between the signal and the setpoint to determine the magnitude of the correction that need to be applied.
This controller can be used to build AIMD (Additive Increase, Multiplicative Decrease) or MIMD style response.

</dd>
<dt>integrator</dt>
<dd>

([V1Integrator](#v1-integrator)) Accumulates sum of signal every tick.

</dd>
<dt>max</dt>
<dd>

([V1Max](#v1-max)) Emits the maximum of the input signals.

</dd>
<dt>min</dt>
<dd>

([V1Min](#v1-min)) Emits the minimum of the input signals.

</dd>
<dt>promql</dt>
<dd>

([V1PromQL](#v1-prom-q-l)) Periodically runs a Prometheus query in the background and emits the result.

</dd>
<dt>rate_limiter</dt>
<dd>

([V1RateLimiter](#v1-rate-limiter)) Rate Limiter provides service protection by applying rate limiter.

</dd>
<dt>sink</dt>
<dd>

([V1Sink](#v1-sink)) Sink is a sink component that does nothing.

</dd>
<dt>sqrt</dt>
<dd>

([V1Sqrt](#v1-sqrt)) Takes an input signal and emits the square root of the input signal.

</dd>
<dt>switcher</dt>
<dd>

([V1Switcher](#v1-switcher)) Switcher acts as a switch that emits one of the two signals based on third signal.

</dd>
</dl>

### v1ConcurrencyLimiter {#v1-concurrency-limiter}

Concurrency Limiter is an actuator component that regulates flows in order to provide active service protection

:::info
See also [Concurrency Limiter overview](/concepts/flow-control/concurrency-limiter.md).
:::

It is based on the actuation strategy (e.g. load actuator) and workload scheduling which is based on Weighted Fair Queuing principles.
Concurrency is calculated in terms of total tokens which translate to (avg. latency \* in-flight requests), i.e. Little's Law.

ConcurrencyLimiter configuration is split into two parts: An actuation
strategy and a scheduler. Right now, only `load_actuator` strategy is available.

#### Properties

<dl>
<dt>flow_selector</dt>
<dd>

([V1FlowSelector](#v1-flow-selector), `required`) Flow Selector decides the service and flows at which the concurrency limiter is applied.

</dd>
<dt>load_actuator</dt>
<dd>

([V1LoadActuator](#v1-load-actuator)) Actuator based on limiting the accepted concurrency under incoming concurrency \* load multiplier.

Actuation strategy defines the input signal that will drive the scheduler.

</dd>
<dt>scheduler</dt>
<dd>

([V1Scheduler](#v1-scheduler), `required`) Configuration of Weighted Fair Queuing-based workload scheduler.

Contains configuration of per-agent scheduler, and also defines some
output signals.

</dd>
</dl>

### v1Constant {#v1-constant}

Component that emits a constant value as an output signal

#### Properties

<dl>
<dt>out_ports</dt>
<dd>

([V1ConstantOuts](#v1-constant-outs)) Output ports for the Constant component.

</dd>
<dt>value</dt>
<dd>

(float64) The constant value to be emitted.

</dd>
</dl>

### v1ConstantOuts {#v1-constant-outs}

Outputs for the Constant component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) The constant value is emitted to the output port.

</dd>
</dl>

### v1ControllerDynamicConfig {#v1-controller-dynamic-config}

Dynamic Configuration for a Controller

#### Properties

<dl>
<dt>manual_mode</dt>
<dd>

(bool) Decides whether the controller runs in "manual_mode".
In manual mode, the controller does not adjust the control variable I.E. emits the same output as the control variable input.

</dd>
</dl>

### v1Decider {#v1-decider}

Type of combinator that computes the comparison operation on lhs and rhs signals

The comparison operator can be greater-than, less-than, greater-than-or-equal, less-than-or-equal, equal, or not-equal.

This component also supports time-based response, i.e. the output
transitions between 1.0 or 0.0 signal if the decider condition is
true or false for at least "true_for" or "false_for" duration. If
`true_for` and `false_for` durations are zero then the transitions are
instantaneous.

#### Properties

<dl>
<dt>false_for</dt>
<dd>

(string, default: `0s`) Duration of time to wait before a transition to false state.
If the duration is zero, the transition will happen instantaneously.

</dd>
<dt>in_ports</dt>
<dd>

([V1DeciderIns](#v1-decider-ins)) Input ports for the Decider component.

</dd>
<dt>operator</dt>
<dd>

(string, `oneof=gt lt gte lte eq neq`) Comparison operator that computes operation on lhs and rhs input signals.

</dd>
<dt>out_ports</dt>
<dd>

([V1DeciderOuts](#v1-decider-outs)) Output ports for the Decider component.

</dd>
<dt>true_for</dt>
<dd>

(string, default: `0s`) Duration of time to wait before a transition to true state.
If the duration is zero, the transition will happen instantaneously.

</dd>
</dl>

### v1DeciderIns {#v1-decider-ins}

Inputs for the Decider component.

#### Properties

<dl>
<dt>lhs</dt>
<dd>

([V1InPort](#v1-in-port)) Left hand side input signal for the comparison operation.

</dd>
<dt>rhs</dt>
<dd>

([V1InPort](#v1-in-port)) Right hand side input signal for the comparison operation.

</dd>
</dl>

### v1DeciderOuts {#v1-decider-outs}

Outputs for the Decider component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Selected signal (1.0 or 0.0).

</dd>
</dl>

### v1Differentiator {#v1-differentiator}

Differentiator calculates rate of change per tick.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1DifferentiatorIns](#v1-differentiator-ins)) Input ports for the Differentiator component.

</dd>
<dt>out_ports</dt>
<dd>

([V1DifferentiatorOuts](#v1-differentiator-outs)) Output ports for the Differentiator component.

</dd>
<dt>window</dt>
<dd>

(string, default: `5s`) The window of time over which differentiator operates.

</dd>
</dl>

### v1DifferentiatorIns {#v1-differentiator-ins}

Inputs for the Differentiator component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1InPort](#v1-in-port))

</dd>
</dl>

### v1DifferentiatorOuts {#v1-differentiator-outs}

Outputs for the Differentiator component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port))

</dd>
</dl>

### v1EMA {#v1-e-m-a}

Exponential Moving Average (EMA) is a type of moving average that applies exponentially more weight to recent signal readings

At any time EMA component operates in one of the following states:

1. Warm up state: The first warm_up_window samples are used to compute the initial EMA.
   If an invalid reading is received during the warm_up_window, the last good average is emitted and the state gets reset back to beginning of Warm up state.
2. Normal state: The EMA is computed using following formula.

The EMA for a series $Y$ is calculated recursively as:

$$
\text{EMA} _t =
\begin{cases}
  Y_0, &\text{for } t = 0 \\
  \alpha Y_t + (1 - \alpha) \text{EMA} _{t-1}, &\text{for }t > 0
\end{cases}
$$

The coefficient $\alpha$ represents the degree of weighting decrease, a constant smoothing factor between 0 and 1.
A higher $\alpha$ discounts older observations faster.
The $\alpha$ is computed using ema_window:

$$
\alpha = \frac{2}{N + 1} \quad\text{where } N = \frac{\text{ema\_window}}{\text{evaluation\_period}}
$$

The EMA filter also employs a min-max-envelope logic during warm up stage, explained [here](#v1-e-m-a-ins).

#### Properties

<dl>
<dt>correction_factor_on_max_envelope_violation</dt>
<dd>

(float64, `gte=0,lte=1.0`, default: `1`) Correction factor to apply on the output value if its in violation of the max envelope.

</dd>
<dt>correction_factor_on_min_envelope_violation</dt>
<dd>

(float64, `gte=1.0`, default: `1`) Correction factor to apply on the output value if its in violation of the min envelope.

</dd>
<dt>ema_window</dt>
<dd>

(string, default: `5s`) Duration of EMA sampling window.

</dd>
<dt>in_ports</dt>
<dd>

([V1EMAIns](#v1-e-m-a-ins)) Input ports for the EMA component.

</dd>
<dt>out_ports</dt>
<dd>

([V1EMAOuts](#v1-e-m-a-outs)) Output ports for the EMA component.

</dd>
<dt>valid_during_warmup</dt>
<dd>

(bool) Whether the output is valid during the warm up stage.

</dd>
<dt>warm_up_window</dt>
<dd>

(string, default: `0s`) Duration of EMA warming up window.

The initial value of the EMA is the average of signal readings received during the warm up window.

</dd>
</dl>

### v1EMAIns {#v1-e-m-a-ins}

Inputs for the EMA component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1InPort](#v1-in-port)) Input signal to be used for the EMA computation.

</dd>
<dt>max_envelope</dt>
<dd>

([V1InPort](#v1-in-port)) Upper bound of the moving average.

Used during the warm-up stage: if the signal would exceed `max_envelope`
it's multiplied by `correction_factor_on_max_envelope_violation` **once per tick**.

:::note
If the signal deviates from `max_envelope` faster than the correction
faster, it might end up exceeding the envelope.
:::

:::note
The envelope logic is **not** used outside the warm-up stage!
:::

</dd>
<dt>min_envelope</dt>
<dd>

([V1InPort](#v1-in-port)) Lower bound of the moving average.

Used during the warm-up stage analogously to `max_envelope`.

</dd>
</dl>

### v1EMAOuts {#v1-e-m-a-outs}

Outputs for the EMA component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Exponential moving average of the series of reading as an output signal.

</dd>
</dl>

### v1EqualsMatchExpression {#v1-equals-match-expression}

Label selector expression of the equal form "label == value".

#### Properties

<dl>
<dt>label</dt>
<dd>

(string, `required`) Name of the label to equal match the value.

</dd>
<dt>value</dt>
<dd>

(string) Exact value that the label should be equal to.

</dd>
</dl>

### v1Extractor {#v1-extractor}

Defines a high-level way to specify how to extract a flow label value given http request metadata, without a need to write rego code

There are multiple variants of extractor, specify exactly one.

#### Properties

<dl>
<dt>address</dt>
<dd>

([V1AddressExtractor](#v1-address-extractor)) Display an address as a single string - `<ip>:<port>`.

</dd>
<dt>from</dt>
<dd>

(string) Use an attribute with no conversion

Attribute path is a dot-separated path to attribute.

Should be either:

- one of the fields of [Attribute Context][attribute-context], or
- a special "request.http.bearer" pseudo-attribute.
  Eg. "request.http.method" or "request.http.header.user-agent"

Note: The same attribute path syntax is shared by other extractor variants,
wherever attribute path is needed in their "from" syntax.

Example:

```yaml
from: request.http.headers.user-agent
```

[attribute-context]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto

</dd>
<dt>json</dt>
<dd>

([V1JSONExtractor](#v1-json-extractor)) Deserialize a json, and extract one of the fields.

</dd>
<dt>jwt</dt>
<dd>

([V1JWTExtractor](#v1-j-w-t-extractor)) Parse the attribute as JWT and read the payload.

</dd>
<dt>path_templates</dt>
<dd>

([V1PathTemplateMatcher](#v1-path-template-matcher)) Match HTTP Path to given path templates.

</dd>
</dl>

### v1Extrapolator {#v1-extrapolator}

Extrapolates the input signal by repeating the last valid value during the period in which it is invalid

It does so until `maximum_extrapolation_interval` is reached, beyond which it emits invalid signal unless input signal becomes valid again.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1ExtrapolatorIns](#v1-extrapolator-ins)) Input ports for the Extrapolator component.

</dd>
<dt>max_extrapolation_interval</dt>
<dd>

(string, default: `10s`) Maximum time interval to repeat the last valid value of input signal.

</dd>
<dt>out_ports</dt>
<dd>

([V1ExtrapolatorOuts](#v1-extrapolator-outs)) Output ports for the Extrapolator component.

</dd>
</dl>

### v1ExtrapolatorIns {#v1-extrapolator-ins}

Inputs for the Extrapolator component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1InPort](#v1-in-port)) Input signal for the Extrapolator component.

</dd>
</dl>

### v1ExtrapolatorOuts {#v1-extrapolator-outs}

Outputs for the Extrapolator component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Extrapolated signal.

</dd>
</dl>

### v1FirstValid {#v1-first-valid}

Picks the first valid input signal from the array of input signals and emits it as an output signal

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1FirstValidIns](#v1-first-valid-ins)) Input ports for the FirstValid component.

</dd>
<dt>out_ports</dt>
<dd>

([V1FirstValidOuts](#v1-first-valid-outs)) Output ports for the FirstValid component.

</dd>
</dl>

### v1FirstValidIns {#v1-first-valid-ins}

Inputs for the FirstValid component.

#### Properties

<dl>
<dt>inputs</dt>
<dd>

([[]V1InPort](#v1-in-port)) Array of input signals.

</dd>
</dl>

### v1FirstValidOuts {#v1-first-valid-outs}

Outputs for the FirstValid component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) First valid input signal as an output signal.

</dd>
</dl>

### v1FlowMatcher {#v1-flow-matcher}

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
    - label: user_agent
      regex: ^(?!.*Chrome).*Safari
```

#### Properties

<dl>
<dt>control_point</dt>
<dd>

(string, `required`) [Control Point](/concepts/flow-control/flow-control.md#control-point)
identifies the location of a Flow within a Service. For an SDK based insertion, a Control Point can represent a particular feature or execution
block within a Service. In case of Service Mesh or Middleware insertion, a Control Point can identify ingress vs egress calls or distinct listeners
or filter chains.

</dd>
<dt>label_matcher</dt>
<dd>

([V1LabelMatcher](#v1-label-matcher)) Label matcher allows to add _additional_ condition on
[flow labels](/concepts/flow-control/flow-label.md)
must also be satisfied (in addition to service+control point matching)

:::info
See also [Label Matcher overview](/concepts/flow-control/flow-selector.md#label-matcher).
:::

:::note
[Classifiers](#v1-classifier) _can_ use flow labels created by some other
classifier, but only if they were created at some previous control point
(and propagated in baggage).

This limitation doesn't apply to selectors of other entities, like
Flux Meters or Actuators. It's valid to create a flow label on a control
point using classifier, and immediately use it for matching on the same
control point.
:::

</dd>
</dl>

### v1FlowSelector {#v1-flow-selector}

Describes which flow in which service a [flow control
component](/concepts/flow-control/flow-control.md#components) should apply
to

:::info
See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).
:::

#### Properties

<dl>
<dt>flow_matcher</dt>
<dd>

([V1FlowMatcher](#v1-flow-matcher), `required`)

</dd>
<dt>service_selector</dt>
<dd>

([V1ServiceSelector](#v1-service-selector), `required`)

</dd>
</dl>

### v1FluxMeter {#v1-flux-meter}

Flux Meter gathers metrics for the traffic that matches its selector.
The histogram created by Flux Meter measures the workload latency by default.

:::info
See also [Flux Meter overview](/concepts/flow-control/flux-meter.md).
:::

Example of a selector that creates a histogram metric for all HTTP requests
to particular service:

```yaml
selector:
  service_selector:
    service: myservice.mynamespace.svc.cluster.local
  flow_selector:
    control_point: ingress
```

#### Properties

<dl>
<dt>attribute_key</dt>
<dd>

(string, default: `workload_duration_ms`) Key of the attribute in access log or span from which the metric for this flux meter is read.

:::info
For list of available attributes in Envoy access logs, refer
[Envoy Filter](/get-started/installation/agent/envoy/istio.md#envoy-filter)
:::

</dd>
<dt>exponential_buckets</dt>
<dd>

([FluxMeterExponentialBuckets](#flux-meter-exponential-buckets))

</dd>
<dt>exponential_buckets_range</dt>
<dd>

([FluxMeterExponentialBucketsRange](#flux-meter-exponential-buckets-range))

</dd>
<dt>flow_selector</dt>
<dd>

([V1FlowSelector](#v1-flow-selector)) The selection criteria for the traffic that will be measured.

</dd>
<dt>linear_buckets</dt>
<dd>

([FluxMeterLinearBuckets](#flux-meter-linear-buckets))

</dd>
<dt>static_buckets</dt>
<dd>

([FluxMeterStaticBuckets](#flux-meter-static-buckets))

</dd>
</dl>

### v1GradientController {#v1-gradient-controller}

Gradient controller is a type of controller which tries to adjust the
control variable proportionally to the relative difference between setpoint
and actual value of the signal

The `gradient` describes a corrective factor that should be applied to the
control variable to get the signal closer to the setpoint. It is computed as follows:

$$
\text{gradient} = \left(\frac{\text{signal}}{\text{setpoint}}\right)^{\text{slope}}
$$

`gradient` is then clamped to [min_gradient, max_gradient] range.

The output of gradient controller is computed as follows:

$$
\text{output} = \text{gradient}_{\text{clamped}} \cdot \text{control\_variable} + \text{optimize}.
$$

Note the additional `optimize` signal, that can be used to "nudge" the
controller into desired idle state.

The output can be _optionally_ clamped to desired range using `max` and
`min` input.

#### Properties

<dl>
<dt>default_config</dt>
<dd>

([V1ControllerDynamicConfig](#v1-controller-dynamic-config)) Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

(string) Configuration key for DynamicConfig

</dd>
<dt>in_ports</dt>
<dd>

([V1GradientControllerIns](#v1-gradient-controller-ins)) Input ports of the Gradient Controller.

</dd>
<dt>max_gradient</dt>
<dd>

(float64, default: `1.7976931348623157e+308`) Maximum gradient which clamps the computed gradient value to the range, [min_gradient, max_gradient].

</dd>
<dt>min_gradient</dt>
<dd>

(float64, default: `-1.7976931348623157e+308`) Minimum gradient which clamps the computed gradient value to the range, [min_gradient, max_gradient].

</dd>
<dt>out_ports</dt>
<dd>

([V1GradientControllerOuts](#v1-gradient-controller-outs)) Output ports of the Gradient Controller.

</dd>
<dt>slope</dt>
<dd>

(float64, `required`) Slope controls the aggressiveness and direction of the Gradient Controller.

Slope is used as exponent on the signal to setpoint ratio in computation
of the gradient (see the [main description](#v1-gradient-controller) for
exact equation). Good intuition for this parameter is "What should the
Gradient Controller do to the control variable when signal is too high",
eg.:

- $\text{slope} = 1$: when signal is too high, increase control variable,
- $\text{slope} = -1$: when signal is too high, decrease control variable,
- $\text{slope} = -0.5$: when signal is to high, decrease control variable more slowly.

The sign of slope depends on correlation between the signal and control variable:

- Use $\text{slope} < 0$ if signal and control variable are _positively_
  correlated (eg. Per-pod CPU usage and total concurrency).
- Use $\text{slope} > 0$ if signal and control variable are _negatively_
  correlated (eg. Per-pod CPU usage and number of pods).

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
variable more slowly.

We recommend setting $|\text{slope}| < 1$ (eg. $\pm0.8$).
If you experience overshooting, consider lowering the magnitude even more.
Values of $|\text{slope}| > 1$ are not recommended.

:::note
Remember that the gradient and output signal can be (optionally) clamped,
so the _slope_ might not fully describe aggressiveness of the controller.
:::

</dd>
</dl>

### v1GradientControllerIns {#v1-gradient-controller-ins}

Inputs for the Gradient Controller component.

#### Properties

<dl>
<dt>control_variable</dt>
<dd>

([V1InPort](#v1-in-port)) Actual current value of the control variable.

This signal is multiplied by the gradient to produce the output.

</dd>
<dt>max</dt>
<dd>

([V1InPort](#v1-in-port)) Maximum value to limit the output signal.

</dd>
<dt>min</dt>
<dd>

([V1InPort](#v1-in-port)) Minimum value to limit the output signal.

</dd>
<dt>optimize</dt>
<dd>

([V1InPort](#v1-in-port)) Optimize signal is added to the output of the gradient calculation.

</dd>
<dt>setpoint</dt>
<dd>

([V1InPort](#v1-in-port)) Setpoint to be used for the gradient computation.

</dd>
<dt>signal</dt>
<dd>

([V1InPort](#v1-in-port)) Signal to be used for the gradient computation.

</dd>
</dl>

### v1GradientControllerOuts {#v1-gradient-controller-outs}

Outputs for the Gradient Controller component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Computed desired value of the control variable.

</dd>
</dl>

### v1InPort {#v1-in-port}

Components receive input from other components via InPorts

#### Properties

<dl>
<dt>constant_value</dt>
<dd>

(float64) Constant value to be used for this InPort instead of a signal.

</dd>
<dt>signal_name</dt>
<dd>

(string) Name of the incoming Signal on the InPort.

</dd>
</dl>

### v1Integrator {#v1-integrator}

Accumulates sum of signal every tick.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1IntegratorIns](#v1-integrator-ins)) Input ports for the Integrator component.

</dd>
<dt>out_ports</dt>
<dd>

([V1IntegratorOuts](#v1-integrator-outs)) Output ports for the Integrator component.

</dd>
</dl>

### v1IntegratorIns {#v1-integrator-ins}

Inputs for the Integrator component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1InPort](#v1-in-port)) The input signal.

</dd>
<dt>max</dt>
<dd>

([V1InPort](#v1-in-port)) The maximum output when reset is not set.

</dd>
<dt>min</dt>
<dd>

([V1InPort](#v1-in-port)) The minimum output when reset is not set.

</dd>
<dt>reset</dt>
<dd>

([V1InPort](#v1-in-port)) Resets the integrator output to zero when reset signal is valid and non-zero.

</dd>
</dl>

### v1IntegratorOuts {#v1-integrator-outs}

Outputs for the Integrator component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port))

</dd>
</dl>

### v1JSONExtractor {#v1-json-extractor}

Deserialize a json, and extract one of the fields

Example:

```yaml
from: request.http.body
pointer: /user/name
```

#### Properties

<dl>
<dt>from</dt>
<dd>

(string, `required`) Attribute path pointing to some strings - eg. "request.http.body".

</dd>
<dt>pointer</dt>
<dd>

(string) Json pointer represents a parsed json pointer which allows to select a specified field from the json payload.

Note: Uses [json pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
eg. `/foo/bar`. If the pointer points into an object, it'd be stringified.

</dd>
</dl>

### v1JWTExtractor {#v1-j-w-t-extractor}

Parse the attribute as JWT and read the payload

Specify a field to be extracted from payload using "json_pointer".

Note: The signature is not verified against the secret (we're assuming there's some
other parts of the system that handles such verification).

Example:

```yaml
from: request.http.bearer
json_pointer: /user/email
```

#### Properties

<dl>
<dt>from</dt>
<dd>

(string, `required`) Jwt token can be pulled from any input attribute, but most likely you'd want to use "request.http.bearer".

</dd>
<dt>json_pointer</dt>
<dd>

(string) Json pointer allowing to select a specified field from the json payload.

Note: Uses [json pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
eg. `/foo/bar`. If the pointer points into an object, it'd be stringified.

</dd>
</dl>

### v1K8sLabelMatcherRequirement {#v1-k8s-label-matcher-requirement}

Label selector requirement which is a selector that contains values, a key, and an operator that relates the key and values.

#### Properties

<dl>
<dt>key</dt>
<dd>

(string, `required`) Label key that the selector applies to.

</dd>
<dt>operator</dt>
<dd>

(string, `oneof=In NotIn Exists DoesNotExists`) Logical operator which represents a key's relationship to a set of values.
Valid operators are In, NotIn, Exists and DoesNotExist.

</dd>
<dt>values</dt>
<dd>

([]string) An array of string values that relates to the key by an operator.
If the operator is In or NotIn, the values array must be non-empty.
If the operator is Exists or DoesNotExist, the values array must be empty.

</dd>
</dl>

### v1LabelMatcher {#v1-label-matcher}

Allows to define rules whether a map of
[labels](/concepts/flow-control/flow-label.md)
should be considered a match or not

It provides three ways to define requirements:

- matchLabels
- matchExpressions
- arbitrary expression

If multiple requirements are set, they are all ANDed.
An empty label matcher always matches.

#### Properties

<dl>
<dt>expression</dt>
<dd>

([V1MatchExpression](#v1-match-expression)) An arbitrary expression to be evaluated on the labels.

</dd>
<dt>match_expressions</dt>
<dd>

([[]V1K8sLabelMatcherRequirement](#v1-k8s-label-matcher-requirement)) List of k8s-style label matcher requirements.

Note: The requirements are ANDed.

</dd>
<dt>match_labels</dt>
<dd>

(map of string) A map of {key,value} pairs representing labels to be matched.
A single {key,value} in the matchLabels requires that the label "key" is present and equal to "value".

Note: The requirements are ANDed.

</dd>
</dl>

### v1LoadActuator {#v1-load-actuator}

Takes the load multiplier input signal and publishes it to the schedulers in the data-plane

#### Properties

<dl>
<dt>alerter_config</dt>
<dd>

([V1AlerterConfig](#v1-alerter-config)) Configuration for embedded alerter.

</dd>
<dt>default_config</dt>
<dd>

([V1LoadActuatorDynamicConfig](#v1-load-actuator-dynamic-config)) Default configuration.

</dd>
<dt>dynamic_config_key</dt>
<dd>

(string) Configuration key for DynamicConfig.

</dd>
<dt>in_ports</dt>
<dd>

([V1LoadActuatorIns](#v1-load-actuator-ins)) Input ports for the Load Actuator component.

</dd>
</dl>

### v1LoadActuatorDynamicConfig {#v1-load-actuator-dynamic-config}

Dynamic Configuration for LoadActuator

#### Properties

<dl>
<dt>dry_run</dt>
<dd>

(bool) Decides whether to run the load actuator in dry-run mode. Dry run mode ensures that no traffic gets dropped by this load actuator.
Useful for observing the behavior of Load Actuator without disrupting any real traffic.

</dd>
</dl>

### v1LoadActuatorIns {#v1-load-actuator-ins}

Input for the Load Actuator component.

#### Properties

<dl>
<dt>load_multiplier</dt>
<dd>

([V1InPort](#v1-in-port)) Load multiplier is ratio of [incoming
concurrency](#v1-scheduler-outs) that needs to be accepted.

</dd>
</dl>

### v1MatchExpression {#v1-match-expression}

Defines a [map<string, string> → bool] expression to be evaluated on labels

MatchExpression has multiple variants, exactly one should be set.

Example:

```yaml
all:
  of:
    - label_exists: foo
    - label_equals: { label = app, value = frobnicator }
```

#### Properties

<dl>
<dt>all</dt>
<dd>

([MatchExpressionList](#match-expression-list)) The expression is true when all subexpressions are true.

</dd>
<dt>any</dt>
<dd>

([MatchExpressionList](#match-expression-list)) The expression is true when any subexpression is true.

</dd>
<dt>label_equals</dt>
<dd>

([V1EqualsMatchExpression](#v1-equals-match-expression)) The expression is true when label value equals given value.

</dd>
<dt>label_exists</dt>
<dd>

(string, `required`) The expression is true when label with given name exists.

</dd>
<dt>label_matches</dt>
<dd>

([V1MatchesMatchExpression](#v1-matches-match-expression)) The expression is true when label matches given regex.

</dd>
<dt>not</dt>
<dd>

([V1MatchExpression](#v1-match-expression)) The expression negates the result of subexpression.

</dd>
</dl>

### v1MatchesMatchExpression {#v1-matches-match-expression}

Label selector expression of the matches form "label matches regex".

#### Properties

<dl>
<dt>label</dt>
<dd>

(string, `required`) Name of the label to match the regular expression.

</dd>
<dt>regex</dt>
<dd>

(string, `required`) Regular expression that should match the label value.
It uses [golang's regular expression syntax](https://github.com/google/re2/wiki/Syntax).

</dd>
</dl>

### v1Max {#v1-max}

Takes a list of input signals and emits the signal with the maximum value

Max: output = max([]inputs).

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1MaxIns](#v1-max-ins)) Input ports for the Max component.

</dd>
<dt>out_ports</dt>
<dd>

([V1MaxOuts](#v1-max-outs)) Output ports for the Max component.

</dd>
</dl>

### v1MaxIns {#v1-max-ins}

Inputs for the Max component.

#### Properties

<dl>
<dt>inputs</dt>
<dd>

([[]V1InPort](#v1-in-port)) Array of input signals.

</dd>
</dl>

### v1MaxOuts {#v1-max-outs}

Output for the Max component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Signal with maximum value as an output signal.

</dd>
</dl>

### v1Min {#v1-min}

Takes an array of input signals and emits the signal with the minimum value
Min: output = min([]inputs).

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1MinIns](#v1-min-ins)) Input ports for the Min component.

</dd>
<dt>out_ports</dt>
<dd>

([V1MinOuts](#v1-min-outs)) Output ports for the Min component.

</dd>
</dl>

### v1MinIns {#v1-min-ins}

Inputs for the Min component.

#### Properties

<dl>
<dt>inputs</dt>
<dd>

([[]V1InPort](#v1-in-port)) Array of input signals.

</dd>
</dl>

### v1MinOuts {#v1-min-outs}

Output ports for the Min component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Signal with minimum value as an output signal.

</dd>
</dl>

### v1OutPort {#v1-out-port}

Components produce output for other components via OutPorts

#### Properties

<dl>
<dt>signal_name</dt>
<dd>

(string) Name of the outgoing Signal on the OutPort.

</dd>
</dl>

### v1PathTemplateMatcher {#v1-path-template-matcher}

Matches HTTP Path to given path templates

HTTP path will be matched against given path templates.
If a match occurs, the value associated with the path template will be treated as a result.
In case of multiple path templates matching, the most specific one will be chosen.

#### Properties

<dl>
<dt>template_values</dt>
<dd>

(map of string, `required`) Template value keys are OpenAPI-inspired path templates.

- Static path segment `/foo` matches a path segment exactly
- `/{param}` matches arbitrary path segment.
  (The param name is ignored and can be omitted (`{}`))
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

### v1Policy {#v1-policy}

Policy expresses reliability automation workflow that automatically protects services

:::info
See also [Policy overview](/concepts/policy/policy.md).
:::

Policy specification contains a circuit that defines the controller logic and resources that need to be setup.

#### Properties

<dl>
<dt>circuit</dt>
<dd>

([V1Circuit](#v1-circuit)) Defines the control-loop logic of the policy.

</dd>
<dt>resources</dt>
<dd>

([V1Resources](#v1-resources)) Resources (Flux Meters, Classifiers etc.) to setup.

</dd>
</dl>

### v1PromQL {#v1-prom-q-l}

Component that runs a Prometheus query periodically and returns the result as an output signal

#### Properties

<dl>
<dt>evaluation_interval</dt>
<dd>

(string, default: `10s`) Describes the interval between successive evaluations of the Prometheus query.

</dd>
<dt>out_ports</dt>
<dd>

([V1PromQLOuts](#v1-prom-q-l-outs)) Output ports for the PromQL component.

</dd>
<dt>query_string</dt>
<dd>

(string) Describes the Prometheus query to be run.

:::caution
TODO we should describe how to construct the query, eg. how to employ the
fluxmeters here or link to appropriate place in docs.
:::

</dd>
</dl>

### v1PromQLOuts {#v1-prom-q-l-outs}

Output for the PromQL component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) The result of the Prometheus query as an output signal.

</dd>
</dl>

### v1RateLimiter {#v1-rate-limiter}

Limits the traffic on a control point to specified rate

:::info
See also [Rate Limiter overview](/concepts/flow-control/rate-limiter.md).
:::

Ratelimiting is done separately on per-label-value basis. Use _label_key_
to select which label should be used as key.

#### Properties

<dl>
<dt>default_config</dt>
<dd>

([V1RateLimiterDynamicConfig](#v1-rate-limiter-dynamic-config)) Default configuration

</dd>
<dt>dynamic_config_key</dt>
<dd>

(string) Configuration key for DynamicConfig

</dd>
<dt>flow_selector</dt>
<dd>

([V1FlowSelector](#v1-flow-selector), `required`) Which control point to apply this ratelimiter to.

</dd>
<dt>in_ports</dt>
<dd>

([V1RateLimiterIns](#v1-rate-limiter-ins), `required`)

</dd>
<dt>label_key</dt>
<dd>

(string, `required`) Specifies which label the ratelimiter should be keyed by.

Rate limiting is done independently for each value of the
[label](/concepts/flow-control/flow-label.md) with given key.
Eg., to give each user a separate limit, assuming you have a _user_ flow
label set up, set `label_key: "user"`.

</dd>
<dt>lazy_sync</dt>
<dd>

([RateLimiterLazySync](#rate-limiter-lazy-sync)) Configuration of lazy-syncing behaviour of ratelimiter

</dd>
<dt>limit_reset_interval</dt>
<dd>

(string, default: `60s`) Time after which the limit for a given label value will be reset.

</dd>
</dl>

### v1RateLimiterDynamicConfig {#v1-rate-limiter-dynamic-config}

Dynamic Configuration for the rate limiter

#### Properties

<dl>
<dt>overrides</dt>
<dd>

([[]RateLimiterOverride](#rate-limiter-override)) Allows to specify different limits for particular label values.

</dd>
</dl>

### v1RateLimiterIns {#v1-rate-limiter-ins}

Inputs for the RateLimiter component

#### Properties

<dl>
<dt>limit</dt>
<dd>

([V1InPort](#v1-in-port), `required`) Number of flows allowed per _limit_reset_interval_ per each label.
Negative values disable the ratelimiter.

:::tip
Negative limit can be useful to _conditionally_ enable the ratelimiter
under certain circumstances. [Decider](#v1-decider) might be helpful.
:::

</dd>
</dl>

### v1Resources {#v1-resources}

Resources that need to be setup for the policy to function

:::info
See also [Resources overview](/concepts/policy/resources.md).
:::

Resources are typically Flux Meters, Classifiers, etc. that can be used to create on-demand metrics or label the flows.

#### Properties

<dl>
<dt>classifiers</dt>
<dd>

([[]V1Classifier](#v1-classifier)) Classifiers are installed in the data-plane and are used to label the requests based on payload content.

The flow labels created by Classifiers can be matched by Flux Meters to create metrics for control purposes.

</dd>
<dt>flux_meters</dt>
<dd>

(map of [V1FluxMeter](#v1-flux-meter)) Flux Meters are installed in the data-plane and form the observability leg of the feedback loop.

Flux Meter created metrics can be consumed as input to the circuit via the PromQL component.

</dd>
</dl>

### v1Rule {#v1-rule}

Rule describes a single Flow Classification Rule

Flow classification rule extracts a value from request metadata.
More specifically, from `input`, which has the same spec as [Envoy's External Authorization Attribute Context][attribute-context].
See https://play.openpolicyagent.org/p/gU7vcLkc70 for an example input.
There are two ways to define a flow classification rule:

- Using a declarative extractor – suitable from simple cases, such as directly reading a value from header or a field from json body.
- Rego expression.

Performance note: It's recommended to use declarative extractors where possible, as they may be slightly performant than Rego expressions.

Example of Declarative JSON extractor:

```yaml
extractor:
  json:
    from: request.http.body
    pointer: /user/name
```

Example of Rego module which also disables telemetry visibility of label:

```yaml
rego:
  query: data.user_from_cookie.user
  source: |
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
telemetry: false
```

[attribute-context]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto

#### Properties

<dl>
<dt>extractor</dt>
<dd>

([V1Extractor](#v1-extractor)) High-level declarative extractor.

</dd>
<dt>rego</dt>
<dd>

([RuleRego](#rule-rego)) Rego module to extract a value from.

</dd>
<dt>telemetry</dt>
<dd>

(bool, `required`) Decides if the created flow label should be available as an attribute in OLAP telemetry and
propagated in [baggage](/concepts/flow-control/flow-label.md#baggage)

:::note
The flow label is always accessible in Aperture Policies regardless of this setting.
:::

:::caution
When using [FluxNinja Cloud plugin](cloud/plugin.md), telemetry enabled
labels are sent to FluxNinha Cloud for observability. Telemetry should be disabled for
sensitive labels.
:::

</dd>
</dl>

### v1Scheduler {#v1-scheduler}

Weighted Fair Queuing-based workload scheduler

:::note
Each Agent instantiates an independent copy of the scheduler, but output
signals for accepted and incoming concurrency are aggregated across all agents.
:::

See [ConcurrencyLimiter](#v1-concurrency-limiter) for more context.

#### Properties

<dl>
<dt>auto_tokens</dt>
<dd>

(bool, default: `true`) Automatically estimate the size of a request in each workload, based on
historical latency. Each workload's `tokens` will be set to average
latency of flows in that workload during last few seconds (exact duration
of this average can change).

</dd>
<dt>default_workload_parameters</dt>
<dd>

([SchedulerWorkloadParameters](#scheduler-workload-parameters), `required`) WorkloadParameters to be used if none of workloads specified in `workloads` match.

</dd>
<dt>max_timeout</dt>
<dd>

(string, default: `0.49s`) Max Timeout is the value with which the flow timeout calculated by `timeout_factor` is capped

:::caution
This timeout needs to be strictly less than the timeout set on the
client for the whole GRPC call:

- in case of envoy, timeout set on `grpc_service` used in `ext_authz` filter,
- in case of libraries, timeout configured... TODO.

We're using fail-open logic in integrations, so if the GRPC timeout
fires first, the flow will end up being unconditionally allowed while
it're still waiting on the scheduler.

To avoid such cases, the end-to-end GRPC timeout should also contain
some headroom for constant overhead like serialization, etc. Default
value for GRPC timeouts is 500ms, giving 50ms of headeroom, so when
tweaking this timeout, make sure to adjust the GRPC timeout accordingly.
:::

</dd>
<dt>out_ports</dt>
<dd>

([V1SchedulerOuts](#v1-scheduler-outs)) Output ports for the Scheduler component.

</dd>
<dt>timeout_factor</dt>
<dd>

(float64, `gte=0.0`, default: `0.5`) Timeout as a factor of tokens for a flow in a workload

If a flow is not able to get tokens within `timeout_factor` \* `tokens` of duration,
it will be rejected.

This value impacts the prioritization and fairness because the larger the timeout the higher the chance a request has to get scheduled.

</dd>
<dt>workloads</dt>
<dd>

([[]SchedulerWorkload](#scheduler-workload)) List of workloads to be used in scheduler.

Categorizing [flows](/concepts/flow-control/flow-control.md#flow) into workloads
allows for load-shedding to be "smarter" than just "randomly deny 50% of
requests". There are two aspects of this "smartness":

- Scheduler can more precisely calculate concurrency if it understands
  that flows belonging to different classes have different weights (eg.
  inserts vs lookups).
- Setting different priorities to different workloads lets the scheduler
  avoid dropping important traffic during overload.

Each workload in this list specifies also a matcher that's used to
determine which flow will be categorized into which workload.
In case of multiple matching workloads, the first matching one will be used.
If none of workloads match, `default_workload` will be used.

:::info
See also [workload definition in the concepts
section](/concepts/flow-control/concurrency-limiter.md#workload).
:::

</dd>
</dl>

### v1SchedulerOuts {#v1-scheduler-outs}

Output for the Scheduler component.

#### Properties

<dl>
<dt>accepted_concurrency</dt>
<dd>

([V1OutPort](#v1-out-port)) Accepted concurrency is the number of accepted tokens per second.

:::info
**Accepted tokens** are tokens associated with
[flows](/concepts/flow-control/flow-control.md#flow) that were accepted by
this scheduler. Number of tokens for a flow is determined by a
[workload parameters](#scheduler-workload-parameters) that the flow was assigned to (either
via `auto_tokens` or explicitly by `Workload.tokens`).
:::

Value of this signal is the sum across all the relevant schedulers.

</dd>
<dt>incoming_concurrency</dt>
<dd>

([V1OutPort](#v1-out-port)) Incoming concurrency is the number of incoming tokens/sec.
This is the same as `accepted_concurrency`, but across all the flows
entering scheduler, including rejected ones.

</dd>
</dl>

### v1ServiceSelector {#v1-service-selector}

Describes which service a [flow control or observability
component](/concepts/flow-control/flow-control.md#components) should apply
to

:::info
See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).
:::

#### Properties

<dl>
<dt>agent_group</dt>
<dd>

(string, default: `default`) Which [agent-group](/concepts/service.md#agent-group) this
selector applies to.

</dd>
<dt>service</dt>
<dd>

(string) The Fully Qualified Domain Name of the
[service](/concepts/service.md) to select.

In kubernetes, this is the FQDN of the Service object.

Empty string means all services within an agent group (catch-all).

:::note
One entity may belong to multiple services.
:::

</dd>
</dl>

### v1Sink {#v1-sink}

Sink is a component that consumes input signals and does nothing with them

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1SinkIns](#v1-sink-ins)) Input ports for the Sink component.

</dd>
</dl>

### v1SinkIns {#v1-sink-ins}

Inputs for the Sink component.

#### Properties

<dl>
<dt>inputs</dt>
<dd>

([[]V1InPort](#v1-in-port)) Array of input signals.

</dd>
</dl>

### v1Sqrt {#v1-sqrt}

Takes an input signal and emits the square root of it multiplied by scale as an output

$$
\text{output} = \text{scale} \sqrt{\text{input}}
$$

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1SqrtIns](#v1-sqrt-ins)) Input ports for the Sqrt component.

</dd>
<dt>out_ports</dt>
<dd>

([V1SqrtOuts](#v1-sqrt-outs)) Output ports for the Sqrt component.

</dd>
<dt>scale</dt>
<dd>

(float64, default: `1`) Scaling factor to be multiplied with the square root of the input signal.

</dd>
</dl>

### v1SqrtIns {#v1-sqrt-ins}

Inputs for the Sqrt component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1InPort](#v1-in-port)) Input signal.

</dd>
</dl>

### v1SqrtOuts {#v1-sqrt-outs}

Outputs for the Sqrt component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Output signal.

</dd>
</dl>

### v1Switcher {#v1-switcher}

Type of combinator that switches between `on_true` and `on_false` signals based on switch input

`on_true` will be returned if switch input is valid and not equal to 0.0 ,
otherwise `on_false` will be returned.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1SwitcherIns](#v1-switcher-ins)) Input ports for the Switcher component.

</dd>
<dt>out_ports</dt>
<dd>

([V1SwitcherOuts](#v1-switcher-outs)) Output ports for the Switcher component.

</dd>
</dl>

### v1SwitcherIns {#v1-switcher-ins}

Inputs for the Switcher component.

#### Properties

<dl>
<dt>on_false</dt>
<dd>

([V1InPort](#v1-in-port)) Output signal when switch is invalid or 0.0.

</dd>
<dt>on_true</dt>
<dd>

([V1InPort](#v1-in-port)) Output signal when switch is valid and not 0.0.

</dd>
<dt>switch</dt>
<dd>

([V1InPort](#v1-in-port)) Decides whether to return on_true or on_false.

</dd>
</dl>

### v1SwitcherOuts {#v1-switcher-outs}

Outputs for the Switcher component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1OutPort](#v1-out-port)) Selected signal (on_true or on_false).

</dd>
</dl>

<!---
Generated File Ends
-->
