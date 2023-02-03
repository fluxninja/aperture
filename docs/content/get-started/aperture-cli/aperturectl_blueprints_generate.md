---
title: Aperturectl Blueprints Generate
description: Aperturectl Blueprints Generate
keywords:
  - aperturectl
  - aperturectl_blueprints_generate
sidebar_position: 13
---

## aperturectl blueprints generate

Generate Aperture Policy related resources from Aperture Blueprints

### Synopsis

Use this command to generate Aperture Policy related resources like Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.

```
aperturectl blueprints generate [flags]
```

### Examples

```
aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml

aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml --version v0.22.0

aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml --apply

aperturectl blueprints generate --custom-blueprint-path=/path/to/blueprint/ --values-file=values.yaml
```

### Options

```
      --apply                          Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed
      --custom-blueprint-path string   Path to the directory containing custom Blueprints which has 'config.libsonnet' and 'bundle.libsonnet' files
  -h, --help                           help for generate
      --kube-config string             Path to the Kubernets cluster config. Defaults to '~/.kube/config'
      --name string                    Name of the Aperture Blueprint to generate Aperture Policy resources for. Can be skipped when '--custom-blueprint-path' is provided
      --output-dir string              Directory path where the generated Policy resources will be stored. If not provided, will use current directory
      --values-file string             Path to the values file for Blueprint's input
```

### Options inherited from parent commands

```
      --uri string       URI of Aperture Blueprints, could be a local path or a remote git repository (default "github.com/fluxninja/aperture/blueprints")
      --version string   Version of Aperture Blueprints (default "main")
```

### SEE ALSO

- [aperturectl blueprints](aperturectl_blueprints.md) - Aperture Blueprints
