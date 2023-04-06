---
title: RabbitMQ Queue Buildup Policy
---

## Introduction

This policy detects RabbitMQ queue buildup by looking at number of messages in
"ready" state. Gradient controller is then used to calculate a proportional
response that limits the accepted concurrency. Concurrency is increased
additively when the overload is no longer detected.

## Build Instructions

Make sure that the aperture binary is [built](reference/aperturectl/build/agent/) with `rabbitmqreceiver` extension enabled.

```yaml
bundled_extensions:
  - integrations/otel/rabbitmqreceiver
```

Use the following agent configuration to allow agent to collect RabbitMQ metrics.

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

## Configuration

<!-- Configuration Marker -->
