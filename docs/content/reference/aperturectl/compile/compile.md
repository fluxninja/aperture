---
sidebar_label: Compile
hide_title: true
keywords:
  - aperturectl
  - aperturectl_compile
---

## aperturectl compile

Compile circuit from Aperture Policy file

### Synopsis

Use this command to compile the Aperture Policy circuit from a file to validate the circuit.
You can also generate the DOT and Mermaid graphs of the compiled Aperture Policy circuit to visualize it.

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
      --cr string        Path to Aperture Policy custom resource file
      --dot string       Path to store the dot file
  -h, --help             help for compile
      --mermaid string   Path to store the mermaid file
      --policy string    Path to Aperture Policy file
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
