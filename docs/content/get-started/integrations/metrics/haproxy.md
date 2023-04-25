---
title: HAProxy
description: Integrating HAProxy Metrics
keywords:
  - haproxy
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [haproxyreceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

:::note

The `haproxyreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/haproxyreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for HAProxy using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    haproxy:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - haproxy
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        haproxy: [haproxyreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/haproxyreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
