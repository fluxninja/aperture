---
title: SSH Check
description: Integrating SSH Check Metrics
keywords:
  - sshcheck
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [sshcheckreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `sshcheckreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/sshcheckreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for SSH Check using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    sshcheck:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - sshcheck
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        sshcheck: [sshcheckreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/sshcheckreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
