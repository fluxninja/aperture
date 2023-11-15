---
title: PostgreSQL
description: Integrating PostgreSQL Metrics
keywords:
  - postgresql
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [postgresqlreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `postgresqlreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/postgresqlreceiver` to the `bundled_extensions` list to make
the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
PostgreSQL as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      postgresql:
        agent_group: default
        per_agent_group: true
        receivers:
          postgresql: [postgresqlreceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/postgresqlreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
