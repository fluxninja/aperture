---
title: Configuration
keywords:
  - cli
sidebar_position: 2
---

Configure aperturectl to point to your Aperture Cloud endpoint: Save the
following as `~/.aperturectl/config`:

```toml
[controller]
url = "ORGANIZATION_NAME.app.fluxninja.com:443"
project_name = "PROJECT_NAME"
access_token = "PERSONAL_ACCESS_TOKEN"
```

Replace `ORGANIZATION_NAME` with the Aperture Cloud organization name and
`PERSONAL_ACCESS_TOKEN` with the Personal Access Token linked to the user. If a
Personal Access Token has not been created, generate a new one through the
Aperture Cloud UI. Refer to [Personal Access Tokens][access-tokens] for
step-by-step instructions.

:::info

See also [aperturectl configuration file format reference][aperturectl-config].

:::

:::note Self-hosted Aperture Controller

With a [self-hosted][self-hosted] Aperture Controller, if the Controller is at
the cluster pointed at by `~/.kube/config` or `KUBECONFIG`, no configuration
file nor flags are needed at all. Otherwise, you need the `--controller` flag.

:::

[self-hosted]: /aperture-for-infra/aperture-for-infra.md
[aperturectl-config]: /reference/configuration/aperturectl.md
[access-tokens]: /reference/cloud-ui/personal-access-tokens.md
