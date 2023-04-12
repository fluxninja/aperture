---
title: Snowflake
description: Integrating Snowflake Metrics
keywords:
  - snowflake
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`snowflakereceiver` extension enabled, so that [snowflakereceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for Snowflake using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    snowflake:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - snowflake
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        snowflake: [snowflakereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/snowflakereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
