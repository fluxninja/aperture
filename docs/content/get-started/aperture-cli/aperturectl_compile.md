---
title: Aperturectl Compile
description: Aperturectl Compile
keywords:
  - aperturectl
  - aperturectl_compile
sidebar_position: 8
---

## aperturectl compile

Compile circuit from policy file

### Synopsis

Use this command to compile the Aperture policy circuit from a policy file to validate the circuit.
You can also generate the DOT and Mermaid graphs of the compiled policy circuit to visualize it.

```
aperturectl compile [flags]
```

### Examples

```
aperturectl compile --cr=policy-cr.yaml --mermaid --dot

aperturectl compile --policy=policy.yaml --mermaid --dot
```

### Options

```
      --cr string        Path to policy custom resource file
      --dot string       Path to store the dot file
  -h, --help             help for compile
      --mermaid string   Path to store the mermaid file
      --policy string    Path to policy file
```

### SEE ALSO

- [aperturectl](aperturectl.md) - aperturectl - CLI tool to interact with Aperture
