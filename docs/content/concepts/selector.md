---
title: Selector
sidebar_label: Selector
sidebar_position: 4
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

:::info See also

[_Selector_ configuration specification.](/reference/configuration/spec.md#selector)

:::

_Selectors_ are used by flow control and observability components instantiated
by Aperture Agents like [_Classifiers_][classifier], [_Flux Meters_][flux-meter]
and [_Load Schedulers_][load-scheduler]. _Selectors_ define scoping rules that
decide how these components should select flows for their operations.

A Selector consists of the following parameters:

- [_Control Point_](#control-point) (required)
- [_Label Matcher_](#label-matcher) (optional)
- [_Agent Group_](#agent-group) (optional)
- [_Service_](#service) (optional)

**Example:**

```yaml
service: checkout.myns.svc.cluster.local # Service
control_point: ingress # Control Point
agent_group: default # Agent Group
label_matcher: # Label Matcher
match_labels:
  user_tier: premium
  http.method: GET
match_expressions:
  - key: query
    operator: In
    values:
      - insert
      - delete
expression: # Using Label Matcher with expression
  label_matches:
    - label: user_agent
      regex: ^(?!.*Chrome).*Safari
```

## Control Point {#control-point}

Control point is the only required parameter in the selector that identifies
either a feature in the code or an interception point within a proxy or
middleware. For more details, refer to the
[control points concept](control-point.md).

## Label Matcher {#label-matcher}

The label matcher is used to narrow down the selected flows using conditions
defined on [Labels][label]. It allows for precise filtering of flows based on
specific criteria.

There are multiple ways to define a label matcher. If multiple match criteria
are defined simultaneously, then they all must match for a flow to be selected.

- **Exact Match**: It is the simplest way to match a label. It matches the label
  value exactly.

  ```yaml
  label_matcher:
    match_labels:
      http.method: GET
  ```

- **Matching Expressions**: It allows for more complex matching conditions using
  operators such as `In`, `NotIn`, `Exists`, and `DoesNotExists`.

  ```yaml
  label_matcher:
    match_expressions:
      - key: http.method
        operator: In
        values:
          - GET
          - POST
  ```

- **Arbitrary Expression**: This allows for defining complex matching
  conditions, including regular expression matching.

  ```yaml
  label_matcher:
    expression:
      label_matches:
        - label: user_agent
          regex: ^(?!.*Chrome).*Safari
  ```

Refer to [Label Matcher Reference][label-matcher] for further details on each of
these matching methods.

## Agent Group {#agent-group}

The agent group parameter identifies the Agents where the selector gets applied.
For more details, refer to the [agent group concept](./advanced/agent-group.md).

In the example below, the agent group is `prod-cluster`.

```yaml
agent_group: prod-cluster # Agent Group
control_point: ingress
label_matcher:
  match_labels:
    user_tier: gold
```

## Service {#service}

The service parameter limits the selector to flows that belong to the specified
service. For more details, refer to the
[service concept](./advanced/service.md).

:::tip Special Service Names

- `any`: Can be used in a policy to match all services

:::

In the example below, the service name is `checkout.myns.svc.cluster.local`.

```yaml
service: checkout.myns.svc.cluster.local #Service Name
agent_group: default
control_point: ingress
label_matcher:
  match_labels:
    user_tier: gold
```

## Selectors for different scenarios

### Features defined through SDK

Control point alone is sufficient to identify flows belonging to a unique
feature.

An example flow:

| Agent Group | Control Point                | Service                                 | Flow Labels     |
| ----------- | ---------------------------- | --------------------------------------- | --------------- |
| default     | smart-recommendation-feature | checkout-service.prod.svc.cluster.local | user_type:guest |

### DaemonSet deployment of Aperture Agents

In this installation mode, control points alone might not be sufficient to
identify flows. For instance, in the case of Envoy interception, the default
Envoy Filter assigns `ingress` and `egress` as the control points to all
listener chains, depending on whether they are for ingress or egress traffic. In
this case, the service parameter can be used to distinguish among flows
belonging to different services.

An example flow:

<!-- vale off -->

| Agent Group | Control Point | Service                                 | Flow Labels                                                |
| ----------- | ------------- | --------------------------------------- | ---------------------------------------------------------- |
| default     | ingress       | checkout-service.prod.svc.cluster.local | http.request.header.user_type:guest, http.target:/checkout |

<!-- vale on -->

### Aperture Agent as a sidecar container

In this installation mode, Kubernetes service discovery is disabled by default.
It is recommended to configure each service's deployment with a unique Aperture
agent group. Therefore, the selector in this scenario can identify the flows for
a service based on the agent group and control point.

Example flows:

<!-- vale off -->

| Agent Group      | Control Point                | Service | Flow Labels                                                |
| ---------------- | ---------------------------- | ------- | ---------------------------------------------------------- |
| checkout-service | ingress                      |         | http.request.header.user_type:guest, http.target:/checkout |
| checkout-service | smart-recommendation-feature |         | user_type:guest                                            |

<!-- vale on -->

### Standalone deployment of Aperture Agents

In this installation mode, Aperture Agents are installed as a standalone load
management service. Clients call into this service to get load management
functionality. The selectors in this scenario rely on control points as the
distinguishing factor to identify unique features and services.

An example flow:

<!-- vale off -->

| Agent Group | Control Point                | Service | Flow Labels     |
| ----------- | ---------------------------- | ------- | --------------- |
| default     | smart-recommendation-feature |         | user_type:guest |

<!-- vale on -->

### Gateways Integration {#gateways-integration}

Aperture can be integrated with [Gateways][gateway] to control traffic before it
is routed to the upstream service. Gateways are configured to send flow control
requests to Aperture for every incoming request.

As the requests to Aperture are sent from the Gateway, the selector has to be
configured to match the Gateway's service. For example, if the Gateway
controller is running with service name `nginx-server` in namespace `nginx`, for
upstream service having location or route as `/search-service`, the selector
should be configured as follows:

```yaml
service: nginx-server.nginx.svc.cluster.local
agent_group: default
control_point: search-service
label_matcher:
  match_labels:
    http.target: "/search-service"
```

An example flow:

<!-- vale off -->

| Agent Group | Control Point  | Service                              | Flow Labels                 |
| ----------- | -------------- | ------------------------------------ | --------------------------- |
| default     | search-service | nginx-server.nginx.svc.cluster.local | http.target:/search-service |

<!-- vale on -->

Also, if the control point is configured uniquely for each location or route,
the `control_point` alone can be used to match the upstream service and the rest
of the parameters can be omitted:

```yaml
agent_group: default
control_point: search-service
```

An example flow:

<!-- vale off -->

| Agent Group | Control Point  | Service                              | Flow Labels                 |
| ----------- | -------------- | ------------------------------------ | --------------------------- |
| default     | search-service | nginx-server.nginx.svc.cluster.local | http.target:/search-service |

<!-- vale on -->

## Filtering out health and metrics endpoints

Liveness and health probes are essential for checking the health of the
application, and metrics endpoints are necessary for monitoring its performance.
However, these endpoints do not usually represent the intended workload in an
Aperture policy. If included in a _Flux Meter_, they can reduce the accuracy of
latency calculations. If included in an actuation component like _Load
Scheduler_, they might cause these requests to be rejected under load, leading
to unnecessary pod restarts.

To prevent these issues, traffic to these endpoints can be filtered out by
matching expressions. In the example below, flows with `http.target` starting
with `/health`, `/live`, or `/ready`, and User Agent starting with
`kube-probe/1.23` are filtered out.

```yaml
service: checkout.myns.svc.cluster.local
agent_group: default
control_point: ingress
label_matcher:
  match_expressions:
    - key: http.target
      operator: NotIn
      values:
        - /health
        - /live
        - /ready
        - /metrics
    - key: http.user_agent
      operator: NotIn
      values:
        - kube-probe/1.23
```

[label]: ./flow-label.md
[flux-meter]: ./advanced/flux-meter.md
[load-scheduler]: ./scheduler/load-scheduler.md
[classifier]: ./advanced/classifier.md
[label-matcher]: /reference/configuration/spec.md#label-matcher
[gateway]: /aperture-for-infra/integrations/gateway/gateway.md
