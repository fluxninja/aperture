---
title: Manually setting feature control points
slug: manually-setting-feature-control-points-using-dotnet-sdk
sidebar_position: 1
keywords:
  - csharp
  - dotnet
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

[Aperture C# SDK](https://www.nuget.org/packages/ApertureSDK/) can be used to
manually set feature control points within a .NET service.

To do so, first create an instance of ApertureClient:

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

```cpp
var sdk = ApertureSdk
    .Builder()
    .SetAddress("ORGANIZATION.app.fluxninja.com:443")
    .SetAgentApiKey("API_KEY")
    .Build();
```

The created instance can then be used to start a flow:

<CodeSnippet lang="cs" snippetName="handleRequest" highlightLanguage="cpp"/>

For more context on using the Aperture C# SDK to set feature control points,
refer to the [example app][example] available in the repository.

[example]: https://github.com/fluxninja/aperture-csharp/tree/main/Examples
[api-keys]: /reference/cloud-ui/api-keys.md
