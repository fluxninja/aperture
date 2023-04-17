---
title: Pure Storage FlashBlade
description: Integrating Pure Storage FlashBlade Metrics
keywords:
  - purefb
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`purefbreceiver` extension enabled, so that [purefbreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for Pure Storage FlashBlade
using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    purefb:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - purefb
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        purefb: [purefbreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/purefbreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
