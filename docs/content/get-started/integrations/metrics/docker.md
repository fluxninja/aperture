---
title: Docker Stats
description: Integrating Docker Stats Metrics
keywords:
  - docker_stats
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`dockerstatsreceiver` extension enabled, so that [dockerstatsreceiver][receiver]
is available.

You can configure [Custom metrics][custom-metrics] for Docker Stats using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    docker_stats:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - docker_stats
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        docker_stats: [dockerstatsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/dockerstatsreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
