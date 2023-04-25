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

:::info

See also [purefbreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

:::note

The `purefbreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/purefbreceiver` to the `bundled_extensions` list to make [the
receiver][receiver] available.

:::

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
