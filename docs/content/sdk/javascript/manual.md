---
title: Manually setting feature control points
slug: manually-setting-feature-control-points-using-javascript-sdk
sidebar_position: 1
keywords:
  - js
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'
```

[Aperture JavaScript SDK](https://www.npmjs.com/package/@fluxninja/aperture-js)
can be used to manually set feature control points within a JavaScript service.

To do so, first create an instance of ApertureClient:

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

<CodeSnippet lang="ts" snippetName="clientConstructor" />

The created instance can then be used to start a flow:

<CodeSnippet
    lang="ts"
    snippetName="handleRequestRateLimit"
 />

For more context on using the Aperture JavaScript SDK to set feature control
points, refer to the [example app][example] available in the repository.

[example]: https://github.com/fluxninja/aperture-js/tree/main/example
[api-keys]: /reference/cloud-ui/api-keys.md
