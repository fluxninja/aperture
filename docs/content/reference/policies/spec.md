---
title: Policy Language Specification
sidebar_position: 1
sidebar_label: Specification
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

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `start`     |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `factor`    |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `count`     |
| Type          | _int32_     |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### FluxMeterExponentialBucketsRange {#flux-meter-exponential-buckets-range}

ExponentialBucketsRange creates `count` number of buckets where the lowest bucket is `min` and the highest
bucket is `max`. The final +inf bucket is not counted.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `min`       |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `max`       |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `count`     |
| Type          | _int32_     |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### FluxMeterLinearBuckets {#flux-meter-linear-buckets}

LinearBuckets creates `count` number of buckets, each `width` wide, where the lowest bucket has an
upper bound of `start`. The final +inf bucket is not counted.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `start`     |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `width`     |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `count`     |
| Type          | _int32_     |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### FluxMeterStaticBuckets {#flux-meter-static-buckets}

StaticBuckets holds the static value of the buckets where latency histogram will be stored.

| <!-- -->      | <!-- -->                                                                           |
| ------------- | ---------------------------------------------------------------------------------- |
| Property      | `buckets`                                                                          |
| Type          | _[]float64_                                                                        |
| Default Value | `, default: `[5.0,10.0,25.0,50.0,100.0,250.0,500.0,1000.0,2500.0,5000.0,10000.0]`` |
| Description   | Lorem Ipsum                                                                        |

### HorizontalPodScalerScaleActuator {#horizontal-pod-scaler-scale-actuator}

