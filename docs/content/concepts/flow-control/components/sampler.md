---
title: Sampler
sidebar_position: 6
---

:::info

See also [_Sampler_ reference][flow-sampler]

:::

The _Sampler_ component in Aperture is designed to manage the load at a
[_Control Point_][control-point] by allowing only a specified percentage of
flows at random or by sticky sessions. This is useful for controlling the load
at a feature inside your service for performing controlled tests. It achieves
this by either acting as a stateless filter or a sticky filter, depending on the
provided configuration.

## Stateless Filter {#stateless-filter}

When the label key is not specified, the _Sampler_ acts as a stateless filter.
In this mode, it selects a percentage of flows randomly. This is useful when you
want to control the overall load on a control point without focusing on specific
sessions or users.

## Sticky Filter {#sticky-filter}

If a label key is specified, the _Sampler_ acts as a sticky filter. In this
mode, it ensures that a series of flows with the same value of the label key get
the same decision, provided the accept percentage remains the same or higher.
This is useful when you want to maintain consistent behavior for specific
sessions or users.

## Dynamic Configuration {#dynamic-configuration}

The _Sampler_ supports dynamic configuration by specifying certain label values
to be accepted by the flow filter, regardless of the accept percentage. This
allows you to make exceptions for specific sessions or users, granting them
unhindered access to a control point.

## Flow Selector {#flow-selectors}

The Flow Selector is responsible for determining the service and flows at which
the _Sampler_ is applied. By configuring the Flow Selector, you can control
which parts of your service are affected by the _Sampler_.

## Accept Percentage Signal {#accept-percentage-signal}

The _Sampler_ component takes an accept-percentage input signal, based on which
it decides the percentage of flows to accept. This gives you fine-grained
control over the flow of requests to a control point, enabling you to achieve
the desired balance between load and performance.

[flow-sampler]: /reference/configuration/spec.md#flow-sampler
[control-point]: /concepts/flow-control/selector.md/#control-point
