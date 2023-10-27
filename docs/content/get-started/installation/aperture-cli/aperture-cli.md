---
title: Install CLI (aperturectl)
description: Aperture CLI for interacting with Aperture seamlessly.
keywords:
  - cli
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../../apertureVersion.js';
import {DownloadScript} from '../../self-hosting/agent/bare-metal.md';
```

```mdx-code-block
export const BinaryDownload = ({}) => (
<CodeBlock language="bash">
{`# Substitute BIN for your bin directory.
VERSION="${apertureVersionWithOutV}"
BIN="/usr/local/bin"
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64";;
  aarch64) ARCH="arm64";;
  *) echo "Unsupported architecture: $ARCH"; exit 1;;
esac
echo "Will download $OS package version $VERSION compiled for $ARCH machine"
url="https://github.com/fluxninja/aperture/releases/download/v\${VERSION}/aperturectl-$\{VERSION}-$\{OS}-$\{ARCH}"
curl --fail --location --remote-name "\${url}"
mv aperturectl* "\${BIN}/aperturectl"
chmod +x "\${BIN}/aperturectl"
`}</CodeBlock>
);
```

The Aperture CLI provides a convenient way to interact with Aperture. It can be
used to generate a policy, apply a policy, list policies, and perform other
operations such as listing connected agents, listing flow control points and
preview samples of flow labels or HTTP requests on control points. Its use cases
extend beyond the aforementioned operations. Detailed documentation is available
on the [aperturectl](/reference/aperturectl/aperturectl.md) reference page.

The Aperture CLI, compatible with all major platforms, is accessible as an
executable program. Binaries for the CLI can be obtained from the GitHub <a
href={`https://github.com/fluxninja/aperture/releases/tag/${apertureVersion}`}>Release
Page</a>.

Alternatively, it can also be downloaded using the following script:

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg">
    <DownloadScript packager="deb" arch="amd64" archSeparator="_" versionSeparator="_" component="aperturectl" />
  </TabItem>
  <TabItem value="rpm" label="rpm">
    <DownloadScript packager="rpm" arch="x86_64" archSeparator="." versionSeparator="-" component="aperturectl" />
  </TabItem>
  <TabItem value="binary" label="binary">
    <BinaryDownload  />
  </TabItem>
</Tabs>

## Installation

:::info

Skip the following steps if you have obtained the binary file directly using the
steps mentioned above.

:::

<!-- vale off -->

<Tabs groupId="setup" queryString>
<TabItem value="macOS" label="macOS">
With Homebrew:
<CodeBlock language="bash">
brew install fluxninja/aperture/aperturectl
</CodeBlock>
</TabItem>
<TabItem value="Linux" label="Linux">
With Homebrew:
<CodeBlock language="bash">
brew install fluxninja/aperture/aperturectl
</CodeBlock>
With dpkg:
<CodeBlock language="bash">
{`sudo dpkg -i aperturectl_${apertureVersionWithOutV}*.deb`}
</CodeBlock>
With rpm:
<CodeBlock language="bash">
{`sudo rpm -i aperturectl-${apertureVersionWithOutV}*.rpm`}
</CodeBlock>
</TabItem>
</Tabs>

<!-- vale on -->

## Enable shell autocompletion

To configure your shell to load `aperturectl`
[bash completions](/reference/aperturectl/completion/completion.md), add to your
profile:

<Tabs>
<TabItem value="bash" label="bash">
<CodeBlock language="bash">
source &lt;(aperturectl completion bash)
</CodeBlock>
</TabItem>
<TabItem value="zsh" label="zsh">
<CodeBlock language="zsh">
source &lt;(aperturectl completion zsh); compdef _aperturectl aperturectl
</CodeBlock>
</TabItem>
<TabItem value="fish" label="fish">
<CodeBlock language="fish">
aperturectl completion fish | source
</CodeBlock>
</TabItem>
<TabItem value="powershell" label="powershell">
<CodeBlock language="powershell">
aperturectl completion powershell | Out-String | Invoke-Expression
</CodeBlock>
</TabItem>
</Tabs>

## Uninstall

<!-- vale off -->
<Tabs groupId="setup" queryString>
<TabItem value="macOS" label="macOS">
With Homebrew:
<CodeBlock language="bash">
{`brew uninstall aperturectl
brew untap fluxninja/aperture`}
</CodeBlock>
</TabItem>
<TabItem value="Linux" label="Linux">
With Homebrew:
<CodeBlock language="bash">
{`brew uninstall aperturectl
brew untap fluxninja/aperture`}
</CodeBlock>
With dpkg:
<CodeBlock language="bash">
sudo dpkg -r aperturectl
</CodeBlock>
With rpm:
<CodeBlock language="bash">
sudo rpm -e aperturectl
</CodeBlock>
</TabItem></Tabs>

<!-- vale on -->

---

## [aperturectl](/reference/aperturectl/aperturectl.md)
