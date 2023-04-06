---
title: RabbitMQ
description: Integrating RabbitMQ Metrics
keywords:
  - rabbitmq
  - otel
  - opentelemetry
  - collector
  - metrics
sidebar_position: 1
---

First, make sure you've [built][build] agent with `rabbitmqreceiver` extension
enabled, so that [rabbitmqreceiver][receiver] is available.

Then, you can use the following [custom metrics][custom-metrics] configuration
in [agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    rabbitmq:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - rabbitmq
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        rabbitmq:
          collection_interval: 1s
          endpoint: http://rabbitmq.rabbitmq.svc.cluster.local:15672
          password: secretpassword
          username: admin
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/rabbitmqreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#custom-metrics-config
