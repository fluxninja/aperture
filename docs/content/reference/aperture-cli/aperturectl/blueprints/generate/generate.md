---
sidebar_label: Generate
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints_generate
---

<!-- markdownlint-disable -->

## aperturectl blueprints generate

Generate Aperture Policy related resources from Aperture Blueprints

### Synopsis

Use this command to generate Aperture Policy related resources such as Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.

```
aperturectl blueprints generate [flags]
```

### Examples

```
aperturectl blueprints generate --values-file=rate-limiting.yaml

aperturectl blueprints generate --values-file=rate-limiting.yaml --apply
```

### Options

```
      --apply                  Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -f, --force                  Force apply policy even if it already exists
      --graph-depth int        Max depth of the graph when generating DOT and Mermaid files (default 1)
  -h, --help                   help for generate
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --no-validation          Do not validate values.yaml file
      --no-yaml-modeline       Do not add YAML language server modeline to generated YAML files
      --output-dir string      Directory path where the generated Policy resources will be stored. If not provided, will use current directory
      --overwrite              Overwrite existing output directory
      --select-all             Select all blueprints
      --skip-verify            Skip TLS certificate verification while connecting to controller
      --values-dir string      Directory path to the values file(s)
      --values-file string     Path to the values file
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the blueprints update.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperture-cli/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
