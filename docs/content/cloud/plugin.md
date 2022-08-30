---
title: FluxNinja Cloud Plugin
sidebar_position: 1
keywords:
  - cloud
  - plugin
---

# FluxNinja cloud plugin

If you are also a FluxNinja Cloud user you can enhance your Aperture experience by enabling FluxNinja plugins.
Plugins provide us with the ability to enrich your logs and traces with additional information.
It implements Open Telemetry standard and allows us to batch/rollup metrics so that they do not take a lot of bandwidth.
By using them Aperture can make better flow control decisions.
We are also sending heartbeats from Aperture Agents and Controllers alike to track their liveness, policy allocation and history.

## Configuration

If you enable plugins then you can configure for example certificates for a secure connection with FluxNinja cloud or a heartbeats interval.

```yaml
plugins:
  disable_plugins: false

fluxninja_plugin:
  fluxninja_endpoint: "test"
  heartbeat_interval: "10s"
  client_grpc:
    insecure: true
    tls:
      insecure_skip_verify: true
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
