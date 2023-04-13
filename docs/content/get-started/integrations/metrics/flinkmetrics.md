---
title: FlinkMetrics
description: Integrating FlinkMetrics Metrics
keywords:
  - flinkmetrics
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`flinkmetricsreceiver` extension enabled, so that
[flinkmetricsreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for FlinkMetrics using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    flinkmetrics:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - flinkmetrics
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        flinkmetrics: [flinkmetricsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/flinkmetricsreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
