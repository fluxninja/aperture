---
title: Simple Prometheus
description: Integrating Simple Prometheus Metrics
keywords:
  - prometheus_simple
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`simpleprometheusreceiver` extension enabled, so that
[simpleprometheusreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for Simple Prometheus using
the following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    prometheus_simple:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - prometheus_simple
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        prometheus_simple: [simpleprometheusreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/simpleprometheusreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
