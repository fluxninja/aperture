# Aperture Controller Configuration Reference

## Table of contents

### POLICIES CONFIGURATION

| Key | Reference         |
| --- | ----------------- |
|     | [Policy](#policy) |

### Object Index

- [MatchExpressionList](#match-expression-list) – List of MatchExpressions that is used for all/any matching.
  eg. {any: {of: [expr…
- [RateLimiterLazySyncConfig](#rate-limiter-lazy-sync-config)
- [RateLimiterOverrideConfig](#rate-limiter-override-config)
- [SchedulerWorkloadConfig](#scheduler-workload-config) – Configuration that provides scheduling parameters such as priority for the given…
- [WorkloadConfigDefaultWorkload](#workload-config-default-workload) – Workload defines a class of requests that preferably have similar properties suc…
- [WorkloadConfigWorkload](#workload-config-workload) – Workload defines a class of requests that preferably have similar properties suc…
- [v1ArithmeticCombinator](#v1-arithmetic-combinator) – Type of combinator that computes the arithmetic operation on the operand signals…
- [v1ArithmeticCombinatorIns](#v1-arithmetic-combinator-ins) – Inputs for the Arithmetic Combinator component.
- [v1ArithmeticCombinatorOuts](#v1-arithmetic-combinator-outs) – Outputs for the Arithmetic Combinator component.
- [v1Component](#v1-component) – Computational blocks that form the circuit.
  Signals flow into the components via…
- [v1ConcurrencyLimiter](#v1-concurrency-limiter) – Concurrency Limiter is an actuator component that regulates flows in order to pr…
- [v1Constant](#v1-constant) – Component that emits a constant value as an output signal.
- [v1ConstantOuts](#v1-constant-outs) – Outputs for the Constant component.
- [v1ControlPoint](#v1-control-point) – Identifies control point within a service that the rule or policy should apply t…
- [v1Decider](#v1-decider) – Type of combinator that computes the comparison operation on lhs and rhs signals…
- [v1DeciderIns](#v1-decider-ins) – Inputs for the Decider component.
- [v1DeciderOuts](#v1-decider-outs) – Outputs for the Decider component.
- [v1EMA](#v1-e-m-a) – Exponential Moving Average (EMA) is a type of moving average that applies expone…
- [v1EMAIns](#v1-e-m-a-ins) – Inputs for the EMA component.
- [v1EMAOuts](#v1-e-m-a-outs) – Outputs for the EMA component.
- [v1EqualsMatchExpression](#v1-equals-match-expression) – Label selector expression of the equal form "label == value".
- [v1Extrapolator](#v1-extrapolator) – Extrapolates the input signal by repeating the last valid value during the perio…
- [v1ExtrapolatorIns](#v1-extrapolator-ins) – Inputs for the Extrapolator component.
- [v1ExtrapolatorOuts](#v1-extrapolator-outs) – Outputs for the Extrapolator component.
- [v1FluxMeter](#v1-flux-meter) – FluxMeter gathers metrics for the traffic that matches its selector.
- [v1GradientController](#v1-gradient-controller) – Describes the gradient values which is computed as follows: gradient = (setpoint…
- [v1GradientControllerIns](#v1-gradient-controller-ins) – Inputs for the Gradient Controller component.
- [v1GradientControllerOuts](#v1-gradient-controller-outs) – Outputs for the Gradient Controller component.
- [v1K8sLabelMatcherRequirement](#v1-k8s-label-matcher-requirement) – Label selector requirement which is a selector that contains values, a key, and …
- [v1LabelMatcher](#v1-label-matcher) – Allows to define rules whether a map of labels should be considered a match or n…
- [v1LoadShedActuator](#v1-load-shed-actuator) – Takes the load shed factor input signal and publishes it to the schedulers in th…
- [v1LoadShedActuatorIns](#v1-load-shed-actuator-ins) – Input for the Load Shed Actuator component.
- [v1MatchExpression](#v1-match-expression) – Defines a [map<string, string> → bool] expression to be evaluated on labels.
  …
- [v1MatchesMatchExpression](#v1-matches-match-expression) – Label selector expression of the matches form "label matches regex".
- [v1Max](#v1-max) – Takes a list of input signals and emits the signal with the maximum value.
  Max: …
- [v1MaxIns](#v1-max-ins) – Inputs for the Max component.
- [v1MaxOuts](#v1-max-outs) – Output for the Max component.
- [v1Min](#v1-min) – Takes an array of input signals and emits the signal with the minimum value.
  Min…
- [v1MinIns](#v1-min-ins) – Inputs for the Min component.
- [v1MinOuts](#v1-min-outs) – Output ports for the Min component.
- [v1Policy](#v1-policy) – Policy is defined as a dataflow graph (circuit) of inter-connected components.
  S…
- [v1Port](#v1-port) – Components are interconnected with each other via Ports.
- [v1PromQL](#v1-prom-q-l) – Component that runs a Prometheus query periodically and returns the result as an…
- [v1PromQLOuts](#v1-prom-q-l-outs) – Output for the PromQL component.
- [v1RateLimiter](#v1-rate-limiter)
- [v1RateLimiterIns](#v1-rate-limiter-ins)
- [v1Scheduler](#v1-scheduler) – Weighted Fair Queuing based workload scheduler.
- [v1SchedulerOuts](#v1-scheduler-outs) – Output for the Scheduler component.
- [v1Selector](#v1-selector) – Describes where a rule or actuation component should apply to.
- [v1Sqrt](#v1-sqrt) – Takes an input signal and emits the square root of it multiplied by scale as an …
- [v1SqrtIns](#v1-sqrt-ins) – Inputs for the Sqrt component.
- [v1SqrtOuts](#v1-sqrt-outs) – Outputs for the Sqrt component.

## Reference

### <span id="policy"></span> _Policy_

#### Members

<dl>

<dt></dt>
<dd>

Type: [V1Policy](#v1-policy)

</dd>
</dl>

## Objects

### <span id="match-expression-list"></span> MatchExpressionList

List of MatchExpressions that is used for all/any matching.
eg. {any: {of: [expr1, expr2]}}.

#### Properties

<dl>
<dt>of</dt>
<dd>

([[]V1MatchExpression](#v1-match-expression)) List of subexpressions of the match expression.

</dd>
</dl>

### <span id="rate-limiter-lazy-sync-config"></span> RateLimiterLazySyncConfig

#### Properties

<dl>
<dt>enabled</dt>
<dd>

(bool)

</dd>
</dl>
<dl>
<dt>num_sync</dt>
<dd>

(int64) Number of times to lazy sync within the limit_reset_interval

</dd>
</dl>

### <span id="rate-limiter-override-config"></span> RateLimiterOverrideConfig

#### Properties

<dl>
<dt>label_value</dt>
<dd>

(string)

</dd>
</dl>
<dl>
<dt>limit_scale_factor</dt>
<dd>

(float64)

</dd>
</dl>

### <span id="scheduler-workload-config"></span> SchedulerWorkloadConfig

Configuration that provides scheduling parameters such as priority for the given workloads.

#### Properties

<dl>
<dt>auto_tokens</dt>
<dd>

(bool)

</dd>
</dl>
<dl>
<dt>default_workload</dt>
<dd>

([WorkloadConfigDefaultWorkload](#workload-config-default-workload))

</dd>
</dl>
<dl>
<dt>label_key</dt>
<dd>

(string) TODO: this entire section is being reworked by @TanveerGill

</dd>
</dl>
<dl>
<dt>workloads</dt>
<dd>

([[]WorkloadConfigWorkload](#workload-config-workload)) list of workloads
workload can describe priority, tokens (if auto_tokens are set to false) and timeout

</dd>
</dl>

### <span id="workload-config-default-workload"></span> WorkloadConfigDefaultWorkload

Workload defines a class of requests that preferably have similar properties such as response latency.

#### Properties

<dl>
<dt>priority</dt>
<dd>

(int64, `gte=0,lte=255`) Describes priority level of the requests within the workload.
Priority level ranges from 0 to 255.
Higher numbers means higher priority level.

</dd>
</dl>
<dl>
<dt>timeout</dt>
<dd>

(string, default: `0.005s`) Timeout override decides how long a request in the workload can wait for tokens. This value impacts the fairness because the larger the timeout the higher the chance a request has to get scheduled.

</dd>
</dl>
<dl>
<dt>tokens</dt>
<dd>

(string, default: `1`) Tokens determines the cost of admitting a single request the workload, which is typically defined as milliseconds of response latency. This override is applicable only if auto_tokens is set to false.

</dd>
</dl>

### <span id="workload-config-workload"></span> WorkloadConfigWorkload

Workload defines a class of requests that preferably have similar properties such as response latency.

#### Properties

<dl>
<dt>label_value</dt>
<dd>

(string) Value of label for specified workload label key.

</dd>
</dl>
<dl>
<dt>priority</dt>
<dd>

(int64, `gte=0,lte=255`) Describes priority level of the requests within the workload.
Priority level ranges from 0 to 255.
Higher numbers means higher priority level.

</dd>
</dl>
<dl>
<dt>timeout</dt>
<dd>

(string, default: `0.005s`) Timeout override decides how long a request in the workload can wait for tokens. This value impacts the fairness because the larger the timeout the higher the chance a request has to get scheduled.

</dd>
</dl>
<dl>
<dt>tokens</dt>
<dd>

(string, default: `1`) Tokens determines the cost of admitting a single request the workload, which is typically defined as milliseconds of response latency. This override is applicable only if auto_tokens is set to false.

</dd>
</dl>

### <span id="v1-arithmetic-combinator"></span> v1ArithmeticCombinator

Type of combinator that computes the arithmetic operation on the operand signals.
The arithmetic operation can be addition, subtraction, multiplication, division, XOR, right bit shift or left bit shift.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1ArithmeticCombinatorIns](#v1-arithmetic-combinator-ins)) Input ports for the Arithmetic Combinator component.

</dd>
</dl>
<dl>
<dt>operator</dt>
<dd>

(string, `oneof=add sub mul div xor lshift rshift`) Operator of the arithmetic operation.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1ArithmeticCombinatorOuts](#v1-arithmetic-combinator-outs)) Output ports for the Arithmetic Combinator component.

</dd>
</dl>

### <span id="v1-arithmetic-combinator-ins"></span> v1ArithmeticCombinatorIns

Inputs for the Arithmetic Combinator component.

#### Properties

<dl>
<dt>lhs</dt>
<dd>

([V1Port](#v1-port)) Left hand side of the arithmetic operation.

</dd>
</dl>
<dl>
<dt>rhs</dt>
<dd>

([V1Port](#v1-port)) Right hand side of the arithmetic operation.

</dd>
</dl>

### <span id="v1-arithmetic-combinator-outs"></span> v1ArithmeticCombinatorOuts

Outputs for the Arithmetic Combinator component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Result of arithmetic operation.

</dd>
</dl>

### <span id="v1-component"></span> v1Component

Computational blocks that form the circuit.
Signals flow into the components via input ports and results are emitted on output ports.
Components are wired to each other based on signal names forming an execution graph of the circuit.
Loops are broken by the runtime at the earliest component index that is part of the loop.
The looped signals are saved in the tick they are generated and served in the subsequent tick.

#### Properties

<dl>
<dt>arithmetic_combinator</dt>
<dd>

([V1ArithmeticCombinator](#v1-arithmetic-combinator)) Applies the given operator on input operands (signals) and emits the result.

</dd>
</dl>
<dl>
<dt>concurrency_limiter</dt>
<dd>

([V1ConcurrencyLimiter](#v1-concurrency-limiter)) Concurrency Limiter provides service protection by applying prioritized load shedding of flows using a network scheduler (e.g. Weighted Fair Queuing).

</dd>
</dl>
<dl>
<dt>constant</dt>
<dd>

([V1Constant](#v1-constant)) Emits a constant signal.

</dd>
</dl>
<dl>
<dt>decider</dt>
<dd>

([V1Decider](#v1-decider)) Decider acts as a switch that emits one of the two signals based on the binary result of comparison operator on two operands.

</dd>
</dl>
<dl>
<dt>ema</dt>
<dd>

([V1EMA](#v1-e-m-a)) Exponential Moving Average filter.

</dd>
</dl>
<dl>
<dt>extrapolator</dt>
<dd>

([V1Extrapolator](#v1-extrapolator)) Takes an input signal and emits the extrapolated value; either mirroring the input value or repeating the last known value up to the maximum extrapolation interval.

</dd>
</dl>
<dl>
<dt>gradient_controller</dt>
<dd>

([V1GradientController](#v1-gradient-controller)) Gradient controller basically calculates the ratio between the signal and the setpoint to determine the magnitude of the correction that need to be applied. This controller can be used to build AIMD (Additive Increase, Multiplicative Decrease) or MIMD style response.

</dd>
</dl>
<dl>
<dt>max</dt>
<dd>

([V1Max](#v1-max)) Emits the maximum of the input siganls.

</dd>
</dl>
<dl>
<dt>min</dt>
<dd>

([V1Min](#v1-min)) Emits the minimum of the input signals.

</dd>
</dl>
<dl>
<dt>promql</dt>
<dd>

([V1PromQL](#v1-prom-q-l)) Periodically runs a Prometheus query in the background and emits the result.

</dd>
</dl>
<dl>
<dt>rate_limiter</dt>
<dd>

([V1RateLimiter](#v1-rate-limiter)) Rate Limiter provides service protection by applying rate limiter.

</dd>
</dl>
<dl>
<dt>sqrt</dt>
<dd>

([V1Sqrt](#v1-sqrt)) Takes an input signal and emits the square root of the input signal.

</dd>
</dl>

### <span id="v1-concurrency-limiter"></span> v1ConcurrencyLimiter

Concurrency Limiter is an actuator component that regulates flows in order to provide active service protection. It is based on the actuation strategy (e.g. load shed) and workload scheduling which is based on Weighted Fair Queuing principles. Concurrency is calculated in terms of total tokens which translate to (avg. latency \* inflight requests), i.e. Little's Law.

#### Properties

<dl>
<dt>load_shed_actuator</dt>
<dd>

([V1LoadShedActuator](#v1-load-shed-actuator)) Actuator based on load shedding a portion of requests.

</dd>
</dl>
<dl>
<dt>scheduler</dt>
<dd>

([V1Scheduler](#v1-scheduler), `required`) Weighted Fair Queuing based workfload scheduler.

</dd>
</dl>

### <span id="v1-constant"></span> v1Constant

Component that emits a constant value as an output signal.

#### Properties

<dl>
<dt>out_ports</dt>
<dd>

([V1ConstantOuts](#v1-constant-outs)) Output ports for the Constant component.

</dd>
</dl>
<dl>
<dt>value</dt>
<dd>

(float64) The constant value to be emitted.

</dd>
</dl>

### <span id="v1-constant-outs"></span> v1ConstantOuts

Outputs for the Constant component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) The constant value is emitted to the output port.

</dd>
</dl>

### <span id="v1-control-point"></span> v1ControlPoint

Identifies control point within a service that the rule or policy should apply to.
Controlpoint is either a library feature name or one of "ingress/egress" traffic control point.

#### Properties

<dl>
<dt>feature</dt>
<dd>

(string, `required`) Name of FlunxNinja library's feature.
Feature corresponds to a block of code that can be "switched off" which usually is a "named opentelemetry's Span".

Note: Flowcontrol only.

</dd>
</dl>
<dl>
<dt>traffic</dt>
<dd>

(string, `required,oneof=ingress egress`) Type of traffic service, either "ingress" or "egress".
Apply the policy to the whole incoming/outgoing traffic of a service.
Usually powered by integration with a proxy (like envoy) or a web framework.

- Flowcontrol: Blockable atom here is a single HTTP-transaction.
- Classification: Apply the classification rules to every incoming/outgoing request and attach the resulting flow labels to baggage and telemetry.

</dd>
</dl>

### <span id="v1-decider"></span> v1Decider

Type of combinator that computes the comparison operation on lhs and rhs signals and switches between on_true and on_false signals based on the result of the comparison.
The comparison operator can be greater-than, less-than, greater-than-or-equal, less-than-or-equal, equal, or not-equal.
This component also supports time-based response, i.e. the output transitions between on_true or on_false signal if the decider condition is true or false for at least "positive_for" or "negative_for" duration.
If true_for and false_for durations are zero then the transitions are instantaneous.

#### Properties

<dl>
<dt>false_for</dt>
<dd>

(string, default: `0s`) Duration of time to wait before a transition to false state.
If the duration is zero, the transition will happen instantaneously.

</dd>
</dl>
<dl>
<dt>in_ports</dt>
<dd>

([V1DeciderIns](#v1-decider-ins)) Input ports for the Decider component.

</dd>
</dl>
<dl>
<dt>operator</dt>
<dd>

(string, `oneof=gt lt gte lte eq neq`) Comparison operator that computes operation on lhs and rhs input signals.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1DeciderOuts](#v1-decider-outs)) Output ports for the Decider component.

</dd>
</dl>
<dl>
<dt>true_for</dt>
<dd>

(string, default: `0s`) Duration of time to wait before a transition to true state.
If the duration is zero, the transition will happen instantaneously.

</dd>
</dl>

### <span id="v1-decider-ins"></span> v1DeciderIns

Inputs for the Decider component.

#### Properties

<dl>
<dt>lhs</dt>
<dd>

([V1Port](#v1-port)) Left hand side input signal for the comparison operation.

</dd>
</dl>
<dl>
<dt>on_false</dt>
<dd>

([V1Port](#v1-port)) Output signal when the result of the operation is false.

</dd>
</dl>
<dl>
<dt>on_true</dt>
<dd>

([V1Port](#v1-port)) Output signal when the result of the operation is true.

</dd>
</dl>
<dl>
<dt>rhs</dt>
<dd>

([V1Port](#v1-port)) Right hand side input signal for the comparison operation.

</dd>
</dl>

### <span id="v1-decider-outs"></span> v1DeciderOuts

Outputs for the Decider component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Selected signal (on_true or on_false).

</dd>
</dl>

### <span id="v1-e-m-a"></span> v1EMA

Exponential Moving Average (EMA) is a type of moving average that applies exponenially more weight to recent signal readings.
At any time EMA component operates in one of the following states:

1. Warm up state: The first warm_up_window samples are used to compute the initial EMA.
   If an invalid reading is received during the warm_up_window, the last good average is emitted and the state gets reset back to beginning of Warm up state.
2. Normal state: The EMA is computed using following formula.
   If an invalid reading is received continuously for ema_window during the EMA stage, the last good EMA is emitted and the state gets reset back to Warm up state.
   The EMA for a series Y is calculated recursively as:
   EMAt = Y0 , t = 0
   EMAt = alpha \* Yt + (1 - alpha) \* EMAt-1 , t > 0
   The coefficient alpha represents the degree of weighting decrease, a constant smoothing factor between 0 and 1.
   A higher alpha discounts older observations faster.
   The alpha is computed using ema_window: alpha = 2 / (N + 1). Where, N = ema_window / evalutation_period.

#### Properties

<dl>
<dt>correction_factor_on_max_envelope_violation</dt>
<dd>

(float64, `gte=0,lte=1.0`, default: `1`) Correction factor to apply on the output value if its in violation of the max envelope.

</dd>
</dl>
<dl>
<dt>correction_factor_on_min_envelope_violation</dt>
<dd>

(float64, `gte=1.0`, default: `1`) Correction factor to apply on the output value if its in violation of the min envelope.

</dd>
</dl>
<dl>
<dt>ema_window</dt>
<dd>

(string, default: `5s`) Duration of EMA sampling window.

</dd>
</dl>
<dl>
<dt>in_ports</dt>
<dd>

([V1EMAIns](#v1-e-m-a-ins)) Input ports for the EMA component.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1EMAOuts](#v1-e-m-a-outs)) Output ports for the EMA component.

</dd>
</dl>
<dl>
<dt>warm_up_window</dt>
<dd>

(string, default: `0s`) Duration of EMA warming up window.
The initial value of the EMA is the average of signal readings received during the warm up window.

</dd>
</dl>

### <span id="v1-e-m-a-ins"></span> v1EMAIns

Inputs for the EMA component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1Port](#v1-port)) Input signal to be used for the EMA computation.

</dd>
</dl>
<dl>
<dt>max_envelope</dt>
<dd>

([V1Port](#v1-port)) Upper bound of the moving average.

</dd>
</dl>
<dl>
<dt>min_envelope</dt>
<dd>

([V1Port](#v1-port)) Lower bound of the moving average.

</dd>
</dl>

### <span id="v1-e-m-a-outs"></span> v1EMAOuts

Outputs for the EMA component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Exponential moving average of the series of reading as an output signal.

</dd>
</dl>

### <span id="v1-equals-match-expression"></span> v1EqualsMatchExpression

Label selector expression of the equal form "label == value".

#### Properties

<dl>
<dt>label</dt>
<dd>

(string, `required`) Name of the label to equal match the value.

</dd>
</dl>
<dl>
<dt>value</dt>
<dd>

(string) Exact value that the label should be equal to.

</dd>
</dl>

### <span id="v1-extrapolator"></span> v1Extrapolator

Extrapolates the input signal by repeating the last valid value during the period in which it is invalid. It does so until maximum_extrapolation_interval is reached, beyond which it emits invalid signal unless input signal becomes valid again.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1ExtrapolatorIns](#v1-extrapolator-ins)) Input ports for the Extrapolator component.

</dd>
</dl>
<dl>
<dt>max_extrapolation_interval</dt>
<dd>

(string, default: `10s`) Maximum time interval to repeat the last valid value of input signal.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1ExtrapolatorOuts](#v1-extrapolator-outs)) Output ports for the Extrapolator component.

</dd>
</dl>

### <span id="v1-extrapolator-ins"></span> v1ExtrapolatorIns

Inputs for the Extrapolator component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1Port](#v1-port)) Input signal for the Extrapolator component.

</dd>
</dl>

### <span id="v1-extrapolator-outs"></span> v1ExtrapolatorOuts

Outputs for the Extrapolator component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Extrapolated signal.

</dd>
</dl>

### <span id="v1-flux-meter"></span> v1FluxMeter

FluxMeter gathers metrics for the traffic that matches its selector.

#### Properties

<dl>
<dt>histogram_buckets</dt>
<dd>

([]float64, default: `[5.0,10.0,25.0,50.0,100.0,250.0,500.0,1000.0,2500.0,5000.0,10000.0]`) Latency histogram buckets (in ms) for this FluxMeter.

</dd>
</dl>
<dl>
<dt>name</dt>
<dd>

(string) Name of the flux meter.

</dd>
</dl>
<dl>
<dt>selector</dt>
<dd>

([V1Selector](#v1-selector)) Policies are only applied to flows that are matched based on the fields in the selector.

</dd>
</dl>

### <span id="v1-gradient-controller"></span> v1GradientController

Describes the gradient values which is computed as follows: gradient = (setpoint)/(signal) \* tolerance.
Limits gradient to range [min_gradient, max_gradient].
Output: (gradient \* control_variable) + optimize.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1GradientControllerIns](#v1-gradient-controller-ins)) Input ports of the Gradient Controller.

</dd>
</dl>
<dl>
<dt>max_gradient</dt>
<dd>

(float64, default: `1.7976931348623157e+308`) Maximum gradient which clamps the computed gradient value to the range, [min_gradient, max_gradient].

</dd>
</dl>
<dl>
<dt>min_gradient</dt>
<dd>

(float64, default: `-1.7976931348623157e+308`) Minimum gradient which clamps the computed gradient value to the range, [min_gradient, max_gradient].

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1GradientControllerOuts](#v1-gradient-controller-outs)) Output ports of the Gradient Controller.

</dd>
</dl>
<dl>
<dt>tolerance</dt>
<dd>

(float64, `gte=0.0`) Tolerance of the gradient controller beyond which the correction is made.

</dd>
</dl>

### <span id="v1-gradient-controller-ins"></span> v1GradientControllerIns

Inputs for the Gradient Controller component.

#### Properties

<dl>
<dt>control_variable</dt>
<dd>

([V1Port](#v1-port)) Control variable is multiplied by the gradient to produce the output.

</dd>
</dl>
<dl>
<dt>max</dt>
<dd>

([V1Port](#v1-port)) Maximum value to limit the gradient.

</dd>
</dl>
<dl>
<dt>min</dt>
<dd>

([V1Port](#v1-port)) Minimum value to limit the gradient.

</dd>
</dl>
<dl>
<dt>optimize</dt>
<dd>

([V1Port](#v1-port)) Optimize signal is added to the output of the gradient calculation.

</dd>
</dl>
<dl>
<dt>setpoint</dt>
<dd>

([V1Port](#v1-port)) Setpoint to be used for the gradient computation.

</dd>
</dl>
<dl>
<dt>signal</dt>
<dd>

([V1Port](#v1-port)) Signal to be used for the gradient computation.

</dd>
</dl>

### <span id="v1-gradient-controller-outs"></span> v1GradientControllerOuts

Outputs for the Gradient Controller component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Computed desired value of the control variable.

</dd>
</dl>

### <span id="v1-k8s-label-matcher-requirement"></span> v1K8sLabelMatcherRequirement

Label selector requirement which is a selector that contains values, a key, and an operator that relates the key and values.

#### Properties

<dl>
<dt>key</dt>
<dd>

(string, `required`) Label key that the selector applies to.

</dd>
</dl>
<dl>
<dt>operator</dt>
<dd>

(string, `oneof=In NotIn Exists DoesNotExists`) Logical operator which represents a key's relationship to a set of values.
Valid operators are In, NotIn, Exists and DoesNotExist.

</dd>
</dl>
<dl>
<dt>values</dt>
<dd>

([]string) An array of string values that relates to the key by an operator.
If the operator is In or NotIn, the values array must be non-empty.
If the operator is Exists or DoesNotExist, the values array must be empty.

</dd>
</dl>

### <span id="v1-label-matcher"></span> v1LabelMatcher

Allows to define rules whether a map of labels should be considered a match or not.
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
</dl>
<dl>
<dt>match_expressions</dt>
<dd>

([[]V1K8sLabelMatcherRequirement](#v1-k8s-label-matcher-requirement)) List of k8s-style label matcher requirements.

Note: The requirements are ANDed.

</dd>
</dl>
<dl>
<dt>match_labels</dt>
<dd>

(map of string) A map of {key,value} pairs representing labels to be matched.
A single {key,value} in the matchLabels requires that the label "key" is present and equal to "value".

Note: The requirements are ANDed.

</dd>
</dl>

### <span id="v1-load-shed-actuator"></span> v1LoadShedActuator

Takes the load shed factor input signal and publishes it to the schedulers in the data-plane.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1LoadShedActuatorIns](#v1-load-shed-actuator-ins)) Input ports for the Load Shed Actuator component.

</dd>
</dl>

### <span id="v1-load-shed-actuator-ins"></span> v1LoadShedActuatorIns

Input for the Load Shed Actuator component.

#### Properties

<dl>
<dt>load_shed_factor</dt>
<dd>

([V1Port](#v1-port)) Load shedding factor is a fraction of incoming concurrency (tokens \* requests) that needs to be dropped.

</dd>
</dl>

### <span id="v1-match-expression"></span> v1MatchExpression

Defines a [map<string, string> → bool] expression to be evaluated on labels.
MatchExpression has multiple variants, exactly one should be set.

#### Properties

<dl>
<dt>all</dt>
<dd>

([MatchExpressionList](#match-expression-list)) The expression is true when all subexpressions are true.

</dd>
</dl>
<dl>
<dt>any</dt>
<dd>

([MatchExpressionList](#match-expression-list)) The expression is true when any subexpression is true.

</dd>
</dl>
<dl>
<dt>label_equals</dt>
<dd>

([V1EqualsMatchExpression](#v1-equals-match-expression)) The expression is true when label value equals given value.

</dd>
</dl>
<dl>
<dt>label_exists</dt>
<dd>

(string, `required`) The expression is true when label with given name exists.

</dd>
</dl>
<dl>
<dt>label_matches</dt>
<dd>

([V1MatchesMatchExpression](#v1-matches-match-expression)) The expression is true when label matches given regex.

</dd>
</dl>
<dl>
<dt>not</dt>
<dd>

([V1MatchExpression](#v1-match-expression)) The expression negates the result of subexpression.

</dd>
</dl>

### <span id="v1-matches-match-expression"></span> v1MatchesMatchExpression

Label selector expression of the matches form "label matches regex".

#### Properties

<dl>
<dt>label</dt>
<dd>

(string, `required`) Name of the label to match the regular expression.

</dd>
</dl>
<dl>
<dt>regex</dt>
<dd>

(string, `required`) Regular expression that should match the label value.
It uses [golang's regular expression syntax](https://github.com/google/re2/wiki/Syntax).

</dd>
</dl>

### <span id="v1-max"></span> v1Max

Takes a list of input signals and emits the signal with the maximum value.
Max: output = max([]inputs).

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1MaxIns](#v1-max-ins)) Input ports for the Max component.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1MaxOuts](#v1-max-outs)) Output ports for the Max component.

</dd>
</dl>

### <span id="v1-max-ins"></span> v1MaxIns

Inputs for the Max component.

#### Properties

<dl>
<dt>inputs</dt>
<dd>

([[]V1Port](#v1-port)) Array of input signals.

</dd>
</dl>

### <span id="v1-max-outs"></span> v1MaxOuts

Output for the Max component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Signal with maximum value as an output signal.

</dd>
</dl>

### <span id="v1-min"></span> v1Min

Takes an array of input signals and emits the signal with the minimum value.
Min: output = min([]inputs).

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1MinIns](#v1-min-ins)) Input ports for the Min component.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1MinOuts](#v1-min-outs)) Output ports for the Min component.

</dd>
</dl>

### <span id="v1-min-ins"></span> v1MinIns

Inputs for the Min component.

#### Properties

<dl>
<dt>inputs</dt>
<dd>

([[]V1Port](#v1-port)) Array of input signals.

</dd>
</dl>

### <span id="v1-min-outs"></span> v1MinOuts

Output ports for the Min component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Signal with minimum value as an output signal.

</dd>
</dl>

### <span id="v1-policy"></span> v1Policy

Policy is defined as a dataflow graph (circuit) of inter-connected components.
Signals flow between components via ports.
As signals traverse the circuit, they get processed, stored within components or get acted upon (e.g. load shed, rate-limit, auto-scale etc.).
Policies are evaluated periodically in order to respond to changes in signal readings.

#### Properties

<dl>
<dt>circuit</dt>
<dd>

([[]V1Component](#v1-component)) Defines a signal processing graph as a list of components.

</dd>
</dl>
<dl>
<dt>evaluation_interval</dt>
<dd>

(string, default: `0.5s`) Evaluation interval (ticks) is the time period between consecutive runs of the policy circuit.
This interval is typically aligned with how often the corrective action (actuation) needs to be taken.

</dd>
</dl>
<dl>
<dt>flux_meters</dt>
<dd>

([[]V1FluxMeter](#v1-flux-meter)) FluxMeters are installed in the data-plane and form the observability leg of the feedback loop.

</dd>
</dl>

### <span id="v1-port"></span> v1Port

Components are interconnected with each other via Ports.

#### Properties

<dl>
<dt>signal_name</dt>
<dd>

(string) Name of the incoming or outgoing Signal on the Port.

</dd>
</dl>

### <span id="v1-prom-q-l"></span> v1PromQL

Component that runs a Prometheus query periodically and returns the result as an output signal.

#### Properties

<dl>
<dt>evaluation_interval</dt>
<dd>

(string, default: `10s`) Describes the interval between successive evaluations of the Prometheus query.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1PromQLOuts](#v1-prom-q-l-outs)) Output ports for the PromQL component.

</dd>
</dl>
<dl>
<dt>query_string</dt>
<dd>

(string) Describes the Prometheus query to be run.

</dd>
</dl>

### <span id="v1-prom-q-l-outs"></span> v1PromQLOuts

Output for the PromQL component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) The result of the Prometheus query as an output signal.

</dd>
</dl>

### <span id="v1-rate-limiter"></span> v1RateLimiter

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1RateLimiterIns](#v1-rate-limiter-ins))

</dd>
</dl>
<dl>
<dt>label_key</dt>
<dd>

(string)

</dd>
</dl>
<dl>
<dt>lazy_sync_config</dt>
<dd>

([RateLimiterLazySyncConfig](#rate-limiter-lazy-sync-config))

</dd>
</dl>
<dl>
<dt>limit_reset_interval</dt>
<dd>

(string)

</dd>
</dl>
<dl>
<dt>overrides</dt>
<dd>

([[]RateLimiterOverrideConfig](#rate-limiter-override-config))

</dd>
</dl>
<dl>
<dt>selector</dt>
<dd>

([V1Selector](#v1-selector))

</dd>
</dl>

### <span id="v1-rate-limiter-ins"></span> v1RateLimiterIns

#### Properties

<dl>
<dt>limit</dt>
<dd>

([V1Port](#v1-port)) negative limit means no limit is applied

</dd>
</dl>

### <span id="v1-scheduler"></span> v1Scheduler

Weighted Fair Queuing based workload scheduler.

#### Properties

<dl>
<dt>fairness_key</dt>
<dd>

(string) TODO: this is getting reworked.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1SchedulerOuts](#v1-scheduler-outs)) Output ports for the Scheduler component.

</dd>
</dl>
<dl>
<dt>selector</dt>
<dd>

([V1Selector](#v1-selector)) Selector decides for which service or flows the scheduler will be applied.

</dd>
</dl>
<dl>
<dt>workload_config</dt>
<dd>

([SchedulerWorkloadConfig](#scheduler-workload-config))

</dd>
</dl>

### <span id="v1-scheduler-outs"></span> v1SchedulerOuts

Output for the Scheduler component.

#### Properties

<dl>
<dt>accepted_concurrency</dt>
<dd>

([V1Port](#v1-port)) Accepted concurrency is the number of accepted tokens/sec.

</dd>
</dl>
<dl>
<dt>incoming_concurrency</dt>
<dd>

([V1Port](#v1-port)) Incoming concurrency is the number of incoming tokens/sec.

</dd>
</dl>

### <span id="v1-selector"></span> v1Selector

Describes where a rule or actuation component should apply to.

#### Properties

<dl>
<dt>agent_group</dt>
<dd>

(string) Describes where this selector applies to.

</dd>
</dl>
<dl>
<dt>control_point</dt>
<dd>

([V1ControlPoint](#v1-control-point), `required`) Describes control point Within the entity where the policy should apply to.

</dd>
</dl>
<dl>
<dt>label_matcher</dt>
<dd>

([V1LabelMatcher](#v1-label-matcher)) Allows to add _additional_ condition on labels that must also be satisfied (in addition to namespace+service+control point matching).
The label matcher allows to match on infra labels, flow labels and request labels.
Arbitrary label matcher can be used to match infra labels.
For flowcontrol policies, the matcher can be used to match flow labels.

Note: For classification we can only match flow labels that were created at some **previous** control point.

In case of k8s, infra labels are labels on entities (note: there might exist some additional labels).
Flow label names are always prefixed with `flow_`
Request labels are always prefixed with `request_`.
Available request labels are `id` (available as `request_id`), `method`, `path`, `host`, `scheme`, `size`, `protocol`
(mapped from fields of [HttpRequest](https://github.com/envoyproxy/envoy/blob/637a92a56e2739b5f78441c337171968f18b46ee/api/envoy/service/auth/v3/attribute_context.proto#L102)).
Also, (non-pseudo) headers are available as `request_header_<headername>`.

Note: Request headers are only available for "traffic" control points.

</dd>
</dl>
<dl>
<dt>namespace</dt>
<dd>

(string, `required`) The namespace of the entities to select.
In k8s, this is the k8s' namespace.

</dd>
</dl>
<dl>
<dt>service</dt>
<dd>

(string, `required`) The service (name) of the entities.
In k8s, this is a name of the Service object.

Note: Entity may belong to multiple services.

</dd>
</dl>

### <span id="v1-sqrt"></span> v1Sqrt

Takes an input signal and emits the square root of it multiplied by scale as an output.
Sqrt: output = sqrt(input) \* scale.

#### Properties

<dl>
<dt>in_ports</dt>
<dd>

([V1SqrtIns](#v1-sqrt-ins)) Input ports for the Sqrt component.

</dd>
</dl>
<dl>
<dt>out_ports</dt>
<dd>

([V1SqrtOuts](#v1-sqrt-outs)) Output ports for the Sqrt component.

</dd>
</dl>
<dl>
<dt>scale</dt>
<dd>

(float64, default: `1`) Scaling factor to be multiplied with the square root of the input signal.

</dd>
</dl>

### <span id="v1-sqrt-ins"></span> v1SqrtIns

Inputs for the Sqrt component.

#### Properties

<dl>
<dt>input</dt>
<dd>

([V1Port](#v1-port)) Input signal.

</dd>
</dl>

### <span id="v1-sqrt-outs"></span> v1SqrtOuts

Outputs for the Sqrt component.

#### Properties

<dl>
<dt>output</dt>
<dd>

([V1Port](#v1-port)) Output signal.

</dd>
</dl>

<!---
Generated File Ends
-->
