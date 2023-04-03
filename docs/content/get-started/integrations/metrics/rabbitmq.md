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
      receivers:
        rabbitmq:
          endpoint: http://<rabbitmq-service-address>:15672
          username: <username>
          password: <password>
          collection_interval: 1s
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/rabbitmqreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#custom-metrics-config
