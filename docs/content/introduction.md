---
title: Introduction
slug: /
sidebar_position: 1
sidebar_class_name: introduction
description:
  Introduction to FluxNinja Aperture, an intelligent load management platform
  for modern cloud applications.
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

If you have questions about how to use Aperture, feel free to reach out. We are
happy to help!

<!-- vale off -->

[**💬 Ask the expert**](https://calendly.com/desaijai/fluxninja-meeting) |
[**👥 Join our Slack Community**](https://join.slack.com/t/fluxninja-aperture/shared_invite/zt-1vm2t2yjb-AG8rzKkB5TpPmqihJB6YYw)
| ✉️ Email: [**support@fluxninja.com**](mailto:support@fluxninja.com)

<!-- vale on -->

## ⚙️ Load management capabilities

Aperture offers a suite of intelligent load management capabilities that are
applicable to a wide range of cloud-native applications. These capabilities
ensure the reliability and stability of applications, and include:

- 🔀
  [**Intelligent Auto Scaling**](./applying-policies/auto-scale/auto-scale.md):
  Aperture adjusts resource allocation based on demand and performance to ensure
  that the application can scale up or down as needed; However, it is different
  from traditional auto-scaling as it is based on the policies defined by the
  user which take multiple factors into consideration.
- 📝 [**Declarative policies**](./concepts/policy/policy.md): Aperture provides
  a policy language that enables teams to define how to react to different
  situations, such as when there is a deviation from service-level objectives.
  These policies are expressed as a signal processing circuit that enables
  Aperture to go from telemetry to appropriate actions.
- 🚀
  [**Dark Launch (aka Feature Flag Rollout)**](./applying-policies/feature-rollout/feature-rollout.md):
  Aperture enables teams to gradually release new features to a subset of users,
  without impacting the rest of the system, using dark launch.
- ⏱️
  [**Distributed Rate-Limiting**](./applying-policies/rate-limiting/rate-limiting.md):
  Safeguard APIs from potential abuse with Aperture's high-performance,
  distributed rate limiter. This feature enforces per-key limits based on
  fine-grained labels, ensuring precise control and prevention of excessive
  usage.
- 🛡️
  [**Adaptive Service Protection**](./applying-policies/service-protection/basic-service-protection.md):
  Enhance resource utilization and safeguard against abrupt service overloads
  with an intelligent queue at the entry point of services. This queue
  dynamically adjusts the rate of requests based on live service health, thereby
  mitigating potential service disruptions and ensuring optimal performance
  under all load conditions.
- 🎯
  [**Workload Prioritization**](./applying-policies/service-protection/workload-prioritization.md):
  Safeguard crucial user experience pathways and ensure prioritized access to
  external APIs even during high-load conditions by strategically prioritizing
  workloads. This is achieved through the use of declarative policies that label
  and prioritize workload requests, such as API calls. By employing
  [weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
  for scheduling, Aperture ensures a fair distribution of resources that aligns
  with the business value and urgency of requests.
- 📊
  [**Intelligent quota management**](./applying-policies/quota-scheduler/quota-scheduler.md):
  Maintain compliance with external API quotas with a global token bucket and
  smart request queuing. This feature regulates requests aimed at external
  services, ensuring that the usage remains within prescribed rate limits and
  avoids penalties or additional costs.
- 🔍 [**Monitoring and Telemetry**](./reference/observability/observability.md):
  Aperture continuously monitors service performance and request attributes
  using an in-built telemetry system, which enables the agent and controller to
  make informed decisions about how to handle requests and prioritize workloads.

## 🛠️ How it works

At its core, load management involves the control loop of observing, analyzing,
and actuating workloads to ensure the stability and reliability of cloud-native
applications. This control loop is applied to both flow control and auto-scaling
use cases. In flow control, the control loop is used to manage workloads and
ensure the system remains within capacity. In auto-scaling, the control loop is
used to adjust resource allocation based on demand and performance.

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)

## ✨ Get started

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

## 📖 Learn

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
