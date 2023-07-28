---
title: Introduction – Aperture
sidebar_label: Introduction
sidebar_position: 1
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
  - load ramp
  - feature flag
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
import { Cards } from '@site/src/components/Cards';
```

<!-- vale off -->

## What is Aperture?

<!-- vale on -->

[Aperture](https://github.com/fluxninja/aperture) is an open source load
management platform designed for classifying, scheduling, and rate-limiting API
traffic in cloud applications. Built upon a foundation of observability and a
global control plane, it offers a comprehensive suite of load management
capabilities. These capabilities enhance the reliability and performance of
cloud applications while also optimizing resource utilization.

Aperture can seamlessly integrate with existing control points such as gateways,
service meshes, and application middlewares. Moreover, it offers SDKs for
developers who need to establish control points around specific features or code
sections inside applications.

## ⚙️ Load management capabilities {#load-management-capabilities}

Aperture provides a variety of advanced load management features:

- 🛡️
  [**Adaptive Service Protection**](./use-cases/adaptive-service-protection/adaptive-service-protection.md):
  Enhance resource utilization and safeguard against abrupt service overloads
  with an intelligent queue at the entry point of services. This queue
  dynamically adjusts the rate of requests based on live service health, thereby
  mitigating potential service disruptions and ensuring optimal performance
  under all load conditions.
- 📊
  [**Global Quota Management**](./use-cases/managing-quotas/managing-quotas.md):
  Maintain compliance with external API quotas with a global token bucket and
  smart request queuing. This feature regulates requests aimed at external
  services, ensuring that the usage remains within prescribed rate limits and
  avoids penalties or additional costs.
- 🎯
  [**Workload Prioritization**](./use-cases/adaptive-service-protection/workload-prioritization.md):
  Safeguard crucial user experience pathways and ensure prioritized access to
  external APIs even during high-load conditions by strategically prioritizing
  workloads. This is achieved through the use of declarative policies that label
  and prioritize workload requests, such as API calls. By employing
  [weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
  for scheduling, Aperture ensures a fair distribution of resources that aligns
  with the business value and urgency of requests.
- 🔀
  [**Load-based Auto Scaling**](./use-cases/auto-scaling/load-based-auto-scaling.md):
  Eliminate the need for costly over-provisioning and enhance efficiency with
  Aperture's load-based auto-scaling. Aperture's policies are expressed as
  circuit graphs that continuously track deviations from service-level
  objectives and calculate recovery or escalation actions. Auto-scaling can be
  implemented as an escalation that triggers based on load throttling signal.
- ⏱️
  [**Distributed Rate-Limiting**](./use-cases/rate-limiting/rate-limiting.md):
  Safeguard APIs from potential abuse with Aperture's high-performance,
  distributed rate limiter. This feature enforces per-key limits based on
  fine-grained labels, ensuring precise control and prevention of excessive
  usage.
- 🚀
  [**Percentage Rollouts**](./use-cases/percentage-rollouts/percentage-rollouts.md):
  Enable teams to gradually release new features to a subset of users, without
  impacting the rest of the system. Aperture provides automated load ramping
  functionality, allowing for a safe and controlled increment of load to new
  features or API endpoints. This feature continuously monitors for potential
  performance issues and includes an automatic response mechanism to dial back
  load in case of a performance regression. This proactive approach minimizes
  service disruptions and maintains consistent performance, even when rolling
  out new features.

## 🛠️ How it works {#how-it-works}

Load management, at its core, consists of a control loop that observes,
analyzes, and actuates workloads to ensure the stability and reliability of
cloud-native applications.

This control loop is pivotal in both flow control and auto-scaling use cases. In
flow control, the loop manages workloads to maintain the system within its
capacity. In auto-scaling scenarios, the control loop adjusts resource
allocation in response to demand and performance fluctuations.

During the observation phase, an in-built telemetry system continuously monitors
service performance and request attributes, allowing the Agent and Controller to
make informed decisions about request handling and workload prioritization.

The analysis and actuation phases use
[**Declarative policies**](./concepts/advanced/policy.md) that facilitate teams
in defining responses to different situations, such as deviations from
service-level objectives.

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)

## ✨ Get started {#get-started}

- [**Setting up your application**](/get-started/setting-up-application/setting-up-application.md)
- [**Install Aperture**](/get-started/installation/installation.md)
- [**Your first policy**](/get-started/policies/policies.md)
- [**Use cases**](/use-cases/use-cases.md)

For an in-depth understanding of how Aperture interacts with applications and
its various integral components, explore the
[Architecture](/architecture/architecture.md) section.

## 📖 Learn {#learn}

The [Concepts section](/concepts/concepts.md) provides detailed insights into
essential elements of Aperture's system and policies, offering a comprehensive
understanding of their key components.

## Additional Support

Navigating Aperture's capabilities might bring up questions, and we understand
that. Don't hesitate to engage with us for any queries or clarifications. We are
here to assist and ensure that your experience with Aperture is smooth and
beneficial.

<!-- vale off -->

[**💬 Consult with an expert**](https://calendly.com/desaijai/fluxninja-meeting)
|
[**👥 Join our Slack Community**](https://join.slack.com/t/fluxninja-aperture/shared_invite/zt-1vm2t2yjb-AG8rzKkB5TpPmqihJB6YYw)
| ✉️ Email: [**support@fluxninja.com**](mailto:support@fluxninja.com)

<!-- vale on -->
