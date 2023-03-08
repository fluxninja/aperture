---
title: Circuit
sidebar_position: 1
---

:::info

See also [_Circuit_ reference][circuit-reference]

:::

_Circuit_ describes a [control system][control-system] as an execution graph.
_Circuit_ is defined as a dataflow graph of inter-connected components.
_Signals_ flow between components via ports. As signals traverse the circuit,
they get processed, stored within components or get acted upon (e.g. load-shed,
rate-limit, auto-scale etc.). _Circuit_ is evaluated periodically in order to
respond to changes in signal readings.

## Component

Building blocks of a circuit are components. Each component has input ports
(`in_ports`) and output ports (`out_ports`). The exact ports available are
determined by the [type of component][components]. Each port can be associated
with a [signal][signal]. Components get chained to one another based on name of
the signal.

## Signal

_Signal_ represents a `float64` value that updates with every [tick][tick] of
circuit execution. Every signal must have a name to uniquely identify it within
a circuit.

Output port on a component may emit a signal. No other port (on any component)
in the circuit can emit a signal with the same name.

In order to receive a named signal at a component it must be defined exactly
once as an output at some component in the circuit. Once defined, a signal may
be received at multiple components.

## Circuit Runtime

_Circuit_ evaluates at a constant _tick_ frequency. Each round of evaluation is
called a tick. The `evaluation_interval` parameter in [policy
spec][policy-reference] configures how often the circuit evaluates (ticks).

On every tick, each component in the circuit gets executed exactly once.
components get executed as they become ready. A component is ready if all of its
input signals are available.

During execution, the input signals are processed and output signals are emitted
by the component. Any [looping signals][looping-signals] are saved and consumed
by circuit in the next tick.

_Circuit_ runtime provides very predictable execution semantics. Any timed
operations like [_PromQL_][promql-reference] queries are synchronized to execute
on multiples of ticks. All _PromQL_ queries in a circuit are centrally
synchronized to ensure that all the queries that fire in the same tick return
results together in a future tick.

### Looping Signals

Loops are allowed in the circuit execution graph. In fact they enable expression
of powerful paradigms such as integration using basic arithmetic components.

In reality, the execution is still performed on a
[directed acyclic graph](https://en.wikipedia.org/wiki/Directed_acyclic_graph).
Before execution, loops are detected in the circuit. Each loop is un-linked at
the component with the smallest index (in list of components). The un-linked
component ports consume looping signals. A looping signal has the value of the
un-linked signal from the previous tick.

## Example Components {#components}

The exhaustive list of the built-in components can be found in the
[policy reference](reference/policies/spec.md#component).

Examples of built-in components include:

- **Query**: These components emit signals into the circuit from outside.
  - [PromQL][promql-reference]: Converts results from a PromQL query into a
    signal.
- **Signal Processors**: These components transform input signal(s) into output
  signal(s).
  - **Arithmetic**: These components perform basic arithmetic operations on
    signal(s).
    - [Arithmetic Combinator](/reference/policies/spec.md#arithmetic-combinator):
      This component takes two input signals and performs a basic arithmetic
      operation to generate an output signal.
    - [Max](/reference/policies/spec.md#max) and
      [Min](/reference/policies/spec.md#min): These components take multiple
      input or output signals and emit maximum or minimum of those signals.
      signal.
  - **Transformers**: These components statefully transform an input signal in
    an output signal.
    - [EMA](/reference/policies/spec.md#e-m-a): Exponential moving average.
  - [Decider and Switcher](/reference/policies/spec.md#decider): These
    components work in tandem to make the circuit adapt based on conditions.
- **Controllers**: Controllers are an essential part of a closed loop control
  system. A controller take as input a signal, a setpoint and emits the
  suggested value of Control Variable as output. The aim of the controller is to
  make the signal achieve the setpoint.
  - [Gradient Controller](/reference/policies/spec.md#gradient-controller): This
    controller acts on the ratio of setpoint and signal.
- **Actuators**: Actuators are components which act on signals to make real
  changes like shed traffic, change rate limits etc.
  - [Concurrency Limiter](/reference/policies/spec.md#concurrency-limiter):
    Takes load multiplier as a signal which determines the proportion of Flow
    concurrency to accept.
  - [Rate Limiter](/reference/policies/spec.md#rate-limiter): Take rate limit as
    a signal which determines the rate of flows handled by that rate limiter.

[control-system]: https://en.wikipedia.org/wiki/Control_system
[tick]: #runtime
[signal]: #signal
[looping-signals]: #looping-signals
[components]: #components
[policy-reference]: /reference/policies/spec.md#policy
[circuit-reference]: /reference/policies/spec.md#circuit
[promql-reference]: /reference/policies/spec.md#prom-q-l
[scheduler-reference]: /reference/policies/spec.md#scheduler
