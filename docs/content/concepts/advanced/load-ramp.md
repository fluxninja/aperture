---
title: Load Ramp
sidebar_position: 6
---

:::info See also

[_Load Ramp_ reference][load-ramp]

:::

The _Load Ramp_ is a high-level component in the policy [circuit](circuit.md).
It is specifically designed to facilitate controlled traffic ramp-up scenarios.
This component allows for gradual increases in traffic by specifying rollout
steps, each defined by a time duration and target percentage. The target
percentage transitions linearly from one step to the next.

To achieve this functionality, the _Load Ramp_ utilizes input ports named
`forward`, `backward`, and `reset` to control its behavior. These input signals
drive the internal logic of the component. Additionally, the _Load Ramp_
internally leverages the [_Sampler_](#sampler) component, providing it with an
`accept_percentage` signal and setting its `ramp_mode` flag.

## Sampler

:::info See also

[_Sampler_ reference][sampler]

:::

The _Sampler_ component in Aperture is designed to manage the load at a
[_Control Point_][control-point] by allowing only a specified percentage of
flows at random or by sticky sessions. _Sampler_ is useful for ramping the load
in a controlled way.

The following example creates a _Sampler_ at the `ingress` control point of
`user-service.default.svc.cluster.local` service. It accepts 50% of the flows
with stickiness based on the value of `http.request.header.user_id` flow label:

```yaml
circuit:
components:
  - flow_control:
      sampler:
        in_ports:
          accept_percentage:
            constant_signal:
              value: 50
        parameters:
          label_key: http.request.header.user_id
          selectors:
            - control_point: ingress
              service: user-service.default.svc.cluster.local
          ramp_mode: true
```

### Stateless Filter {#stateless-filter}

When the label key is not specified, the _Sampler_ acts as a stateless filter.
In this mode, it selects a percentage of flows randomly. This is useful for
controlling the overall load at a control point without focusing on specific
sessions or users.

### Sticky Filter {#sticky-filter}

If a label key is specified, the _Sampler_ acts as a sticky filter. In this
mode, it ensures that a series of flows with the same value of the label key get
the same decision, provided the accept percentage remains the same or higher.
This is useful for maintaining consistent behavior for specific sessions or
users.

### Dynamic Configuration {#dynamic-configuration}

The _Sampler_ supports dynamic configuration by specifying certain label values
to be accepted by the flow filter, regardless of the accept percentage. This
enables making exceptions for specific sessions or users, granting them
unhindered access to a control point.

### Selector {#selectors}

The Flow Selector is responsible for determining the service and flows at which
the _Sampler_ is applied. By configuring the Flow Selector, one can control
which parts of a service are affected by the _Sampler_.

### Accept Percentage Signal {#accept-percentage-signal}

The _Sampler_ component takes an accept-percentage input signal, based on which
it decides the percentage of flows to accept, allowing for a fine-grained
control.

### Ramp Mode {#ramp-mode}

Flows with `ramp_mode` flag set require at least one ramp component to accept
them, resulting in flow rejection if Aperture Agent is disconnected or policy is
not loaded. When _Sampler_'s `ramp_mode` flag is set, the _Sampler_ becomes such
a ramp component and can accept ramp mode flows.

[sampler]: /reference/configuration/spec.md#sampler
[control-point]: /concepts/control-point.md
[load-ramp]: /reference/configuration/spec.md#load-ramp
