---
title: Flow Control
sidebar_position: 2
keywords:
  - flows
  - tracing
  - opentracing
  - opentelemetry
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
import {apertureVersion} from '../../apertureVersion.js';
```

Many applications consist of intricate networks of interconnected services that
drive essential features. While both monolithic and decoupled architectures
offer various advantages, they also introduce new challenges and complexities.
During periods of high traffic, a critical service may experience queue buildup,
triggering a detrimental positive feedback loop and leading to
[cascading failures](https://sre.google/sre-book/addressing-cascading-failures/).
As a result, the application becomes unresponsive, causing disruptions to
crucial end-user transactions.

![Absence of flow control](assets/img/no-flow-control.png#gh-light-mode-only)
![Absence of flow control](assets/img/no-flow-control-dark.png#gh-dark-mode-only)

Applications are governed by
[Little’s Law](https://en.wikipedia.org/wiki/Little%27s_law), which describes
the relationship between concurrent requests in the system, arrival rate of
requests, and response times. For the application to remain stable, the
concurrent requests in the system must be limited. Indirect techniques to
stabilize applications such as rate-limiting and auto-scaling fall short in
enabling good user experiences or business outcomes. Rate-limiting individual
users is insufficient in protecting services. Auto-scaling is slow to respond
and can be cost-prohibitive. As the number of services scales, these techniques
get harder to deploy.

![Reliability with flow control](assets/img/active-flow-control.png#gh-light-mode-only)
![Reliability with flow control](assets/img/active-flow-control-dark.png#gh-dark-mode-only)

This is where flow control comes in. Applications can degrade gracefully in
real-time when using flow control techniques with Aperture, by prioritizing
high-importance features over others. Reliable operations at web-scale are
impossible without effective flow control.

Aperture splits the process of flow control into two layers:

- Governing the flow control process and making high-level decisions. This is
  done by the Aperture Controller through [_Policies_][policies].
- Actual execution of flow control is performed by Aperture Agent through [_Load
  Regulators_][regulator], [_Load Schedulers_][load-scheduler] and [_Rate
  Limiters_][rate-limiter]. Additionally, the Agent handles other flow-control
  related tasks, like gathering metrics through [_Flux Meters_][flux-meter] and
  classifying traffic through [_Classifiers_][classifier]. This chapter
  describes flow control capabilities at the Agent.

## Insertion {#insertion}

For Aperture to be able to act at any of the [_Control Points_][control-point],
you need to install integrations that will communicate with the Aperture Agent.

- _HTTP_ _Control Points_: Web framework and service-mesh based integrations
  expose _Control Points_ at in the traffic path of a service.

  In principle, any web proxy or web framework can be integrated with Aperture
  in this way. These integrations use [Envoy's External Authorization
  API][ext-authz]. Integrations with several popular web frameworks are
  available.

  Integration instructions for [Istio/Envoy][istio] are provided, and the
  Control Point can be named to identify a particular filter chain in Envoy. If
  insertion is done through Istio, the
  [default filter configuration](/get-started/integrations/flow-control/envoy/istio.md#envoy-filter)
  assigns _ingress_ and _egress_ Control Points as identified by
  [Istio][istio-patch-context].

- _Feature_ _Control Points_:
  [Aperture SDKs](/get-started/integrations/flow-control/sdk/sdk.md) are
  available for popular programming languages. Aperture SDK wraps any function
  call or code snippet inside the service code as a _Feature_ _Control Point_.
  Every invocation of the feature is a flow from the perspective of Aperture.

  The SDK provides an API to begin a flow, which translates to a
  [`flowcontrol.v1.Check`][flowcontrol-proto] call into Agent. The response of
  this call contains a decision on whether to allow or reject the flow. The
  execution of a feature might be gated based on this decision. There is an API
  to end a flow, which sends an [OpenTelemetry span][span] representing the flow
  to the _Agent_ as telemetry.

[policies]: /concepts/policy/policy.md
[control-point]: ./selector.md#control-point
[load-scheduler]: ./components/load-scheduler.md
[regulator]: ./components/regulator.md
[rate-limiter]: ./components/rate-limiter.md
[flux-meter]: ./resources/flux-meter.md
[classifier]: ./resources/classifier.md
[span]: https://opentelemetry.io/docs/reference/specification/trace/api/#span
[istio]: /get-started/integrations/flow-control/envoy/istio.md
[ext-authz]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#authorization-service-proto
[flowcontrol-proto]:
  https://buf.build/fluxninja/aperture/docs/main:aperture.flowcontrol.check.v1
[istio-patch-context]:
  https://istio.io/latest/docs/reference/config/networking/envoy-filter/#EnvoyFilter-PatchContext
