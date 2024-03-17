---
title: Getting Started with Aperture
keywords:
  - setup
  - getting started
sidebar_position: 2
sidebar_label: Get Started
hide_table_of_contents: true
---

```mdx-code-block
import { Cards } from '@site/src/components/Cards';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import CodeSnippet from '../codeSnippet.js'

```

:::info Sign up

FluxNinja has been acquired by [CodeRabbit](https://www.coderabbit.ai). New
sign-ups are temporarily disabled. Existing users can continue to use the
service as usual.

:::

Aperture is available as a managed service, [Aperture Cloud][cloud], or can be
self-hosted within your infrastructure.

<!-- markdownlint-disable -->

## 1. Sign up to Aperture Cloud

For signing up, head over to
[Aperture Cloud](https://app.fluxninja.com/sign-up). For detailed instructions,
refer to our [step-by-step][sign-up] guide.

## 2. Pick your Integration

These are two main modes on how to get started with Aperture.

<Tabs>

<TabItem value="Aperture Serverless">

The Serverless mode is the quickest way to start using Aperture.

![Aperture Serverless Architecture](./assets/architecture/saas-dark.svg#gh-dark-mode-only)
![Aperture Serverless Architecture](./assets/architecture/saas-light.svg#gh-light-mode-only)

## 3. Get your API key

In Aperture Cloud, authentication for SDK integrations is handled using API
keys. These keys can be found in the Aperture Cloud UI. For detailed
instructions on locating API Keys, refer to the [API Keys][api-keys] section.

## 4. Install SDK and Define Control Points

Aperture provides SDKs for various languages. You can find the SDKs and
installation instructions [here][sdks].

```mdx-code-block
<Tabs groupId="language" queryString>
<TabItem value="TypeScript">
```

```bash
npm install @fluxninja/aperture-js
```

```mdx-code-block
</TabItem>
<TabItem value="Python">
```

```bash
pip install aperture-py
```

```mdx-code-block
</TabItem>
<TabItem value="Go">
```

```bash
go get github.com/fluxninja/aperture-go/v2
```

```mdx-code-block
</TabItem>
</Tabs>
```

Connect to Aperture Cloud by creating an Aperture Client with your
organization's address and API Key.

```mdx-code-block
<Tabs groupId="language" queryString>
<TabItem value="TypeScript">
```

<CodeSnippet lang="ts" snippetName="clientConstructor" />

```mdx-code-block
</TabItem>
<TabItem value="Python">
```

<CodeSnippet lang="py" snippetName="clientConstructor" />

```mdx-code-block
</TabItem>
<TabItem value="Go">
```

<CodeSnippet lang="go" snippetName="grpcOptions" />

<CodeSnippet lang="go" snippetName="clientConstructor" />

```mdx-code-block
</TabItem>
</Tabs>
```

**Define a Control Point**

Once the SDK is connected to Aperture Cloud, you can create a feature control
point anywhere within your code. For example, before executing the business
logic of a specific API, you can create a feature control point that can control
the execution flow of the API and can reject the request based on the policy
defined in Aperture.

Let's create a feature control point in the following code snippet.

```mdx-code-block
<Tabs groupId="language" queryString>
<TabItem value="TypeScript">
```

The code snippet below shows how to wrap your
[Control Point](/concepts/control-point.md) within the `StartFlow` call while
also passing [labels](/concepts/flow-label.md) to Aperture Agents.

- The function `Flow.ShouldRun()` checks if the flow allows the request.
- The `Flow.End()` function is responsible for sending telemetry.

<CodeSnippet
    lang="ts"
    snippetName="handleRequestRateLimit"
 />

```mdx-code-block
</TabItem>
<TabItem value="Python">
```

The code snippet below shows how to wrap your
[Control Point](/concepts/control-point.md) within the `start_flow` call while
passing [labels](/concepts/flow-label.md) to Aperture.

<CodeSnippet
    lang="py"
    snippetName="syncFlow"
 />

```mdx-code-block
</TabItem>
<TabItem value="Go">
```

<CodeSnippet
    lang="go"
    snippetName="manualFlowNoCaching"
 />

```mdx-code-block
</TabItem>
</Tabs>
```

This is how you can create a feature control point in your code. The complete
example is available
[here](https://github.com/fluxninja/aperture-js/blob/main/example/routes/use_aperture.ts).

## 5. Create Your Policy

Within the Aperture UI, navigate to the policy in tab in the sidebar menu, and
click the `Create Policy` button in the upper-right corner. There you can pick
the blueprint that best aligns with your needs. After a few clicks, you'll be
directed to a screen where you can input the necessary parameters to generate
the policy.

![Rate Limiter Blueprint](./assets/rate-limiter-blueprint-dark.png#gh-dark-mode-only)
![Rate Limiter Blueprint](./assets/rate-limiter-blueprint-light.png#gh-light-mode-only)
![Rate Limiter Dashboard](./assets/rate-limiter-dashboard-dark.png#gh-dark-mode-only)
![Rate Limiter Dashboard](./assets/rate-limiter-dashboard-light.png#gh-light-mode-only)

For more details on how to create a policy, follow our [step-by-step
guide][policies].

</TabItem>
<TabItem value="Aperture for Infrastructure">

There are two ways to integrate Aperture into your infrastructure. The first
option includes a cloud-hosted controller and self-managed agents to ensure
minimal performance impact. The second option involves hosting Aperture entirely
within your infrastructure, making it perfect for air-gapped environments.

## 3. Set Up the Environment

The Aperture Agent can be installed in various modes. For installation steps,
see [Agent][agent-docs] docs under [Aperture For
Infrastructure][aperture-for-infra].

:::info

For more details on fully self-hosted installation, refer to the
[Self-hosted][aperture-for-infra] section.

:::

## 4. Integrate with Aperture

Here are various [Integrations][integrations] methods with Aperture

- [SDKs](../sdk/sdk.md)
- [Istio](/aperture-for-infra/integrations/istio/istio.md)
- [Gateways](/aperture-for-infra/integrations/gateway/gateway.md)
- [Consul](/aperture-for-infra/integrations/consul/consul.md)
- [Auto Scale](/aperture-for-infra/integrations/auto-scale/auto-scale.md)
- [Metrics](/aperture-for-infra/integrations/metrics/metrics.md)

## 5. Map to Aperture SaaS Controller

Aperture Cloud authenticates requests from integrations using API keys, which
are created for your project and can be found within the Aperture Cloud UI. Copy
the API key and save it in a secure location. This key will be used during the
configuration of [Self-hosted][aperture-for-infra] Agents. For detailed
instructions on locating API Keys, refer to the [API Keys][api-keys] section.

:::info

Using the API key, you can map your integration to the Aperture Cloud. See
[FluxNinja Cloud Extension][cloud-extension] for more details.

:::

## 6. Create Your Policy

For more details on how to create a policy, follow our [step-by-step
guide][policies].

</TabItem>
</Tabs>

[cloud]: https://www.fluxninja.com
[aperture-for-infra]: /aperture-for-infra/aperture-for-infra.md
[sign-up]: /reference/cloud-ui/sign-up.md
[policies]: /reference/cloud-ui/policies/policy-creation.md
[cloud-extension]: /reference/fluxninja.md
[agent-docs]: /aperture-for-infra/agent/agent.md
[integrations]: /aperture-for-infra/integrations/integrations.md
[sdks]: /sdk/sdk.md
[api-keys]: /reference/cloud-ui/api-keys.md
