---
title: OpenCensus
description: Integrating OpenCensus Metrics
keywords:
  - opencensus
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [opencensusreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `opencensusreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/opencensusreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for OpenCensus using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    opencensus:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - opencensus
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        opencensus: [opencensusreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/opencensusreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
