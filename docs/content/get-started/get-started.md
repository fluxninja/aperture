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

Get stated guide will provide you quick overview on how to integrate Aperture
with your application. Once integrated with Aperture, you are ready to enforce
`rate limiting as a service`.

Below are two different ways to get started with Aperture, ranging from
Serverless to Self-Hosted. If you are unsure which one to choose, **we recommend
starting with the Serverless option**, as it's the fastest way to get started
and involves fewer teams.

<Tabs>

<TabItem value="Aperture for Serverless">

This mode is tailored for developers who prefer to use our SDKs directly without
delving into infrastructure components. It allows developers to focus on their
application and requires minimal permissions to get started.

1. **Sign up for an account**: Get started with Aperture by
   [creating an account](/get-started/sign-up.md), completing the simple
   onboarding - giving you an introduction to dashboard - and inviting your team
   members so you can all collaborate on creating your first policy.
2. **Map to Cloud Agent**: Aperture Agent authenticates requests from SDK
   integrations using Agent API keys, which can be created for your project
   within the Aperture UI; for more information, see
   [Agent API Keys](/get-started/agent-api-keys/agent-api-keys.md) for more
   details.

   Once you have the Agent API key, you can use it to initialize the SDK.

3. **Integrate the SDK**: Add the necessary lines of code to your desired
   application where you want Aperture to take action. Refer to the available
   [SDKs](../sdk/sdk.md) for additional guidance.

4. **Create Your Policy**: Deploy your first policy. See
   [Policies](/get-started/policies/policies.md) for more details.

</TabItem>

<TabItem value="Aperture for Infrastructure">

Aperture offers two modes for infrastructure needs. The first is ideal for teams
favoring a cloud-hosted controller and self-managed agents, eliminating the need
to self-host Prometheus and etcd, and ensuring minimal performance impact. The
second mode, suited for teams comfortable with hosting their own Prometheus and
etcd instances, needs complete control over the Aperture Controller and Agent,
ideal for situations like air-gapped environments.

1. **Sign up for an account**: Get started with Aperture by
   [creating an account](/get-started/sign-up.md), completing the simple
   onboarding - giving you an introduction to the dashboard - and inviting your
   team members so you can all collaborate on creating your first policy.

2. **Set Up the Environment**: The Aperture Agent can be installed in various
   modes. For installation steps, please refer to the following
   [page](/self-hosting/agent/agent.md).

   :::info

   For more details on fully self-hosted solutions, please refer to the
   following [page](/self-hosting/self-hosting.md).

   :::

3. **Integrate with Aperture**: Here are various
   [Integrations](/self-hosting/integrations/integrations.md) methods with
   Aperture

   - [SDKs](../sdk/sdk.md)
   - [Istio](/self-hosting/integrations/istio/istio.md)
   - [Gateways](/self-hosting/integrations/gateway/gateway.md)
   - [Consul](/self-hosting/integrations/consul/consul.md)
   - [Auto Scale](/self-hosting/integrations/auto-scale/auto-scale.md)
   - [Metrics](/self-hosting/integrations/metrics/metrics.md)

4. **Map to Aperture SaaS Controller**: Aperture Agent authenticates requests
   from integrations using Agent API keys, which can be created for your project
   within the Aperture UI; for more information, see
   [Agent API Keys](/get-started/agent-api-keys/agent-api-keys.md) for more
   details.

   Using the API key, you can map your integration to the Aperture Cloud. See
   [FluxNinja Cloud Extension](/reference/fluxninja.md) for more details.

5. **Create Your Policy**: Deploy your first policy. See
   [Policies](/get-started/policies/policies.md) for more details.

</TabItem>

</Tabs>
