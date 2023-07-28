---
title: Configure CLI
keywords:
  - cli
sidebar_position: 2
---

Configure aperturectl to point to your FluxNinja endpoint: Save the following as
`~/.aperturectl/config`:

```toml
[controller]
url = "ORGANIZATION_NAME.app.fluxninja.com:443"
api_key = "API_KEY"
```

:::tip

You can create multiple configuration files and use `APERTURE_CONFIG`
environment variable to switch between them.

:::

:::note Self-Hosted Aperture Controller

With a [Self-Hosted][self-hosted] Aperture Controller, if the Controller is at
the cluster pointed at by `~/.kube/config` or `KUBECONFIG`, no configuration
file nor flags are needed at all. Otherwise, you need the `--controller` flag.
See [aperturectl][] reference for details.

:::

[self-hosted]: /self-hosting/self-hosting.md
[aperturectl]: /reference/aperturectl/aperturectl.md
