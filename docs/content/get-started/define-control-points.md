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
**`Aperture`** tab in the sidebar menu. From there, select **`Agent API Keys`**
in the top bar. This is where you can find and copy the **`AGENT_API_KEY`** and
use it in the SDK as follows:

:::

```mdx-code-block
<Tabs>
<TabItem value="TypeScript">
```

```typescript
import { ApertureClient, FlowStatusEnum } from "@fluxninja/aperture-js";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  agentAPIKey: "AGENT_API_KEY",
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
showcases how to define policy in Aperture. The code snippet below shows how to
wrap your [Control Point](/concepts/control-point.md) within the `StartFlow`
call and pass [labels](/concepts/flow-label.md) that will be matched with
policy. The function `Flow.ShouldRun()` checks if the flow allows the request.
The `Flow.End()` function is responsible for sending telemetry, and updating the
specified cache entry within Aperture.

Let's create a feature control point in the following code snippet.

```mdx-code-block
<Tabs>
<TabItem value="TypeScript">
```

```typescript
async function handleFlow() {
  const flow = await apertureClient.StartFlow("archimedes-service", {
    labels: {
      label_key: "api_key",
      interval: "60",
    },
    grpcCallOptions: {
      deadline: Date.now() + 1200000, // 20 minutes deadline
    },
  });

  if (flow.ShouldRun()) {
    // Do Actual Work
  } else {
    // Handle flow rejection
    flow.SetStatus(FlowStatusEnum.Error);
  }

  if (flow) {
    flow.End();
  }
}
```

```mdx-code-block
</TabItem>
</Tabs>
```

This is how you can create a feature control point in your code. The complete
example is available
[here](https://github.com/fluxninja/aperture-js/blob/main/example/routes/use_aperture.ts).

:::info

Aperture SDKs are available for multiple languages, you choose the one that fits
your needs. [See all SDKs][sdks].

:::

<!-- vale off -->

[control-points]: /concepts/control-point.md
[sdks]: /sdk/sdk.md
