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

In the next few minutes, you should have your application integrated with
Aperture and ready to enforce `rate limiting as a service`.

There are multiple integration methods, allowing you to choose the one that best
suits your application. Below are the three different ways to get started with
Aperture, ranging from Managed SaaS to Self-Hosted. If you are unsure which one
to choose, _we recommend starting with the Managed SaaS option_, as it is the
fastest way to get started.

<Tabs>

<TabItem value="Aperture SaaS - Fully Cloud-Hosted">

This mode is tailored for developers who prefer to directly use our SDKs without
delving into infrastructure components. Allowing developers to focus on their
application and requires minimal level permissions to get started.

1. **Sign up for an account**: Sign up process is simple and takes less than a
   minute. Complete the sign up process [here](/get-started/sign-up.md).
2. **Set Up the Environment**: With Aperture SaaS, you don't need to worry about
   `Aperture Controller and Aperture agent are both hosted within the FluxNinja SaaS platform.`
3. **Integrate the SDK**: Add the necessary few lines of code to your desired
   place in code, where you want aperture to act. See available
   [SDKs](../sdk/sdk.md).
4. **Map to Cloud Agent**: Aperture Cloud uses API keys to authenticate requests
   coming from SDK integrations. You can create API keys for your project in the
   Aperture Cloud UI. See
   [API Keys](/get-started/agent-api-keys/agent-api-keys.md) for more details.

   Once you have the API key, you can use it to initialize the SDK. See
   [SDKs](../sdk/sdk.md) for more details.

5. **Create Your Policy**: Deploy your first policy. See
   [Policies](/get-started/policies/policies.md) for more details.

</TabItem>

<TabItem value="Cloud-Only Controller">

Ideal for teams preferring a cloud-hosted controller and self-managed agents,
this option caters to developers and infrastructure teams overseeing diverse
platforms. It offers the convenience of not self-hosting Prometheus and etcd,
ensuring a lightweight impact on performance from the Aperture agent.

1. **Sign up for an account**: Sign up process is simple and takes less than a
   minute. See the Sign up process [here](/get-started/sign-up.md).

2. **Set Up the Environment**:
   `Aperture Controller is hosted within the FluxNinja SaaS platform and Aperture Agent is installed within your infrastcture.`
   Aperture Agents need to be installed by yourself, checkout
   [Agent Installation](/self-hosting/agent/agent.md) for more details.

3. **Integrate with Aperture**: Apart from the SDKs, you can also integrate with

   - ServiceMesh
   - Gateways

   See [SDKs](/sdk/sdk.md) and
   [Integrations](/self-hosting/integrations/integrations.md) Docs for more
   details.

4. **Map to Aperture SaaS Controller**: Aperture Cloud uses API keys to
   authenticate requests coming from SDK integrations. You can create API keys
   for your project in the Aperture Cloud UI. See
   [API Keys](/get-started/agent-api-keys/agent-api-keys.md) for more details.

   Using the API key, you can map your integration to the Aperture Cloud. See
   [FluxNinja Cloud Extension](/reference/fluxninja.md) for more details.

5. **Create Your Policy**: Deploy your first policy. See
   [Policies](/get-started/policies/policies.md) for more details.

</TabItem>
<TabItem value="Self-Hosted - Open Source">

This mode is tailored for teams deploying applications on a variety of
infrastructure platforms. It's well-suited for those who are
`comfortable hosting their own Prometheus and etcd instances` and aim for
minimal performance impact from the Aperture agent. Those who need complete
control over the Aperture Controller and Aperture Agent, in cases like
air-gapped environments.

1. **Set Up the Environment**:
   `Aperture Controller and Aperture Agent installation within your infrastructure.`
   You will have to managed the Aperture Controller and Aperture Agent yourself.
   Checkout [Self Hosting Section](/self-hosting/self-hosting.md) for more
   details.
2. **Integrate with Aperture**: All the integrations are available -

   - SDKs
   - ServiceMesh
   - Gateways

   See [SDKs](/sdk/sdk.md) and
   [Integrations](/self-hosting/integrations/integrations.md) Docs for more
   details.

3. **Create Your Policy**: Deploy your first policy. See
   [Policies](/get-started/policies/policies.md) for more details.

</TabItem>

</Tabs>
