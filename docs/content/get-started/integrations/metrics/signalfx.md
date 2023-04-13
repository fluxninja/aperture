---
title: SignalFx
description: Integrating SignalFx Metrics
keywords:
  - signalfx
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`signalfxreceiver` extension enabled, so that [signalfxreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for SignalFx using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    signalfx:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - signalfx
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        signalfx: [signalfxreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/signalfxreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
