---
title: HTTP Check
description: Integrating HTTP Check Metrics
keywords:
  - httpcheck
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [httpcheckreceiver docs][receiver] in opentelemetry-collect-contrib repo.

:::

::: note

The httpcheckreceiver extension is available in default agent image, but if you're [building][build] your own Aperture Agent, make sure to add `integrations/otel/httpcheckreceiver` to `bundled_extensions` list.

:::

You can configure [Custom metrics][custom-metrics] for HTTP Check using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    httpcheck:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - httpcheck
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        httpcheck: [httpcheckreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/httpcheckreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
