---
title: Pulsar
description: Integrating Pulsar Metrics
keywords:
  - pulsar
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [pulsarreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `pulsarreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/pulsarreceiver` to the `bundled_extensions` list to make [the
receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Pulsar as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    infra_meters:
      pulsar:
        agent_group: default
        per_agent_group: true
        receivers:
          pulsar: [pulsarreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/pulsarreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
