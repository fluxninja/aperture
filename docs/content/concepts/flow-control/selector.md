---
title: Flow Selector
sidebar_label: Selector
sidebar_position: 2
keywords:
  - flows
  - services
  - discovery
  - labels
---

:::info

See also [Selector reference](/references/configuration/policies.md#v1-selector)

:::

Flow Selectors are used by all flow observability and control components
instantiated on Aperture Agents ([Classifiers][classifier], [Flux
Meters][flux-meter] and Limiters). Flow Selectors define scoping rules – how
these components should select [Flows][flow] for their operations.

A Selector consists of:

- [agent group][agent-group] name,
- [service][service] name,
- [control point][control-point], and
- optional [flow label matcher](#label-matcher).

### Service

_Agent group_ name together with _service_ name determine the [service][service]
to select flows from.

:::tip Default Agent Group

The default Agent Group is called `default`. If you're using this group, you can
skip the _agent group_ field.

:::

:::tip Catch-all service

If the agent group is already logically a single service or you simply want to
select all services within the agent group, you can skip the service name.

:::

### Control Point

Flow Selector selects flows from only one [Control Point][control-point] within
a service.

### Label Matcher

Label matcher allows to optionally narrow down the selected flow based on
conditions on [Flow Labels][label].

There are multiple ways to define a label matcher. The simplest way is to
provide a map of labels for exact-match:

```yaml
label_matcher:
  match_labels:
    http.method: GET
```

You can also provide a matching-expression-tree, which allows for arbitrary
conditions, including regex matching. Refer to [LabelMatcher][label-matcher] for further details.

### Example

```yaml
service: checkout.myns.svc.cluster.local
agent_group: default
control_point:
  traffic: ingress
label_matcher:
  match_labels:
    user_tier: gold
```

[flow]: /concepts/flow-control/flow-control.md#flow
[label]: /concepts/flow-control/flow-label.md
[control-point]: /concepts/flow-control/flow-control.md#control-point
[service]: /concepts/service.md
[agent-group]: /concepts/service.md#agent-group
[flux-meter]: /concepts/flow-control/flux-meter.md
[classifier]: /concepts/flow-control/flow-classifier.md
[label-matcher]: /references/configuration/policies.md#v1-label-matcher
