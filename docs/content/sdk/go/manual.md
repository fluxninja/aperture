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

<!-- TODO: Convert into Generated Code Snippet -->

```go
  // business logic produces labels
  labels := map[string]string{
      "key": "value",
  }

  rampMode := false

  // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow object.
  flow := apertureClient.StartFlow(ctx, "featureName", labels, rampMode, 200 * time.Millisecond)

  // See whether flow was accepted by Aperture Agent.
  if flow.ShouldRun() {
      // do actual work
  } else {
      // handle flow rejection by Aperture Agent
      flow.SetStatus(aperture.Error)
  }
  _ = flow.End()
```

For more context on using Aperture Go SDK to set feature control points, refer
to the [example app][example] available in the repository.

## HTTP Middleware

You can also automatically set middleware for your HTTP server using the SDK. To
do so, after creating an instance of ApertureClient, use the middleware on your
router:

```go
  mux.Use(aperturemiddlewares.NewHTTPMiddleware(apertureClient, "awesomeFeature", labels, nil, false, 200 * time.Millisecond).Handle)
```

For simplicity, you can also pass a list of regexp patterns to match against the
request path, for which the middleware will pass through. This is especially
useful for endpoints like `/healthz`:

```go
  mux.Use(aperturemiddlewares.NewHTTPMiddleware(apertureClient, "awesomeFeature", labels,  []regexp.Regexp{regexp.MustCompile("/health.*")}, false, 200 * time.Millisecond).Handle)
```

[example]: https://github.com/fluxninja/aperture-go/tree/main/example
[api-keys]: /reference/cloud-ui/api-keys.md
