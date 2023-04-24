---
title: Aerospike
description: Integrating Aerospike Metrics
keywords:
  - aerospike
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [aerospikereceiver docs][receiver] in opentelemetry-collect-contrib repo.

:::

::: note

The aerospikereceiver extension is available in default agent image, but if you're [building][build] your own Aperture Agent, make sure to add `integrations/otel/aerospikereceiver` to `bundled_extensions` list.

:::

You can configure [Custom metrics][custom-metrics] for Aerospike using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    aerospike:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - aerospike
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        aerospike: [aerospikereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/aerospikereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
