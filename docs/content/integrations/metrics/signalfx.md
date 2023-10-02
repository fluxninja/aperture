---
title: SignalFx
description: Integrating SignalFx Metrics
keywords:
  - signalfx
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [signalfxreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `signalfxreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/signalfxreceiver` to the `bundled_extensions` list to make
the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
SignalFx as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      signalfx:
        agent_group: default
        per_agent_group: true
        receivers:
          signalfx: [signalfxreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/signalfxreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
