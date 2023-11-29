---
title: Manually setting feature control points
sidebar_position: 1
slug: manually-setting-feature-control-points-using-python-sdk
keywords:
  - python
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

[Aperture Python SDK][pythonsdk] can be used to manually set feature control
points within a Python service.

To do so, first create an instance of ApertureClient:

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, please refer to the
[API Keys](/reference/cloud-ui/api-keys.md) section.

:::

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'

<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="clientConstructor"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

The created instance can then be used to start a flow:

```mdx-code-block
<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="manualFlow"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

You can also use the flow as a context manager:

```mdx-code-block
<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="contextManagerFlow"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

Additionally, you can decorate any function with aperture client. This will skip
running the function if the flow is rejected by Aperture Agent. This might be
helpful to handle specific routes in your service.

```mdx-code-block
<Tabs>
<TabItem value="Python">
```

<CodeSnippet
lang="py"
snippetName="apertureDecorator"
/>

```mdx-code-block
</TabItem>
</Tabs>
```

For more context on using the Aperture Python SDK to set feature control points,
refer to the [example app][example] available in the repository.

[pythonsdk]: https://pypi.org/project/aperture-py/
[example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-py/example
