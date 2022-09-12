---
title: FluxNinja Cloud Plugin
sidebar_label: Plugin
sidebar_position: 1
keywords:
  - cloud
  - plugin
---

# FluxNinja cloud plugin

If you are a FluxNinja Cloud customer you can enhance your Aperture experience
by enabling FluxNinja plugins. Plugins provide us with the ability to enrich
your logs and traces with additional information. The plugins process logs and
traces according to the Open Telemetry standard. We batch/rollup metrics so that
they do not take a lot of bandwidth. In FluxNinja Cloud, you can visualize
policies and use the information from plugins to refine your policies. We also
send heartbeats from Aperture Agents and Controllers to track their liveness,
policy allocation and history.

## Configuration

If you enable plugins then you can configure certificates for a secure
connection with FluxNinja Cloud or a heartbeat interval.

```yaml
plugins:
  disable_plugins: false

fluxninja_plugin:
  fluxninja_endpoint: "placeholder.local.dev.fluxninja.com:443"
  heartbeat_interval: "10s"
  client:
    grpc:
      insecure: false
      tls:
        insecure_skip_verify: true
        cert_file: "cert.pem"
        key_file: "key.pem"
        ca_file: test
    http:
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

- [Agent](/reference/configuration/agent.md#plugins)
- [Controller](/reference/configuration/controller.md#plugins)

## See also

How various components interact with the plugin:

- [Flow labels](/concepts/flow-control/selector/flow-label.md#plugin)
