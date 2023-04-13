---
title: Azure Event Hub
description: Integrating Azure Event Hub Metrics
keywords:
  - azureeventhub
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`azureeventhubreceiver` extension enabled, so that
[azureeventhubreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for Azure Event Hub using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    azureeventhub:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - azureeventhub
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        azureeventhub: [azureeventhubreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/azureeventhubreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
