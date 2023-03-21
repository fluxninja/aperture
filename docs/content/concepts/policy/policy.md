---
title: Policy
sidebar_position: 1
---

:::info

See also [Policy reference](/reference/policies/spec.md#policy)

:::

Aperture's policy language enables expression of closed-loop [control
systems][control-system] in a declarative manner. Aperture comes with
pre-packaged [blueprints][blueprints] that can be used both as a guide for
creating new policies, or used as-is.

Policies provide a framework for defining and managing reliability criteria and
conditions as code. They allow service operators to define and enforce
reliability policies programmatically, running in a continuous control loop. In
an application reliability context, policies codify the capability of the
application to modify its operational state to achieve the best possible mode of
operation despite overload and failures.

Policy specification consists of two parts:

## Circuit

A [circuit][circuit] in Aperture's policy language represents the signal
processing circuit of the control system as an execution graph. It captures the
decision-making process and describes the actions to be taken based on the
current state of the system. A circuit is made up of nodes, which represent the
various components of the control system, including signal processing
components, and edges, which represent the flow of signals between the nodes.

Observability-driven control is an important aspect of Aperture's policy
language. By monitoring signals such as request latency, error rate, and
saturation, Aperture's circuits can detect deviations from service-level
objectives (SLOs) and trigger appropriate actions to restore system stability
and reliability. The circuit is the heart of the policy specification and is
responsible for the logic of the control system.

## Resources

A list of [Resources][resources] which need to be set up in order to support the
circuit.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />

[circuit]: /concepts/policy/circuit.md
[resources]: /concepts/policy/resources.md
[blueprints]: /reference/policies/bundled-blueprints/bundled-blueprints.md
[control-system]: https://en.wikipedia.org/wiki/Control_system
