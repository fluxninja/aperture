---
title: Cloud Foundry
description: Integrating Cloud Foundry Metrics
keywords:
  - cloudfoundry
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [cloudfoundryreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `cloudfoundryreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/cloudfoundryreceiver` to the `bundled_extensions` list to
make the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Cloud Foundry as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      cloudfoundry:
        agent_group: default
        per_agent_group: true
        receivers:
          cloudfoundry: [cloudfoundryreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/cloudfoundryreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
