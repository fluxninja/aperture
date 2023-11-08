---
sidebar_label: Dashboard
hide_title: true
keywords:
  - aperturectl
  - aperturectl_dashboard
---

<!-- markdownlint-disable -->

## aperturectl dashboard

Generate dashboards for Aperture

### Synopsis

Generate dashboards for Aperture

```
aperturectl dashboard [flags]
```

### Options

```
      --datasource-name string   Name of the datasource to use (default "controller-prometheus")
  -h, --help                     help for dashboard
      --output-dir string        Output directory for the generated dashboards (default "dashboards")
      --overwrite                Overwrite the output directory if it already exists
      --policy-file string       Path to the policy file to use
      --skip-pull                Skip pulling the dashboards from the remote repository
      --uri string               URI of Custom Dashboards, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/dashboards/grafana@latest. This field should not be provided when the Version is provided.
      --version string           Version of official Aperture Dashboards, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
