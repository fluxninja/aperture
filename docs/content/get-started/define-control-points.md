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

Let's create a feature control point in code. To begin, you first need to
install the Aperture SDK([See all SDKs][sdks]) for your chosen application
language.

### Install Aperture SDK

The following code snippets are in the TypeScript language. To install the SDK,
run the following command:

```mdx-code-block
<Tabs>
<TabItem value="TypeScript">
```

```bash
npm install @fluxninja/aperture-js
```

```mdx-code-block
</TabItem>
</Tabs>
```

### Define Aperture Client

:::note Required

In the organization created within Aperture Cloud, navigate to the
**`Aperture`** tab in the sidebar menu. From there, select **`API Keys`** in the
top bar. This is where you can find and copy the **`API_KEY`** and use it in the
SDK. For detailed instructions on locating API Keys, refer to the [API
Keys][api-keys] section.

:::

```mdx-code-block
import CodeSnippet from '../codeSnippet.js'

<Tabs>
<TabItem value="TypeScript">
```

<CodeSnippet
    lang="ts"
    snippetName="clientConstructor"
 />

```mdx-code-block
</TabItem>
</Tabs>
```

### Create Feature Control Point

Once you have configured Aperture SDK, you can create a feature control point
anywhere within your code. Before executing the business logic of a specific
API, you can create a feature control point that can control the execution flow
of the API and can reject the request based on the policy defined in Aperture.
The [Create Your First Policy](./policies/policies.md) section showcases how to
define policy in Aperture.

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

```mdx-code-block
<Tabs>
<TabItem value="TypeScript">
```

<CodeSnippet
    lang="ts"
    snippetName="handleRequest"
 />

```mdx-code-block
</TabItem>
</Tabs>
```

This is how you can create a feature control point in your code. The complete
example is available
[here](https://github.com/fluxninja/aperture-js/blob/main/example/routes/use_aperture.ts).

<!-- vale off -->

[control-points]: /concepts/control-point.md
[sdks]: /sdk/sdk.md
[api-keys]: /reference/cloud-ui/api-keys.md
