---
title: Metrics
description: Integrating custom metrics pipelines
keywords:
  - setup
  - otel
  - opentelemetry
  - collector
sidebar_position: 3
---

Aperture enables feeding custom metrics to the Aperture Controller's Prometheus
instance by adding custom metric receivers into the OpenTelemetry Collector in
the Aperture Policy. See configuration [reference][config].

## Configuring OpenTelemetry Collectors

For all the OpenTelemetry collectors, configuration can be passed using
Environment variables in the Aperture Controller.

For example, environment variables for [RabbitMQ][rabbitmq] can be passed as
below:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
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

You can use a Secret and a ConfigMap to pass environment variables to the
Aperture Controller, as shown below:

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
  namespace: aperture-controller
data:
  RABBITMQ_USERNAME: <rabbitmq-username>
  RABBITMQ_PASSWORD: <rabbitmq-password>
```

To use these secrets and ConfigMap during the
[Aperture Controller Installation](/get-started/installation/controller/controller.md#installation),
refer to them in the values.yaml file, as shown below:

```yaml
controller:
  extraEnvVarsSecret: "rabbitmq-creds"
  extraEnvVarsCM: "rabbitmq"
```

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />

[config]: /reference/policies/spec.md#resources
[rabbitmq]: /get-started/integrations/metrics/rabbitmq.md
