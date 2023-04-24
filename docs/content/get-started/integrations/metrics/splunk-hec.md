---
title: Splunk HEC
description: Integrating Splunk HEC Metrics
keywords:
  - splunk_hec
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [splunkhecreceiver docs][receiver] in opentelemetry-collect-contrib repo.

:::

::: note

The splunkhecreceiver extension is available in default agent image, but if you're [building][build] your own Aperture Agent, make sure to add `integrations/otel/splunkhecreceiver` to `bundled_extensions` list.

:::

You can configure [Custom metrics][custom-metrics] for Splunk HEC using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    splunk_hec:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - splunk_hec
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        splunk_hec: [splunkhecreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/splunkhecreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
