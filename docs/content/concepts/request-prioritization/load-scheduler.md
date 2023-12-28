---
title: Load Scheduler
keywords:
  - scheduler
  - service protection
  - queuing
sidebar_position: 3
---

:::info See Also

Load Scheduler [Reference](../../reference/configuration/spec.md#load-scheduler)

:::

The _Load Scheduler_ is used to throttle request rates dynamically during high
load, therefore protecting services from overloads and cascading failures. It
uses a local token bucket for estimating the allowed token rate. The fill rate
of the token bucket gets adjusted by the controller based on the specified
policy. Since this component builds upon the [_Scheduler_](../scheduler.md), it
allows defining workloads along with their priority and tokens. The scheduler
employs weighted fair queuing of requests to achieve graceful degradation of
applications.

The _Load Scheduler_'s throttling behavior is controlled by the signal at its
`load_multiplier` input port. As the policy circuit adjusts the signal at the
load multiplier port, it gets translated to the token refill rate at the Agents.
At each Agent, the adjusted token rate is determined by multiplying the past
token rate with the load multiplier. The past 30 seconds of data is used for
finding the past token rate.

$$
adjusted\_token\_rate = past\_token\_rate * load\_multiplier
$$

If the incoming request rate surpasses the adjusted rate, the scheduler starts
queuing requests. The queued requests get admitted as tokens become available in
an order determined by the scheduler based on the weighted fair queuing
algorithm. Any request that fails to be scheduled within its designated timeout
is rejected.

Additionally, by defining workloads with varying priorities and weights
(tokens), the load scheduler can prioritize certain requests over others,
facilitating graceful service degradation during high-traffic periods.

## Load Schedulers {#load-schedulers}

### Additive Increase Multiplicative Decrease Load Scheduler {#aimd-load-scheduler}

:::info Specification

AIMD Load Scheduler
[Reference](../../reference/configuration/spec.md#a-i-m-d-load-scheduler)

:::

_AIMD Load Scheduler_ is a high-level [circuit](../advanced/circuit.md)
component that uses the _Load Scheduler_ internally. In addition, it employs a
[_Gradient Controller_](../../reference/configuration/spec.md#gradient-controller)
and an [Integrator](../../reference/configuration/spec.md#integrator) for
computing the load multiplier.

Offering a more high-level interface, this component has `signal`, `setpoint`,
and `overload_confirmation` ports. The core function of the _AIMD Load
Scheduler_ is its ability to modify the accepted token rate based on the
deviation of the input `signal` from the `setpoint`. This scheduler reduces
token rate proportionally (or any arbitrary power) based on deviation of the
`signal` from `setpoint`. During recovery, it increases the token rate linearly
until the system is not overloaded. It allows for the translation of health
signals into adjustments in token rate, thereby providing an active defense
mechanism for the service.

### Additive Increase Additive Decrease Load Scheduler {#aiad-load-scheduler}

:::info Specification

AIAD Load Scheduler
[Reference](../../reference/configuration/spec.md#a-i-a-d-load-scheduler)

:::

_AIAD Load Scheduler_ is a high-level [circuit](../advanced/circuit.md)
component that uses the _Load Scheduler_ internally.

This component has `signal`, `setpoint`, and `overload_confirmation` ports. It
uses `overload_condition` to compare `signal` and `setpoint` to determine if the
service is overloaded. _AIAD Load Scheduler_ reduces the token rate linearly
over time while in overload state. During recovery, it increases the token rate
linearly until the system is not overloaded.

### Range-Driven Load Scheduler {#range-driven-load-scheduler}

:::info Specification

Range-Driven Load Scheduler
[Reference](../../reference/configuration/spec.md#range-driven-load-scheduler)

:::

_Range-Driven Load Scheduler_ is a high-level [circuit](../advanced/circuit.md)
component that uses the _Load Scheduler_ internally.

This component has `signal` and `overload_confirmation` ports. It uses the
[polynomial range function](../../reference/configuration/spec.md#polynomial-range-function)
to throttle the token rate based on the range of the `signal`, attempting to
keep it between `low_throttle_threshold` and `high_throttle_threshold`.
