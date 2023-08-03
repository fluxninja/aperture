---
title: FluxNinja Extension
sidebar_label: FluxNinja Extension
sidebar_position: 8
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

FluxNinja extension enables [FluxNinja Cloud][] integration for Aperture Agents
(and [self-hosted][Self-Hosting] Controllers). It enriches logs and
traces collected by Aperture and sends them to FluxNinja Cloud. This data is
batched and rolled up to optimize bandwidth usage. The extension also sends
periodic heartbeats to FluxNinja Cloud to track health and configuration.
This allows you to monitor your policies and analyze flows in FluxNinja Cloud.

FluxNinja extension also provides the possibility to use the FluxNinja Cloud
Controller.

## FluxNinja Cloud Controller {#cloud-controller}

Without [Aperture Controller][], Aperture Agents won't be able to work.
While it's possible to [self-host][Self-Hosting] Aperture Controller,
FluxNinja Cloud Controller can be used instead.

FluxNinja Cloud Controller is an [Aperture Controller] hosted by FluxNinja Cloud.
The Cloud Controller is available for every FluxNinja Cloud Organization in the
`default` project.

## Configuration

Configure the following parameters in the `values.yaml` file generated during
installation of the Aperture Controller or Agent:

<Tabs>
  <TabItem value="FluxNinja Cloud Controller">
    <Tabs>
      <TabItem value="Agent">
        <CloudExtensionConfig />
      </TabItem>
    </Tabs>
  </TabItem>
  <TabItem value="Self-Hosted Controller">
    <Tabs>
      <TabItem value="Controller">
        <ExtensionConfig component="controller" />
      </TabItem>
      <TabItem value="Agent">
        <ExtensionConfig component="agent" />
      </TabItem>
    </Tabs>
  </TabItem>
</Tabs>

Replace the values of `ORGANIZATION_NAME` and `API_KEY` with the actual values
of the organization on FluxNinja Cloud and API Key generated on it.

More specific details about particular agent installation modes could be found
in [Get Started: Installation](/get-started/installation/agent/agent.md).

Configuration parameters for the FluxNinja Extension are as follows:

- [Aperture Agent](/reference/configuration/agent.md#flux-ninja-extension-config)
- [Aperture Controller](/reference/configuration/controller.md/#flux-ninja-extension-config)

## See also

How various components interact with the extension:

- [Flow labels](/concepts/flow-label.md#extension)

[Self-Hosting]: /self-hosting/self-hosting.md
[FluxNinja Cloud]: /introduction.md
[Aperture Controller]: /architecture/architecture.md#aperture-controller
