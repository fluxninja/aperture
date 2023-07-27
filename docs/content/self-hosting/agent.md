---
title: Agent Configuration
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
```

## Installation

Installation process of Aperture Agents doesn't change when using self-hosted
Controller, so follow the steps from [agent installation guide][install-agent].
The places when configuration with self-hosted Controller differs will be marked
as such.

## Configuration

When using self-hosted controller instead of FluxNinja Cloud Controller, you
need to turn off the `enable_cloud_controller` flag and configure controller,
etcd and Prometheus endpoints directly, for example:

```mdx-code-block
<Tabs>
  <TabItem value="aperturectl or helm">
```

`values.yaml`:

```yaml
agent:
  config:
    fluxninja:
      enable_cloud_controller: false
      endpoint: ORGANIZATION_NAME.app.fluxninja.com:443
    etcd:
      endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
    prometheus:
      address: "http://controller-prometheus-server.default.svc.cluster.local:80"
    agent_functions:
      endpoints: ["aperture-controller.default.svc.cluster.local:8080"]
  secrets: ...
```

The values above assume that you have installed the
[Aperture Controller](/self-hosting/controller/controller.md) on the same
cluster in `default` namespace, with etcd and Prometheus using `controller` as
release name. If your setup is different, adjust these endpoints accordingly.

```mdx-code-block
  </TabItem>

  <TabItem value="Docker or Bare Metal">
```

`agent.yaml`:

```yaml
fluxninja:
  enable_cloud_controller: false
  endpoint: ORGANIZATION_NAME.app.fluxninja.com:443
etcd:
  endpoints: ["http://etcd:2379"]
prometheus:
  address: "http://prometheus:80"
agent_functions:
  endpoints: ["aperture-controller:8080"]
```

You may need to adjust the endpoints, depending on your exact setup.

```mdx-code-block
  </TabItem>
</Tabs>
```

:::info

If you're not using [FluxNinja][] at all, simply remove the `fluxninja` and
`secrets` sections.

:::

:::note

`agent_functions.endpoints` is optional. If you skip it, some `aperturectl`
subcommands (like `flow-control`) won't work.

:::

[FluxNinja]: /fluxninja/introduction.md
[install-agent]: /get-started/installation/agent/agent.md
