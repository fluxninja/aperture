# Istio Configuration

This Chart inserts Envoy filters that integrate with Aperture Agent.

## Parameters

### Envoy Filter Parameters

| Name                                          | Description                                                                                                             | Value            |
| --------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | ---------------- |
| `envoyFilter.name`                            | Name of service running `aperture-agent`                                                                                | `aperture-agent` |
| `envoyFilter.namespace`                       | Namespace where `aperture-agent` is running                                                                             | `aperture-agent` |
| `envoyFilter.port`                            | Port serving external authorization API and for streaming access logs                                                   | `8080`           |
| `envoyFilter.authzGrpcTimeout`                | Timeout in seconds to authorization requests made to `aperture-agent`.                                                  | `0.25s`          |
| `envoyFilter.enableAuthzRequestBodyBuffering` | Enable buffering request body that is sent over external authorization API. Note: This might break some streaming APIs. | `true`           |
| `envoyFilter.maxRequestBytes`                 | Maximum size of request that is sent over external authorization API                                                    | `8192`           |
| `envoyFilter.sidecarMode`                     | Aperture Agent installed using the Sidecar mode                                                                         | `false`          |
