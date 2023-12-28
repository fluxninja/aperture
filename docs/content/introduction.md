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
sections inside applications. The following diagram depicts the role of Aperture
in a cloud application:

![Unified Load Management (dark)](./assets/img/unified-load-management-dark.svg#gh-dark-mode-only)
![Unified Load Management (light)](./assets/img/unified-load-management-light.svg#gh-light-mode-only)

Aperture is available as a managed service, [Aperture Cloud][cloud], or can be
[self-hosted][self-hosted] within your infrastructure. Visit the
[Architecture][architecture] page for more details.

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
- 📊
  [**API Quota Management**](concepts/request-prioritization/quota-scheduler.md):
  Maintain compliance with external API quotas with a global token bucket and
  smart request queuing. This feature regulates requests aimed at external
  services, ensuring that the usage remains within prescribed rate limits and
  avoids penalties or additional costs. Refer to the
  [API Quota Management](guides/api-quota-management.md) guide for more details.
- 🛡️ [**Adaptive Queuing**](concepts/request-prioritization/load-scheduler.md):
  Enhance resource utilization and safeguard against abrupt service overloads
  with an intelligent queue at the entry point of services. This queue
  dynamically adjusts the rate of requests based on live service health, thereby
  mitigating potential service disruptions and ensuring optimal performance
  under all load conditions. Refer to the
  [Service Load Management](aperture-for-infra/guides/service-load-management/service-load-management.md)
  and
  [Database Load Management](aperture-for-infra/guides/database-load-management/database-load-management.md)
  guides for more details.
- 🎯 [**Workload Prioritization**](concepts/scheduler.md): Safeguard crucial
  user experience pathways and ensure prioritized access to external APIs by
  strategically prioritizing workloads. With
  [weighted fair queuing](https://en.wikipedia.org/wiki/Weighted_fair_queueing),
  Aperture aligns resource distribution with business value and urgency of
  requests. Workload prioritization applies to API Quota Management and Adaptive
  Queuing use cases.
- 💾 **Caching**: Boost application performance and reduce costs by caching
  costly operations, preventing duplicate requests to pay-per-use services, and
  easing the load on constrained services.

## ✨ Get started {#get-started}

- [**Get Started**](get-started/get-started.md)
- [**Guides**](guides/guides.md)

## 📖 Learn {#learn}

The [Concepts](concepts/concepts.md) section provides detailed insights into
essential elements of Aperture's system and policies, offering a comprehensive
understanding of their key components.

## Additional Support

Don't hesitate to engage with us for any queries or clarifications. Our team is
here to assist and ensure that your experience with Aperture is smooth and
beneficial.

<!-- vale off -->

[**💬 Consult with an expert**](https://calendly.com/fluxninja/fluxninja-meeting)
|
[**👥 Join our Slack Community**](https://join.slack.com/t/fluxninja-aperture/shared_invite/zt-1vm2t2yjb-AG8rzKkB5TpPmqihJB6YYw)
| ✉️ Email: [**support@fluxninja.com**](mailto:support@fluxninja.com)

<!-- vale on -->

[cloud]: https://www.fluxninja.com
[sign-up]: https://app.fluxninja.com/sign-up
[architecture]: /aperture-for-infra/architecture.md
[self-hosted]: /aperture-for-infra/aperture-for-infra.md
