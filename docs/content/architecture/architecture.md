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

Aperture is built on a distributed architecture that provides a unified
observability and controllability platform for cloud-native applications. The
architecture is designed to ensure high availability, scalability, and
reliability.

<Zoom>

```mermaid
{@include: ../assets/diagrams/architecture/architecture_simple.mmd}
```

</Zoom>

### Aperture Controller

The Aperture Controller is the central component of the platform. The controller
monitors the system using an in-built telemetry system and collects metrics on
service performance and workloads, including information on customer tiers,
request types, and other relevant attributes.

he controller uses declarative policies, expressed as a control circuit, to
analyze the collected metrics and make decisions on load shedding, flow control,
and auto scaling to ensure that the application operates within the specified
SLOs. The controller's policies are based on the principles of
Observability-driven closed-loop automation, which continuously track deviations
from service-level objectives (SLOs) and calculate recovery or escalation
actions.

For example, a gradient control circuit component can be used to implement
[AIMD](https://en.wikipedia.org/wiki/Additive_increase/multiplicative_decrease)
(Additive Increase, Multiplicative Decrease) style closed-loop automation that
limits the concurrency on a service when response times deteriorate. Advanced
control components like
[PID controller](https://en.wikipedia.org/wiki/PID_controller) can be used to
further tune the concurrency limits based on specific service requirements.

The controller's policies are stored in a policy database and are managed using
the Kubernetes Custom Resource Definition (CRD) API, allowing users to easily
configure and modify policies as needed. The controller interacts with Aperture
Agents, which run alongside service instances as sidecars, to enforce the
policies and ensure the reliable operation of cloud-native applications.

### Aperture Agents

Aperture Agents are the building blocks of the Aperture platform, residing next
to service instances as a sidecar. They provide powerful flow control
components, such as a weighted fair queuing scheduler for prioritized load
shedding, a distributed rate limiter for abuse prevention, and intelligent
autoscaling to handle changes in demand. Agents are responsible for observing,
analyzing, and acting on workloads, providing a foundation for intelligent load
management capabilities.

Aperture Agents schedule workloads based on their priorities, helping maximize
user experience or revenue even during overload scenarios. They prioritize
critical application features over background workloads, much like when boarding
an aircraft, business class passengers get priority over other passengers; every
application has workloads with varying priorities.

Agents monitor golden signals using an in-built telemetry system and a
programmable, high-fidelity flow classifier used to label requests based on
attributes such as customer tier or request type. These metrics are analyzed by
the Aperture Controller to make informed decisions on flow control and auto
scaling.

Aperture Agents are capable of working with different infrastructures such as
Kubernetes, VM, or bare-metal. They integrate with Service Meshes or can be used
with SDKs depending on your requirements. Additionally, agents work with
auto-scaling integration to help scale infrastructure when needed.
