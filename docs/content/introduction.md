---
title: Introduction
slug: /
sidebar_position: 1
sidebar_class_name: introduction
description:
  Introduction to FluxNinja Aperture, an open-source flow control and
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
project, an open-source platform designed to empower cloud-native reliability
engineering. The platform provides a unified controllability layer that enables
platform and reliability engineering teams to manage complex microservices-based
applications with ease.

## Simplify Cloud-Native Load Management

With Aperture, teams can automate load management processes, including flow
control and auto scaling, to ensure the reliability and stability of
cloud-native applications. These capabilities improve the overall user
experience, while optimizing resources and reducing costs.

## Declarative Policy Language

Aperture's declarative policy language allows teams to effortlessly develop and
manage policies that dictate their applications' behavior under various
circumstances. This offers a visual depiction of their policies, enabling them
to grasp the system's behavior and self-correction intuitively.

## Load Management Capabilities

Aperture's intelligent load management capabilities, such as fine-grained rate
limiting, prioritized load shedding and auto scaling, can be applied to a wide
range of cloud-native applications. These capabilities, ensure the reliability
and stability of applications.

- **Prioritized load shedding**: Aperture enables organizations to gracefully
  degrade application performance by dropping traffic that's deemed less
  important, ensuring that the most critical traffic is served.
- **Distributed rate-limiting**: Aperture includes a distributed rate-limiter to
  prevent abuse and protect the service from excessive requests by users.
- **Intelligent auto scaling**: Aperture adjusts resource allocation based on
  demand and performance to ensure that the application can scale up or down as
  needed.
- **Monitoring and telemetry**: Aperture continuously monitors service
  performance and request attributes using an in-built telemetry system, which
  enables the agent and controller to make informed decisions about how to
  handle requests and prioritize workloads.
- **Declarative policies**: Aperture provides a policy language that enables
  teams to define how to react to different situations, such as when there is a
  deviation from service-level objectives. These policies are expressed as a
  signal processing circuit that enables Aperture to go from telemetry to
  actions within minutes.

## How it works

At its core, load management involves the control loop of observing, analyzing,
and actuating workloads to ensure the stability and reliability of cloud-native
applications. This control loop is applied to both flow control and auto scaling
use cases. In flow control, the control loop is used to manage workloads and
ensure the system remains within capacity. In auto scaling, the control loop is
used to adjust resource allocation based on demand and performance.

Learn more about how Aperture interfaces with your application in the
[Architecture](/architecture/architecture.md) section.

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)
