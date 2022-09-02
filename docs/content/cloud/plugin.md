---
title: FluxNinja Cloud Plugin
sidebar_label: Plugin
sidebar_position: 1
keywords:
  - cloud
  - plugin
---

# FluxNinja cloud plugin

If you are also a FluxNinja Cloud user you can enhance your Aperture experience by enabling FluxNinja plugins.
Plugins provide us with the ability to enrich your logs and traces with additional information.
Plugin is processing logs and traces according to Open Telemetry standard.
It allows us to batch/rollup metrics so that they do not take a lot of bandwidth.
In FluxNinja Cloud you can view and use information from plugins to refine your policies.
We are also sending heartbeats from Aperture Agents and Controllers alike to track their liveness, policy allocation and history.

## Configuration

If you enable plugins then you can configure for example certificates for a secure connection with FluxNinja cloud or a heartbeats interval.

```yaml
plugins:
  disable_plugins: false

fluxninja_plugin:
  fluxninja_endpoint: "placeholder.local.dev.fluxninja.com:443"
  heartbeat_interval: "10s"
  client_grpc:
    insecure: false
    tls:
      insecure_skip_verify: true
      cert_file: "cert.pem"
      key_file: "key.pem"
      ca_file: test
  client_http:
    tls:
      insecure_skip_verify: true
      ca_file: test
```

You can also selectively disable plugins like this:

```yaml
plugins:
  disable_plugins: false
  disabled_plugins:
    - aperture-plugin-fluxninja
```

For more details, see plugin configuration reference:

- [Agent](/reference/configuration/agent#plugin-configuration)
- [Controller](/reference/configuration/controller#plugin-configuration)

## See also

How various components interact with the plugin:

- [Flow labels](/concepts/flow-control/label/label.md#plugin)
