---
sidebar_label: Bash
hide_title: true
keywords:
  - aperturectl
  - aperturectl_completion_bash
---

<!-- markdownlint-disable -->

## aperturectl completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it using your operating system's package manager.

To load completions in your current shell session:

    source <(aperturectl completion bash)

To load completions for every new session, execute once:

#### Linux:

    aperturectl completion bash > /etc/bash_completion.d/aperturectl

#### macOS:

    aperturectl completion bash > $(brew --prefix)/etc/bash_completion.d/aperturectl

You will need to start a new shell for this setup to take effect.

```
aperturectl completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

- [aperturectl completion](/reference/aperturectl/completion/completion.md) - Generate the autocompletion script for the specified shell
