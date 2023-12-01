---
title: Manually setting feature control points
sidebar_position: 1
slug: manually-setting-feature-control-points-using-golang-sdk
keywords:
  - go
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

<a href={`https://pkg.go.dev/github.com/fluxninja/aperture-go/v2`}>Aperture Go
SDK</a> can be used to manually set feature control points within a Go service.

To do so, first create an instance of ApertureClient:

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

```bash
go get github.com/fluxninja/aperture-go/v2
```

<CodeSnippet lang="go" snippetName="grpcOptions" />

<CodeSnippet lang="go" snippetName="clientConstructor" />

The created instance can then be used to start a flow:

<CodeSnippet
    lang="go"
    snippetName="manualFlowNoCaching"
 />

For more context on using Aperture Go SDK to set feature control points, refer
to the [example app][example] available in the repository.

## HTTP Middleware

You can also configure middleware for your HTTP server using the SDK. To do
this, after creating an instance of ApertureClient, apply the middleware to your
router as demonstrated in the example below.

For added convenience, you can specify a list of regular expression patterns.
The middleware will only be applied to request paths that match these patterns.
This feature is particularly beneficial for endpoints such as `/health`,
`/connected` which may not require Aperture intervention.

<CodeSnippet
    lang="go"
    snippetName="middleware"
 />

<!-- TODO: Fix Link -->

[example]: https://github.com/fluxninja/aperture-go/tree/main/example
[api-keys]: /reference/cloud-ui/api-keys.md
