# Aperture-Go SDK

`aperture-go` is an SDK to interact with Aperture Agent. It allows flow control
functionality on fine-grained features inside service code.

## Usage

### ApertureClient Interface

`ApertureClient` maintains a gRPC connection with Aperture Agent.

```go
 options := aperture.Options{
  ClientConn: client,
  // checkTimeout is the time that the client will wait for a response from Aperture Agent.
  // if not provided, the default value of 200 milliseconds will be used.
  CheckTimeout: 200 * time.Millisecond,
 }

 // initialize Aperture Client with the provided options.
 apertureClient, err := aperture.NewClient(options)
 if err != nil {
  log.Fatalf("failed to create client: %v", err)
 }
```

### Flow Interface

`Flow` is created every time `ApertureClient.BeginFlow` is called.

```go
 // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
 flow, err := a.apertureClient.StartFlow(ctx, "awesomeFeature", labels)
 if err != nil {
  log.Printf("Aperture flow control got error. Returned flow defaults to Allowed. flow.ShouldRun(): %t", flow.ShouldRun())
 }

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
