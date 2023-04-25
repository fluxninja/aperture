---
title: Apache Web Server
description: Integrating Apache Web Server Metrics
keywords:
  - apache
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [apachereceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

::: note

The `apachereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/apachereceiver` to the `bundled_extensions` list to make [the
receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for Apache Web Server using
the following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    apache:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - apache
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        apache: [apachereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/apachereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
