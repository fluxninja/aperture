---
title: Microsoft Azure Event Hub
description: Integrating Azure Event Hub Metrics
keywords:
  - azureeventhub
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [azureeventhubreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

:::note

The `azureeventhubreceiver` extension is available in the default agent image.
If you're [building][build] your own Aperture Agent, add
`integrations/otel/azureeventhubreceiver` to the `bundled_extensions` list to
make [the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for Microsoft Azure Event Hub
using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    azureeventhub:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - azureeventhub
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        azureeventhub: [azureeventhubreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/azureeventhubreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
