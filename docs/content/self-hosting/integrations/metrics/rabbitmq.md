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

:::info

See also [rabbitmqreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `rabbitmqreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/rabbitmqreceiver` to the `bundled_extensions` list to make
the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
RabbitMQ as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      rabbitmq:
        agent_group: default
        per_agent_group: true
        receivers:
          rabbitmq: [rabbitmqreceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/rabbitmqreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
