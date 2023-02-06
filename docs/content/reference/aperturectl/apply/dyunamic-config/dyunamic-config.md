---
sidebar_label: Dyunamic-Config
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply_dyunamic-config
---

## aperturectl apply dyunamic-config

Apply Aperture DynamicConfig to a Policy

### Synopsis

Use this command to apply the Aperture DynamicConfig to a Policy.

```
aperturectl apply dyunamic-config [flags]
```

### Examples

```
aperturectl apply dynamic-config --policy=static-rate-limiting
```

### Options

```
      --file string     Path to the dynamic config file
  -h, --help            help for dyunamic-config
      --policy string   Name of the Policy to apply the DynamicConfig to
```

### SEE ALSO

- [aperturectl apply](/reference/aperturectl/apply/apply.md) - Apply Aperture Policy to the cluster
