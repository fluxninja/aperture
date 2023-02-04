---
sidebar_label: Generate
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints_generate
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
      --apply                Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed
  -h, --help                 help for generate
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --name string          Name of the Aperture Blueprint to generate Aperture Policy resources for. Can be skipped when '--custom-blueprint-path' is provided
      --output-dir string    Directory path where the generated Policy resources will be stored. If not provided, will use current directory
      --values-file string   Path to the values file for Blueprint's input
```

### Options inherited from parent commands

```
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@main. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
