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

Reliable operations at web-scale are impossible without effective flow control.
Aperture splits the process of flow control in two layers:

- Governing the flow control process and making high-level decisions. This is
  done by Aperture Controller through [_Policies_][policies].
- Actual execution of flow control is performed by Aperture Agent via
  [_Concurrency Limiters_][cl] or [_Rate Limiters_][rate-limiter]. Additionally
  the Agent handles other flow-control related tasks, like gathering metrics via
  [_Flux Meters_][flux-meter] and classifying traffic via
  [_Classifiers_][classifier]. This chapter describes flow control capabilities
  at the Agent.

## Insertion {#insertion}

For Aperture to be able to act at any of the [_Control Points_][control-point],
you need to install integrations that will communicate with the Aperture Agent.

- _HTTP_ _Control Points_: Web framework and service-mesh based integrations
  expose _Control Points_ at in the traffic path of a service.

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

- _Feature_ _Control Points_: We provide
  [Aperture SDKs](/get-started/integrations/flow-control/sdk/sdk.md) for popular
  languages. Aperture SDK wraps any function call or code snippet inside the
  service code as a _Feature_ _Control Point_. Every invocation of th feature is
  a flow from the perspective of Aperture.

  The SDK provides API to begin a flow which translates to a
  [flowcontrol.v1.Check][flowcontrol-proto] call into Agent. Response of this
  call contains a decision on whether to allow or reject the flow. The execution
  of a feature may be gated based on this decision. There is an API to end a
  flow which sends an [OpenTelemetry span][span] representing the flow to the
  _Agent_ as telemetry.

[policies]: /concepts/policy/policy.md
[control-point]: ./flow-selector.md#control-point
[cl]: ./components/concurrency-limiter.md
[rate-limiter]: ./components/rate-limiter.md
[flux-meter]: ./resources/flux-meter.md
[classifier]: ./resources/classifier.md
[span]: https://opentelemetry.io/docs/reference/specification/trace/api/#span
[istio]: /get-started/integrations/flow-control/envoy/istio.md
[ext-authz]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#authorization-service-proto
[aperture-go]: https://github.com/FluxNinja/aperture-go
[flowcontrol-proto]:
  https://buf.build/fluxninja/aperture/docs/main:aperture.flowcontrol.check.v1
[istio-patch-context]:
  https://istio.io/latest/docs/reference/config/networking/envoy-filter/#EnvoyFilter-PatchContext
