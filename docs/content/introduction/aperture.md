---
id: aperture
title: Aperture Open Source
slug: /introduction/aperture
description: Aperture Introduction page
keywords:
  - observe
  - analyze
  - actuate
  - how-it-works
---

# Introduction

Welcome to the official guide for Aperture!

Aperture is an open-source flow control and reliability management platform for
modern web application.

## Why flow control is needed?

Modern web-scale apps are a complex network of inter-connected microservices
that implement features such as account management, search, payments & more.
This decoupled architecture has advantages but introduces new complex failure
modes. When traffic surges, it can result in a queue buildup on a critical
service, kick-starting a positive feedback loop and causing
[cascading failures](https://sre.google/sre-book/addressing-cascading-failures/).
The application stops serving responses in a timely manner and critical end-user
transactions are interrupted.

![Absence of flow control](../assets/img/no-flow-control.jpg)

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

![Reliability with flow control](../assets/img/active-flow-control.jpg)

This is where flow control comes in. Applications can degrade gracefully in
real-time when using flow control techniques with Aperture, by prioritizing
high-importance features over others.

## How Aperture Works

At the fundamental level, Aperture enables flow control through observing,
analyzing, and actuating, facilitated by agents and a controller.

![Aperture Control Loop](../assets/img/control-loop.svg)

### Aperture Agents

Aperture Agents live next to your service instances as a sidecar and provide
powerful flow control components such as a weighted fair queuing scheduler for
prioritized load-shedding and a distributed rate-limiter for abuse prevention. A
flow is the fundamental unit of work from the perspective of an Aperture Agent.
It could be an API call, a feature, or even a database query.

Graceful degradation of services is achieved by prioritizing critical
application features over background workloads. Much like when boarding an
aircraft, business class passengers get priority over other passengers; every
application has workloads with varying priorities. A video streaming service
might view a request to play a movie by a customer as a higher priority than
running an internal machine learning workload. A SaaS product might prioritize
features used by paid users over those being used by free users. Aperture Agents
schedule workloads based on their priorities helping maximize user experience or
revenue even during overload scenarios.

Aperture Agents monitor golden signals using an in-built telemetry system and a
programmable, high-fidelity flow classifier used to label requests based on
attributes such as customer tier or request type. These metrics are analyzed by
the controller.

### Aperture Controller

The controller is powered by always-on, dataflow-driven policies that
continuously track deviations from service-level objectives (SLOs) and calculate
recovery or escalation actions. The policies running in the controller are
expressed as circuits, much like circuit networks in the game
[Factorio](https://wiki.factorio.com/Circuit_network).

For example, a gradient control circuit component can be used to implement
[AIMD](https://en.wikipedia.org/wiki/Additive_increase/multiplicative_decrease)
(Additive Increase, Multiplicative Decrease) style counter-measure that limits
the concurrency on a service when response times deteriorate. Advanced control
components like [PID](https://en.wikipedia.org/wiki/PID_controller) can be used
to further tune the concurrency limits.

Aperture’s Controller is comparable in capabilities to autopilot in aircraft or
adaptive cruise control in some automobiles.

### Deploying Aperture

Aperture can be inserted into service instances with either Service Meshes or
SDKs:

- Service Mesh: Aperture can be deployed with no changes to application code,
  using [Envoy](https://www.envoyproxy.io/). It latches onto Envoy’s
  [External Authorization API](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter)
  for control purposes and collects access logs for telemetry purposes. On each
  request, Envoy sends request metadata to the Aperture Agent for a flow control
  decision. Inside the Aperture Agent, the request traverses classifiers,
  rate-limiters, and schedulers, before the decision to accept or drop the
  request is sent back to Envoy. Aperture participates in the
  [OpenTelemetry](https://opentelemetry.io/) tracing protocol as it inserts flow
  classification labels into requests, enabling visualization in tracing tools
  such as [Jaeger](https://www.jaegertracing.io/).
- Aperture SDKs: In addition to service mesh insertion, Aperture provides SDKs
  that can be used by developers to achieve fine-grained flow control at the
  feature level inside service code. For example, an e-commerce app may
  prioritize users in the checkout flow over new sessions when the application
  is experiencing an overload. The Aperture Controller can be programmed to
  degrade features as an escalated recovery action when basic load shedding is
  triggered for several minutes.
