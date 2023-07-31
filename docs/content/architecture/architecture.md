---
title: Architecture
sidebar_position: 2
description:
  Discover the core components of Aperture architecture and learn how they work
  together to provide powerful and efficient load management.
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

## Aperture Controller

The Aperture Controller is a centralized control system, equipped with a
comprehensive global perspective. It is programmed using declarative policies
that are stored in a policy database that can be managed using the Kubernetes
Custom Resource Definition (CRD) API, allowing users to configure and modify
policies as needed.

A policy represents a closed-loop control circuit that is executed periodically.
The control circuit draws input signals from metrics aggregated across Aperture
Agents, providing the Controller with a holistic view of the application's
health and performance. Service-level objectives (SLOs) are defined against
these health and performance signals. The policies continuously track deviations
from SLOs and calculate recovery or escalation actions that are translated as
adjustments to the Agents.

After computing the adjustments, the Aperture Controller synchronizes them with
the relevant Aperture Agents. These adjustments encompass load throttling,
workload prioritization, and auto-scaling actions, among others. By
disseminating the calculated adjustments to the Agents, the Controller ensures
that the Agents take localized actions in line with the global state of the
system.

## Aperture Agents

Aperture Agents are the workhorses of the platform, providing powerful flow
control components such as a weighted fair queuing scheduler for workload
prioritization and a distributed rate-limiter for abuse prevention. Agents
integrate with service meshes, gateways and HTTP middlewares. Alternately,
developers can use SDKs to get flow control around specific features or code
sections inside services.

The Agents monitor service and infrastructure health signals using an in-built
telemetry system. In addition, a programmable, high-fidelity flow classifier is
used to label requests based on attributes such as customer tier or request
type. These metrics are then analyzed by the Aperture Controller.

Aperture Agents schedule workloads based on their priorities, helping prioritize
critical features over less important workloads during overload scenarios. For
example, a video streaming service might prioritize a request to play a movie by
a customer over a recommended movies API. A SaaS product might prioritize
features used by paid users over those being used by free users.

Aperture Agents can be installed on a variety of infrastructure such as
Kubernetes, VMs, or bare-metal. In addition to flow control capabilities, Agents
work with auto-scaling APIs for platforms such as Kubernetes, to help scale
infrastructure when needed.

## Aperture Databases

Aperture uses two databases to store configuration, telemetry, and flow control
information: [Prometheus](https://prometheus.io) and [etcd](https://etcd.io).
Prometheus enables Aperture to monitor the system and detect deviations from the
service-level objectives (SLOs) defined in the declarative policies. Aperture
Controller uses etcd (distributed key-value store) to persist the declarative
policies that define the control circuits and their components, as well as the
adjustments synchronized between the Controller and Agents.

Users can optionally reuse their existing etcd and
[scalable Prometheus](https://promlabs.com/blog/2021/10/14/promql-vendor-compatibility-round-three)
installations to minimize operational overhead and use their existing monitoring
infrastructure.
