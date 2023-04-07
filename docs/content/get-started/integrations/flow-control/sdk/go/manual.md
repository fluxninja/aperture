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
SDK</a> can be used to manually set feature Control Points within a Go service.

To do so, first create an instance of ApertureClient:

```go
agentHost := "localhost"
agentPort := "8089"

apertureAgentGRPCClient, err := grpcClient(ctx, net.JoinHostPort(agentHost, agentPort))
if err != nil {
  log.Fatalf("failed to create flow control client: %v", err)
}

opts := aperture.Options{
  ApertureAgentGRPCClientConn: apertureAgentGRPCClient,
  CheckTimeout:                200 * time.Millisecond,
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

    // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
    flow, err := a.apertureClient.StartFlow(ctx, "featureName", labels)
    if err != nil {
        log.Printf("Aperture flow control got error. Returned flow defaults to Allowed. flow.Accepted(): %t", flow.Accepted())
    }

    // See whether flow was accepted by Aperture Agent.
    if flow.Accepted() {
        // do actual work
        _ = flow.End(aperture.OK)
    } else {
        // handle flow rejection by Aperture Agent
        _ = flow.End(aperture.Error)
    }
```

For more context on how to use Aperture Go SDK to set feature Control Points,
you can take a look at the [example app][example] available in our repository.

[example]: https://github.com/fluxninja/aperture-go/tree/v1.0.0/example
