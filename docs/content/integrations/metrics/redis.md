---
title: Redis
description: Integrating Redis Metrics
keywords:
  - redis
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [redisreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `redisreceiver` extension is available in the default agent image. If you're
[building][build] your own Aperture Agent, add `integrations/otel/redisreceiver`
to the `bundled_extensions` list to make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Redis as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          redis:
            per_agent_group: true
            pipeline:
              processors:
                - batch
              receivers:
                - redis
            processors:
              batch:
                send_batch_size: 10
                timeout: 10s
            receivers:
              redis: [redisreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/redisreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
