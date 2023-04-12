---
title: HTTP Check
description: Integrating HTTP Check Metrics
keywords:
  - httpcheck
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`httpcheckreceiver` extension enabled, so that [httpcheckreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for HTTP Check using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    httpcheck:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - httpcheck
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        httpcheck: [httpcheckreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/httpcheckreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
