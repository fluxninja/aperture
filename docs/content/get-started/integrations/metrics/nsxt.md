---
title: NSX-T
description: Integrating NSX-T Metrics
keywords:
  - nsx-t
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [nsxtreceiver docs][receiver] in opentelemetry-collect-contrib repo.

:::

::: note

The nsxtreceiver extension is available in default agent image, but if you're [building][build] your own Aperture Agent, make sure to add `integrations/otel/nsxtreceiver` to `bundled_extensions` list.

:::

You can configure [Custom metrics][custom-metrics] for NSX-T using the following
configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    nsxt:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - nsxt
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        nsxt: [nsxtreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/nsxtreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
