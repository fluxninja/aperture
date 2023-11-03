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

<a href={`https://pkg.go.dev/github.com/fluxninja/aperture-go`}>Aperture Go
SDK</a> can be used to manually set feature control points within a Go service.

To do so, first create an instance of ApertureClient:

:::info Agent API Key

You can create an Agent API key for your project in the Aperture Cloud UI. For
more information, refer to
[Agent API Keys](/get-started/aperture-cloud/agent-api-keys.md).

:::

```go
  agentAddress = "ORGANIZATION.app.fluxninja.com:443"
  agentAPIKey = "AGENT_API_KEY"

  opts := aperture.Options{
      Address:                     agentAddress,
      APIKey:                      agentAPIKey,
  }

  // initialize Aperture Client with the provided options.
  apertureClient, err := aperture.NewClient(ctx, opts)
  if err != nil {
      log.Fatalf("failed to create client: %v", err)
  }
```

The created instance can then be used to start a flow:

```go
  // business logic produces labels
  labels := map[string]string{
      "key": "value",
  }

  rampMode := false

  // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
  flow, err := a.apertureClient.StartFlow(ctx, "featureName", labels, rampMode, 200 * time.Millisecond)
  if err != nil {
      log.Printf("Aperture flow control got error. Returned flow defaults to Allowed. flow.ShouldRun(): %t", flow.ShouldRun())
  }

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
