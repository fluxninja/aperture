---
title: Podman Stats
description: Integrating Podman Stats Metrics
keywords:
  - podman_stats
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`podmanreceiver` extension enabled, so that [podmanreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for Podman Stats using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    podman_stats:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - podman_stats
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        podman_stats: [podmanreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/podmanreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
