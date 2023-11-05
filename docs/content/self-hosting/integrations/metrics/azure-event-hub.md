---
title: Microsoft Azure Event Hub
description: Integrating Microsoft Azure Event Hub Metrics
keywords:
  - azureeventhub
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [azureeventhubreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `azureeventhubreceiver` extension is available in the default agent image.
If you're [building][build] your own Aperture Agent, add
`integrations/otel/azureeventhubreceiver` to the `bundled_extensions` list to
make the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Microsoft Azure Event Hub as part of [Policy resources][policy-resources] while
applying the policy:

```yaml
policy:
  resources:
    infra_meters:
      azureeventhub:
        agent_group: default
        per_agent_group: true
        receivers:
          azureeventhub: [azureeventhubreceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/azureeventhubreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
