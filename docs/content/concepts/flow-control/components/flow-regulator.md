---
title: Flow Regulator
sidebar_position: 6
---

:::info

See also [_Flow Regulator_ reference][reference]

:::

The Flow Regulator component in Aperture is designed to manage the flow of
requests to a service by allowing only a specified percentage of requests or
sticky sessions. This is useful for controlling the load on your service or for
performing controlled tests. It achieves this by either acting as a stateless
filter or a sticky filter, depending on the provided configuration.

## Stateless Filter {#stateless-filter}

When the label key is not specified, the Flow Regulator acts as a stateless
filter. In this mode, it selects a percentage of flows randomly. This is useful
when you want to control the overall load on your service without focusing on
specific sessions or users.

## Sticky Filter {#sticky-filter}

If a label key is specified, the Flow Regulator acts as a sticky filter. In this
mode, it ensures that a series of flows with the same value of the label key get
the same decision, provided the percentage of accepted requests remains the same
or higher. This is useful when you want to maintain consistent behavior for
specific sessions or users.

## Dynamic Configuration {#dynamic-configuration}

The Flow Regulator supports dynamic configuration by specifying certain label
values to be accepted by the flow filter, regardless of the accept percentage.
This allows you to make exceptions for specific sessions or users, granting them
unhindered access to your service.

## Flow Selector {#flow-selector}

The Flow Selector is responsible for determining the service and flows at which
the Flow Regulator is applied. By configuring the Flow Selector, you can control
which parts of your service are affected by the Flow Regulator.

## Accept Percentage Signal {#accept-percentage-signal}

The Flow Regulator component takes an accept-percentage input signal, based on
which it decides the percentage of requests to accept. This gives you
fine-grained control over the flow of requests to your service, enabling you to
achieve the desired balance between load and performance.

[reference]: /reference/policies/spec.md#flow-regulator
