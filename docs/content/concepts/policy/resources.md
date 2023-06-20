---
title: Resources
sidebar_position: 2
---

:::info

See also [Resources reference](/reference/configuration/spec.md#resources)

:::

Resources configuration in the policy specification describes resources needed
to set up a [circuit][circuit]. It's possible but not recommended to share
resources across policies since resources are always defined in the global
scope. Resources might be referenced inside the circuit, the exact reference
mechanism depends on the type of resource.

Examples:

1. [_Flux Meters_][flux-meter]: The metrics generated by a _Flux Meter_ can be
   referenced inside [PromQL components][promql-reference]
2. [_Classifiers_][classifier]: The labels generated by a _Classifier_ can be
   referred inside a [`FlowSelector`][selector-reference].

[circuit]: circuit.md
[flux-meter]: /concepts/flow-control/resources/flux-meter.md
[classifier]: /concepts/flow-control/resources/classifier.md
[promql-reference]: /reference/configuration/spec.md#prom-q-l
[selector-reference]: /reference/configuration/spec.md#flow-selectors