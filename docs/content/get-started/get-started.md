---
title: Get Started
keywords:
  - setup
  - getting started
sidebar_position: 2
sidebar_label: Get Started
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
```

Aperture is available as a managed service, [Aperture Cloud][cloud], or can be
self-hosted within your infrastructure.

The easiest way to get started with Aperture is to sign up on Aperture Cloud and
integrate your application using SDKs.

These are two main modes on how to get started with Aperture.

<Tabs>

<TabItem value="Aperture for Serverless">

This mode is tailored for developers who prefer to use Aperture SDKs directly
without delving into infrastructure components. It allows developers to focus on
their application and requires minimal permissions to get started.

![Aperture Serverless Architecture](./assets/architecture/saas-dark.svg#gh-dark-mode-only)
![Aperture Serverless Architecture](./assets/architecture/saas-light.svg#gh-light-mode-only)

1. **Sign up for an account**: Get started with Aperture by [creating an
   account][sign-up], completing the simple onboarding - giving you an
   introduction to the dashboard - and inviting your team members, so you can
   collaborate on load management.
2. **Connect to Aperture Cloud**: Aperture Cloud authenticates requests from SDK
   integrations using Agent API keys, which can be copied within the Aperture
   UI; more information on where to find and use the Agent API Keys is provided
   in the [Define Control Points][define-control-points] section.

   Once you have the Agent API key, you can use it to initialize the SDK.

3. **Integrate the SDK**: Add the necessary lines of code to your desired
   application where you want Aperture to take action. Refer to the available
   [SDKs][sdks] for additional guidance.

4. **Create Your Policy**: Deploy your first policy. See [Policies][policies]
   for more details.

</TabItem>

<TabItem value="Aperture for Infrastructure">

Aperture offers two options for infrastructure needs. The first is ideal for
teams favoring a cloud-hosted controller and self-managed agents, eliminating
the need to self-host Prometheus and etcd, and ensuring minimal performance
impact. The second option, suited for teams comfortable hosting their own
Prometheus and etcd instances, needs complete control over the Aperture
Controller and Agent, ideal for situations like air-gapped environments.

1. **Sign up for an account**: Get started with Aperture by [creating an
   account][sign-up], completing the simple onboarding - giving you an
   introduction to the dashboard - and inviting your team members, so you can
   collaborate on load management.

2. **Set Up the Environment**: The Aperture Agent can be installed in various
   modes. For installation steps, see [Agent][agent-docs] docs under [Aperture
   For Infra section][aperture-for-infra].

   :::info

   For more details on fully self-hosted installation, please refer to the
   [Self-hosted][aperture-for-infra] section.

   :::

3. **Integrate with Aperture**: Here are various [Integrations][integrations]
   methods with Aperture

   - [SDKs](../sdk/sdk.md)
   - [Istio](/aperture-for-infra/integrations/istio/istio.md)
   - [Gateways](/aperture-for-infra/integrations/gateway/gateway.md)
   - [Consul](/aperture-for-infra/integrations/consul/consul.md)
   - [Auto Scale](/aperture-for-infra/integrations/auto-scale/auto-scale.md)
   - [Metrics](/aperture-for-infra/integrations/metrics/metrics.md)

4. **Map to Aperture SaaS Controller**: Aperture Cloud authenticates requests
   from integrations using Agent API keys, which are created for your project
   within the Aperture UI; navigate to the **`Aperture`** tab in the sidebar
   menu and then select **`Aperture Agent Keys`** in the top bar. From there you
   can either copy the existing key or create a new one by clicking on
   **`Create Agent API Key`**.

   Using the API key, you can map your integration to the Aperture Cloud. See
   [FluxNinja Cloud Extension][cloud-extension] for more details.

5. **Create Your Policy**: Deploy your first policy. See [Policies][policies]
   for more details.

</TabItem>

</Tabs>

[cloud]: https://www.fluxninja.com/product
[aperture-for-infra]: /aperture-for-infra/aperture-for-infra.md
[sign-up]: /get-started/sign-up.md
[policies]: /get-started/policies/policies.md
[cloud-extension]: /reference/fluxninja.md
[agent-docs]: /aperture-for-infra/agent/agent.md
[integrations]: /aperture-for-infra/integrations/integrations.md
[sdks]: /sdk/sdk.md
[define-control-points]: /get-started/define-control-points.md
