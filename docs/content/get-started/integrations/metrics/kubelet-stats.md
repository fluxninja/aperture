---
title: kubelet metrics
description: Integrating kubelet metrics
keywords:
  - kubeletstats
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`kubeletstatsreceiver` extension enabled, so that
[kubeletstatsreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for kubelet using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    kubeletstats:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - kubeletstats
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        kubeletstats: [kubeletstatsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kubeletstatsreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
