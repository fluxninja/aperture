---
title: Microsoft SQL Server
description: Integrating Microsoft SQL Server Metrics
keywords:
  - sqlserver
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [sqlserverreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `sqlserverreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/sqlserverreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Microsoft SQL Server as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          sqlserver:
            per_agent_group: true
            receivers:
              sqlserver: [sqlserverreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/sqlserverreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
