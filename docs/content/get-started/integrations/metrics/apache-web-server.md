---
title: Apache Web Server
description: Integrating Apache Web Server Metrics
keywords:
  - apache
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`apachereceiver` extension enabled, so that [apachereceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for Apache Web Server using
the following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    apache:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - apache
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        apache: [apachereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/apachereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
