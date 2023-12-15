---
title: Define Control Points
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
define feature control points within a .NET service.

Run the command below to install the SDK:

```bash
dotnet add package ApertureSDK --version 2.23.1
```

The next step is to create an ApertureClient instance, for which, the address of
the organization created in Aperture Cloud and API key are needed. You can
locate both these details by clicking on the Aperture tab in the sidebar menu of
Aperture Cloud.

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

The above code snippets starts by defining business critical labels that can be
passed to Aperture as `FeatureFlowParams` when making `StartFlow` calls. Labels
will be matched to the labels set policy created in Aperture Cloud, and a
decision will be returned on whether a flow `ShouldRun` or not. In this example,
we only see log returns, but in a production environment, actual business logic
could be executed when a request is allowed. It is important to make the `End`
call made after processing each request, in order to send telemetry data that
would provide granular visibility for each flow.

For more context on using the Aperture C# SDK to set feature control points,
refer to the [example app][example] available in the repository.

[example]: https://github.com/fluxninja/aperture-csharp/tree/main/Examples
[api-keys]: /reference/cloud-ui/api-keys.md
