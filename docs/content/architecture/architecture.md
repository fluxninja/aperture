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
  # TODO
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
```

The diagram below shows the interaction between the main components of
[Aperture][]-powered [FluxNinja][] platform: FluxNinja Cloud, Aperture Agents,
and various integrations.

![FluxNinja Architecture](../assets/img/FluxNinja-arc-dark.svg#gh-dark-mode-only)
![FluxNinja Architecture](../assets/img/FluxNinja-arc-light.svg#gh-light-mode-only)

FluxNinja Cloud is the brain of the system and its role is collecting data and
evaluating policies. Policy evaluation is performed by Aperture Controller and
results in high-level decisions, which are then sent down to Aperture Agents.

Aperture Agents are part of the system that's much closer to the infrastructure
– they're installed on every node. The Agents are where the actual execution of
policies takes place. Note that while the Agents by themselves are able to
collect some metrics and perform limited actions like auto-scaling, they need
[integrations][] to actually control the traffic.

Aperture provides [integrations][] for service meshes and gateways. It's also
possible to instrument your application directly with [Aperture SDKs][]. When
integration is enabled, it will ask the Agent on the local node to make a
decision for every request or flow. Note that this RPC call never leaves the
node, so its overhead and impact on latency are minimized.

## FluxNinja Cloud

FluxNinja Cloud is a centralized platform that provides tools for policy
management and observability. There are two significant components of FluxNinja
Cloud worth mentioning: The analytics database and the Aperture Controller.

### Analytics database

FluxNinja uses a real-time analytics database to support FluxNinja observability
capabilities. All the logs and traces collected by Aperture Agents are batched
and rolled up and sent to FluxNinja. Thanks to the use of rollup, similar events
are aggregated to reduce the traffic, but no data is lost (as it would with
usage of sampling-based solutions).

### Aperture Controller

FluxNinja Cloud [provides a per-project Aperture
Controller][FluxNinja Cloud Controller] for every organization.

The Aperture Controller is a centralized control system, equipped with a
comprehensive global perspective. It is programmed using declarative policies.
Policies can be applied by configuring a [pre-defined blueprint][Use Cases].
It's also possible to build a policy [from scratch from policy
components][Policy].

A policy represents a closed-loop control circuit that is executed periodically.
The control circuit draws input signals from [metrics](#metrics) aggregated
across Aperture Agents, providing the Controller with a holistic view of the
application's health and performance. Service-level objectives (SLOs) are
defined against these health and performance signals. The policies continuously
track deviations from SLOs and calculate recovery or escalation actions that are
translated as adjustments to the Agents.

After computing the adjustments, the Aperture Controller synchronizes them with
the relevant Aperture Agents. These adjustments encompass load throttling,
workload prioritization, and auto-scaling actions, among others. By
disseminating the calculated adjustments to the Agents, the Controller ensures
that the Agents take localized actions in line with the global state of the
system.

:::note

Here the Aperture Controller is shown as part of FluxNinja Cloud Platform, but
it's also possible to [self-host it][Self-Hosting].

:::

## Aperture Agents

Aperture Agents are the workhorses of the platform, providing powerful flow
control components such as a weighted fair queuing scheduler for workload
prioritization and a distributed rate-limiter for abuse prevention. Agents
integrate with service meshes, gateways and HTTP middlewares. Alternately,
developers can use SDKs to get flow control around specific features or code
sections inside services.

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
infrastructure][Install Agents] such as Kubernetes, VMs, or bare-metal. In
addition to flow control capabilities, Agents work with auto-scaling APIs for
platforms such as Kubernetes, to help scale infrastructure when needed.

### Metrics

Aperture Agents use metrics to provide input signals to policies in the Aperture
Controller. These metrics can either be defined based on existing traffic using
[Flux Meters](/concepts/flux-meter.md) or using [any OpenTelemetry Collector
receiver][Metrics]. These metrics can then be used in policies using [PromQL
syntax][].

:::info

For more details about the interaction between Aperture Controller and Agents
and the exact databases, see [Architecture of Self-Hosted Aperture][].

:::

[FluxNinja]: /introduction.md
[Aperture]: https://github.com/fluxninja/aperture
[FluxNinja Cloud Controller]: /reference/fluxninja.md#cloud-controller
[Architecture of Self-Hosted Aperture]: /self-hosting/architecture.md
[Use Cases]: /use-cases/use-cases.md
[Policy]: /concepts/advanced/policy.md
[integrations]: /integrations/integrations.md
[Aperture SDKs]: /integrations/sdk/sdk.md
[Metrics]: /integrations/metrics/metrics.md
[Install Agents]: /get-started/installation/agent/agent.md
[Self-Hosting]: /self-hosting/self-hosting.md
[PromQL syntax]: https://prometheus.io/docs/prometheus/latest/querying/basics/
