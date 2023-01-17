# sdk-validtor

`sdk-validtor` is a CLI tool to validate Aperture SDKs. Following are the
supported flags that can be passed to the binary.

```shell
$ go run main.go --help
Usage of sdk-validator:
  -port string
     Port to start sdk-validator's grpc server on. (default "8089")
  -rejects int
     Number of requests (out of 'requests') to reject. (default 5)
  -requests int
     Number of requests to make to SDK example server. (default 10)
  -sdk-docker-image string
     Docker image of SDK example to run.
  -sdk-port string
     Port to expose on SDK's example container. (default "8080")
```

In order to validate an SDK, you need to pass the docker image of the SDK's
example server. For example, to validate the [Go SDK](../../sdks/aperture-go),
you can run the following command:

```shell
# Change your current directory to SDK's directory
$ cd fluxninja/aperture/sdks/aperture-go
# Build SDK example's Docker image
$ DOCKER_BUILDKIT=1 docker build -t aperture-go-example:0.1.0 .
# Change your current directory to SDK Validator's directory
$ cd ../../cmd/sdk-validator
# Run SDK Validator
$ go run main.go --sdk-docker-image aperture-go-example:0.1.0
```

The above command will start the SDK's example server in a Docker container and
start the SDK Validator's gRPC server. The SDK Validator will make sure the
example server running inside the Docker container is exposed on the host
machine. It will then make requests to the SDK's example server and validate the
SDK.

In order for the SDK Validator to validate the SDK, it is required that the
SDK's example server exposes the following HTTP endpoints:

1. `/super` where a feature is protected using flow control.
2. `/connected` checks the gRPC client connection to the Aperture server, and
   responds with `StatusServiceUnavailable 503` status if the connection is not
   in `READY` state.
3. `/health` responds with `StatusOK 200` status and "Healthy" message.

The SDK Validator will check that the SDK's example server is exposing the above
endpoints and that the SDK is working as expected.

Refer to the [Go SDK's example server](../../sdks/aperture-go/example/main.go)
for an example of how to expose the above endpoints.
