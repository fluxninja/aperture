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

Use this command to generate Aperture Policy related resources like Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.

```
aperturectl blueprints generate [flags]
```

### Examples

```
aperturectl blueprints generate --name=rate-limiting/base --values-file=rate-limiting.yaml

aperturectl blueprints generate --name=rate-limiting/base --values-file=rate-limiting.yaml --apply
```

### Options

```
      --api-key string         FluxNinja API Key to be used when using Cloud Controller
      --apply                  Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed
      --config string          Path to the Aperture config file
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -f, --force                  Force apply policy even if it already exists
      --graph-depth int        Max depth of the graph when generating DOT and Mermaid files (default 1)
  -h, --help                   help for generate
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --name string            Name of the Aperture Blueprint to generate Aperture Policy resources for
      --no-validation          Do not validate values.yaml file
      --no-yaml-modeline       Do not add YAML language server modeline to generated YAML files
      --output-dir string      Directory path where the generated Policy resources will be stored. If not provided, will use current directory
      --overwrite              Overwrite existing output directory
  -s, --select-all             Apply all the generated Policies
      --skip-verify            Skip TLS certificate verification while connecting to controller
      --values-file string     Path to the values file for Blueprint's input
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the blueprints update.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
