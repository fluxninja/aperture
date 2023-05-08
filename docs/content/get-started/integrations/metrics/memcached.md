---
title: Memcached
description: Integrating Memcached Metrics
keywords:
  - memcached
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [memcachedreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `memcachedreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/memcachedreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Memcached as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          memcached:
            per_agent_group: true
            pipeline:
              processors:
                - batch
              receivers:
                - memcached
            processors:
              batch:
                send_batch_size: 10
                timeout: 10s
            receivers:
              memcached: [memcachedreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/memcachedreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
