---
title: Oracle DB
description: Integrating Oracle DB Metrics
keywords:
  - oracledb
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [oracledbreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `oracledbreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/oracledbreceiver` to the `bundled_extensions` list to make
the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Oracle DB as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      oracledb:
        agent_group: default
        per_agent_group: true
        receivers:
          oracledb: [oracledbreceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/oracledbreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
