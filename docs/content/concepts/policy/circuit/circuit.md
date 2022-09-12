---
title: Circuit
sidebar_position: 1
---

<!-- How and what circuits can express? How the controller interacts with the agents based on signals from a circuit. -->
<!-- Section: How to express a circuit -->
<!-- Sub-section Component -->
<!-- Sub-section Signal -->
<!-- Section: Runtime: how the circuit executes circuit per tick, handles loops. Cover how the circuit executes circuit per tick, loops -->

:::info

See also [Circuit reference][circuit-reference]

:::

Circuit describes a [Control System][control-system] as an execution graph.

# Component

Building blocks of a Circuit are Components. Each Component has Input Ports
(`in_ports`) and Output Ports (`out_ports`). The exact Ports available are
determined by the [type of Component][component-types]. Each Port can be
associated with a [Signal][signal]. Components get chained to one another based
on name of the Signal.

# Signal

Signal represents a `float64` value that updates with every [Tick][tick] of
Circuit execution. Every Signal must have a name to uniquely identify it within
a Circuit.

Output Port on a Component may emit a Signal. No other Port (on any Component)
in the Circuit can emit a Signal with the same name.

In order to receive a named Signal at a Component it must be defined exactly
once as an Output at some Component in the Circuit. Once defined, a Signal may
be received at multiple Components.

# Circuit Runtime

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

## Looping Signals

Loops are allowed in the Circuit execution graph. In fact they enable expression
of powerful paradigms such as integration using basic Arithmetic components.

In reality, the execution is still performed on a Directed Acyclic Graph. Before
execution, loops are detected in the Circuit. Each Loop is un-linked at the
Component with the smallest index (in list of Components). The un-linked
Component Ports consume Looping Signals. A Looping Signal has the value of the
un-linked Signal from the previous Tick.

[control-system]: https://en.wikipedia.org/wiki/Control_system
[component-types]: /reference/configuration/policies.md#v1-component
[tick]: #runtime
[signal]: #signal
[policy-reference]: /reference/configuration/policies.md#v1-policy
[circuit-reference]: /reference/configuration/policies.md#v1-circuit
