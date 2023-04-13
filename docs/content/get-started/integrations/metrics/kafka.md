---
title: Kafka Metrics
description: Integrating Kafka Metrics Metrics
keywords:
  - kafkametrics
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`kafkametricsreceiver` extension enabled, so that
[kafkametricsreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for Kafka Metrics using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    kafkametrics:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - kafkametrics
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        kafkametrics: [kafkametricsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkametricsreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
