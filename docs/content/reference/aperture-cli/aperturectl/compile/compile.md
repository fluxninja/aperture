---
sidebar_label: Compile
hide_title: true
keywords:
  - aperturectl
  - aperturectl_compile
---

<!-- markdownlint-disable -->

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
      --depth int        Maximum depth to expand the graph. Use -1 for maximum possible depth (default 1)
      --dot string       Path to store the dot file
      --graph string     Path to store the graph file
  -h, --help             help for compile
      --mermaid string   Path to store the mermaid file
      --policy string    Path to Aperture Policy file
      --tree string      Path to store the graph file
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
