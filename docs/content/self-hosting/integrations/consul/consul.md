---
title: Consul
keywords:
  - install
  - setup
  - service mesh
  - consul
  - service defaults
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
```

![Consul](./assets/consul-light.svg#gh-light-mode-only)

![Consul](./assets/consul-dark.svg#gh-dark-mode-only)

## Supported Versions

Aperture supports the following versions of Consul:

| Platform          | Extent of Support |
| ----------------- | ----------------- |
| Consul            | 1.17 and above    |
| Consul Data Plane | 1.2.2 and above   |

:::info

This integration is currently only supported with the
[self-hosted](/self-hosting/agent/agent.md) version of Aperture Agent.

:::

## Envoy Extensions {#envoy-extensions}

The
[Envoy Extensions](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions)
in Consul is used to modify the behavior of Envoy.

The Aperture Agent requires additional details and needs the
[External Authorization](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/ext-authz)
and
[OpenTelemetry Access Logging](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/otel-access-logging)
Envoy Extensions to be added through the
[Service Defaults](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-defaults).

:::info

Consul also supports
[Proxy Defaults](https://developer.hashicorp.com/consul/docs/connect/config-entries/proxy-defaults)
but as of now, Aperture Agent does not support it.

The reason being, the Proxy Defaults are applied to all the proxies in the
Consul Data Plane, which includes the proxies running in the Aperture Agent
itself. This will create a loop of requests from the Aperture Agent to itself.

:::

**Note**: In all the below patches, it is presumed that the Aperture Agent is
installed with `DaemonSet` mode and is installed in the `aperture-agent`
namespace, which makes the target service name `aperture-agent` and namespace
`aperture-agent`. If you are running the Aperture Agent in Sidecar mode, use
`localhost` as the target address.

1. As the
   [External Authorization](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/ext-authz)
   and
   [OpenTelemetry Access Logging](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/otel-access-logging)
   support gRPC, the below Service Defaults configures the protocol of the
   Aperture Agent to gRPC.

   ```yaml
   apiVersion: consul.hashicorp.com/v1alpha1
   kind: ServiceDefaults
   metadata:
     name: aperture-agent
     namespace: aperture-agent
   spec:
     protocol: grpc
   ```

2. The below Service Defaults adds
   [OpenTelemetry Access Logging](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/otel-access-logging)
   Envoy Extension for the outbound listener, in the Consul Proxy running with
   the application.

   The OpenTelemetry configuration in the following has extracted values, which
   are forwarded to the Aperture Agent instance using gRPC.

   The prepared log has the request method value as log body and `egress` as the
   log name to differentiate between different access logs coming from the same
   Envoy.

   ```yaml
   name: builtin/otel-access-logging
   required: true
   arguments:
     listenerType: outbound
     config:
       attributes:
         aperture.source: "envoy"
         aperture.check_response: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:aperture.check_response)%"
         http.status_code: "%RESPONSE_CODE%"
         authz_duration: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:ext_authz_duration)%"
         BYTES_RECEIVED: "%BYTES_RECEIVED%"
         BYTES_SENT: "%BYTES_SENT%"
         DURATION: "%DURATION%"
         REQUEST_DURATION: "%REQUEST_DURATION%"
         REQUEST_TX_DURATION: "%REQUEST_TX_DURATION%"
         RESPONSE_DURATION: "%RESPONSE_DURATION%"
         RESPONSE_TX_DURATION: "%RESPONSE_TX_DURATION%"
       body: "%REQ(:METHOD)%"
       logName: "egress"
       grpcService:
         target:
           service:
             name: aperture-agent
             namespace: aperture-agent
   ```

3. The below Service Defaults adds
   [OpenTelemetry Access Logging](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/otel-access-logging)
   Envoy Extension, but for the inbound listener, in the Consul Proxy running
   with the application.

   The OpenTelemetry configuration in the following has extracted values, which
   are forwarded to the Aperture Agent instance using gRPC.

   The prepared log has the request method value as log body and `ingress` as
   the log name to differentiate between different access logs coming from the
   same Envoy.

   ```yaml
   name: builtin/otel-access-logging
   required: true
   arguments:
     listenerType: inbound
     config:
       attributes:
         aperture.source: "envoy"
         aperture.check_response: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:aperture.check_response)%"
         http.status_code: "%RESPONSE_CODE%"
         authz_duration: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:ext_authz_duration)%"
         BYTES_RECEIVED: "%BYTES_RECEIVED%"
         BYTES_SENT: "%BYTES_SENT%"
         DURATION: "%DURATION%"
         REQUEST_DURATION: "%REQUEST_DURATION%"
         REQUEST_TX_DURATION: "%REQUEST_TX_DURATION%"
         RESPONSE_DURATION: "%RESPONSE_DURATION%"
         RESPONSE_TX_DURATION: "%RESPONSE_TX_DURATION%"
       body: "%REQ(:METHOD)%"
       logName: "ingress"
       grpcService:
         target:
           service:
             name: aperture-agent
             namespace: aperture-agent
   ```

4. The below Service Defaults adds
   [External Authorization](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/ext-authz)
   Envoy Extension, for the outbound listener, in the Consul Proxy running with
   the application.

   The External Authorization extension forwards the request to the Aperture
   Agent instance using gRPC with a timeout of `0.25s`, having `egress` value
   for key `control-point` metadata included in the streams initiated to the
   gRPC service. The extension will accept the client request even if the
   communication with the authorization service has failed, or if the
   authorization service has returned an HTTP 5xx error.

   ```yaml
   name: builtin/ext-authz
   required: true
   arguments:
     listenerType: outbound
     config:
       statPrefix: ext_authz
       timeout: 0.25s
       grpcService:
         target:
           service:
             name: aperture-agent
             namespace: aperture-agent
         initialMetadata:
           - key: control-point
             value: "egress"
   ```

5. The below Service Defaults adds
   [External Authorization](https://developer.hashicorp.com/consul/docs/v1.17.x/connect/proxies/envoy-extensions/usage/ext-authz)
   Envoy Extension, but for the inbound listener, in the Consul Proxy running
   with the application.

   The External Authorization extension forwards the request to the Aperture
   Agent instance using gRPC with a timeout of `0.25s`, having `ingress` value
   for key `control-point` metadata included in the streams initiated to the
   gRPC service. The extension will accept the client request even if the
   communication with the authorization service has failed, or if the
   authorization service has returned an HTTP 5xx error.

   ```yaml
   name: builtin/ext-authz
   required: true
   arguments:
     listenerType: inbound
     config:
       statPrefix: ext_authz
       timeout: 0.25s
       grpcService:
         target:
           service:
             name: aperture-agent
             namespace: aperture-agent
         initialMetadata:
           - key: control-point
             value: "ingress"
   ```

## Installation {#installation}

The complete Service Defaults configuration as an example for the Consul Proxy
running with the application is as follows:

<details><summary>service1-demo-app-service-defaults.yaml</summary>
<p>
<CodeBlock language="yaml">
{`apiVersion: consul.hashicorp.com/v1alpha1
kind: ServiceDefaults
metadata:
  name: service1-demo-app
  namespace: demoapp
