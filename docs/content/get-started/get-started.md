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
```

Aperture is available as a managed service, [Aperture Cloud][cloud], or can be
self-hosted within your infrastructure.

<!-- markdownlint-disable -->

## Step 1. Sign up to Aperture Cloud

For signing up, head over to
[Aperture Cloud](https://app.fluxninja.com/sign-up). For detailed instructions,
please refer to our [step-by-step][sign-up] guide.

## Step 2. Select Your Mode

These are two main modes on how to get started with Aperture.

<Tabs>

<TabItem value="Aperture Serverless">

The Serverless mode is the quickest way to get started with using Aperture SDKs.

![Aperture Serverless Architecture](./assets/architecture/saas-dark.svg#gh-dark-mode-only)
![Aperture Serverless Architecture](./assets/architecture/saas-light.svg#gh-light-mode-only)

## Step 3. Get your API key

In Aperture Cloud, authentication for SDK integrations is handled using API
keys. These keys can be found in the Aperture Cloud UI. For detailed
instructions on locating API Keys, please refer to the [API Keys][api-keys]
section.

## Step 4. Grab the SDK and define Control Points

**Install the SDK**

Aperture provides SDKs for various languages. You can find the SDKs and
installation instructions [here][sdks].

Install the Aperture SDK in your project directory

```bash
npm i @fluxninja/aperture-js
```

Configure the Aperture SDK for in your application with the API key and the
address of the Aperture Cloud.

```typescript
import { ApertureClient, FlowStatusEnum } from "@fluxninja/aperture-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  agentAPIKey: "API_KEY",
});
```

**Define a Control Point**

Once you have configured Aperture SDK, you can create a feature control point
anywhere within your code. Before executing the business logic of a specific
API, you can create a feature control point that can control the execution flow
of the API and can reject the request based on the policy defined in Aperture.

The code snippet below shows how to wrap your
[Control Point](/concepts/control-point.md) within the `StartFlow` call and pass
[labels](/concepts/flow-label.md) and `cacheKey` to Aperture Agents.

- The function `Flow.ShouldRun()` checks if the flow allows the request.
- The `Flow.End()` function is responsible for sending telemetry, and updating
  the specified cache entry within Aperture.
- The `flow.CachedValue().GetLookupStatus()` function returns the status of the
  cache lookup. The status can be `Hit` or `Miss`.
- If the status is `Hit`, the `flow.CachedValue().GetValue()` function returns
  the cached response.
- The `flow.SetCachedValue()` function is responsible for setting the response
  in Aperture cache with the specified TTL (time to live).

Let's create a feature control point in the following code snippet.

```typescript
async function handleRequest(req, res) {
  const flow = await apertureClient.StartFlow("archimedes-service", {
    labels: {
      api_key: "some_api_key",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300, // ms
    },
    rampMode: false,
    cacheKey: "cache",
  });

  if (flow.ShouldRun()) {
    // Check if the response is cached in Aperture from a previous request
    if (flow.CachedValue().GetLookupStatus() === LookupStatus.Hit) {
      res.send({ message: flow.CachedValue().GetValue()?.toString() });
    } else {
      // Do Actual Work
      // After completing the work, you can return store the response in cache and return it, for example:
      const resString = "foo";

      // create a new buffer
      const buffer = Buffer.from(resString);

      // set cache value
      const setResult = await flow.SetCachedValue(buffer, {
        seconds: 30,
        nanos: 0,
      });
      if (setResult?.error) {
        console.log(`Error setting cache value: ${setResult.error}`);
      }

      res.send({ message: resString });
    }
  } else {
    // Handle flow rejection
    flow.SetStatus(FlowStatusEnum.Error);
  }

  if (flow) {
    flow.End();
  }
}
```

This is how you can create a feature control point in your code. The complete
example is available
[here](https://github.com/fluxninja/aperture-js/blob/main/example/routes/use_aperture.ts).

## Step 5. Create Your Policy

Deploy your first policy. See [Policies][policies] for more details.

</TabItem>
<TabItem value="Aperture for Infrastructure">

There are two ways to integrate Aperture into your infrastructure. The first
option includes a cloud-hosted controller and self-managed agents to ensure
minimal performance impact. The second option involves hosting Aperture entirely
within your infrastructure, making it perfect for air-gapped environments.

## Step 3. Set Up the Environment

The Aperture Agent can be installed in various modes. For installation steps,
see [Agent][agent-docs] docs under [Aperture For Infra
section][aperture-for-infra].

:::info

For more details on fully self-hosted installation, please refer to the
[Self-hosted][aperture-for-infra] section.

:::

## Step 4. Integrate with Aperture

Here are various [Integrations][integrations] methods with Aperture

- [SDKs](../sdk/sdk.md)
- [Istio](/aperture-for-infra/integrations/istio/istio.md)
- [Gateways](/aperture-for-infra/integrations/gateway/gateway.md)
- [Consul](/aperture-for-infra/integrations/consul/consul.md)
- [Auto Scale](/aperture-for-infra/integrations/auto-scale/auto-scale.md)
- [Metrics](/aperture-for-infra/integrations/metrics/metrics.md)

## Step 5. Map to Aperture SaaS Controller

Aperture Cloud authenticates requests from integrations using API keys, which
are created for your project and can be found within the Aperture Cloud UI. Copy
the API key and save it in a secure location. This key will be used during the
configuration of [Self-hosted][aperture-for-infra] Agents. For detailed
instructions on locating API Keys, please refer to the [API Keys][api-keys]
section.

:::info

Using the API key, you can map your integration to the Aperture Cloud. See
[FluxNinja Cloud Extension][cloud-extension] for more details.

:::

## Step 6. Create Your Policy

Deploy your first policy. See [Policies][policies] for more details.

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
