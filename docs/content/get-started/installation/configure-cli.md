---
title: Configure CLI
keywords:
  - cli
sidebar_position: 3
---

Configure aperturectl to point to your Aperture Cloud endpoint: Save the
following as `~/.aperturectl/config`:

```toml
[controller]
url = "ORGANIZATION_NAME.app.fluxninja.com:443"
api_key = "API_KEY"
```

Replace the `ORGANIZATION_NAME` and `API_KEY` with your Aperture Cloud
organization name and API key associated with your project. If you don't know
the API key, you can create one in the Aperture Cloud UI. See [API
Keys][api-keys] for more details.

:::info

See also [aperturectl configuration file format reference][aperturectl-config].

:::

:::note Self-hosted Aperture Controller

With a [self-hosted][self-hosted] Aperture Controller, if the Controller is at
the cluster pointed at by `~/.kube/config` or `KUBECONFIG`, no configuration
file nor flags are needed at all. Otherwise, you need the `--controller` flag.

:::

[self-hosted]: /self-hosting/self-hosting.md
[aperturectl-config]: /reference/configuration/aperturectl.md
[api-keys]: /get-started/aperture-cloud/api-keys.md
