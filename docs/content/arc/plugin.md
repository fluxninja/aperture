---
title: FluxNinja ARC Plugin
sidebar_label: Plugin
sidebar_position: 1
keywords:
  - cloud
  - plugin
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import CodeBlock from '@theme/CodeBlock';
```

```mdx-code-block
export const PluginConfig = ({children, component}) => (
<CodeBlock language="yaml">
{`${component}:
  config:
    plugins:
      disabled_plugins: []
    fluxninja_plugin:
      fluxninja_endpoint: ORGANIZATION_NAME.app.fluxninja.com:443
      client:
        grpc:
          insecure: false
          tls:
            insecure_skip_verify: true
        http:
          tls:
            insecure_skip_verify: true
  secrets:
    fluxNinjaPlugin:
      create: true
      secretKeyRef:
        name: aperture-${component}-apikey
        key: apiKey
      value: API_KEY
`}</CodeBlock>
);
```

If you are a FluxNinja ARC customer, you can enhance your Aperture experience by
enabling FluxNinja plugin. It enriches logs and traces collected by Aperture
with additional dimensions and batches/rollups metrics to optimize bandwidth
usage. In FluxNinja ARC, you can monitor your policies and analyze flows.
FluxNinja plugin also sends periodic heartbeats from Aperture Agents and
Controllers to track their health and configuration.

## Configuration

Configure below parameters in the `values.yaml` file generated during
installation of the Aperture Controller or Agent:

<Tabs>
  <TabItem value="Controller">
    <PluginConfig component="controller" />
  </TabItem>
  <TabItem value="Agent">
    <PluginConfig component="agent" />
  </TabItem>
</Tabs>

Replace the values of `ORGANIZATION_NAME` and `API_KEY` with the actual values
of the organization on FluxNinja ARC and API Key generated on it.

Configuration parameters for the FluxNinja ARC Plugin are available below:

- [Aperture Controller](/references/configuration/controller.md/#flux-ninja-plugin)
- [Aperture Agent](/references/configuration/agent.md#flux-ninja-plugin)

## See also

How various components interact with the plugin:

- [Flow labels](/concepts/integrations/flow-control/flow-label.md#plugin)
