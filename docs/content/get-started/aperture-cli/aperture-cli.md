---
title: Aperture CLI
description: Aperture CLI for interacting with Aperture Seamlessly.
keywords:
  - cli
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion} from '../../apertureVersion.js';
```

The Aperture CLI is available as a binary executable for all major platforms,
the binaries can be downloaded form GitHub <a
href={`https://github.com/fluxninja/aperture/releases/tag/${apertureVersion}`}>Release
Page</a>.

## Installation

<Tabs>
<TabItem value="macOS" label="macOS">
With Homebrew:
<CodeBlock language="bash">
brew install aperturectl
</CodeBlock>
</TabItem>
<TabItem value="Linux" label="Linux">
With Homebrew:
<CodeBlock language="bash">
brew install aperturectl
</CodeBlock>
With Apt:
<CodeBlock language="bash">
apt-get install aperturectl
</CodeBlock>
With dnf:
<CodeBlock language="bash">
dnf install aperturectl
</CodeBlock>
</TabItem>
</Tabs>

## Enable shell autocompletion

To configure your shell to load `aperturectl`
[bash completions](/get-started/aperture-cli/aperturectl_completion.md) add to
your profile:

<Tabs>
<TabItem value="bash" label="bash">
<CodeBlock language="bash">
. &lt;(aperturectl completion bash)
</CodeBlock>
</TabItem>
<TabItem value="zsh" label="zsh">
<CodeBlock language="zsh">
. &lt;(aperturectl completion zsh)
</CodeBlock>
</TabItem>
<TabItem value="fish" label="fish">
<CodeBlock language="fish">
. &lt;(aperturectl completion fish)
</CodeBlock>
</TabItem>
<TabItem value="powershell" label="powershell">
<CodeBlock language="powershell">
. &lt;(aperturectl completion powershell)
</CodeBlock>
</TabItem>
</Tabs>

---

### [aperturectl](/get-started/aperture-cli/aperturectl.md)
