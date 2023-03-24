---
title: Install CLI (apeturectl)
description: Aperture CLI for interacting with Aperture Seamlessly.
keywords:
  - cli
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion,apertureVersionWithOutV} from '../../apertureVersion.js';
import {DownloadScript} from '../installation/agent/bare_metal.md';
```

The Aperture CLI is available as a binary executable for all major platforms,
the binaries can be downloaded from GitHub <a
href={`https://github.com/fluxninja/aperture/releases/tag/${apertureVersion}`}>Release
Page</a>.

Alternatively download it using following script:

<Tabs groupId="packageManager" queryString>
  <TabItem value="dpkg" label="dpkg">
    <DownloadScript packager="deb" arch="amd64" archSeparator="_" versionSeparator="_" component="aperturectl" />
  </TabItem>
  <TabItem value="rpm" label="rpm">
    <DownloadScript packager="rpm" arch="x86_64" archSeparator="." versionSeparator="-" component="aperturectl" />
  </TabItem>
</Tabs>

## Installation

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

## Enable shell autocompletion

To configure your shell to load `aperturectl`
[bash completions](/reference/aperturectl/completion/completion.md) add to your
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

---

### [aperturectl](/reference/aperturectl/aperturectl.md)
