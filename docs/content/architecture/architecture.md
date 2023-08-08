---
title: Architecture
sidebar_position: 2
description:
  Discover the core components of Aperture architecture and learn how they work
  together to provide powerful and efficient load management.
image: ../assets/img/aperture_logo.png
keywords:
  - reliability
  - overload
  - concurrency
  - aperture
  - fluxninja
  - microservices
  - cloud
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

The diagram below shows the core components of Aperture architecture and how
they integrate with an application.

![Aperture Architecture (dark)](../assets/img/aperture-architecture-dark.svg#gh-dark-mode-only)
![Aperture Architecture (light)](../assets/img/aperture-architecture-light.svg#gh-light-mode-only)

Aperture Cloud is a highly available, fully managed load management platform
offering:

1. Hosted **Aperture Controller**
2. Consoles for managing Aperture policies, **Aperture Agents** and self-hosted
   Aperture Controllers
3. Traffic analytics dashboard
4. Alerting system to notify about actions taken by Aperture Agents

## Aperture Controller (hosted in Aperture Cloud) {#aperture-controller}

:::note

Here the Aperture Controller is shown as part of Aperture Cloud, but it's also
possible to [self-host it][self-hosting].

:::

The Aperture Controller is a centralized control system, equipped with a
comprehensive global perspective. Its role is collecting data and evaluating
policies. Policy evaluation results in high-level adjustments, which are then
sent down to Aperture Agents.

Aperture Cloud [provides a per-project Aperture
Controller][aperture-cloud-controller]. It is programmed using declarative
policies. Policies can be applied by configuring a [pre-defined
blueprint][use-cases]. It's also possible to build a policy [from scratch from
policy components][policy].

## Aperture Agents

Serving as the workhorses of the platform, Aperture Agents provide powerful flow
control components. These include a weighted fair queuing scheduler for workload
prioritization and a distributed rate-limiter for abuse prevention. These agents
are deployed adjacent to services requiring load management and control traffic
flows based on real-time adjustments from the Aperture Controller. They
seamlessly [integrate][integrations] with service meshes, gateways, and HTTP
middlewares. For more specific control, developers can use [SDKs][sdks] to
manage specific features or code sections within services.

The Agents monitor service and infrastructure health signals using an in-built
telemetry system. In addition, a programmable, high-fidelity flow classifier is
used to label requests based on attributes such as customer tier or request
type. These metrics are then analyzed by the Aperture Controller.

Aperture Agents schedule workloads based on their priorities, helping prioritize
critical features over less important workloads during overload scenarios. For
example, a video streaming service might prioritize a request to play a movie by
a customer over a recommended movies API. A SaaS product might prioritize
features used by paid users over those being used by free users.

Aperture Agents can be [installed on a variety of
infrastructure][install-agents] such as Kubernetes, VMs, or bare-metal. In
addition to flow control capabilities, Agents work with auto-scaling APIs for
platforms such as Kubernetes, to help scale infrastructure when needed.

### Metrics

Aperture Agents use metrics to provide input signals to policies in the Aperture
Controller. These metrics can either be defined based on existing traffic using
[Flux Meters](/concepts/flux-meter.md) or using [any OpenTelemetry Collector
receiver][metrics]. These metrics can then be used in policies using [PromQL
syntax][promql-syntax].

:::info

For more details about the interaction between Aperture Controller and Agents
and the exact databases, see [Architecture of Self-Hosted
Aperture][architecture-self-hosted].

:::

[aperture-cloud-controller]: /reference/fluxninja.md#cloud-controller
[architecture-self-hosted]: /self-hosting/architecture.md
[use-cases]: /use-cases/use-cases.md
[policy]: /concepts/advanced/policy.md
[integrations]: /integrations/integrations.md
[sdks]: /integrations/sdk/sdk.md
[metrics]: /integrations/metrics/metrics.md
[install-agents]: /get-started/installation/agent/agent.md
[self-hosting]: /self-hosting/self-hosting.md
[promql-syntax]: https://prometheus.io/docs/prometheus/latest/querying/basics/
