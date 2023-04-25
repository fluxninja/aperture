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

::: info

See also [dockerstatsreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

::: note

The `dockerstatsreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/dockerstatsreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

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
