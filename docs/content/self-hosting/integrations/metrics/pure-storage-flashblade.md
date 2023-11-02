---
title: Pure Storage FlashBlade
description: Integrating Pure Storage FlashBlade Metrics
keywords:
  - purefb
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [purefbreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `purefbreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/purefbreceiver` to the `bundled_extensions` list to make the
[receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Pure Storage FlashBlade as part of [Policy resources][policy-resources] while
applying the policy:

```yaml
policy:
  resources:
    infra_meters:
      purefb:
        agent_group: default
        per_agent_group: true
        receivers:
          purefb: [purefbreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/purefbreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
