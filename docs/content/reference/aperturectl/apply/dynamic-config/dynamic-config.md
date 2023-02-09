---
sidebar_label: Dynamic-Config
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply_dynamic-config
---

## aperturectl apply dynamic-config

Apply Aperture DynamicConfig to a Policy

### Synopsis

Use this command to apply the Aperture DynamicConfig to a Policy.

```
aperturectl apply dynamic-config [flags]
```

### Examples

```
aperturectl apply dynamic-config --name=static-rate-limiting --file=dynamic-config.yaml
```

### Options

```
      --file string     Path to the dynamic config file
  -h, --help            help for dynamic-config
      --policy string   Name of the Policy to apply the DynamicConfig to
```

### SEE ALSO

- [aperturectl apply](/reference/aperturectl/apply/apply.md) - Apply Aperture Policy to the cluster
