---
title: Define Control Points
keywords:
  - Define Control Points
sidebar_position: 3
sidebar_label: Define Control Points
---

```mdx-code-block
import { Cards } from '@site/src/components/Cards';
import Tabs from '@theme/Tabs';
import TabItem from "@theme/TabItem";
```

[Control points][control-points] are used to define where you want to act in
code or at service level. It's important to understand what control points are
because you will be using them many times in your code.

<!-- vale off -->

## What is a Feature Control point?

<!-- vale on -->

A feature control point is essentially a specific point in the codebase where
the execution flow can be controlled using feature flags. Feature flags, also
known as feature toggles, are a programming technique that allows developers to
enable or disable features of their software even after the code has been
deployed to production. This can be extremely useful for testing new features,
performing Blue Green testing, or quickly disabling a feature in response to an
issue or error.

<!-- vale off -->

## How to create a Feature Control point?

<!-- vale on -->

Let's create a feature control point in code. To begin with, you need to
configure the Aperture SDK for your application.

:::note Required

In the organization created within Aperture Cloud, navigate to the
**`Aperture`** tab in the sidebar menu. From there, select
**`Aperture Agent Keys`** in the top bar. This is where you can find and copy
the **`AGENT_API_KEY`** and use it in the SDK as follows:

:::

```mdx-code-block
<Tabs>
<TabItem value="TypeScript">
```

```typescript
import { ApertureClient, FlowStatusEnum } from "@fluxninja/aperture-js";
import grpc from "@grpc/grpc-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address:
    process.env.APERTURE_AGENT_ADDRESS !== undefined
      ? process.env.APERTURE_AGENT_ADDRESS
      : "localhost:8089",
  agentAPIKey: process.env.APERTURE_AGENT_API_KEY || undefined,
  // if process.env.APERTURE_AGENT_INSECURE set channelCredentials to insecure
  channelCredentials:
    process.env.APERTURE_AGENT_INSECURE !== undefined
      ? grpc.credentials.createInsecure()
      : grpc.credentials.createSsl(),
});
```

```mdx-code-block
</TabItem>
</Tabs>
```

Once you have configured Aperture SDK, you can create a feature control point
wherever you want in your code. Before executing the business logic of a
specific API, you can create a feature control point that can control the
execution flow of the API and can reject the request based on the policy defined
in Aperture. The [Create Your First Policy](./policies/policies.md) section
showcases how to define policy in Aperture. For now, some labels have been added
in the code snippet below. These labels will be used while defining a policy.

Let's create a feature control point in the following code snippet.

```mdx-code-block
<Tabs>
<TabItem value="TypeScript">
```

```typescript
let flow: Flow | undefined;

if (apertureClient) {
  try {
    // Start the flow to check rate limiting for the incoming request
    flow = await apertureClient.StartFlow("archimedes-service", {
      labels: {
        label_key: "api_key",
        interval: "60",
      },
      grpcCallOptions: {
        deadline: Date.now() + 1200000, // 20 minutes deadline
      },
    });

    // Check if the flow is allowed by Aperture
    if (flow.ShouldRun()) {
      // Add business logic to process incoming request
      console.log("Request accepted. Processing...");
    } else {
      console.log("Request rate-limited. Try again later.");
    }
  } catch (e) {
    console.error("Error in flow:", e);
    if (flow) {
      flow.SetStatus(FlowStatusEnum.Error);
    }
  } finally {
    if (flow) {
      flow.End();
    }
  }
}
```

```mdx-code-block
</TabItem>
</Tabs>
```

This is how you can create a feature control point in your code. The complete
example is available
[here](https://github.com/fluxninja/aperture-js/tree/main/example).

:::info

Aperture SDKs are available for multiple languages, you choose the one that fits
your needs. [See all SDKs][sdks].

:::

<!-- vale off -->

[control-points]: /concepts/control-point.md
[sdks]: /sdk/sdk.md
