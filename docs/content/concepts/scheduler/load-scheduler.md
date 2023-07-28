---
title: Load Scheduler
keywords:
  - scheduler
  - service protection
  - queuing
sidebar_position: 1
---

:::info See Also

Load Scheduler [Reference](../../reference/configuration/spec.md#load-scheduler)

:::

The _Load Scheduler_ is used to throttle request rates dynamically during high
load, therefore protecting services from overloads and cascading failures. It
uses a local token bucket for estimating the allowed token rate. The fill rate
of the token bucket gets adjusted by the controller based on the specified
policy. Since this component builds upon the [_Scheduler_](./scheduler.md), it
allows defining workloads along with their priority and tokens. The scheduler
employs weighted fair queuing of requests to achieve graceful degradation of
applications.

This diagram illustrates the working of a load scheduler.

![Scheduler](./assets/img/load-scheduler-light.svg#gh-light-mode-only)
![Scheduler](./assets/img/load-scheduler-dark.svg#gh-dark-mode-only)

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

## Adaptive Load Scheduler {#adaptive-load-scheduler}

:::info See Also

Adaptive Load Scheduler
[Reference](../../reference/configuration/spec.md#adaptive-load-scheduler)

:::

_Adaptive Load Scheduler_ is a high-level [circuit](../advanced/circuit.md)
component that uses the _Load Scheduler_ internally. In addition, it employs a
[_Gradient Controller_](../../reference/configuration/spec.md#gradient-controller)
and an [Integrator](../../reference/configuration/spec.md#integrator) for
computing the load multiplier. From the Agents' perspective, the _Load
Scheduler_ and _Adaptive Load Scheduler_ are identical.

Offering a more high-level interface, this component has `signal`, `setpoint`,
and `overload_confirmation` ports. The core function of the _Adaptive Load
Scheduler_ is its ability to modify the accepted token rate based on the
deviation of the input `signal` from the `setpoint`. This functionality allows
for the translation of health signals into adjustments in token rate, thereby
providing an active defense mechanism for the service.

Additionally, by defining workloads with varying priorities and weights
(tokens), the scheduler can prioritize certain requests over others,
facilitating graceful service degradation during high-traffic periods.
