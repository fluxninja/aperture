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

Replace `ORGANIZATION_NAME` with the Aperture Cloud organization name and
`API_KEY` with the API key linked to the project. If an API key has not been
created, generate a new one through the Aperture Cloud UI. Refer to [API
Keys][api-keys] for additional information.

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
