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
  - flow control
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Welcome to the official documentation for
[FluxNinja Aperture](https://github.com/fluxninja/aperture)!

Aperture is an open-source platform that helps manage the flow of traffic and
improve the reliability of modern web applications. It enables the
prioritization of critical functions and prevent issues such as cascading
failures during periods of high traffic, ensuring that the overall performance
of the application is stable and reliable.

## Why is flow control needed?

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

## How Aperture Works

At the fundamental level, Aperture enables flow control through observing,
analyzing, and actuating, facilitated by agents and a controller.

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)

- Observe: Aperture continuously monitors the system and collects metrics on
  service performance and request attributes.
- Analyze: Aperture's agent and controller use the metrics collected to identify
  patterns and trends in the system and make decisions on how to handle requests
  and workloads.
- Actuate: Aperture takes appropriate actions, such as prioritizing critical
  workloads and shedding load on non-critical workloads to ensure the stability
  and reliability of the service in web-scale apps.

## What features does Aperture bring in?

Aperture is a flow control platform that offers several features to help
maintain the stability and reliability of modern web-scale applications:

- **Weighted fair queuing**: Aperture uses a weighted fair queuing scheduler to
  prioritize workloads based on their importance, ensuring that critical
  application features are not affected during overload scenarios.
- **Distributed rate-limiting**: Aperture includes a distributed rate-limiter to
  prevent abuse and protect the service from malicious requests.
- **Prioritization of critical features**: Aperture prioritizes critical
  application features over background workloads to ensure a graceful
  degradation of services during overload scenarios
- **Monitoring and telemetry**: Aperture continuously monitors service
  performance and request attributes using an in-built telemetry system, which
  enables the agent and controller to make informed decisions about how to
  handle requests and prioritize workloads.
- **Dataflow-driven policies**: Aperture's controller uses dataflow-driven
  policies to continuously track service-level indicators and perform recovery
  actions whenever there is a deviation from service-level objectives. This
  ensures that the application remains stable and reliable even in the face of
  failures.
- **Flexibility**: Aperture can be deployed with either Service Meshes or SDKs,
  depending on your infrastructure and requirements, and allows you to customize
  the flow control policies to meet the specific needs of your application
