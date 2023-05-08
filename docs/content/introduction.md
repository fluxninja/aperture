---
title: Introduction
slug: /
sidebar_position: 1
sidebar_class_name: introduction
description:
  Introduction to FluxNinja Aperture, an open source flow control and
  reliability management platform for modern web applications.
image: /assets/img/aperture_logo.png
keywords:
  - reliability
  - overload
  - concurrency
  - aperture
  - fluxninja
  - microservices
  - cloud
  - auto-scale
  - load management
  - flow control
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Welcome to the [FluxNinja Aperture](https://github.com/fluxninja/aperture)
project, an open source platform that empowers platform and reliability
engineering teams. The platform provides a unified controllability layer that
enables platform and reliability engineering teams to manage complex
microservices-based applications with ease.

## Simplify cloud native load management

With Aperture, teams can automate load management processes, including flow
control and auto-scaling, to ensure the reliability and stability of cloud
native applications. These capabilities improve the overall user experience,
while optimizing resources and reducing costs.

## Declarative policy language

Aperture's declarative policy language allows teams to effortlessly develop and
manage policies that dictate their applications' behavior under various
circumstances. This offers a visual depiction of their policies, enabling them
to grasp the system's behavior and self-correction intuitively.

## Load management capabilities

Aperture offers a suite of intelligent load management capabilities that are
applicable to a wide range of cloud-native applications. These capabilities
ensure the reliability and stability of applications, and include:

- **Intelligent Workload Classification**: Aperture provides a weighted fair
  queuing scheduler to ensure that the most critical workloads are served first
  based on the requirements of your application.
- **Circuit-Based policies**: Aperture provides a policy language that enables
  teams to define how to react to different situations, such as when there is a
  deviation from service-level objectives. These policies are expressed as a
  signal processing circuit that enables Aperture to go from telemetry to
  appropriate actions.
- **Dark Launch (aka Feature Flag Rollout)**: Aperture enables teams to
  gradually release new features to a subset of users, without impacting the
  rest of the system, using dark launch.
- **Distributed rate-limiting**: Aperture includes a distributed rate-limiter to
  prevent abuse and protect the service from excessive requests by users.
- **Intelligent Auto Scaling**: Aperture adjusts resource allocation based on
  demand and performance to ensure that the application can scale up or down as
  needed; However, it is different from traditional auto-scaling as it is based
  on the policies defined by the user which take multiple factors into
  consideration.
- **Monitoring and telemetry**: Aperture continuously monitors service
  performance and request attributes using an in-built telemetry system, which
  enables the agent and controller to make informed decisions about how to
  handle requests and prioritize workloads.
- **Quota Scheduling**: It is an important feature of Aperture, as it helps
  manage and allocate quotas for API requests. This feature enables teams to
  manage costs effectively and prevent exceeding API rate limits. By using quota
  scheduling, you can ensure that your application stays within the limits of
  the API and prevent being perceived as a bad actor by the API provider.

## How it works

At its core, load management involves the control loop of observing, analyzing,
and actuating workloads to ensure the stability and reliability of cloud-native
applications. This control loop is applied to both flow control and auto-scaling
use cases. In flow control, the control loop is used to manage workloads and
ensure the system remains within capacity. In auto-scaling, the control loop is
used to adjust resource allocation based on demand and performance.

Learn more about how Aperture interfaces with your application in the
[Architecture](/architecture/architecture.md) section.

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)
