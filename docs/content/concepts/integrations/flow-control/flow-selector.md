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

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

:::info

See also [_Flow Selector_ reference](/reference/policies/spec.md#flow-selector)

:::

_Flow Selectors_ are used by flow control and observability components
instantiated by Aperture Agents like [_Classifiers_][classifier], [_Flux
Meters_][flux-meter] and [_Concurrency Limiters_][cl]. _Flow Selectors_ define
scoping rules – how these components should select flows for their operations.

A _Flow Selector_ consists of:

- _Service Selector_, containing

  - [_Service_](#service) name,
  - [_Agent Group_](#agent-group) name,

- _Flow Matcher_, containing
  - [_Control Point_](#control-point), and
  - [_Flow Matcher_](#label-matcher).

## Service Selector {#service-selector}

:::info

See also
[_Service Selector_ reference](/reference/policies/spec.md#service-selector)

:::

A service in Aperture is similar to services tracked in Kubernetes or Consul.
Services in Aperture are usually referred by their fully qualified domain names
(FQDN).

A service is a collection of entities delivering a common functionality, such as
checkout, billing etc. Aperture maintains a mapping of entity IP addresses to
service names. For each flow control decision request sent by an entity,
Aperture looks up the service name and then decides which flow control
components to execute.

:::note

An entity (K8s Pod, VM, etc) may belong to multiple services.

:::

:::tip Catch-all service

If the agent group is already logically a single service or you simply want to
select all services within the agent group, you can set the service name as
`all`.

:::

:::info Service Discovery

Aperture Agents perform automated discovery of services and entities in
environments such as Kubernetes and watch for any changes. Service and entity
entries can also be created manually via configuration.

:::

## Agent Groups {#agent-group}

_Agent Group_ is a flexible label that defines a collection of agents that
operate as peers. For example, an Agent Group can be a Kubernetes cluster name
in case of DaemonSet deployment or it can be a service name for sidecar
deployments.

_Agent Group_ also defines the scope of **Agent-to-Agent synchronization**.
Agents within their group form a peer-to-peer network to synchronize
fine-grained state such as per-label global counters that are used for [rate
limiting purposes][dc]. Also, all the agents within an _Agent Group_ instantiate
the same set of [flowcontrol components][components], as published by Aperture
Controller.

Services in Aperture are additionally scoped within _Agent Groups_, creating two
level hierarchy, eg.:

<Zoom>

```mermaid
graph TB
    subgraph group2
        s3[search.mynamespace.svc.cluster.local]
        s4[db.mynamespace.svc.cluster.local]
    end
    subgraph group1
        s1[frontend.mynamespace.svc.cluster.local]
        s2[db.mynamespace.svc.cluster.local]
    end
```

</Zoom>

In this example there are two independent _db.mynamespace.svc.cluster.local_
services.

For single-cluster deployments, a single `default` _Agent Group_ can be used:

<Zoom>

```mermaid
graph TB
    subgraph default
        s1[frontend.mynamespace.svc.cluster.local]
        s3[search.mynamespace.svc.cluster.local]
        s2[db.mynamespace.svc.cluster.local]
    end
```

</Zoom>

as an other extreme, if your _Agent Groups_ already group entities into logical
services, you can treat the _Agent Group_ as a service and catch-all-service can
be used to match flows to policies (useful when installing as a sidecar):

<Zoom>

```mermaid
graph TB
    subgraph frontend
        s1[*]
    end
    subgraph search
        s2[*]
    end
    subgraph db
        s3[*]
    end
```

</Zoom>

_Agent group_ name together with _service_ name determine the [service][service]
to select flows from.

## Flow Matcher {#flow-matcher}

:::info

See also [_Flow Matcher_ reference](/reference/policies/spec.md#flow-matcher)

:::

### Control Point {#control-point}

_Flow Matcher_ selects flows from a specific control point within a
[service](#service).

<Zoom>

```mermaid
graph LR
  users(("users"))
  subgraph Frontend Service
    fingress["ingress"]
    recommendations{{"recommendations"}}
    live-update{{"live-update"}}
    fegress["egress"]
  end
  subgraph Checkout Service
    cingress["ingress"]
    cegress["egress"]
  end
  subgraph Database Service
    dbingress["ingress"]
  end
  users -.-> fingress
  fegress -.-> cingress
  cegress -.-> dbingress
```

</Zoom>

In the above diagram, each service has **HTTP** control points. Every incoming
API request to a service is a flow at its **ingress** control point. Likewise
every outgoing request from a service is a flow at its **egress** control point.

In addition, `Frontend` service has **Feature** control points identifying
_recommendations_ and _live-update_ features inside the `Frontend` service's
code.

:::note

_Control Point_ definition doesn't care about which particular entity (like a
pod) is handling particular flow. A single _Control Point_ covers _all_ the
entities belonging to the same service.

:::

:::tip

You can use [`aperturectl flow-control control-points`][aperturectl] to list
active control points:

:::

### Label Matcher {#label-matcher}

_Label Matcher_ allows to optionally narrow down the selected flow based on
conditions on [Flow Labels][label].

There are multiple ways to define a label matcher. The simplest way is to
provide a map of labels for exact-match:

```yaml
label_matcher:
  match_labels:
    http.method: GET
```

You can also provide a matching-expression-tree, which allows for arbitrary
conditions, including regex matching. Refer to [Label Matcher
reference][label-matcher] for further details.

## Example

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

[label]: ./flow-label.md
[flux-meter]: ./resources/flux-meter.md
[cl]: ./components/concurrency-limiter.md
[classifier]: ./resources/classifier.md
[label-matcher]: /reference/policies/spec.md#label-matcher
[dc]: components/rate-limiter.md#distributed-counters
[components]: ./components/components.md
[aperturectl]: /get-started/aperture-cli/aperture-cli.md
