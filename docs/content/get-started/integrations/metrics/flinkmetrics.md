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

::: info

See also [flinkmetricsreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

::: note

The `flinkmetricsreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/flinkmetricsreceiver` to the `bundled_extensions` list to
make [the receiver][receiver] available.

:::

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
