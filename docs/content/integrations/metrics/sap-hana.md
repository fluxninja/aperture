---
title: SAP HANA
description: Integrating SAP HANA Metrics
keywords:
  - saphana
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [saphanareceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `saphanareceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/saphanareceiver` to the `bundled_extensions` list to make the
[receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for SAP
HANA as part of [Policy resources][policy-resources] while applying the policy:

```yaml
policy:
  resources:
    infra_meters:
      saphana:
        agent_group: default
        per_agent_group: true
        receivers:
          saphana: [saphanareceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/saphanareceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
