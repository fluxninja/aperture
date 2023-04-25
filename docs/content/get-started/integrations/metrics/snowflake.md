---
title: Snowflake
description: Integrating Snowflake Metrics
keywords:
  - snowflake
  - otel
  - opentelemetry
  - collector
  - metrics
---

::: info

See also [snowflakereceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

::: note

The `snowflakereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/snowflakereceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for Snowflake using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    snowflake:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - snowflake
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        snowflake: [snowflakereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/snowflakereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
