---
title: Google Cloud Pub/Sub
description: Integrating Google Cloud Pub/Sub Metrics
keywords:
  - googlecloudpubsub
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [googlecloudpubsubreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `googlecloudpubsubreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/googlecloudpubsubreceiver` to the `bundled_extensions` list
to make [the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for Google Cloud Pub/Sub
using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    googlecloudpubsub:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - googlecloudpubsub
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        googlecloudpubsub: [googlecloudpubsubreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/googlecloudpubsubreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
