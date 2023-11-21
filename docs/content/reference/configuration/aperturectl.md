---
title: aperturectl Configuration File Format Reference
sidebar_position: 4
sidebar_label: "Aperturectl"
---

<!-- If our configuration file grows, would be nice to automatically generate
it from corresponding go structs from cmd/aperturectl/cmd/utils/controller.go -->

## Location

To avoid specifying `--controller`, `--api-key` and `--project-name` in every `aperturectl`
invocation, aperturectl can use a configuration file located in
`~/.aperturectl/config`.

The location of this file can be overridden by the `APERTURE_CONFIG` environment
variable and `--config` option (with the command-line option having higher
precedence).

When any explicit flag related to controller location (e.g., `--kube`,
`--controller`, `--api-key` or `--project-name`) is used, the value from the
configuration file is ignored for the flag.

If the configuration file is not specified nor present at the default location,
aperturectl will try to find the controller at the local Kubernetes cluster (as
if the `--kube` flag were passed).

## Format

The aperturectl configuration file uses the following [TOML][toml] syntax:

```toml
[controller]
url = "ORGANIZATION_NAME.app.fluxninja.com:443"
project_name = "PROJECT_NAME"
api_key = "PERSONAL_API_KEY"
```

Replace `ORGANIZATION_NAME` with the Aperture Cloud organization name and
`PERSONAL_API_KEY` with the Personal API key linked to the user. If a Personal
API key has not been created, generate a new one through the Aperture Cloud UI.
Refer to [Personal API Keys][api-keys] for additional information.

All the fields are required (although the file itself is not). See [Configuring
aperturectl][configure-aperturectl] for an example on how to configure
aperturectl with [Aperture Cloud Controller][cloud-controller].

:::tip

You can create multiple configuration files and use `APERTURE_CONFIG`
environment variable to switch between different projects and organizations.

:::

[toml]: https://toml.io/
[configure-aperturectl]: /reference/aperture-cli/configure-cli.md
[cloud-controller]: /reference/fluxninja.md#cloud-controller
[api-keys]: /reference/aperture-cli/personal-access-tokens.md
