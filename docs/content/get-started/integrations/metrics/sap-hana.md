---
title: SAP HANA
description: Integrating SAP HANA Metrics
keywords:
  - saphana
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [saphanareceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `saphanareceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/saphanareceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for SAP HANA using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    saphana:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - saphana
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        saphana: [saphanareceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/saphanareceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
