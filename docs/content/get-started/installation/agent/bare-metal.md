---
title: Bare Metal/VM
description: Install Aperture Agent on Bare Metal or VM
keywords:
  - install
  - setup
  - agent
  - os
  - baremetal
  - vm
sidebar_position: 3
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import CodeBlock from '@theme/CodeBlock';
import {apertureVersion, apertureVersionWithOutV} from '../../../apertureVersion.js';

```

The Aperture Agent can be installed as a system service on any Linux system that
is [supported](../supported-platforms.md).

## Download {#agent-download}

The Aperture Agent can be installed using packages made for your system's
package manager, like `dpkg` or `rpm`.

To install Aperture Agent, first download the package for your manager from the
[Releases Page](https://github.com/fluxninja/aperture/releases/latest).

Alternatively, download it using the following script:

```mdx-code-block
export const DownloadScript = ({children, packager, arch, archSeparator, versionSeparator, component}) => (
<CodeBlock language="bash">
{`VERSION="${apertureVersionWithOutV}"
ARCH="${arch}"
PACKAGER="${packager}"
url="https://github.com/fluxninja/aperture/releases/download/v\${VERSION}/${component}${versionSeparator}\${VERSION}${archSeparator}\${ARCH}.\${PACKAGER}"
echo "Will download \${PACKAGER} package version \${VERSION} compiled for \${ARCH} machine"
curl --fail --location --remote-name "\${url}"
`}</CodeBlock>
);
```

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg">
    <DownloadScript packager="deb" arch="amd64" archSeparator="_" versionSeparator="_" component="aperture-agent" />
  </TabItem>
  <TabItem value="rpm" label="rpm">
    <DownloadScript packager="rpm" arch="x86_64" archSeparator="." versionSeparator="-" component="aperture-agent" />
  </TabItem>
</Tabs>

## Installation {#agent-installation}

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg">
    <CodeBlock language="bash">{`sudo dpkg -i aperture-agent_${apertureVersionWithOutV}*.deb`}</CodeBlock>
  </TabItem>
  <TabItem value="rpm" label="rpm">
    <CodeBlock language="bash">{`sudo rpm -i aperture-agent-${apertureVersionWithOutV}*.rpm`}</CodeBlock>
  </TabItem>
</Tabs>

To point the Aperture Agent to the Aperture Controller, edit the configuration
file located at `/etc/aperture/aperture-agent/config/aperture-agent.yaml`.
Follow the steps for
[FluxNinja Cloud Controller](/reference/fluxninja.md#configuration) or
[Self-Hosted Aperture Controller](/self-hosting/agent.md#configuration),
depending on your setup.

All the configuration parameters for the Aperture Agent are available
[here](/reference/configuration/agent.md).

After installing, you should enable the `aperture-agent` `systemd` service, and
make it start after system boot:

```bash
sudo systemctl enable --now aperture-agent
```

:::caution

Currently, configuration watcher and automatic reload aren't supported. If you
modify the configuration file, do restart the service:

```bash
sudo systemctl restart aperture-agent
```

:::

You can then view service status:

```bash
sudo systemctl status aperture-agent
```

To view the logs, when the default log configuration is used, you can use
`journalctl`:

```bash
journalctl -u aperture-agent --since "15 minutes ago"
```

## Upgrade

[Download](#agent-download) the updated package and follow
[installation](#agent-installation) steps. Remember to restart the service after
installation is complete.

## Uninstall

Stop the Aperture Agent service:

```bash
sudo systemctl stop aperture-agent
```

**Optional**: Remove the agent configuration:

```bash
sudo rm /etc/aperture/aperture-agent/config/aperture-agent.yaml
```

Uninstall the package:

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg">

```bash
sudo dpkg -r aperture-agent
```

  </TabItem>

  <TabItem value="rpm" label="rpm">

```bash
sudo rpm -e aperture-agent
```

  </TabItem>

</Tabs>
