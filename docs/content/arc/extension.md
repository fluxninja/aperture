---
title: FluxNinja ARC Extension
sidebar_label: Extension
sidebar_position: 1
keywords:
  - cloud
  - extension
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import CodeBlock from '@theme/CodeBlock';
```

```mdx-code-block
export const ExtensionConfig = ({children, component}) => (
<CodeBlock language="yaml">
{`${component}:
  config:
    fluxninja:
      endpoint: ORGANIZATION_NAME.app.fluxninja.com:443
      client:
        grpc:
          insecure: false
          tls:
            insecure_skip_verify: true
        http:
          tls:
            insecure_skip_verify: true
  secrets:
    fluxNinjaExtension:
      create: true
      secretKeyRef:
        name: aperture-${component}-apikey
        key: apiKey
      value: API_KEY
`}</CodeBlock>
);
```

```mdx-code-block
export const CloudExtensionConfig = ({children, component}) => (
<CodeBlock language="yaml">
{`agent:
  config:
    fluxninja:
      enable_cloud_controller: true
      endpoint: ORGANIZATION_NAME.app.fluxninja.com:443
  secrets:
    fluxNinjaExtension:
      create: true
      secretKeyRef:
        name: aperture-agent-apikey
        key: apiKey
      value: API_KEY
`}</CodeBlock>
);
```

If you are a FluxNinja ARC customer, you can enhance your Aperture experience by
enabling the FluxNinja extension. It enriches logs and traces collected by
Aperture with additional dimensions and batches and rolls ups metrics to
optimize bandwidth usage. In FluxNinja ARC, you can monitor your policies and
analyze flows. The FluxNinja extension also sends periodic heartbeats from
Aperture Agents and Controllers to track their health and configuration.

## Configuration

Configure the following parameters in the `values.yaml` file generated during
installation of the Aperture Controller or Agent:

<Tabs>
  <TabItem value="On-premise Controller">
    <Tabs>
      <TabItem value="Controller">
        <ExtensionConfig component="controller" />
      </TabItem>
      <TabItem value="Agent">
        <ExtensionConfig component="agent" />
      </TabItem>
    </Tabs>
  </TabItem>
  <TabItem value="ARC Controller">
    <Tabs>
      <TabItem value="Agent">
        <CloudExtensionConfig />
      </TabItem>
    </Tabs>
  </TabItem>
</Tabs>

Replace the values of `ORGANIZATION_NAME` and `API_KEY` with the actual values
of the organization on FluxNinja ARC and API Key generated on it.

Configuration parameters for the FluxNinja ARC Extension are as follows:

- [Aperture Controller](/reference/configuration/controller.md/#flux-ninja)
- [Aperture Agent](/reference/configuration/agent.md#flux-ninja)

## See also

How various components interact with the extension:

- [Flow labels](/concepts/flow-label.md#extension)
