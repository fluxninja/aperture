---
title: Metrics
description: Integrating custom metrics pipelines
keywords:
  - setup
  - otel
  - opentelemetry
  - collector
sidebar_position: 5
---

Aperture enables feeding custom metrics to the Aperture Controller's Prometheus
instance by adding custom metric receivers into the OpenTelemetry Collector in
the Aperture Policy. See configuration [reference][config].

## Configuring OpenTelemetry Collectors

For all the OpenTelemetry collectors, configuration can be passed using
Environment variables in the Aperture Agent.

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
            receivers:
              rabbitmq:
                endpoint: ${RABBITMQ_ENDPOINT}
                username: ${RABBITMQ_USERNAME}
                password: ${RABBITMQ_PASSWORD}
                collection_interval: 1s
```

If you are installing the Aperture Agent on Kubernetes, you can use a Secret and
a ConfigMap to pass environment variables to the Aperture Agent, as shown below:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: rabbitmq
  namespace: aperture-agent
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

To use these Secret and ConfigMap during the
[Aperture Agent Installation](/get-started/installation/agent/agent.md#agent-installation-modes),
refer to them in the values.yaml file, as shown below:

```yaml
agent:
  extraEnvVarsSecret: "rabbitmq-creds"
  extraEnvVarsCM: "rabbitmq"
```

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />

[config]: /reference/policies/spec.md#resources
[rabbitmq]: /integrations/metrics/rabbitmq.md
