---
title: Google Cloud Spanner
description: Integrating Google Cloud Spanner Metrics
keywords:
  - googlecloudspanner
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [googlecloudspannerreceiver docs][receiver] in opentelemetry-collect-contrib repo.

:::

::: note

The googlecloudspannerreceiver extension is available in default agent image, but if you're [building][build] your own Aperture Agent, make sure to add `integrations/otel/googlecloudspannerreceiver` to `bundled_extensions` list.

:::

You can configure [Custom metrics][custom-metrics] for Google Cloud Spanner
using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    googlecloudspanner:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - googlecloudspanner
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        googlecloudspanner: [googlecloudspannerreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/googlecloudspannerreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
