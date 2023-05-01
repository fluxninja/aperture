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

You can configure [Custom metrics][custom-metrics] for SQL Query using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    sqlquery:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - sqlquery
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        sqlquery: [sqlqueryreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/sqlqueryreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
