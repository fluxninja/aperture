---
title: JMX
description: Integrating JMX Metrics
keywords:
  - jmx
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`jmxreceiver` extension enabled, so that [jmxreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for JMX using the following
configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    jmx:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - jmx
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        jmx: [jmxreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/jmxreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
