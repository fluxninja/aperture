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

Aperture Agents are the workhorses of the platform, residing alongside service
instances as sidecars. They provide powerful flow control components such as a
weighted fair queuing scheduler for prioritized load shedding and a distributed
rate-limiter for abuse prevention. A flow is the fundamental unit of work from
the perspective of an Aperture Agent. It could be an API call, a feature, or
even a database query.

The agents monitor golden signals, such as request latency, error rate, and
saturation, using an in-built telemetry system and a programmable, high-fidelity
flow classifier used to label requests based on attributes such as customer tier
or request type. These metrics are then analyzed by the Aperture Controller.

Graceful degradation of services is achieved by prioritizing critical
application features over background workloads. Similar to boarding an aircraft,
business class passengers get priority over other passengers; every application
has workloads with varying priorities. For example, a video streaming service
might prioritize a request to play a movie by a customer over running an
internal machine learning workload. A SaaS product might prioritize features
used by paid users over those being used by free users. Aperture Agents schedule
workloads based on their priorities, helping maximize user experience or revenue
even during overload scenarios.

Aperture Agents are capable of working with different infrastructures such as
Kubernetes, VM, or bare-metal. They integrate with Service Meshes or can be used
with SDKs depending on your requirements. Additionally, agents work with
auto-scaling APIs for platforms such as Kubernetes, to help scale infrastructure
when needed.

## Aperture Databases

Aperture uses two databases to store configuration, telemetry, and flow control
information: [Prometheus](https://prometheus.io) and [etcd](https://etcd.io).
Prometheus is a time-series database used to store and query telemetry data
collected from Aperture Agents. It enables Aperture to monitor the system and
detect deviations from the service-level objectives (SLOs) defined in the
declarative policies.

Etcd is a distributed key-value store used to store configuration and flow
control information. Aperture Controller uses etcd to store the declarative
policies that define the control circuits and their components, as well as the
current system state.

Users can optionally reuse their existing etcd or
[scalable Prometheus](https://promlabs.com/blog/2021/10/14/promql-vendor-compatibility-round-three)
installations to minimize operational overhead and leverage their existing
monitoring infrastructure.