| <!-- -->      | <!-- -->                                                                           |
| ------------- | ---------------------------------------------------------------------------------- |
| Property      | `in_ports`                                                                         |
| Type          | _[HorizontalPodScalerScaleActuatorIns](#horizontal-pod-scaler-scale-actuator-ins)_ |
| Default Value | ``                                                                                 |
| Description   | Lorem Ipsum                                                                        |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `dynamic_config_key` |
| Type          | _string_             |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                                                                                |
| ------------- | ------------------------------------------------------------------------------------------------------- |
| Property      | `default_config`                                                                                        |
| Type          | _[HorizontalPodScalerScaleActuatorDynamicConfig](#horizontal-pod-scaler-scale-actuator-dynamic-config)_ |
| Default Value | ``                                                                                                      |
| Description   | Lorem Ipsum                                                                                             |

### HorizontalPodScalerScaleActuatorDynamicConfig {#horizontal-pod-scaler-scale-actuator-dynamic-config}

Dynamic Configuration for ScaleActuator

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `dry_run`   |
| Type          | _bool_      |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### HorizontalPodScalerScaleActuatorIns {#horizontal-pod-scaler-scale-actuator-ins}

Inputs for the HorizontalPodScaler component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `desired_replicas`        |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### HorizontalPodScalerScaleReporter {#horizontal-pod-scaler-scale-reporter}

| <!-- -->      | <!-- -->                                                                             |
| ------------- | ------------------------------------------------------------------------------------ |
| Property      | `out_ports`                                                                          |
| Type          | _[HorizontalPodScalerScaleReporterOuts](#horizontal-pod-scaler-scale-reporter-outs)_ |
| Default Value | ``                                                                                   |
| Description   | Lorem Ipsum                                                                          |

### HorizontalPodScalerScaleReporterOuts {#horizontal-pod-scaler-scale-reporter-outs}

Outputs for the HorizontalPodScaler component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `actual_replicas`           |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `configured_replicas`       |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### MatchExpressionList {#match-expression-list}

List of MatchExpressions that is used for all/any matching

eg. {any: {of: [expr1, expr2]}}.

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `of`                                          |
| Type          | _[[]V1MatchExpression](#v1-match-expression)_ |
| Default Value | ``                                            |
| Description   | Lorem Ipsum                                   |

### ParametersLazySync {#parameters-lazy-sync}

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `enabled`   |
| Type          | _bool_      |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `num_sync`       |
| Type          | _int64_          |
| Default Value | `, default: `5`` |
| Description   | Lorem Ipsum      |

### RateLimiterOverride {#rate-limiter-override}

| <!-- -->      | <!-- -->      |
| ------------- | ------------- |
| Property      | `label_value` |
| Type          | _string_      |
| Default Value | ``            |
| Description   | Lorem Ipsum   |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `limit_scale_factor` |
| Type          | _float64_            |
| Default Value | `, default: `1``     |
| Description   | Lorem Ipsum          |

### RuleRego {#rule-rego}

Raw rego rules are compiled 1:1 to rego queries

High-level extractor-based rules are compiled into a single rego query.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `source`    |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `query`     |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### SchedulerWorkload {#scheduler-workload}

Workload defines a class of requests that preferably have similar properties such as response latency or desired priority.

| <!-- -->      | <!-- -->                                                        |
| ------------- | --------------------------------------------------------------- |
| Property      | `parameters`                                                    |
| Type          | _[SchedulerWorkloadParameters](#scheduler-workload-parameters)_ |
| Default Value | ``                                                              |
| Description   | Lorem Ipsum                                                     |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `label_matcher`                       |
| Type          | _[V1LabelMatcher](#v1-label-matcher)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

### SchedulerWorkloadParameters {#scheduler-workload-parameters}

Parameters defines parameters such as priority, tokens and fairness key that are applicable to flows within a workload.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `priority`  |
| Type          | _int64_     |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `tokens`         |
| Type          | _string_         |
| Default Value | `, default: `1`` |
| Description   | Lorem Ipsum      |

| <!-- -->      | <!-- -->       |
| ------------- | -------------- |
| Property      | `fairness_key` |
| Type          | _string_       |
| Default Value | ``             |
| Description   | Lorem Ipsum    |

### v1AIMDConcurrencyController {#v1-a-i-m-d-concurrency-controller}

High level concurrency control component. Baselines a signal via exponential moving average and applies concurrency limits based on deviation of signal from the baseline. Internally implemented as a nested circuit.

| <!-- -->      | <!-- -->                                                                   |
| ------------- | -------------------------------------------------------------------------- |
| Property      | `in_ports`                                                                 |
| Type          | _[V1AIMDConcurrencyControllerIns](#v1-a-i-m-d-concurrency-controller-ins)_ |
| Default Value | ``                                                                         |
| Description   | Lorem Ipsum                                                                |

| <!-- -->      | <!-- -->                                                                     |
| ------------- | ---------------------------------------------------------------------------- |
| Property      | `out_ports`                                                                  |
| Type          | _[V1AIMDConcurrencyControllerOuts](#v1-a-i-m-d-concurrency-controller-outs)_ |
| Default Value | ``                                                                           |
| Description   | Lorem Ipsum                                                                  |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `flow_selector`                       |
| Type          | _[V1FlowSelector](#v1-flow-selector)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

| <!-- -->      | <!-- -->                                            |
| ------------- | --------------------------------------------------- |
| Property      | `scheduler_parameters`                              |
| Type          | _[V1SchedulerParameters](#v1-scheduler-parameters)_ |
| Default Value | ``                                                  |
| Description   | Lorem Ipsum                                         |

| <!-- -->      | <!-- -->                                                               |
| ------------- | ---------------------------------------------------------------------- |
| Property      | `gradient_parameters`                                                  |
| Type          | _[V1GradientControllerParameters](#v1-gradient-controller-parameters)_ |
| Default Value | ``                                                                     |
| Description   | Lorem Ipsum                                                            |

| <!-- -->      | <!-- -->                       |
| ------------- | ------------------------------ |
| Property      | `concurrency_limit_multiplier` |
| Type          | _float64_                      |
| Default Value | `, default: `2``               |
| Description   | Lorem Ipsum                    |

| <!-- -->      | <!-- -->                       |
| ------------- | ------------------------------ |
| Property      | `concurrency_linear_increment` |
| Type          | _float64_                      |
| Default Value | `, default: `5``               |
| Description   | Lorem Ipsum                    |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `concurrency_sqrt_increment_multiplier` |
| Type          | _float64_                               |
| Default Value | `, default: `1``                        |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                                        |
| ------------- | ----------------------------------------------- |
| Property      | `alerter_parameters`                            |
| Type          | _[V1AlerterParameters](#v1-alerter-parameters)_ |
| Default Value | ``                                              |
| Description   | Lorem Ipsum                                     |

| <!-- -->      | <!-- -->                     |
| ------------- | ---------------------------- |
| Property      | `dry_run_dynamic_config_key` |
| Type          | _string_                     |
| Default Value | ``                           |
| Description   | Lorem Ipsum                  |

### v1AIMDConcurrencyControllerIns {#v1-a-i-m-d-concurrency-controller-ins}

Inputs for the AIMDConcurrencyController component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `signal`                  |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `setpoint`                |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1AIMDConcurrencyControllerOuts {#v1-a-i-m-d-concurrency-controller-outs}

Outputs for the AIMDConcurrencyController component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `is_overload`               |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `load_multiplier`           |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

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

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `from`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1Alerter {#v1-alerter}

Alerter reacts to a signal and generates alert to send to alert manager.

| <!-- -->      | <!-- -->                          |
| ------------- | --------------------------------- |
| Property      | `in_ports`                        |
| Type          | _[V1AlerterIns](#v1-alerter-ins)_ |
| Default Value | ``                                |
| Description   | Lorem Ipsum                       |

| <!-- -->      | <!-- -->                                        |
| ------------- | ----------------------------------------------- |
| Property      | `parameters`                                    |
| Type          | _[V1AlerterParameters](#v1-alerter-parameters)_ |
| Default Value | ``                                              |
| Description   | Lorem Ipsum                                     |

### v1AlerterIns {#v1-alerter-ins}

Inputs for the Alerter component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `signal`                  |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1AlerterParameters {#v1-alerter-parameters}

Alerter Parameters is a common config for separate alerter components and alerters embedded in other components.

| <!-- -->      | <!-- -->     |
| ------------- | ------------ |
| Property      | `alert_name` |
| Type          | _string_     |
| Default Value | ``           |
| Description   | Lorem Ipsum  |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `severity`          |
| Type          | _string_            |
| Default Value | `, default: `info`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `resolve_timeout`   |
| Type          | _string_            |
| Default Value | `, default: `300s`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `alert_channels` |
| Type          | _[]string_       |
| Default Value | ``               |
| Description   | Lorem Ipsum      |

### v1And {#v1-and}

Logical AND.

Signals are mapped to boolean values as follows:

- Zero is treated as false.
- Any non-zero is treated as true.
- Invalid inputs are considered unknown.

  :::note

  Treating invalid inputs as "unknowns" has a consequence that the result
  might end up being valid even when some inputs are invalid. Eg. `unknown && false == false`,
  because the result would end up false no matter if
  first signal was true or false. On the other hand, `unknown && true == unknown`.

  :::

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `in_ports`                |
| Type          | _[V1AndIns](#v1-and-ins)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `out_ports`                 |
| Type          | _[V1AndOuts](#v1-and-outs)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1AndIns {#v1-and-ins}

Inputs for the And component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `inputs`                    |
| Type          | _[[]V1InPort](#v1-in-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1AndOuts {#v1-and-outs}

Output ports for the And component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1ArithmeticCombinator {#v1-arithmetic-combinator}

Type of combinator that computes the arithmetic operation on the operand signals

| <!-- -->      | <!-- -->                                                     |
| ------------- | ------------------------------------------------------------ |
| Property      | `in_ports`                                                   |
| Type          | _[V1ArithmeticCombinatorIns](#v1-arithmetic-combinator-ins)_ |
| Default Value | ``                                                           |
| Description   | Lorem Ipsum                                                  |

| <!-- -->      | <!-- -->                                                       |
| ------------- | -------------------------------------------------------------- |
| Property      | `out_ports`                                                    |
| Type          | _[V1ArithmeticCombinatorOuts](#v1-arithmetic-combinator-outs)_ |
| Default Value | ``                                                             |
| Description   | Lorem Ipsum                                                    |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `operator`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1ArithmeticCombinatorIns {#v1-arithmetic-combinator-ins}

Inputs for the Arithmetic Combinator component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `lhs`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `rhs`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1ArithmeticCombinatorOuts {#v1-arithmetic-combinator-outs}

Outputs for the Arithmetic Combinator component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1AutoScale {#v1-auto-scale}

AutoScale components are used to scale a service.

| <!-- -->      | <!-- -->                                             |
| ------------- | ---------------------------------------------------- |
| Property      | `horizontal_pod_scaler`                              |
| Type          | _[V1HorizontalPodScaler](#v1-horizontal-pod-scaler)_ |
| Default Value | ``                                                   |
| Description   | Lorem Ipsum                                          |

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

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `evaluation_interval` |
| Type          | _string_              |
| Default Value | `, default: `0.5s``   |
| Description   | Lorem Ipsum           |

| <!-- -->      | <!-- -->                         |
| ------------- | -------------------------------- |
| Property      | `components`                     |
| Type          | _[[]V1Component](#v1-component)_ |
| Default Value | ``                               |
| Description   | Lorem Ipsum                      |

### v1Classifier {#v1-classifier}

Set of classification rules sharing a common selector

:::info

See also [Classifier overview](/concepts/integrations/flow-control/flow-classifier.md).

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

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `flow_selector`                       |
| Type          | _[V1FlowSelector](#v1-flow-selector)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `rules`                     |
| Type          | _map of [V1Rule](#v1-rule)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

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
To do so, use the [InPort](#v1-in_port)'s .withConstantSignal(constant_signal) method.
You can also use it to provide special math values such as NaN and +- Inf.
If You need to provide the same constant signal to multiple components,
You can use the [Variable](#v1-variable) component.

:::

See also [Policy](#v1-policy) for a higher-level explanation of circuits.

| <!-- -->      | <!-- -->                                          |
| ------------- | ------------------------------------------------- |
| Property      | `gradient_controller`                             |
| Type          | _[V1GradientController](#v1-gradient-controller)_ |
| Default Value | ``                                                |
| Description   | Lorem Ipsum                                       |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `ema`                |
| Type          | _[V1EMA](#v1-e-m-a)_ |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                              |
| ------------- | ----------------------------------------------------- |
| Property      | `arithmetic_combinator`                               |
| Type          | _[V1ArithmeticCombinator](#v1-arithmetic-combinator)_ |
| Default Value | ``                                                    |
| Description   | Lorem Ipsum                                           |

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `decider`                  |
| Type          | _[V1Decider](#v1-decider)_ |
| Default Value | ``                         |
| Description   | Lorem Ipsum                |

| <!-- -->      | <!-- -->                     |
| ------------- | ---------------------------- |
| Property      | `switcher`                   |
| Type          | _[V1Switcher](#v1-switcher)_ |
| Default Value | ``                           |
| Description   | Lorem Ipsum                  |

| <!-- -->      | <!-- -->                     |
| ------------- | ---------------------------- |
| Property      | `variable`                   |
| Type          | _[V1Variable](#v1-variable)_ |
| Default Value | ``                           |
| Description   | Lorem Ipsum                  |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `sqrt`               |
| Type          | _[V1Sqrt](#v1-sqrt)_ |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                             |
| ------------- | ------------------------------------ |
| Property      | `extrapolator`                       |
| Type          | _[V1Extrapolator](#v1-extrapolator)_ |
| Default Value | ``                                   |
| Description   | Lorem Ipsum                          |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `max`              |
| Type          | _[V1Max](#v1-max)_ |
| Default Value | ``                 |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `min`              |
| Type          | _[V1Min](#v1-min)_ |
| Default Value | ``                 |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->                          |
| ------------- | --------------------------------- |
| Property      | `first_valid`                     |
| Type          | _[V1FirstValid](#v1-first-valid)_ |
| Default Value | ``                                |
| Description   | Lorem Ipsum                       |

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `alerter`                  |
| Type          | _[V1Alerter](#v1-alerter)_ |
| Default Value | ``                         |
| Description   | Lorem Ipsum                |

| <!-- -->      | <!-- -->                         |
| ------------- | -------------------------------- |
| Property      | `integrator`                     |
| Type          | _[V1Integrator](#v1-integrator)_ |
| Default Value | ``                               |
| Description   | Lorem Ipsum                      |

| <!-- -->      | <!-- -->                                 |
| ------------- | ---------------------------------------- |
| Property      | `differentiator`                         |
| Type          | _[V1Differentiator](#v1-differentiator)_ |
| Default Value | ``                                       |
| Description   | Lorem Ipsum                              |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `and`              |
| Type          | _[V1And](#v1-and)_ |
| Default Value | ``                 |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `or`             |
| Type          | _[V1Or](#v1-or)_ |
| Default Value | ``               |
| Description   | Lorem Ipsum      |

| <!-- -->      | <!-- -->                     |
| ------------- | ---------------------------- |
| Property      | `inverter`                   |
| Type          | _[V1Inverter](#v1-inverter)_ |
| Default Value | ``                           |
| Description   | Lorem Ipsum                  |

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `pulse_generator`                         |
| Type          | _[V1PulseGenerator](#v1-pulse-generator)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

| <!-- -->      | <!-- -->                 |
| ------------- | ------------------------ |
| Property      | `holder`                 |
| Type          | _[V1Holder](#v1-holder)_ |
| Default Value | ``                       |
| Description   | Lorem Ipsum              |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `nested_circuit`                        |
| Type          | _[V1NestedCircuit](#v1-nested-circuit)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                                             |
| ------------- | ---------------------------------------------------- |
| Property      | `nested_signal_ingress`                              |
| Type          | _[V1NestedSignalIngress](#v1-nested-signal-ingress)_ |
| Default Value | ``                                                   |
| Description   | Lorem Ipsum                                          |

| <!-- -->      | <!-- -->                                           |
| ------------- | -------------------------------------------------- |
| Property      | `nested_signal_egress`                             |
| Type          | _[V1NestedSignalEgress](#v1-nested-signal-egress)_ |
| Default Value | ``                                                 |
| Description   | Lorem Ipsum                                        |

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `query`                |
| Type          | _[V1Query](#v1-query)_ |
| Default Value | ``                     |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `flow_control`                      |
| Type          | _[V1FlowControl](#v1-flow-control)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

| <!-- -->      | <!-- -->                        |
| ------------- | ------------------------------- |
| Property      | `auto_scale`                    |
| Type          | _[V1AutoScale](#v1-auto-scale)_ |
| Default Value | ``                              |
| Description   | Lorem Ipsum                     |

### v1ConcurrencyLimiter {#v1-concurrency-limiter}

Concurrency Limiter is an actuator component that regulates flows in order to provide active service protection

:::info

See also [Concurrency Limiter overview](/concepts/integrations/flow-control/components/concurrency-limiter.md).

:::

It is based on the actuation strategy (e.g. load actuator) and workload scheduling which is based on Weighted Fair Queuing principles.
Concurrency is calculated in terms of total tokens which translate to (avg. latency \* in-flight requests), i.e. Little's Law.

ConcurrencyLimiter configuration is split into two parts: An actuation
strategy and a scheduler. Right now, only `load_actuator` strategy is available.

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `flow_selector`                       |
| Type          | _[V1FlowSelector](#v1-flow-selector)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

| <!-- -->      | <!-- -->                       |
| ------------- | ------------------------------ |
| Property      | `scheduler`                    |
| Type          | _[V1Scheduler](#v1-scheduler)_ |
| Default Value | ``                             |
| Description   | Lorem Ipsum                    |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `load_actuator`                       |
| Type          | _[V1LoadActuator](#v1-load-actuator)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

### v1ConstantSignal {#v1-constant-signal}

Special constant input for ports and Variable component. Can provide either a constant value or special Nan/+-Inf value.

| <!-- -->      | <!-- -->        |
| ------------- | --------------- |
| Property      | `special_value` |
| Type          | _string_        |
| Default Value | ``              |
| Description   | Lorem Ipsum     |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `value`     |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1Decider {#v1-decider}

Type of combinator that computes the comparison operation on lhs and rhs signals

The comparison operator can be greater-than, less-than, greater-than-or-equal, less-than-or-equal, equal, or not-equal.

This component also supports time-based response, i.e. the output
transitions between 1.0 or 0.0 signal if the decider condition is
true or false for at least "true_for" or "false_for" duration. If
`true_for` and `false_for` durations are zero then the transitions are
instantaneous.

| <!-- -->      | <!-- -->                          |
| ------------- | --------------------------------- |
| Property      | `in_ports`                        |
| Type          | _[V1DeciderIns](#v1-decider-ins)_ |
| Default Value | ``                                |
| Description   | Lorem Ipsum                       |

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `out_ports`                         |
| Type          | _[V1DeciderOuts](#v1-decider-outs)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `operator`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `true_for`        |
| Type          | _string_          |
| Default Value | `, default: `0s`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `false_for`       |
| Type          | _string_          |
| Default Value | `, default: `0s`` |
| Description   | Lorem Ipsum       |

### v1DeciderIns {#v1-decider-ins}

Inputs for the Decider component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `lhs`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `rhs`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1DeciderOuts {#v1-decider-outs}

Outputs for the Decider component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Differentiator {#v1-differentiator}

Differentiator calculates rate of change per tick.

| <!-- -->      | <!-- -->                                        |
| ------------- | ----------------------------------------------- |
| Property      | `in_ports`                                      |
| Type          | _[V1DifferentiatorIns](#v1-differentiator-ins)_ |
| Default Value | ``                                              |
| Description   | Lorem Ipsum                                     |

| <!-- -->      | <!-- -->                                          |
| ------------- | ------------------------------------------------- |
| Property      | `out_ports`                                       |
| Type          | _[V1DifferentiatorOuts](#v1-differentiator-outs)_ |
| Default Value | ``                                                |
| Description   | Lorem Ipsum                                       |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `window`          |
| Type          | _string_          |
| Default Value | `, default: `5s`` |
| Description   | Lorem Ipsum       |

### v1DifferentiatorIns {#v1-differentiator-ins}

Inputs for the Differentiator component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1DifferentiatorOuts {#v1-differentiator-outs}

Outputs for the Differentiator component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1EMA {#v1-e-m-a}

Exponential Moving Average (EMA) is a type of moving average that applies exponentially more weight to recent signal readings

At any time EMA component operates in one of the following states:

1. Warm up state: The first warmup_window samples are used to compute the initial EMA.
   If an invalid reading is received during the warmup_window, the last good average is emitted and the state gets reset back to beginning of Warm up state.
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

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `in_ports`                  |
| Type          | _[V1EMAIns](#v1-e-m-a-ins)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

| <!-- -->      | <!-- -->                      |
| ------------- | ----------------------------- |
| Property      | `out_ports`                   |
| Type          | _[V1EMAOuts](#v1-e-m-a-outs)_ |
| Default Value | ``                            |
| Description   | Lorem Ipsum                   |

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `parameters`                              |
| Type          | _[V1EMAParameters](#v1-e-m-a-parameters)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

### v1EMAIns {#v1-e-m-a-ins}

Inputs for the EMA component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `max_envelope`            |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `min_envelope`            |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1EMAOuts {#v1-e-m-a-outs}

Outputs for the EMA component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1EMAParameters {#v1-e-m-a-parameters}

Parameters for the EMA component.

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `ema_window`      |
| Type          | _string_          |
| Default Value | `, default: `5s`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `warmup_window`   |
| Type          | _string_          |
| Default Value | `, default: `0s`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `correction_factor_on_min_envelope_violation` |
| Type          | _float64_                                     |
| Default Value | `, default: `1``                              |
| Description   | Lorem Ipsum                                   |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `correction_factor_on_max_envelope_violation` |
| Type          | _float64_                                     |
| Default Value | `, default: `1``                              |
| Description   | Lorem Ipsum                                   |

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `valid_during_warmup` |
| Type          | _bool_                |
| Default Value | ``                    |
| Description   | Lorem Ipsum           |

### v1EqualsMatchExpression {#v1-equals-match-expression}

Label selector expression of the equal form "label == value".

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `label`     |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `value`     |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1Extractor {#v1-extractor}

Defines a high-level way to specify how to extract a flow label value given http request metadata, without a need to write rego code

There are multiple variants of extractor, specify exactly one.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `from`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `json`                                  |
| Type          | _[V1JSONExtractor](#v1-json-extractor)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `address`                                     |
| Type          | _[V1AddressExtractor](#v1-address-extractor)_ |
| Default Value | ``                                            |
| Description   | Lorem Ipsum                                   |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `jwt`                                   |
| Type          | _[V1JWTExtractor](#v1-j-w-t-extractor)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                                             |
| ------------- | ---------------------------------------------------- |
| Property      | `path_templates`                                     |
| Type          | _[V1PathTemplateMatcher](#v1-path-template-matcher)_ |
| Default Value | ``                                                   |
| Description   | Lorem Ipsum                                          |

### v1Extrapolator {#v1-extrapolator}

Extrapolates the input signal by repeating the last valid value during the period in which it is invalid

It does so until `maximum_extrapolation_interval` is reached, beyond which it emits invalid signal unless input signal becomes valid again.

| <!-- -->      | <!-- -->                                    |
| ------------- | ------------------------------------------- |
| Property      | `in_ports`                                  |
| Type          | _[V1ExtrapolatorIns](#v1-extrapolator-ins)_ |
| Default Value | ``                                          |
| Description   | Lorem Ipsum                                 |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `out_ports`                                   |
| Type          | _[V1ExtrapolatorOuts](#v1-extrapolator-outs)_ |
| Default Value | ``                                            |
| Description   | Lorem Ipsum                                   |

| <!-- -->      | <!-- -->                                                  |
| ------------- | --------------------------------------------------------- |
| Property      | `parameters`                                              |
| Type          | _[V1ExtrapolatorParameters](#v1-extrapolator-parameters)_ |
| Default Value | ``                                                        |
| Description   | Lorem Ipsum                                               |

### v1ExtrapolatorIns {#v1-extrapolator-ins}

Inputs for the Extrapolator component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1ExtrapolatorOuts {#v1-extrapolator-outs}

Outputs for the Extrapolator component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1ExtrapolatorParameters {#v1-extrapolator-parameters}

Parameters for the Extrapolator component.

| <!-- -->      | <!-- -->                     |
| ------------- | ---------------------------- |
| Property      | `max_extrapolation_interval` |
| Type          | _string_                     |
| Default Value | `, default: `10s``           |
| Description   | Lorem Ipsum                  |

### v1FirstValid {#v1-first-valid}

Picks the first valid input signal from the array of input signals and emits it as an output signal

| <!-- -->      | <!-- -->                                 |
| ------------- | ---------------------------------------- |
| Property      | `in_ports`                               |
| Type          | _[V1FirstValidIns](#v1-first-valid-ins)_ |
| Default Value | ``                                       |
| Description   | Lorem Ipsum                              |

| <!-- -->      | <!-- -->                                   |
| ------------- | ------------------------------------------ |
| Property      | `out_ports`                                |
| Type          | _[V1FirstValidOuts](#v1-first-valid-outs)_ |
| Default Value | ``                                         |
| Description   | Lorem Ipsum                                |

### v1FirstValidIns {#v1-first-valid-ins}

Inputs for the FirstValid component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `inputs`                    |
| Type          | _[[]V1InPort](#v1-in-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1FirstValidOuts {#v1-first-valid-outs}

Outputs for the FirstValid component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1FlowControl {#v1-flow-control}

FlowControl components are used to regulate requests flow.

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `rate_limiter`                      |
| Type          | _[V1RateLimiter](#v1-rate-limiter)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

| <!-- -->      | <!-- -->                                          |
| ------------- | ------------------------------------------------- |
| Property      | `concurrency_limiter`                             |
| Type          | _[V1ConcurrencyLimiter](#v1-concurrency-limiter)_ |
| Default Value | ``                                                |
| Description   | Lorem Ipsum                                       |

| <!-- -->      | <!-- -->                                                            |
| ------------- | ------------------------------------------------------------------- |
| Property      | `aimd_concurrency_controller`                                       |
| Type          | _[V1AIMDConcurrencyController](#v1-a-i-m-d-concurrency-controller)_ |
| Default Value | ``                                                                  |
| Description   | Lorem Ipsum                                                         |

### v1FlowMatcher {#v1-flow-matcher}

Describes which flows a [flow control
component](/concepts/integrations/flow-control/flow-control.md#components) should apply
to

:::info

See also [FlowSelector overview](/concepts/integrations/flow-control/flow-selector.md).

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

| <!-- -->      | <!-- -->        |
| ------------- | --------------- |
| Property      | `control_point` |
| Type          | _string_        |
| Default Value | ``              |
| Description   | Lorem Ipsum     |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `label_matcher`                       |
| Type          | _[V1LabelMatcher](#v1-label-matcher)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

### v1FlowSelector {#v1-flow-selector}

Describes which flow in which service a [flow control
component](/concepts/integrations/flow-control/flow-control.md#components) should apply
to

:::info

See also [FlowSelector overview](/concepts/integrations/flow-control/flow-selector.md).

:::

| <!-- -->      | <!-- -->                                    |
| ------------- | ------------------------------------------- |
| Property      | `service_selector`                          |
| Type          | _[V1ServiceSelector](#v1-service-selector)_ |
| Default Value | ``                                          |
| Description   | Lorem Ipsum                                 |

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `flow_matcher`                      |
| Type          | _[V1FlowMatcher](#v1-flow-matcher)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

### v1FluxMeter {#v1-flux-meter}

Flux Meter gathers metrics for the traffic that matches its selector.
The histogram created by Flux Meter measures the workload latency by default.

:::info

See also [Flux Meter overview](/concepts/integrations/flow-control/flux-meter.md).

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

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `flow_selector`                       |
| Type          | _[V1FlowSelector](#v1-flow-selector)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

| <!-- -->      | <!-- -->                                               |
| ------------- | ------------------------------------------------------ |
| Property      | `static_buckets`                                       |
| Type          | _[FluxMeterStaticBuckets](#flux-meter-static-buckets)_ |
| Default Value | ``                                                     |
| Description   | Lorem Ipsum                                            |

| <!-- -->      | <!-- -->                                               |
| ------------- | ------------------------------------------------------ |
| Property      | `linear_buckets`                                       |
| Type          | _[FluxMeterLinearBuckets](#flux-meter-linear-buckets)_ |
| Default Value | ``                                                     |
| Description   | Lorem Ipsum                                            |

| <!-- -->      | <!-- -->                                                         |
| ------------- | ---------------------------------------------------------------- |
| Property      | `exponential_buckets`                                            |
| Type          | _[FluxMeterExponentialBuckets](#flux-meter-exponential-buckets)_ |
| Default Value | ``                                                               |
| Description   | Lorem Ipsum                                                      |

| <!-- -->      | <!-- -->                                                                    |
| ------------- | --------------------------------------------------------------------------- |
| Property      | `exponential_buckets_range`                                                 |
| Type          | _[FluxMeterExponentialBucketsRange](#flux-meter-exponential-buckets-range)_ |
| Default Value | ``                                                                          |
| Description   | Lorem Ipsum                                                                 |

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `attribute_key`                     |
| Type          | _string_                            |
| Default Value | `, default: `workload_duration_ms`` |
| Description   | Lorem Ipsum                         |

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

| <!-- -->      | <!-- -->                                                 |
| ------------- | -------------------------------------------------------- |
| Property      | `in_ports`                                               |
| Type          | _[V1GradientControllerIns](#v1-gradient-controller-ins)_ |
| Default Value | ``                                                       |
| Description   | Lorem Ipsum                                              |

| <!-- -->      | <!-- -->                                                   |
| ------------- | ---------------------------------------------------------- |
| Property      | `out_ports`                                                |
| Type          | _[V1GradientControllerOuts](#v1-gradient-controller-outs)_ |
| Default Value | ``                                                         |
| Description   | Lorem Ipsum                                                |

| <!-- -->      | <!-- -->                                                               |
| ------------- | ---------------------------------------------------------------------- |
| Property      | `parameters`                                                           |
| Type          | _[V1GradientControllerParameters](#v1-gradient-controller-parameters)_ |
| Default Value | ``                                                                     |
| Description   | Lorem Ipsum                                                            |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `dynamic_config_key` |
| Type          | _string_             |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                                                      |
| ------------- | ----------------------------------------------------------------------------- |
| Property      | `default_config`                                                              |
| Type          | _[V1GradientControllerDynamicConfig](#v1-gradient-controller-dynamic-config)_ |
| Default Value | ``                                                                            |
| Description   | Lorem Ipsum                                                                   |

### v1GradientControllerDynamicConfig {#v1-gradient-controller-dynamic-config}

Dynamic Configuration for a Controller

| <!-- -->      | <!-- -->      |
| ------------- | ------------- |
| Property      | `manual_mode` |
| Type          | _bool_        |
| Default Value | ``            |
| Description   | Lorem Ipsum   |

### v1GradientControllerIns {#v1-gradient-controller-ins}

Inputs for the Gradient Controller component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `signal`                  |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `setpoint`                |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `optimize`                |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `max`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `min`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `control_variable`        |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1GradientControllerOuts {#v1-gradient-controller-outs}

Outputs for the Gradient Controller component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1GradientControllerParameters {#v1-gradient-controller-parameters}

Gradient Parameters.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `slope`     |
| Type          | _float64_   |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `min_gradient`                          |
| Type          | _float64_                               |
| Default Value | `, default: `-1.7976931348623157e+308`` |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                               |
| ------------- | -------------------------------------- |
| Property      | `max_gradient`                         |
| Type          | _float64_                              |
| Default Value | `, default: `1.7976931348623157e+308`` |
| Description   | Lorem Ipsum                            |

### v1Holder {#v1-holder}

Holds the last valid signal value for the specified duration then waits for next valid value to hold.
If it's holding a value that means it ignores both valid and invalid new signals until the hold_for duration is finished.

| <!-- -->      | <!-- -->                        |
| ------------- | ------------------------------- |
| Property      | `in_ports`                      |
| Type          | _[V1HolderIns](#v1-holder-ins)_ |
| Default Value | ``                              |
| Description   | Lorem Ipsum                     |

| <!-- -->      | <!-- -->                          |
| ------------- | --------------------------------- |
| Property      | `out_ports`                       |
| Type          | _[V1HolderOuts](#v1-holder-outs)_ |
| Default Value | ``                                |
| Description   | Lorem Ipsum                       |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `hold_for`        |
| Type          | _string_          |
| Default Value | `, default: `5s`` |
| Description   | Lorem Ipsum       |

### v1HolderIns {#v1-holder-ins}

Inputs for the Holder component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1HolderOuts {#v1-holder-outs}

Outputs for the Holder component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1HorizontalPodScaler {#v1-horizontal-pod-scaler}

| <!-- -->      | <!-- -->                                                       |
| ------------- | -------------------------------------------------------------- |
| Property      | `kubernetes_object_selector`                                   |
| Type          | _[V1KubernetesObjectSelector](#v1-kubernetes-object-selector)_ |
| Default Value | ``                                                             |
| Description   | Lorem Ipsum                                                    |

| <!-- -->      | <!-- -->                                                                    |
| ------------- | --------------------------------------------------------------------------- |
| Property      | `scale_reporter`                                                            |
| Type          | _[HorizontalPodScalerScaleReporter](#horizontal-pod-scaler-scale-reporter)_ |
| Default Value | ``                                                                          |
| Description   | Lorem Ipsum                                                                 |

| <!-- -->      | <!-- -->                                                                    |
| ------------- | --------------------------------------------------------------------------- |
| Property      | `scale_actuator`                                                            |
| Type          | _[HorizontalPodScalerScaleActuator](#horizontal-pod-scaler-scale-actuator)_ |
| Default Value | ``                                                                          |
| Description   | Lorem Ipsum                                                                 |

### v1InPort {#v1-in-port}

Components receive input from other components via InPorts

| <!-- -->      | <!-- -->      |
| ------------- | ------------- |
| Property      | `signal_name` |
| Type          | _string_      |
| Default Value | ``            |
| Description   | Lorem Ipsum   |

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `constant_signal`                         |
| Type          | _[V1ConstantSignal](#v1-constant-signal)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

### v1Integrator {#v1-integrator}

Accumulates sum of signal every tick.

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `in_ports`                              |
| Type          | _[V1IntegratorIns](#v1-integrator-ins)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `out_ports`                               |
| Type          | _[V1IntegratorOuts](#v1-integrator-outs)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

### v1IntegratorIns {#v1-integrator-ins}

Inputs for the Integrator component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `reset`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `min`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `max`                     |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1IntegratorOuts {#v1-integrator-outs}

Outputs for the Integrator component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Inverter {#v1-inverter}

Logical NOT.

See [And component](#v1-and) on how signals are mapped onto boolean values.

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `in_ports`                          |
| Type          | _[V1InverterIns](#v1-inverter-ins)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `out_ports`                           |
| Type          | _[V1InverterOuts](#v1-inverter-outs)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

### v1InverterIns {#v1-inverter-ins}

Inputs for the Inverter component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1InverterOuts {#v1-inverter-outs}

Output ports for the Inverter component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1JSONExtractor {#v1-json-extractor}

Deserialize a json, and extract one of the fields

Example:

```yaml
from: request.http.body
pointer: /user/name
```

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `from`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `pointer`   |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

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

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `from`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->       |
| ------------- | -------------- |
| Property      | `json_pointer` |
| Type          | _string_       |
| Default Value | ``             |
| Description   | Lorem Ipsum    |

### v1K8sLabelMatcherRequirement {#v1-k8s-label-matcher-requirement}

Label selector requirement which is a selector that contains values, a key, and an operator that relates the key and values.

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `key`       |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `operator`  |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `values`    |
| Type          | _[]string_  |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1KubernetesObjectSelector {#v1-kubernetes-object-selector}

Describes which pods a control or observability
component should apply to.

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `agent_group`          |
| Type          | _string_               |
| Default Value | `, default: `default`` |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `namespace` |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->      |
| ------------- | ------------- |
| Property      | `api_version` |
| Type          | _string_      |
| Default Value | ``            |
| Description   | Lorem Ipsum   |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `kind`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `name`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1LabelMatcher {#v1-label-matcher}

Allows to define rules whether a map of
[labels](/concepts/integrations/flow-control/flow-label.md)
should be considered a match or not

It provides three ways to define requirements:

- matchLabels
- matchExpressions
- arbitrary expression

If multiple requirements are set, they are all ANDed.
An empty label matcher always matches.

| <!-- -->      | <!-- -->        |
| ------------- | --------------- |
| Property      | `match_labels`  |
| Type          | _map of string_ |
| Default Value | ``              |
| Description   | Lorem Ipsum     |

| <!-- -->      | <!-- -->                                                              |
| ------------- | --------------------------------------------------------------------- |
| Property      | `match_expressions`                                                   |
| Type          | _[[]V1K8sLabelMatcherRequirement](#v1-k8s-label-matcher-requirement)_ |
| Default Value | ``                                                                    |
| Description   | Lorem Ipsum                                                           |

| <!-- -->      | <!-- -->                                    |
| ------------- | ------------------------------------------- |
| Property      | `expression`                                |
| Type          | _[V1MatchExpression](#v1-match-expression)_ |
| Default Value | ``                                          |
| Description   | Lorem Ipsum                                 |

### v1LoadActuator {#v1-load-actuator}

Takes the load multiplier input signal and publishes it to the schedulers in the data-plane

| <!-- -->      | <!-- -->                                     |
| ------------- | -------------------------------------------- |
| Property      | `in_ports`                                   |
| Type          | _[V1LoadActuatorIns](#v1-load-actuator-ins)_ |
| Default Value | ``                                           |
| Description   | Lorem Ipsum                                  |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `dynamic_config_key` |
| Type          | _string_             |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                                          |
| ------------- | ----------------------------------------------------------------- |
| Property      | `default_config`                                                  |
| Type          | _[V1LoadActuatorDynamicConfig](#v1-load-actuator-dynamic-config)_ |
| Default Value | ``                                                                |
| Description   | Lorem Ipsum                                                       |

| <!-- -->      | <!-- -->                                        |
| ------------- | ----------------------------------------------- |
| Property      | `alerter_parameters`                            |
| Type          | _[V1AlerterParameters](#v1-alerter-parameters)_ |
| Default Value | ``                                              |
| Description   | Lorem Ipsum                                     |

### v1LoadActuatorDynamicConfig {#v1-load-actuator-dynamic-config}

Dynamic Configuration for LoadActuator

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `dry_run`   |
| Type          | _bool_      |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1LoadActuatorIns {#v1-load-actuator-ins}

Input for the Load Actuator component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `load_multiplier`         |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

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

| <!-- -->      | <!-- -->                                    |
| ------------- | ------------------------------------------- |
| Property      | `not`                                       |
| Type          | _[V1MatchExpression](#v1-match-expression)_ |
| Default Value | ``                                          |
| Description   | Lorem Ipsum                                 |

| <!-- -->      | <!-- -->                                        |
| ------------- | ----------------------------------------------- |
| Property      | `all`                                           |
| Type          | _[MatchExpressionList](#match-expression-list)_ |
| Default Value | ``                                              |
| Description   | Lorem Ipsum                                     |

| <!-- -->      | <!-- -->                                        |
| ------------- | ----------------------------------------------- |
| Property      | `any`                                           |
| Type          | _[MatchExpressionList](#match-expression-list)_ |
| Default Value | ``                                              |
| Description   | Lorem Ipsum                                     |

| <!-- -->      | <!-- -->       |
| ------------- | -------------- |
| Property      | `label_exists` |
| Type          | _string_       |
| Default Value | ``             |
| Description   | Lorem Ipsum    |

| <!-- -->      | <!-- -->                                                 |
| ------------- | -------------------------------------------------------- |
| Property      | `label_equals`                                           |
| Type          | _[V1EqualsMatchExpression](#v1-equals-match-expression)_ |
| Default Value | ``                                                       |
| Description   | Lorem Ipsum                                              |

| <!-- -->      | <!-- -->                                                   |
| ------------- | ---------------------------------------------------------- |
| Property      | `label_matches`                                            |
| Type          | _[V1MatchesMatchExpression](#v1-matches-match-expression)_ |
| Default Value | ``                                                         |
| Description   | Lorem Ipsum                                                |

### v1MatchesMatchExpression {#v1-matches-match-expression}

Label selector expression of the matches form "label matches regex".

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `label`     |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `regex`     |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1Max {#v1-max}

Takes a list of input signals and emits the signal with the maximum value

Max: output = max([]inputs).

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `in_ports`                |
| Type          | _[V1MaxIns](#v1-max-ins)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `out_ports`                 |
| Type          | _[V1MaxOuts](#v1-max-outs)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1MaxIns {#v1-max-ins}

Inputs for the Max component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `inputs`                    |
| Type          | _[[]V1InPort](#v1-in-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1MaxOuts {#v1-max-outs}

Output for the Max component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Min {#v1-min}

Takes an array of input signals and emits the signal with the minimum value
Min: output = min([]inputs).

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `in_ports`                |
| Type          | _[V1MinIns](#v1-min-ins)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `out_ports`                 |
| Type          | _[V1MinOuts](#v1-min-outs)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1MinIns {#v1-min-ins}

Inputs for the Min component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `inputs`                    |
| Type          | _[[]V1InPort](#v1-in-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1MinOuts {#v1-min-outs}

Output ports for the Min component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1NestedCircuit {#v1-nested-circuit}

Nested circuit defines a sub-circuit as a high-level component. It consists of a list of components and a map of input and output ports.

| <!-- -->      | <!-- -->                         |
| ------------- | -------------------------------- |
| Property      | `in_ports_map`                   |
| Type          | _map of [V1InPort](#v1-in-port)_ |
| Default Value | ``                               |
| Description   | Lorem Ipsum                      |

| <!-- -->      | <!-- -->                           |
| ------------- | ---------------------------------- |
| Property      | `out_ports_map`                    |
| Type          | _map of [V1OutPort](#v1-out-port)_ |
| Default Value | ``                                 |
| Description   | Lorem Ipsum                        |

| <!-- -->      | <!-- -->                         |
| ------------- | -------------------------------- |
| Property      | `components`                     |
| Type          | _[[]V1Component](#v1-component)_ |
| Default Value | ``                               |
| Description   | Lorem Ipsum                      |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `name`      |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `short_description` |
| Type          | _string_            |
| Default Value | ``                  |
| Description   | Lorem Ipsum         |

### v1NestedSignalEgress {#v1-nested-signal-egress}

Nested signal egress is a special type of component that allows to extract a signal from a nested circuit.

| <!-- -->      | <!-- -->                                                  |
| ------------- | --------------------------------------------------------- |
| Property      | `in_ports`                                                |
| Type          | _[V1NestedSignalEgressIns](#v1-nested-signal-egress-ins)_ |
| Default Value | ``                                                        |
| Description   | Lorem Ipsum                                               |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `port_name` |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1NestedSignalEgressIns {#v1-nested-signal-egress-ins}

Inputs for the NestedSignalEgress component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `signal`                  |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1NestedSignalIngress {#v1-nested-signal-ingress}

Nested signal ingress is a special type of component that allows to inject a signal into a nested circuit.

| <!-- -->      | <!-- -->                                                      |
| ------------- | ------------------------------------------------------------- |
| Property      | `out_ports`                                                   |
| Type          | _[V1NestedSignalIngressOuts](#v1-nested-signal-ingress-outs)_ |
| Default Value | ``                                                            |
| Description   | Lorem Ipsum                                                   |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `port_name` |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1NestedSignalIngressOuts {#v1-nested-signal-ingress-outs}

Outputs for the NestedSignalIngress component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `signal`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Or {#v1-or}

Logical OR.

See [And component](#v1-and) on how signals are mapped onto boolean values.

| <!-- -->      | <!-- -->                |
| ------------- | ----------------------- |
| Property      | `in_ports`              |
| Type          | _[V1OrIns](#v1-or-ins)_ |
| Default Value | ``                      |
| Description   | Lorem Ipsum             |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `out_ports`               |
| Type          | _[V1OrOuts](#v1-or-outs)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1OrIns {#v1-or-ins}

Inputs for the Or component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `inputs`                    |
| Type          | _[[]V1InPort](#v1-in-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1OrOuts {#v1-or-outs}

Output ports for the Or component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1OutPort {#v1-out-port}

Components produce output for other components via OutPorts

| <!-- -->      | <!-- -->      |
| ------------- | ------------- |
| Property      | `signal_name` |
| Type          | _string_      |
| Default Value | ``            |
| Description   | Lorem Ipsum   |

### v1PathTemplateMatcher {#v1-path-template-matcher}

Matches HTTP Path to given path templates

HTTP path will be matched against given path templates.
If a match occurs, the value associated with the path template will be treated as a result.
In case of multiple path templates matching, the most specific one will be chosen.

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `template_values` |
| Type          | _map of string_   |
| Default Value | ``                |
| Description   | Lorem Ipsum       |

### v1Policy {#v1-policy}

Policy expresses reliability automation workflow that automatically protects services

:::info

See also [Policy overview](/concepts/policy/policy.md).

:::

Policy specification contains a circuit that defines the controller logic and resources that need to be setup.

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `circuit`                  |
| Type          | _[V1Circuit](#v1-circuit)_ |
| Default Value | ``                         |
| Description   | Lorem Ipsum                |

| <!-- -->      | <!-- -->                       |
| ------------- | ------------------------------ |
| Property      | `resources`                    |
| Type          | _[V1Resources](#v1-resources)_ |
| Default Value | ``                             |
| Description   | Lorem Ipsum                    |

### v1PromQL {#v1-prom-q-l}

Component that runs a Prometheus query periodically and returns the result as an output signal

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `out_ports`                         |
| Type          | _[V1PromQLOuts](#v1-prom-q-l-outs)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

| <!-- -->      | <!-- -->       |
| ------------- | -------------- |
| Property      | `query_string` |
| Type          | _string_       |
| Default Value | ``             |
| Description   | Lorem Ipsum    |

| <!-- -->      | <!-- -->              |
| ------------- | --------------------- |
| Property      | `evaluation_interval` |
| Type          | _string_              |
| Default Value | `, default: `10s``    |
| Description   | Lorem Ipsum           |

### v1PromQLOuts {#v1-prom-q-l-outs}

Output for the PromQL component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1PulseGenerator {#v1-pulse-generator}

Generates 0 and 1 in turns.

| <!-- -->      | <!-- -->                                           |
| ------------- | -------------------------------------------------- |
| Property      | `out_ports`                                        |
| Type          | _[V1PulseGeneratorOuts](#v1-pulse-generator-outs)_ |
| Default Value | ``                                                 |
| Description   | Lorem Ipsum                                        |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `true_for`        |
| Type          | _string_          |
| Default Value | `, default: `5s`` |
| Description   | Lorem Ipsum       |

| <!-- -->      | <!-- -->          |
| ------------- | ----------------- |
| Property      | `false_for`       |
| Type          | _string_          |
| Default Value | `, default: `5s`` |
| Description   | Lorem Ipsum       |

### v1PulseGeneratorOuts {#v1-pulse-generator-outs}

Outputs for the PulseGenerator component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Query {#v1-query}

Query components that are query databases such as Prometheus.

| <!-- -->      | <!-- -->                   |
| ------------- | -------------------------- |
| Property      | `promql`                   |
| Type          | _[V1PromQL](#v1-prom-q-l)_ |
| Default Value | ``                         |
| Description   | Lorem Ipsum                |

### v1RateLimiter {#v1-rate-limiter}

Limits the traffic on a control point to specified rate

:::info

See also [Rate Limiter overview](/concepts/integrations/flow-control/components/rate-limiter.md).

:::

Ratelimiting is done separately on per-label-value basis. Use _label_key_
to select which label should be used as key.

| <!-- -->      | <!-- -->                                   |
| ------------- | ------------------------------------------ |
| Property      | `in_ports`                                 |
| Type          | _[V1RateLimiterIns](#v1-rate-limiter-ins)_ |
| Default Value | ``                                         |
| Description   | Lorem Ipsum                                |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `flow_selector`                       |
| Type          | _[V1FlowSelector](#v1-flow-selector)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

| <!-- -->      | <!-- -->                                                 |
| ------------- | -------------------------------------------------------- |
| Property      | `parameters`                                             |
| Type          | _[V1RateLimiterParameters](#v1-rate-limiter-parameters)_ |
| Default Value | ``                                                       |
| Description   | Lorem Ipsum                                              |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `dynamic_config_key` |
| Type          | _string_             |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                                        |
| ------------- | --------------------------------------------------------------- |
| Property      | `default_config`                                                |
| Type          | _[V1RateLimiterDynamicConfig](#v1-rate-limiter-dynamic-config)_ |
| Default Value | ``                                                              |
| Description   | Lorem Ipsum                                                     |

### v1RateLimiterDynamicConfig {#v1-rate-limiter-dynamic-config}

Dynamic Configuration for the rate limiter

| <!-- -->      | <!-- -->                                          |
| ------------- | ------------------------------------------------- |
| Property      | `overrides`                                       |
| Type          | _[[]RateLimiterOverride](#rate-limiter-override)_ |
| Default Value | ``                                                |
| Description   | Lorem Ipsum                                       |

### v1RateLimiterIns {#v1-rate-limiter-ins}

Inputs for the RateLimiter component

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `limit`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1RateLimiterParameters {#v1-rate-limiter-parameters}

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `limit_reset_interval` |
| Type          | _string_               |
| Default Value | `, default: `60s``     |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `label_key` |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

| <!-- -->      | <!-- -->                                      |
| ------------- | --------------------------------------------- |
| Property      | `lazy_sync`                                   |
| Type          | _[ParametersLazySync](#parameters-lazy-sync)_ |
| Default Value | ``                                            |
| Description   | Lorem Ipsum                                   |

### v1Resources {#v1-resources}

Resources that need to be setup for the policy to function

:::info

See also [Resources overview](/concepts/policy/resources.md).

:::

Resources are typically Flux Meters, Classifiers, etc. that can be used to create on-demand metrics or label the flows.

| <!-- -->      | <!-- -->                               |
| ------------- | -------------------------------------- |
| Property      | `flux_meters`                          |
| Type          | _map of [V1FluxMeter](#v1-flux-meter)_ |
| Default Value | ``                                     |
| Description   | Lorem Ipsum                            |

| <!-- -->      | <!-- -->                           |
| ------------- | ---------------------------------- |
| Property      | `classifiers`                      |
| Type          | _[[]V1Classifier](#v1-classifier)_ |
| Default Value | ``                                 |
| Description   | Lorem Ipsum                        |

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

| <!-- -->      | <!-- -->                       |
| ------------- | ------------------------------ |
| Property      | `extractor`                    |
| Type          | _[V1Extractor](#v1-extractor)_ |
| Default Value | ``                             |
| Description   | Lorem Ipsum                    |

| <!-- -->      | <!-- -->                 |
| ------------- | ------------------------ |
| Property      | `rego`                   |
| Type          | _[RuleRego](#rule-rego)_ |
| Default Value | ``                       |
| Description   | Lorem Ipsum              |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `telemetry` |
| Type          | _bool_      |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1Scheduler {#v1-scheduler}

Weighted Fair Queuing-based workload scheduler

:::note

Each Agent instantiates an independent copy of the scheduler, but output
signals for accepted and incoming concurrency are aggregated across all agents.

:::

See [ConcurrencyLimiter](#v1-concurrency-limiter) for more context.

| <!-- -->      | <!-- -->                                |
| ------------- | --------------------------------------- |
| Property      | `out_ports`                             |
| Type          | _[V1SchedulerOuts](#v1-scheduler-outs)_ |
| Default Value | ``                                      |
| Description   | Lorem Ipsum                             |

| <!-- -->      | <!-- -->                                            |
| ------------- | --------------------------------------------------- |
| Property      | `parameters`                                        |
| Type          | _[V1SchedulerParameters](#v1-scheduler-parameters)_ |
| Default Value | ``                                                  |
| Description   | Lorem Ipsum                                         |

### v1SchedulerOuts {#v1-scheduler-outs}

Output for the Scheduler component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `accepted_concurrency`      |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `incoming_concurrency`      |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1SchedulerParameters {#v1-scheduler-parameters}

Scheduler parameters

| <!-- -->      | <!-- -->                                     |
| ------------- | -------------------------------------------- |
| Property      | `workloads`                                  |
| Type          | _[[]SchedulerWorkload](#scheduler-workload)_ |
| Default Value | ``                                           |
| Description   | Lorem Ipsum                                  |

| <!-- -->      | <!-- -->                                                        |
| ------------- | --------------------------------------------------------------- |
| Property      | `default_workload_parameters`                                   |
| Type          | _[SchedulerWorkloadParameters](#scheduler-workload-parameters)_ |
| Default Value | ``                                                              |
| Description   | Lorem Ipsum                                                     |

| <!-- -->      | <!-- -->            |
| ------------- | ------------------- |
| Property      | `auto_tokens`       |
| Type          | _bool_              |
| Default Value | `, default: `true`` |
| Description   | Lorem Ipsum         |

| <!-- -->      | <!-- -->           |
| ------------- | ------------------ |
| Property      | `timeout_factor`   |
| Type          | _float64_          |
| Default Value | `, default: `0.5`` |
| Description   | Lorem Ipsum        |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `max_timeout`        |
| Type          | _string_             |
| Default Value | `, default: `0.49s`` |
| Description   | Lorem Ipsum          |

### v1ServiceSelector {#v1-service-selector}

Describes which service a [flow control or observability
component](/concepts/integrations/flow-control/flow-control.md#components) should apply
to

:::info

See also [FlowSelector overview](/concepts/integrations/flow-control/flow-selector.md).

:::

| <!-- -->      | <!-- -->               |
| ------------- | ---------------------- |
| Property      | `agent_group`          |
| Type          | _string_               |
| Default Value | `, default: `default`` |
| Description   | Lorem Ipsum            |

| <!-- -->      | <!-- -->    |
| ------------- | ----------- |
| Property      | `service`   |
| Type          | _string_    |
| Default Value | ``          |
| Description   | Lorem Ipsum |

### v1Sqrt {#v1-sqrt}

Takes an input signal and emits the square root of it multiplied by scale as an output

$$
\text{output} = \text{scale} \sqrt{\text{input}}
$$

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `in_ports`                  |
| Type          | _[V1SqrtIns](#v1-sqrt-ins)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

| <!-- -->      | <!-- -->                      |
| ------------- | ----------------------------- |
| Property      | `out_ports`                   |
| Type          | _[V1SqrtOuts](#v1-sqrt-outs)_ |
| Default Value | ``                            |
| Description   | Lorem Ipsum                   |

| <!-- -->      | <!-- -->         |
| ------------- | ---------------- |
| Property      | `scale`          |
| Type          | _float64_        |
| Default Value | `, default: `1`` |
| Description   | Lorem Ipsum      |

### v1SqrtIns {#v1-sqrt-ins}

Inputs for the Sqrt component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `input`                   |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1SqrtOuts {#v1-sqrt-outs}

Outputs for the Sqrt component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Switcher {#v1-switcher}

Type of combinator that switches between `on_true` and `on_false` signals based on switch input

`on_true` will be returned if switch input is valid and not equal to 0.0 ,
otherwise `on_false` will be returned.

| <!-- -->      | <!-- -->                            |
| ------------- | ----------------------------------- |
| Property      | `in_ports`                          |
| Type          | _[V1SwitcherIns](#v1-switcher-ins)_ |
| Default Value | ``                                  |
| Description   | Lorem Ipsum                         |

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `out_ports`                           |
| Type          | _[V1SwitcherOuts](#v1-switcher-outs)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

### v1SwitcherIns {#v1-switcher-ins}

Inputs for the Switcher component.

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `on_true`                 |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `on_false`                |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

| <!-- -->      | <!-- -->                  |
| ------------- | ------------------------- |
| Property      | `switch`                  |
| Type          | _[V1InPort](#v1-in-port)_ |
| Default Value | ``                        |
| Description   | Lorem Ipsum               |

### v1SwitcherOuts {#v1-switcher-outs}

Outputs for the Switcher component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

### v1Variable {#v1-variable}

Component that emits a variable value as an output signal, can be defined in dynamic configuration.

| <!-- -->      | <!-- -->                              |
| ------------- | ------------------------------------- |
| Property      | `out_ports`                           |
| Type          | _[V1VariableOuts](#v1-variable-outs)_ |
| Default Value | ``                                    |
| Description   | Lorem Ipsum                           |

| <!-- -->      | <!-- -->             |
| ------------- | -------------------- |
| Property      | `dynamic_config_key` |
| Type          | _string_             |
| Default Value | ``                   |
| Description   | Lorem Ipsum          |

| <!-- -->      | <!-- -->                                                 |
| ------------- | -------------------------------------------------------- |
| Property      | `default_config`                                         |
| Type          | _[V1VariableDynamicConfig](#v1-variable-dynamic-config)_ |
| Default Value | ``                                                       |
| Description   | Lorem Ipsum                                              |

### v1VariableDynamicConfig {#v1-variable-dynamic-config}

| <!-- -->      | <!-- -->                                  |
| ------------- | ----------------------------------------- |
| Property      | `constant_signal`                         |
| Type          | _[V1ConstantSignal](#v1-constant-signal)_ |
| Default Value | ``                                        |
| Description   | Lorem Ipsum                               |

### v1VariableOuts {#v1-variable-outs}

Outputs for the Variable component.

| <!-- -->      | <!-- -->                    |
| ------------- | --------------------------- |
| Property      | `output`                    |
| Type          | _[V1OutPort](#v1-out-port)_ |
| Default Value | ``                          |
| Description   | Lorem Ipsum                 |

<!---
Generated File Ends
-->
