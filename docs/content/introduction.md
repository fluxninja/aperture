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

Aperture's cutting-edge features enable teams to effortlessly implement
intelligent load management strategies, ensuring optimal performance and maximal
infrastructure utilization optimal end-user experience, even during service
failures. This versatility and resilience apply across all types of systems,
including monolithic architectures and distributed microservices environments.

## ⚙️ Load management capabilities

Aperture offers a suite of intelligent load management capabilities that are
applicable to a wide range of cloud-native applications. These capabilities
ensure the reliability and stability of applications, and include:

- 🛡️
  [**Adaptive Service Protection**](./applying-policies/service-protection/basic-service-protection.md):
  Enhance resource utilization and safeguard against abrupt service overloads
  with an intelligent queue at the entry point of services. This queue
  dynamically adjusts the rate of requests based on live service health, thereby
  mitigating potential service disruptions and ensuring optimal performance
  under all load conditions.
- 📊
  [**Intelligent Quota Management**](./applying-policies/quota-scheduler/quota-scheduler.md):
  Maintain compliance with external API quotas with a global token bucket and
  smart request queuing. This feature regulates requests aimed at external
  services, ensuring that the usage remains within prescribed rate limits and
  avoids penalties or additional costs.
- 🎯
  [**Workload Prioritization**](./applying-policies/service-protection/workload-prioritization.md):
  Safeguard crucial user experience pathways and ensure prioritized access to
  external APIs even during high-load conditions by strategically prioritizing
  workloads. This is achieved through the use of declarative policies that label
  and prioritize workload requests, such as API calls. By employing
  [weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing)
  for scheduling, Aperture ensures a fair distribution of resources that aligns
  with the business value and urgency of requests.
- 🔀
  [**Load based Auto Scaling**](./applying-policies/auto-scale/auto-scale.md):
  Eliminate the need for costly over-provisioning and enhance efficiency with
  Aperture's load-based auto-scaling. Aperture's policies are expressed as
  circuit graphs that continuously track deviations from service-level
  objectives and calculate recovery or escalation actions. Auto-scaling can be
  implemented as an escalation that triggers based on load throttling signal.
- ⏱️
  [**Distributed Rate-Limiting**](./applying-policies/rate-limiting/rate-limiting.md):
  Safeguard APIs from potential abuse with Aperture's high-performance,
  distributed rate limiter. This feature enforces per-key limits based on
  fine-grained labels, ensuring precise control and prevention of excessive
  usage.
- 🚀
  [**Automated Load Ramping**](./applying-policies/feature-rollout/feature-rollout.md):
  Aperture enables teams to gradually release new features to a subset of users,
  without impacting the rest of the system, using dark launch.

## 🛠️ How it works

Load management, at its core, consists of a control loop that observes,
analyzes, and actuates workloads to ensure the stability and reliability of
cloud-native applications.

🔍 The observation phase incorporates
[**Monitoring and Telemetry**](./reference/observability/observability.md).
Here, an in-built telemetry system continuously monitors service performance and
request attributes, allowing the agent and controller to make informed decisions
about request handling and workload prioritization.

This control loop is pivotal in both flow control and auto-scaling use cases. In
flow control, the loop manages workloads to maintain the system within its
capacity. In auto-scaling scenarios, the control loop adjusts resource
allocation in response to demand and performance fluctuations.

The analysis and actuation phases use 📝
[**Declarative policies**](./concepts/policy/policy.md). This policy language
facilitates teams in defining responses to different situations, such as
deviations from service-level objectives. Expressed as a signal processing
circuit, these policies streamline action from telemetry data to appropriate
actions.

![Aperture Control Loop](assets/img/oaalight.png#gh-light-mode-only)
![Aperture Control Loop](assets/img/oaadark.png#gh-dark-mode-only)

## 📖 Key concepts

For a high-level overview that explains how Aperture policies work, following
are some key components to be aware of:

- [**Control point**](./concepts/flow-control/selector.md): Control Points
  function as essential signposts in your code or data plane, analogous to
  feature flags, that guide flow control decisions. They can be intuitively
  established using SDKs or seamlessly incorporated during API Gateway or
  Service Mesh integrations.
- [**Selector**](./concepts/flow-control/selector.md): Selectors are similar to
  filters, specifying which flows should be considered by the policy components
  for certain operations.
- [**Classifier**](./concepts/flow-control/resources/classifier.md): The
  Classifier can be used to create additional Flow Labels based on request
  metadata without any service alterations, if the existing flow labels aren't
  sufficient.

<!-- vale off -->

Aperture's [Architecture](/architecture/architecture.md) section provides more
information on how Aperture interfaces with applications as well as the
different components that make up Aperture.

<!-- vale off -->

## Additional Support

<!-- vale on -->

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
