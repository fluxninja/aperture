---
title: SQL Query
description: Integrating SQL Query Metrics
keywords:
  - sqlquery
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [sqlqueryreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `sqlqueryreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/sqlqueryreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for SQL
Query as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          sqlquery:
            per_agent_group: true
            receivers:
              sqlquery: [sqlqueryreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/sqlqueryreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
