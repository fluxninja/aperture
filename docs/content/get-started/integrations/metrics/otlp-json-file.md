---
title: OTLP JSON File
description: Integrating OTLP JSON File Metrics
keywords:
  - otlpjsonfile
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`otlpjsonfilereceiver` extension enabled, so that
[otlpjsonfilereceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for OTLP JSON File using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    otlpjsonfile:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - otlpjsonfile
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        otlpjsonfile: [otlpjsonfilereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/otlpjsonfilereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
