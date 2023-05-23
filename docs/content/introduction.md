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
import { Cards } from '@site/src/components/Cards';
```

Welcome to the [FluxNinja Aperture](https://github.com/fluxninja/aperture)
project, an open source platform designed to empower platform and reliability
engineering teams. The platform offers a unified controllability layer that
simplifies the management of modern applications. It enables teams to
effortlessly implement intelligent load management capabilities, ensuring
optimal performance at any scale and for any infrastructure stack.

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

```mdx-code-block
<Cards data={[
  {
    title: "Intelligent Auto Scaling",
    description: "Aperture adjusts resource allocation based on demand and performance to ensure that the application can scale up or down as needed; However, it is different from traditional auto-scaling as it is based on the policies defined by the user which take multiple factors into consideration.",
    url: "/applying-policies/auto-scale",
  },
  {
    title: "Circuit-Based Policies",
    description: "Aperture provides a policy language that enables teams to define how to react to different situations, such as when there is a deviation from service-level objectives. These policies are expressed as a signal processing circuit that enables Aperture to go from telemetry to appropriate actions.",
    url: "/concepts/policy",
  },
  {
    title: "Dark Launch (aka Feature Flag Rollout)",
    description: "Aperture enables teams to gradually release new features to a subset of users, without impacting the rest of the system, using dark launch.",
    url: "/applying-policies/feature-rollout",
  },
  {
    title: "Distributed Rate-Limiting",
    description: "Aperture includes a distributed rate-limiter to prevent abuse and protect the service from excessive requests by users.",
    url: "/applying-policies/rate-limiting",
  },
  {
    title: "Intelligent Workload Classification",
    description: "Aperture provides a weighted fair queuing scheduler to ensure that the most critical workloads are served first based on the requirements of your application. Requirements are defined in a policy which helps Aperture to classify the workloads into different classes.",
    url: "/applying-policies/service-protection/workload-prioritization",
  },
  {
    title: "Monitoring and Telemetry",
    description: "Aperture continuously monitors service performance and request attributes using an in-built telemetry system, which enables the agent and controller to make informed decisions about how to handle requests and prioritize workloads.",
    url: "/reference/observability",
  },
  {
    title: "Quota Scheduler",
    description: "It is an important feature of Aperture, as it helps manage and allocate quotas for API requests. This feature enables teams to manage costs effectively and prevent exceeding these limits. By using quota scheduler, you can ensure that your application stays within the allowed usage and prevent being perceived as a bad actor by the API provider.",
    url: "/applying-policies/quota-scheduler/",
  },
  {
    title: "Service Protection",
    description: "The Service Protection feature with Adaptive Concurrency Limiting is a powerful mechanism designed to safeguard microservices within a distributed system. By intelligently managing the concurrency of requests, it ensures the stability, reliability, and optimal performance of individual services.",
    url: "/applying-policies/service-protection",
  }
  ]}/>
```

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
