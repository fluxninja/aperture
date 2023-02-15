---
title: Architecture
sidebar_position: 2
description:
  Discover the core components of Aperture Architecture and learn how they work
  together to provide powerful and efficient reliability automation.
image: ../assets/img/aperture_logo.png
keywords:
  - reliability
  - overload
  - concurrency
  - aperture
  - fluxninja
  - microservices
  - cloud
  - TODO
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

FluxNinja Aperture architecture in detail

<Zoom>

```mermaid
{@include: ../assets/diagrams/architecture/architecture_simple.mmd}
```

</Zoom>

FluxNinja Aperture consists of multiple components that build up the
architecture (As shown in above diagram). Some of the core components are list
below.

- [Aperture Controller](#aperture-controller)
  - Controller Circuit
- [Aperture Agent](#aperture-agents)
  - Flow Controller
  - Auto Scaler

### Aperture Controller

The Aperture Controller is powered by always-on, dataflow-driven policies that
continuously track deviations from service-level objectives (SLOs) and calculate
recovery or escalation actions. The policies running in the Aperture Controller
are expressed as circuits, much like circuit networks in the game
[Factorio](https://wiki.factorio.com/Circuit_network). Aperture Controller store
certain list of metrics in Prometheus and Etcd store information about decision
it took.

For example, a gradient control circuit component can be used to implement
[AIMD](https://en.wikipedia.org/wiki/Additive_increase/multiplicative_decrease)
(Additive Increase, Multiplicative Decrease) style counter-measure that limits
the concurrency on a service when response times deteriorate. Advanced control
components like [PID](https://en.wikipedia.org/wiki/PID_controller) can be used
to further tune the concurrency limits.

Aperture Controller is comparable in capabilities to autopilot in aircraft or
adaptive cruise control in some automobiles.

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

Aperture also includes an integration for auto-scaling, which can help you scale
your infrastructure as needed by countermeasures.