spec:
  protocol: http
  envoyExtensions:
  - name: builtin/ext-authz
    required: true
    arguments:
      listenerType: inbound
      config:
        statPrefix: ext_authz
        timeout: 0.25s
        grpcService:
          target:
            service:
              name: aperture-agent
              namespace: aperture-agent
          initialMetadata:
          - key: control-point
            value: "ingress"
  - name: builtin/otel-access-logging
    required: true
    arguments:
      listenerType: inbound
      config:
        attributes:
          aperture.source: "envoy"
          aperture.check_response: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:aperture.check_response)%"
          http.status_code: "%RESPONSE_CODE%"
          authz_duration: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:ext_authz_duration)%"
          BYTES_RECEIVED: "%BYTES_RECEIVED%"
          BYTES_SENT: "%BYTES_SENT%"
          DURATION: "%DURATION%"
          REQUEST_DURATION: "%REQUEST_DURATION%"
          REQUEST_TX_DURATION: "%REQUEST_TX_DURATION%"
          RESPONSE_DURATION: "%RESPONSE_DURATION%"
          RESPONSE_TX_DURATION: "%RESPONSE_TX_DURATION%"
        body: "%REQ(:METHOD)%"
        logName: "ingress"
        grpcService:
          target:
            service:
              name: aperture-agent
              namespace: aperture-agent
  - name: builtin/ext-authz
    required: true
    arguments:
      listenerType: outbound
      config:
        statPrefix: ext_authz
        timeout: 0.25s
        grpcService:
          target:
            service:
              name: aperture-agent
              namespace: aperture-agent
          initialMetadata:
          - key: control-point
            value: "egress"
  - name: builtin/otel-access-logging
    required: true
    arguments:
      listenerType: outbound
      config:
        attributes:
          aperture.source: "envoy"
          aperture.check_response: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:aperture.check_response)%"
          http.status_code: "%RESPONSE_CODE%"
          authz_duration: "%DYNAMIC_METADATA(envoy.filters.http.ext_authz:ext_authz_duration)%"
          BYTES_RECEIVED: "%BYTES_RECEIVED%"
          BYTES_SENT: "%BYTES_SENT%"
          DURATION: "%DURATION%"
          REQUEST_DURATION: "%REQUEST_DURATION%"
          REQUEST_TX_DURATION: "%REQUEST_TX_DURATION%"
          RESPONSE_DURATION: "%RESPONSE_DURATION%"
          RESPONSE_TX_DURATION: "%RESPONSE_TX_DURATION%"
        body: "%REQ(:METHOD)%"
        logName: "egress"
        grpcService:
          target:
            service:
              name: aperture-agent
              namespace: aperture-agent
`}
</CodeBlock>
</p>
</details>

To install the above Service Defaults, run the following command:

```bash
kubectl apply -f service1-demo-app-service-defaults.yaml
```

## Verifying the Installation

To verify the installation, run the following command:

```bash
kubectl get servicedefaults -n demoapp service1-demo-app -o yaml
```

## Uninstall

To uninstall the Service Defaults, run the following command:

```bash
kubectl delete servicedefaults -n demoapp service1-demo-app
```
