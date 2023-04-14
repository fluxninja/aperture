---
sidebar_label: Zsh
hide_title: true
keywords:
  - aperturectl
  - aperturectl_completion_zsh
---

## aperturectl completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion isn't already enabled in your environment you will need to
enable it. You can execute the following once:

    echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

    source <(aperturectl completion zsh); compdef _aperturectl aperturectl

To load completions for every new session, execute once:

#### Linux:

    aperturectl completion zsh > "${fpath[1]}/_aperturectl"

#### macOS:

    aperturectl completion zsh > $(brew --prefix)/share/zsh/site-functions/_aperturectl

You will need to start a new shell for this setup to take effect.

```
aperturectl completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

- [aperturectl completion](/reference/aperturectl/completion/completion.md) -
  Generate the autocompletion script for the specified shell
