---
title: Introduction
slug: /
sidebar_position: 1
sidebar_label: Introduction
sidebar_class_name: introduction
keywords:
  - cloud
  - enterprise
  - platform
  - fluxninja
  - aperture
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

[Aperture](https://github.com/fluxninja/aperture) is an open source load
management platform designed for classifying, rate limiting, queuing and
prioritizing API traffic in cloud applications. Built upon a foundation of
observability and a global control plane, it offers a comprehensive suite of
load management capabilities. These capabilities enhance the reliability and
performance of cloud applications while also optimizing resource utilization.

Aperture can seamlessly integrate with existing control points such as gateways,
service meshes, and application middlewares. Moreover, it offers SDKs for
developers who need to establish control points around specific features or code
sections inside applications.

Aperture's control plane is available as a managed service, [Aperture
Cloud][cloud], or can be [self-hosted][self-hosted] within your infrastructure.

Here's a simplified diagram of how Aperture Cloud (managed by FluxNinja)
interacts with your infrastructure. Visit the [Architecture][architecture] page
for more details.

![Aperture Architecture (dark)](./assets/img/aperture-architecture-dark.svg#gh-dark-mode-only)
![Aperture Architecture (light)](./assets/img/aperture-architecture-light.svg#gh-light-mode-only)

:::info Sign-up

To sign-up to Aperture Cloud, [click here][sign-up].

:::

## ⚙️ Load management capabilities {#load-management-capabilities}

- ⏱️ [**Global Rate-Limiting**](concepts/rate-limiter.md): Safeguard APIs and
  features against excessive usage with Aperture's high-performance, distributed
  rate limiter. Identify individual users or entities by fine-grained labels.
  Create precise rate limiters controlling burst-capacity and fill-rate tailored
  to business-specific labels. Refer to the
  [Rate Limiting](guides/per-user-rate-limiting.md) guide for more details.
- 📊 [**API Quota Management**](concepts/scheduler/quota-scheduler.md): Maintain
  compliance with external API quotas with a global token bucket and smart
  request queuing. This feature regulates requests aimed at external services,
  ensuring that the usage remains within prescribed rate limits and avoids
  penalties or additional costs. Refer to the
  [API Quota Management](guides/api-quota-management/api-quota-management.md)
  guide for more details.
- 🛡️ [**Adaptive Queuing**](concepts/scheduler/load-scheduler.md): Enhance
  resource utilization and safeguard against abrupt service overloads with an
  intelligent queue at the entry point of services. This queue dynamically
  adjusts the rate of requests based on live service health, thereby mitigating
  potential service disruptions and ensuring optimal performance under all load
  conditions. Refer to the
  [Service Load Management](guides/service-load-management/service-load-management.md)
  and
  [Database Load Management](guides/database-load-management/database-load-management.md)
  guides for more details.
- 🎯 [**Workload Prioritization**](concepts/scheduler/scheduler.md): Safeguard
  crucial user experience pathways and ensure prioritized access to external
  APIs by strategically prioritizing workloads. With
  [weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing),
  Aperture aligns resource distribution with business value and urgency of
  requests. Workload prioritization applies to API Quota Management and Adaptive
  Queuing use cases.

## 🛠️ How it works {#how-it-works}

Load management, at its core, consists of a control loop that observes,
analyzes, and actuates workloads to ensure the stability and reliability of
cloud-native applications. This control loop is pivotal in flow control use
cases where it manages workloads to maintain the system within its capacity.

During the observation phase, an in-built telemetry system continuously monitors
service performance and request attributes, allowing the Agent and Controller to
make informed decisions about request handling and workload prioritization.

The analysis and actuation phases use
[**Declarative policies**](concepts/advanced/policy.md) that facilitate teams in
defining responses to different situations, such as deviations from
service-level objectives.

![Aperture Control Loop](./assets/img/oaalight.svg#gh-light-mode-only)
![Aperture Control Loop](./assets/img/oaadark.svg#gh-dark-mode-only)

## ✨ Get started {#get-started}

- [**Setting up your application**](get-started/set-up-application/set-up-application.md)
- [**Install Aperture**](get-started/installation/installation.md)
- [**Your first policy**](get-started/policies/policies.md)
- [**Guides**](guides/guides.md)

For an in-depth understanding of how Aperture interacts with applications and
its various integral components, explore the
[Architecture](architecture/architecture.md) section.

## 📖 Learn {#learn}

The [Concepts](concepts/concepts.md) section provides detailed insights into
essential elements of Aperture's system and policies, offering a comprehensive
understanding of their key components.

## Additional Support

Don't hesitate to engage with us for any queries or clarifications. Our team is
here to assist and ensure that your experience with Aperture is smooth and
beneficial.

<!-- vale off -->

[**💬 Consult with an expert**](https://calendly.com/desaijai/fluxninja-meeting)
|
[**👥 Join our Slack Community**](https://join.slack.com/t/fluxninja-aperture/shared_invite/zt-1vm2t2yjb-AG8rzKkB5TpPmqihJB6YYw)
| ✉️ Email: [**support@fluxninja.com**](mailto:support@fluxninja.com)

<!-- vale on -->

[cloud]: https://www.fluxninja.com/product
[sign-up]: https://app.fluxninja.com/sign-up
[architecture]: /architecture/architecture.md
[self-hosted]: /get-started/self-hosting/self-hosting.md
