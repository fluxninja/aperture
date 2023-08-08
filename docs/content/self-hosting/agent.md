---
title: Agent Configuration
sidebar_position: 3
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
```

## Installation

Installation process of Aperture Agent doesn't change when using the self-hosted
Aperture Controller, so follow the steps from the [Aperture Agent installation
guide][install-agent].

Below are the configuration options that are specific to the self-hosted
Aperture Controller.

## Configuration

When using the self-hosted Aperture Controller instead of the Aperture Cloud
Controller, you need to turn off the `enable_cloud_controller` flag and
configure Controller, etcd and Prometheus endpoints directly, for example:

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
  secrets:
    fluxNinjaExtension:
      create: true
      secretKeyRef:
        name: aperture-agent-apikey
        key: apiKey
      value: API_KEY
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

You might need to adjust the endpoints, depending on your exact setup.

```mdx-code-block
  </TabItem>
</Tabs>
```

:::info

If you're not using [Aperture Cloud][aperture-cloud], simply remove the
`fluxninja` and `secrets` sections.

:::

:::note

`agent_functions.endpoints` is optional. If you skip it, some `aperturectl`
commands (like `flow-control`) won't work.

:::

[aperture-cloud]: /introduction.md
[install-agent]: /get-started/installation/agent/agent.md
