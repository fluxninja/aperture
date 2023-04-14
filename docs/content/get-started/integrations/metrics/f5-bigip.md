---
title: F5 Big-IP
description: Integrating F5 Big-IP Metrics
keywords:
  - f5
  - big-ip
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`bigipreceiver` extension enabled, so that [bigipreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for F5 Big-IP using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    bigip:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - bigip
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        bigip: [bigipreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/bigipreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
