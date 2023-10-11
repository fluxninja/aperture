---
title: aperturectl Configuration File Format Reference
sidebar_position: 4
sidebar_label: "Aperturectl"
---

<!-- If our configuration file grows, would be nice to automatically generate
it from corresponding go structs from cmd/aperturectl/cmd/utils/controller.go -->

## Location

To avoid specifying `--controller` and `--api-key` in every `aperturectl`
invocation, aperturectl can use a configuration file located in
`~/.aperturectl/config`.

The location of this file can be overridden by the `APERTURE_CONFIG` environment
variable and `--config` option (with the command-line option having higher
precedence).

When any explicit flag related to controller location (e.g., `--kube`,
`--controller`, or `--api-key`) is used, the _entire_ configuration file is
ignored.

If the configuration file is not specified nor present at the default location,
aperturectl will try to find the controller at the local Kubernetes cluster (as
if the `--kube` flag were passed).

## Format

The aperturectl configuration file uses the following [TOML][toml] syntax:

```toml
[controller]
url = "controller hostname:port"
api_key = "api key for the controller"
```

All the fields are required (although the file itself is not). See [Configuring
aperturectl][configure-aperturectl] for an example on how to configure
aperturectl with [Aperture Cloud Controller][cloud-controller].

:::tip

You can create multiple configuration files and use `APERTURE_CONFIG`
environment variable to switch between different projects and organizations.

:::

[toml]: https://toml.io/
[configure-aperturectl]: /get-started/installation/configure-cli.md
[cloud-controller]: /reference/fluxninja.md#cloud-controller
