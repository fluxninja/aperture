---
title: kubelet metrics
description: Integrating kubelet metrics
keywords:
  - kubeletstats
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [kubeletstatsreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

::: note

The `kubeletstatsreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/kubeletstatsreceiver` to the `bundled_extensions` list to
make [the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for kubelet using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    kubeletstats:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - kubeletstats
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        kubeletstats: [kubeletstatsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kubeletstatsreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
