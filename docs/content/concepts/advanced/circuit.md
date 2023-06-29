---
title: Circuit
sidebar_position: 2
---

:::info See also

[_Circuit_ reference][circuit-reference]

:::

The _Circuit_ describes a [control system][control-system] as an execution
graph. The _Circuit_ is defined as a graph of interconnected signal processing
components. _Signals_ flow between components through ports. As signals traverse
the circuit, they get processed, stored within components or get acted upon (for
example: load-throttling, rate-limit, auto-scale and so on). The _Circuit_ is
evaluated periodically to respond to changes in signal readings.

## Component

The building blocks of a circuit are components. Each component has input ports
(`in_ports`) and output ports (`out_ports`). The exact ports available are
determined by the [type of component][components]. Each port can be associated
with a [signal][signal]. Components get chained to one another based on the name
of the signal.

## Signal

The _Signal_ represents a `float64` value that updates with every [tick][tick]
of circuit execution. Every signal must have a name to uniquely identify it
within a circuit.

The output port on a component can emit a signal. No other port (on any
component) in the circuit can emit a signal with the same name.

To receive a named signal at a component, it must be defined exactly once as an
output at some component in the circuit. Once defined, a signal can be received
at other components that have an input port with the same name.

## Circuit Runtime {#runtime}

The _Circuit_ evaluates at a constant _tick_ frequency. Each round of evaluation
is called a tick. The `evaluation_interval` parameter in the [policy
spec][policy-reference] configures how often the circuit evaluates (ticks).

On every tick, each component in the circuit gets executed exactly once.
Components get executed as they become ready. A component is ready if all its
input signals are available.

During execution, the input signals are processed, and output signals are
emitted by the component. Any [looping signals][looping-signals] are saved and
consumed by the circuit in the next tick.

_Circuit_ runtime provides highly predictable execution semantics. Any timed
operations like [_PromQL_][promql-reference] queries are synchronized to execute
on multiples of ticks. For _PromQL_ queries, this ensures that all the queries
that fire in the same tick return results together in a future tick.

### Looping Signals

Circuit execution graph can have loops. Looping signals enable expression of
powerful paradigms such as integration using basic arithmetic components.

Despite loops, the execution is still performed on a
[directed acyclic graph](https://en.wikipedia.org/wiki/Directed_acyclic_graph).
Before execution, loops are detected in the circuit. Each loop is unlinked at
the component with the smallest index (in list of components). The unlinked
component ports consume looping signals. A looping signal has the value of the
unlinked signal from the previous tick.

## Example Components {#components}

The exhaustive list of the built-in components can be found in the
[policy reference](/reference/configuration/spec.md#component).

Examples of built-in components include:

- **Query**: These components emit signals into the circuit from outside.
  - [PromQL][promql-reference]: Converts results from a PromQL query into a
    signal.
- **Signal Processors**: These components transform input signal(s) into output
  signal(s).
  - **Arithmetic**: These components perform basic arithmetic operations on
    signal(s).
    - [Arithmetic Combinator](/reference/configuration/spec.md#arithmetic-combinator):
      This component takes two input signals and performs a basic arithmetic
      operation to generate an output signal.
    - [Max](/reference/configuration/spec.md#max) and
      [Min](/reference/configuration/spec.md#min): These components take
      multiple input or output signals and emit the maximum or minimum of those
      signals.
  - **Transformers**: These components transform an input signal into an output
    signal based on past state.
    - [EMA](/reference/configuration/spec.md#e-m-a): Exponential moving average.
  - [Decider and Switcher](/reference/configuration/spec.md#decider): These
    components work in tandem to make the circuit adapt based on conditions.
- **Controllers**: Controllers are an essential part of a closed loop control
  system. A controller takes as input a signal, a setpoint and emits the
  suggested value of the control variable as output. The aim of the controller
  is to make the signal achieve the setpoint.
  - [Gradient Controller](/reference/configuration/spec.md#gradient-controller):
    This controller acts on the ratio of setpoint and signal.
- **Actuators**: Actuators are components which act on signals and interface
  with external systems to perform actions such as throttling and queuing
  traffic, changing rate limits or auto-scaling.
  - [Load Scheduler](/reference/configuration/spec.md#load-scheduler): Takes
    load multiplier as a signal which determines the proportion of tokens to
    accept relative to past window.
  - [Rate Limiter](/reference/configuration/spec.md#rate-limiter): Take fill
    rate and bucket capacity as signals which determines the API rate limit for
    each user or label value.

[control-system]: https://en.wikipedia.org/wiki/Control_system
[tick]: #runtime
[signal]: #signal
[looping-signals]: #looping-signals
[components]: #components
[policy-reference]: /reference/configuration/spec.md#policy
[circuit-reference]: /reference/configuration/spec.md#circuit
[promql-reference]: /reference/configuration/spec.md#prom-q-l
