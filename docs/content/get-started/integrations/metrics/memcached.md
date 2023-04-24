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

::: info

See also [memcachedreceiver docs][receiver] in opentelemetry-collect-contrib repo.

:::

::: note

The memcachedreceiver extension is available in default agent image, but if you're [building][build] your own Aperture Agent, make sure to add `integrations/otel/memcachedreceiver` to `bundled_extensions` list.

:::

You can configure [Custom metrics][custom-metrics] for Memcached using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
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
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
