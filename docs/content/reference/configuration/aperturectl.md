---
title: aperturectl Configuration File Format Reference
sidebar_position: 5
sidebar_label: .aperturectl/config
---

<!-- If our config grows, would be nice to automatically generate it from
corresponding go structs from cmd/aperturectl/cmd/utils/controller.go -->

## Location

To avoid specifying `--controller` and `--api-key` in every `aperturectl` invocation,
aperturectl can use a configuration file located in `~/.aperturectl/config`.

Location of this file can be overridden by `APERTURE_CONFIG` environment
variable and `--config` option (with the commandline option having higher
precedence).

Config file is ignored when any explicit flag related to controller
location (such as `--kube`, `--controller` or `--api-key`) is passed.

If the config file is not specified nor present at the default location,
aperturectl will try to find the controller at the local kubernetes cluster (as
if `--kube` flag was passed).

## Format

The aperturectl configuration file uses following [TOML][] syntax:

```toml
[controller]
url = "controller hostname:port"
api_key = "api key for the controller"
```

All the fields are required. See [Configuring aperturectl][] for an example
how to configure aperturectl with [FluxNinja Cloud Controller][].

:::tip

You can create multiple configuration files and use `APERTURE_CONFIG`
environment variable to switch between different projects and organizations.

:::

[TOML]: https://toml.io/
[Configuring aperturectl]: /get-started/installation/configure-cli.md
[FluxNinja Cloud Controller]: /reference/fluxninja.md#cloud-controller
