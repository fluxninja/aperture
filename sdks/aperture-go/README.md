# Aperture-Go SDK

`aperture-go` is an SDK to interact with Aperture Agent. It allows flow control
functionality on fine-grained features inside service code.

## Usage

```bash
go get github.com/fluxninja/aperture-go/v2
```

### ApertureClient Interface

`ApertureClient` maintains a gRPC connection with Aperture Agent.

```go
options := aperture.Options{
   DialOptions: grpcOptions,
   Address:     "ORGANIZATION.app.fluxninja.com",
   AgentAPIKey: "AGENT_API_KEY",
}

// initialize Aperture Client with the provided options.
apertureClient, err := aperture.NewClient(options)
if err != nil {
   log.Fatalf("failed to create client: %v", err)
}
```

### HTTP Middleware

`aperture-go` provides an HTTP middleware to be used with routers.

```go
// Create a new mux router
router := mux.NewRouter()

superRouter := mux.PathPrefix("/super").Subrouter()
superRouter.HandleFunc("", a.SuperHandler)

superRouter.Use(aperturegomiddleware.NewHTTPMiddleware(apertureClient, "awesomeFeature", nil, nil, false, 2000*time.Millisecond).Handle)
```

### gRPC Unary Interceptor

`aperture-go` provides a gRPC unary interceptor to be used with gRPC clients.

```go
// Create a new gRPC interceptor
interceptor := aperturegomiddleware.NewGRPCUnaryInterceptor(apertureClient, "awesomeFeature", nil, false, 2000*time.Millisecond)

// Create a new gRPC server
s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
```

### Flow Interface

`Flow` is created every time `ApertureClient.StartFlow` is called.

```go
// StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow object.
flow := apertureClient.StartFlow(ctx, "awesomeFeature", labels, false, 200 * time.Millisecond)

// See whether flow was accepted by Aperture Agent.
if flow.ShouldRun() {
   // Simulate work being done
   time.Sleep(5 * time.Second)
} else {
   // Flow has been rejected by Aperture Agent.
   flow.SetStatus(aperture.Error)
}

// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop. SetStatus() method of Flow object can be used to capture whether the Flow was successful or resulted in an error. If not set, status defaults to OK.
_ = flow.End()
```

## Relevant Resources

[FluxNinja Aperture](https://github.com/fluxninja/aperture)
