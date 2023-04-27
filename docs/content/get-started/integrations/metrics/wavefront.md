---
title: Wavefront
description: Integrating Wavefront Metrics
keywords:
  - wavefront
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [wavefrontreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

:::note

The `wavefrontreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/wavefrontreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for Wavefront using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    wavefront:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - wavefront
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        wavefront: [wavefrontreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/wavefrontreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
