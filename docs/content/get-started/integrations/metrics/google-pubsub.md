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

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`googlecloudpubsubreceiver` extension enabled, so that
[googlecloudpubsubreceiver][receiver] is available.

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
