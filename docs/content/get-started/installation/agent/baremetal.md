---
title: Baremetal/VM
description: Install Aperture Agent on Baremetal or VM
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
import {apertureVersion} from '../../../introduction.md';
```

The Aperture Agent can be installed as a system service.

## Download {#agent-download}

The Aperture Agent can be installed using packages made for your system's
package manager like `dpkg`<!-- or `rpm` -->.

To install agent, first download package for your manager from
[Releases Page](https://github.com/fluxninja/aperture/releases/latest).

Alternatively download it using following script:

```mdx-code-block
export const DownloadScript = ({children, packager}) => (
<CodeBlock language="bash">
{`VERSION="${apertureVersion}"
ARCH="amd64"
PACKAGER="${packager}"
url="https://github.com/fluxninja/aperture/releases/download/v\${VERSION}/aperture-agent_\${VERSION}_\${ARCH}.\${PACKAGER}"
echo "Will download \${PACKAGER} package version \${VERSION} compiled for \${ARCH} machine"
curl --fail --location --remote-name "\${url}"
`}</CodeBlock>
);
```

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg"><DownloadScript packager="deb"/></TabItem>
</Tabs>

## Installation {#agent-installation}

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg">

```bash
sudo dpkg -i "aperture-agent_${VERSION}_${ARCH}.${PACKAGER}"
```

  </TabItem>
</Tabs>

You should then point agent at etcd and prometheus deployed by the Controller,
by editing `/etc/aperture/aperture-agent/config/aperture-agent.yaml`.

:::info

The default config disables the FluxNinja Cloud Plugin for the Aperture Agent.
If you want to keep it enabled, add parameters provided
[here](/cloud/plugin.md#configuration).

:::

After installing, you should enable the `aperture-agent` systemd service, and
make it start after system boot:

```bash
sudo systemctl enable --now aperture-agent
```

:::caution

Currently configuration watcher and automatic reload doesn't work. If you modify
the configuration file, make sure to restart the service:

```bash
sudo systemctl restart aperture-agent
```

:::

You can then view service status:

```bash
systemctl status aperture-service
```

To view the logs, when default log configuration is used, you can use
`journalctl`:

```bash
journalctl -u aperture-agent --since "15 minutes ago"
```

## Upgrade

[Download](#agent-download) updated package and follow
[installation](#agent-installation) steps. Remember to restart the service after
installation is complete.

## Uninstall

1. **Optional**: Remove the agent configuration:

   ```bash
   sudo rm /etc/aperture/aperture-agent/config/aperture-agent.yaml
   ```

2. Uninstall the package:

  <Tabs groupId="packageManager" queryString>
    <TabItem value="dpkg" label="dpkg">

    ```bash
    sudo dpkg -r aperture-agent
    ```

    </TabItem>

  </Tabs>
