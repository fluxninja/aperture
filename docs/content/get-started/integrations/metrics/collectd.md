---
title: CollectD write_http plugin JSON
description: Integrating CollectD write_http plugin JSON Metrics
keywords:
  - collectd
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`collectdreceiver` extension enabled, so that [collectdreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for CollectD write_http
plugin JSON using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    collectd:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - collectd
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        collectd: [collectdreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/collectdreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
