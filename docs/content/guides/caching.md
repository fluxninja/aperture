---
title: Caching
sidebar_position: 3
keywords:
  - guides
  - caching
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
import {apertureVersion} from '../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from "@theme/TabItem";
import {BashTab, TabContent} from './blueprintsComponents.js';
import CodeSnippet from '../codeSnippet.js'

```

## Overview

Caching is a critical technology that significantly speeds up application
performance by storing frequently accessed data in memory for quick retrieval.
This approach reduces the need to repeatedly fetch data that might be slow to
process. A common example is web caching, where content from often-visited web
pages is stored and served to users. This method avoids the repeated loading of
the entire webpage, thereby enhancing the user experience and application
performance.

Caching can be integrated into various application layers, for example, the
database layer. By caching frequently requested data, it reduces the load on the
database and enables faster request processing. In the context of generative AI
and Large Language Models (LLMs), where computations can be intensive and
costly, effective caching is essential. It facilitates rapid request processing,
leading to significant cost savings and performance improvements.

<Zoom>

```mermaid
{@include: ./assets/caching/caching.mmd}
```

The diagram above shows how developers using the Aperture SDK can connect to
Aperture Cloud, to set the cache or lookup a stored response before processing
an incoming request.

</Zoom>

:::note Pre-Requisites

Before exploring Aperture's caching capabilities, make sure that you have signed
up to [Aperture Cloud](https://app.fluxninja.com/sign-up) and set up an
organization. For more information on how to sign up, follow our
[step-by-step guide](/reference/cloud-ui/sign-up.md).

:::

## Caching with Aperture SDK

The first step to using the Aperture SDK is to import and set up Aperture
Client:

```mdx-code-block
<Tabs>
  <TabItem value="TypeScript">
```

<CodeSnippet lang="ts" snippetName="clientConstructor" />

```mdx-code-block
  </TabItem>
</Tabs>
```

You can obtain your organization address and API Key within the Aperture Cloud
UI by clicking the `Aperture` tab in the sidebar menu.

The next step is making a `startFlow` call to Aperture. For this call, it's
crucial to designate the `control point` (`caching-example` in our example) and
the `resultCacheKey`, which facilitates access to the cache in Aperture Cloud.
Additionally, to obtain detailed telemetry data for each Aperture request,
include the labels related to business logic.

```mdx-code-block
<Tabs>
  <TabItem value="TypeScript">
```

<CodeSnippet lang="ts" snippetName="CStartFlow" />

```mdx-code-block
  </TabItem>
</Tabs>
```

After making a `startFlow` call, we check for cached responses in Aperture Cloud
using `flow.resultCache().getLookupStatus()`matching it to (`LookupStatus.Hit`).
Otherwise, in the case of a cache miss, developers can store a new response in
the cache. This is where setting the `ttl` (Time to Live) becomes important, as
it dictates how long the response will be stored in the cache. A longer TTL is
ideal for stable data that doesn't change often, ensuring it's readily available
for frequent access. Conversely, a shorter TTL is more suitable for dynamic data
that requires regular updates, maintaining the cache's relevance and accuracy.
It is important to make the `end` call made after processing each request, in
order to send telemetry data that would provide granular visibility for each
flow.

```mdx-code-block
<Tabs>
  <TabItem value="TypeScript">
```

<CodeSnippet lang="ts" snippetName="CacheLookup" />

```mdx-code-block
  </TabItem>
</Tabs>
```

## Caching in Action

Begin by cloning the
[Aperture JS SDK](https://github.com/fluxninja/aperture-js).

Switch to the example directory and follow these steps to run the example:

1. Install the necessary packages:
   - Run `npm install` to install the base dependencies.
   - Run `npm install @fluxninja/aperture-js` to install the Aperture SDK.
2. Run `npx tsc` to compile the TypeScript example.
3. Run `node dist/caching_example.js` to start the compiled example.

Once the example is running, it will prompt you for your Organization address
and API Key. In the Aperture Cloud UI, select the Aperture tab from the sidebar
menu. Copy and enter both your Organization address and API Key to establish a
connection between the SDK and Aperture Cloud.

Aperture will cache and serve the response for the duration specified by the
TTL. Once the TTL expires, and the cache lookup returns a miss, Aperture will
reset the response in the cache.

Using Aperture's caching feature, developers can enhance application performance
by storing commonly requested data, thereby reducing system load.
