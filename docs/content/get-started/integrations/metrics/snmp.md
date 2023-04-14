---
title: SNMP
description: Integrating SNMP Metrics
keywords:
  - snmp
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`snmpreceiver` extension enabled, so that [snmpreceiver][receiver] is available.

You can configure [Custom metrics][custom-metrics] for SNMP using the following
configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    snmp:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - snmp
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        snmp: [snmpreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/snmpreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
