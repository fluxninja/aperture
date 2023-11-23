---
title: Getting Started with Aperture
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

## Sign up to Aperture Cloud

For signing up, head over to
[Aperture Cloud](https://app.fluxninja.com/sign-up). For detailed instructions,
please refer to our [step-by-step][sign-up] guide.

These are two main modes on how to get started with Aperture.

<Tabs>

<TabItem value="Aperture Serverless">

![Aperture Serverless Architecture](./assets/architecture/saas-dark.svg#gh-dark-mode-only)
![Aperture Serverless Architecture](./assets/architecture/saas-light.svg#gh-light-mode-only)

1. **Connect to Aperture Cloud**: In Aperture Cloud, authentication for SDK
   integrations is handled using API keys. These keys can be found in the
   Aperture Cloud UI. For detailed instructions on locating API Keys, please
   refer to the [API Keys][api-keys] section.

2. **Define Control Points**: Check out the [Define Control
   Points][define-control-points] and learn how to use Aperture within your
   code. Explore our [SDKs][sdks] to pick the one that aligns with your
   requirements.

3. **Create Your Policy**: Deploy your first policy. See [Policies][policies]
   for more details.

</TabItem>

<TabItem value="Aperture for Infrastructure">

There are two ways to integrate Aperture into your infrastructure. The first
option includes a cloud-hosted controller and self-managed agents to ensure
minimal performance impact. The second option involves hosting Aperture entirely
within your infrastructure, making it perfect for air-gapped environments.

1. **Set Up the Environment**: The Aperture Agent can be installed in various
   modes. For installation steps, see [Agent][agent-docs] docs under [Aperture
   For Infra section][aperture-for-infra].

   :::info

   For more details on fully self-hosted installation, please refer to the
   [Self-hosted][aperture-for-infra] section.

   :::

2. **Integrate with Aperture**: Here are various [Integrations][integrations]
   methods with Aperture

   - [SDKs](../sdk/sdk.md)
   - [Istio](/aperture-for-infra/integrations/istio/istio.md)
   - [Gateways](/aperture-for-infra/integrations/gateway/gateway.md)
   - [Consul](/aperture-for-infra/integrations/consul/consul.md)
   - [Auto Scale](/aperture-for-infra/integrations/auto-scale/auto-scale.md)
   - [Metrics](/aperture-for-infra/integrations/metrics/metrics.md)

3. **Map to Aperture SaaS Controller**: Aperture Cloud authenticates requests
   from integrations using API keys, which are created for your project and can
   be found within the Aperture Cloud UI. Copy the API key and save it in a
   secure location. This key will be used during the configuration of
   [Self-hosted][aperture-for-infra] Agents. For detailed instructions on
   locating API Keys, please refer to the [API Keys][api-keys] section.

   :::info

   Using the API key, you can map your integration to the Aperture Cloud. See
   [FluxNinja Cloud Extension][cloud-extension] for more details.

   :::

4. **Create Your Policy**: Deploy your first policy. See [Policies][policies]
   for more details.

</TabItem>

</Tabs>

[cloud]: https://www.fluxninja.com/product
[aperture-for-infra]: /aperture-for-infra/aperture-for-infra.md
[sign-up]: /reference/cloud-ui/sign-up.md
[policies]: /get-started/policies/policies.md
[cloud-extension]: /reference/fluxninja.md
[agent-docs]: /aperture-for-infra/agent/agent.md
[integrations]: /aperture-for-infra/integrations/integrations.md
[sdks]: /sdk/sdk.md
[api-keys]: /reference/cloud-ui/api-keys.md
[define-control-points]: /get-started/define-control-points.md
