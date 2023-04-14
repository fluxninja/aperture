---
title: RabbitMQ
description: Integrating RabbitMQ Metrics
keywords:
  - rabbitmq
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`rabbitmqreceiver` extension enabled, so that [rabbitmqreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for RabbitMQ using the
following configuration in the [Aperture Agent's
configuration][agent-configuration]:

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
          endpoint: http://<rabbitmq-service-address>:15672
          username: <username>
          password: <password>
          collection_interval: 1s
```

You can also pass these configurations using environment variables as shown
below:

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
          endpoint: ${RABBITMQ_ENDPOINT}
          username: ${RABBITMQ_USERNAME}
          password: ${RABBITMQ_PASSWORD}
          collection_interval: 1s
```

If you're installing the Aperture Agent on Kubernetes, you can use a Secret and
a ConfigMap to pass environment variables to the Aperture Agent, as shown below:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: rabbitmq
data:
  RABBITMQ_ENDPOINT: http://rabbitmq.rabbitmq.svc.cluster.local:15672
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: rabbitmq-creds
  namespace: aperture-agent
data:
  RABBITMQ_USERNAME: <rabbitmq-username>
  RABBITMQ_PASSWORD: <rabbitmq-password>
```

To use these secrets and ConfigMap during the
[Aperture Agent Installation](/get-started/installation/agent/agent.md#agent-installation-modes),
refer to them in the values.yaml file, as shown below:

```yaml
agent:
  config:
    etcd:
      endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
    prometheus:
      address: "http://controller-prometheus-server.default.svc.cluster.local:80"
    agent_functions:
      endpoints: ["aperture-controller.default.svc.cluster.local:8080"]
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
              endpoint: ${RABBITMQ_ENDPOINT}
              username: ${RABBITMQ_USERNAME}
              password: ${RABBITMQ_PASSWORD}
              collection_interval: 1s
  extraEnvVarsSecret: "rabbitmq-creds"
  extraEnvVarsCM: "rabbitmq"
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/rabbitmqreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-configuration]: /reference/configuration/agent.md#custom-metrics-config
