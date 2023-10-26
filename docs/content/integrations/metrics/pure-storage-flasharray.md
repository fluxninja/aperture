---
title: Pure Storage FlashArray
description: Integrating Pure Storage FlashArray Metrics
keywords:
  - purefa
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [purefareceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `purefareceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/purefareceiver` to the `bundled_extensions` list to make the
[receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Pure Storage FlashArray as part of [Policy resources][policy-resources] while
applying the policy:

```yaml
policy:
  resources:
    infra_meters:
      purefa:
        agent_group: default
        per_agent_group: true
        receivers:
          purefa: [purefareceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/purefareceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
