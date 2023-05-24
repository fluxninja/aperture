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
  - dark launch
  - workload prioritization
  - rate limiting
  - observability
  - feature rollout
  - feature flag
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
import { Cards } from '@site/src/components/Cards';
```

## What is Aperture?

The [FluxNinja Aperture](https://github.com/fluxninja/aperture) project, an open
source Intelligent Load Management platform that seamlessly integrates into any
tech stack. This innovative platform is built to empower developers, platform
engineers, and reliability engineers, providing them with an advanced control
mechanism dovetailed with an observability layer. Aperture streamlines the task
of handling diverse traffic load intensities, spanning from low throughput
instances to conditions necessitating web-scale capacities.

Its core functionality facilitates teams to effortlessly implement intelligent
load management strategies, fostering optimized performance and promoting
maximal infrastructure utilization. Crucially, Aperture also ensures the
maintenance of an optimal end-user experience, even during service failures.
This versatility and resilience apply across all types of systems, including
monolithic architectures and distributed microservices environments.

<!-- vale off -->

## Where to get help?

<!-- vale on -->

If you have questions about how to use Aperture, feel free to reach out:

| [Ask the expert](https://calendly.com/desaijai/fluxninja-meeting) |
[Join our Slack Community](https://join.slack.com/t/fluxninja-aperture/shared_invite/zt-1vm2t2yjb-AG8rzKkB5TpPmqihJB6YYw)
| Email: [support@fluxninja.com](mailto:support@fluxninja.com) |

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
    title: "Declarative policies",
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
    title: "Workload Prioritization",
    description: "Aperture provides a weighted fair queuing scheduler to ensure that the most critical workloads are served first based on the requirements of your application.",
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

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)

## Get started

```mdx-code-block

<Cards data={[{
  title: "Setting up your application",
  url: "/getting-started/setting-up-application/",
},
{
  title: "Install Aperture",
  url: "/getting-started/installation/",
},
{
  title: "Your First Policy",
  url: "/getting-started/policies",
},
{
  title: "Applying Policies",
  url: "/applying-policies/",
}
]}/>

```

## Learn

For a high-level overview that explains how Aperture works, check out the
Concepts section:

```mdx-code-block
<Cards data={[{
  title: "Concepts",
  url: "/concepts",
}
]}/>
```

<!-- vale off -->

To understand how Aperture interfaces with your application, take a look at the
[Architecture](/architecture/architecture.md) section.

```mdx-code-block
<Cards data={[{
  title: "Architecture",
  url: "/architecture",
}
]}/>
```
