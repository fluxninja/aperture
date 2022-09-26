---
title: Circuit
sidebar_position: 2
---

:::info

See also [Circuit reference][circuit-reference]

:::

Circuit describes a [Control System][control-system] as an execution graph.
Circuit is defined as a dataflow graph of inter-connected Components. Signals
flow between Components via Ports. As Signals traverse the Circuit, they get
processed, stored within Components or get acted upon (e.g. load-shed,
rate-limit, auto-scale etc.). Circuit is evaluated periodically in order to
respond to changes in Signal readings.

## Component

Building blocks of a Circuit are Components. Each Component has Input Ports
(`in_ports`) and Output Ports (`out_ports`). The exact Ports available are
determined by the [type of Component][components]. Each Port can be associated
with a [Signal][signal]. Components get chained to one another based on name of
the Signal.

## Signal

Signal represents a `float64` value that updates with every [Tick][tick] of
Circuit execution. Every Signal must have a name to uniquely identify it within
a Circuit.

Output Port on a Component may emit a Signal. No other Port (on any Component)
in the Circuit can emit a Signal with the same name.

In order to receive a named Signal at a Component it must be defined exactly
once as an Output at some Component in the Circuit. Once defined, a Signal may
be received at multiple Components.

## Circuit Runtime

The Circuit evaluates at a constant _Tick_ frequency. Each round of evaluation
is called a Tick. The `evaluation_interval` parameter in [Policy
spec][policy-reference] configures how often the Circuit evaluates (Ticks).

On every Tick, each Component in the Circuit gets executed exactly once.
Components get executed as they become ready. A Component is ready if all of its
Input Signals are available.

During execution, the Input Signals are processed and Output Signals are emitted
by the Component. Any [Looping Signals][looping-signals] are saved and consumed
by Circuit in the next Tick.

Circuit runtime provides very predictable execution semantics. Any timed
operations like PromQL queries are synchronized to execute on multiples of
Ticks. All PromQL queries in a circuit are centrally synchronized to ensure that
all the queries that fire in the same Tick return results together in a future
Tick.

### Looping Signals

Loops are allowed in the Circuit execution graph. In fact they enable expression
of powerful paradigms such as integration using basic Arithmetic Components.

In reality, the execution is still performed on a Directed Acyclic Graph. Before
execution, loops are detected in the Circuit. Each loop is un-linked at the
Component with the smallest index (in list of Components). The un-linked
Component Ports consume Looping Signals. A Looping Signal has the value of the
un-linked Signal from the previous Tick.

## Exploring Components {#components}

Components are divided into following functional categories to aid exploration:

- **Sources**: These Components emit Signals into the Circuit from outside.
  - [PromQL][promql-reference]: Converts results from a PromQL query into a
    Signal.
  - [ConcurrencyLimiter.Scheduler][scheduler-reference]: Emits Signals
    representing incoming and accepted Flow concurrencies observed by the
    specified [Flow Control Scheduler][flow-control-scheduler] at Aperture
    Agents.
- **Signal Processors**: These Components transform input Signal(s) into output
  Signal(s).
  - **Arithmetic**: These Components perform basic Arithmetic operations on
    Signal(s).
    - [Arithmetic Combinator](/references/configuration/policies.md#v1-arithmetic-combinator):
      This Component takes two input Signals and performs a basic arithmetic
      operation to generate an output Signal.
    - [Max](/references/configuration/policies.md#v1-max) and
      [Min](/references/configuration/policies.md#v1-min): These Components take
      multiple input or output Signals and emit maximum or minimum of those
      Signals.
    - [Sqrt](/references/configuration/policies.md#v1-sqrt): This Component
      square roots a Signal.
    - [Constant](/references/configuration/policies.md#v1-constant): A constant
      Signal value.
  - **Transformers**: These Components statefully transform an input Signal in
    an output Signal.
    - [EMA](/references/configuration/policies.md#v1-e-m-a): Exponential Moving
      Average.
    - [Extrapolator](/references/configuration/policies.md#v1-extrapolator):
      Extrapolate a Signal based on past values.
  - [Decider and Switcher](/references/configuration/policies.md#v1-decider):
    These Components work in tandem to make the Circuit adapt based on
    conditions.
- **Controllers**: Controllers are an essential part of a closed loop control
  system. A Controller take as input a signal, a setpoint and emits the
  suggested value of Control Variable as output. The aim of the Controller is to
  make the Signal achieve the Setpoint.
  - [Gradient Controller](/references/configuration/policies.md#v1-gradient-controller):
    This Controller acts on the ratio of Setpoint and Signal.
- **Actuators**: Actuators are Components which act on Signals to make real
  changes like shed traffic, change rate limits etc.
  - [Concurrency Limiter](/references/configuration/policies.md#v1-concurrency-limiter):
    Takes load shed factor as a Signal which determines the proportional of Flow
    concurrency to shed.
  - [Rate Limiter](/references/configuration/policies.md#v1-rate-limiter): Take
    rate limit as a Signal which determines the rate of flows handled by that
    Rate Limiter.

[control-system]: https://en.wikipedia.org/wiki/Control_system
[tick]: #runtime
[signal]: #signal
[looping-signals]: #looping-signals
[components]: #components
[policy-reference]: /references/configuration/policies.md#v1-policy
[circuit-reference]: /references/configuration/policies.md#v1-circuit
[promql-reference]: /references/configuration/policies.md#v1-prom-q-l
[scheduler-reference]: /references/configuration/policies.md#v1-scheduler
[flow-control-scheduler]: ../flow-control/concurrency-limiter.md#scheduler
