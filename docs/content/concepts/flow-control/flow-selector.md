---
title: Flow Selector
sidebar_label: Flow Selector
sidebar_position: 2
keywords:
  - flows
  - services
  - discovery
  - labels
---

:::info

See also
[FlowSelector reference](/references/configuration/policy.md#v1-flow-selector)

:::

Flow Selectors are used by flow control and observability components
instantiated by Aperture Agents like [Classifiers][classifier], [Flux
Meters][flux-meter] and [Concurrency Limiters][concurrency-limiter]. Flow
Selectors define scoping rules – how these components should select
[Flows][flow] for their operations.

A Flow Selector consists of:

- Service Selector, containing

  - [agent group][agent-group] name,
  - [service][service] name,

- FlowMatcher, containing
  - [control point][control-point], and
  - optional [flow label matcher](#label-matcher).

### Service Selector

:::info

See also
[ServiceSelector reference](/references/configuration/policy.md#v1-service-selector)

:::

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

### Flow Matcher

:::info

See also
[FlowMatcher reference](/references/configuration/policy.md#v1-flow-matcher)

:::

#### Control Point

Flow Selector selects flows from only one [Control Point][control-point] within
a service.

#### Label Matcher

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
conditions, including regex matching. Refer to [LabelMatcher][label-matcher] for
further details.

### Example

```yaml
service_selector:
  service: checkout.myns.svc.cluster.local
  agent_group: default
flow_selector:
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
[concurrency-limiter]: /concepts/flow-control/concurrency-limiter.md
[classifier]: /concepts/flow-control/flow-classifier.md
[label-matcher]: /references/configuration/policy.md#v1-label-matcher
