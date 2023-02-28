---
title: Flow Control
sidebar_position: 1
keywords:
  - flows
  - tracing
  - opentracing
  - opentelemetry
---

# Flow Control

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
import {apertureVersion} from '../../../apertureVersion.js';
```

Modern web-scale apps are a complex network of inter-connected microservices
that implement features such as account management, search, payments & more.
This decoupled architecture has advantages but introduces new complex failure
modes. When traffic surges, it can result in a queue buildup on a critical
service, kick-starting a positive feedback loop and causing
[cascading failures](https://sre.google/sre-book/addressing-cascading-failures/).
The application stops serving responses in a timely manner and critical end-user
transactions are interrupted.

![Absence of flow control](assets/img/no-flow-control.png#gh-light-mode-only)
![Absence of flow control](assets/img/no-flow-control-dark.png#gh-dark-mode-only)

Applications are governed by
[Little’s Law](https://en.wikipedia.org/wiki/Little%27s_law), which describes
the relationship between concurrent requests in the system, arrival rate of
requests, and response times. For the application to remain stable, the
concurrent requests in the system must be limited. Indirect techniques to
stabilize applications such as rate-limiting and auto-scaling fall short in
enabling good user experiences or business outcomes. Rate-limiting individual
users are insufficient in protecting services. Autoscaling is slow to respond
and can be cost-prohibitive. As the number of services scales, these techniques
get harder to deploy.

![Reliability with flow control](assets/img/active-flow-control.png#gh-light-mode-only)
![Reliability with flow control](assets/img/active-flow-control-dark.png#gh-dark-mode-only)

This is where flow control comes in. Applications can degrade gracefully in
real-time when using flow control techniques with Aperture, by prioritizing
high-importance features over others.

A Flow is the fundamental unit of work from the perspective of an Aperture
Agent. It could be an API call, a feature, or even a database query. A Flow in
Aperture is similar to [OpenTelemetry span][span] and contains [flow
labels][flow-label].

Reliable operations at web-scale are impossible without effective flow control.
Aperture splits the process of flow control in two layers:

- Governing the flow control process and making high-level decisions. This is
  done by Aperture Controller through _policies_. You can read more about
  policies in [Policies chapter][policies].
- Actual execution of flow control is performed by Aperture Agent via
  [Concurrency Limiters][cl] or [Rate Limiters][rate-limiter]. Additionally the
  Agent handles other flow-control related tasks, like gathering metrics via
  [Flux Meters][flux-meter] and classifying traffic via
  [Classifiers][classifier]. This chapter describes flow control capabilities at
  the Agent.

## Insertion {#integrations}

For Aperture to be able to act at any of the Control Points, you need to install
integrations that will communicate with the Aperture Agent.

- _HTTP_ Control Points: Web framework and service-mesh based integrations
  expose control points at in the traffic path of a service.

  In principle, any web proxy or web framework can be integrated with Aperture
  in this way. These integrations use [Envoy's External Authorization
  API][ext-authz]. Integrations with several popular web frameworks are
  available.

  We provide integration instructions for [Istio/Envoy][istio]. The user can
  name the Control Point to identify a particular filter chain in Envoy. In case
  of insertion via Istio, a
  [default filter config](/get-started/integrations/flow-control/envoy/istio.md#envoy-filter),
  assigns _ingress_ and _egress_ Control Points as [identified by
  Istio][istio-patch-context].

- _Feature_ Control Points: We provide
  [Aperture SDKs](/get-started/integrations/flow-control/sdk/sdk.md) for popular
  languages. Aperture SDK wraps any function call or code snippet inside the
  Service code as a Feature Control Point. Every invocation of the Feature is a
  Flow from the perspective of Aperture.

  The SDK provides API to begin a flow which translates to a
  [flowcontrol.v1.Check][flowcontrol-proto] call into Agent. Response of this
  call contains a decision on whether to allow or reject the flow. The execution
  of a feature may be gated based on this decision. There is an API to end a
  flow which sends an OpenTelemetry span representing the flow to the Agent as
  telemetry.

:::note

Exact instructions on custom proxies / web frameworks / SDK integrations will be
added in the future.

:::

## Flow Control Components {#components}

Agent uses the following observability and control components (in order of
execution):

- [Classifiers][classifier]
- [Rate Limiter][rate-limiter]
- [Concurrency Limiter][cl]
- [Flux Meters][flux-meter]

You can learn more about each of the components in the subsequent sections, but
we recommend to start with concepts like [services][service] and
[labels][flow-label] first.

[policies]: /concepts/policy/policy.md
[cl]: ./components/concurrency-limiter.md
[rate-limiter]: components/rate-limiter.md
[flux-meter]: /concepts/integrations/flow-control/flux-meter.md
[classifier]: /concepts/integrations/flow-control/flow-classifier.md
[span]: https://opentelemetry.io/docs/reference/specification/trace/api/#span
[istio]: /get-started/integrations/flow-control/envoy/istio.md
[ext-authz]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#authorization-service-proto
[aperture-go]: https://github.com/FluxNinja/aperture-go
[service]: /concepts/integrations/flow-control/service.md
[flow-label]: /concepts/integrations/flow-control/flow-label.md
[flowcontrol-proto]:
  https://buf.build/fluxninja/aperture/docs/main:aperture.flowcontrol.v1
[istio-patch-context]:
  https://istio.io/latest/docs/reference/config/networking/envoy-filter/#EnvoyFilter-PatchContext
