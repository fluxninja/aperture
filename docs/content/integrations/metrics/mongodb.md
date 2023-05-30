---
title: MongoDB
description: Integrating MongoDB Metrics
keywords:
  - mongodb
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [mongodbreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `mongodbreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/mongodbreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
MongoDB as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          mongodb:
            per_agent_group: true
            receivers:
              mongodb: [mongodbreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/mongodbreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/policies/spec.md#resources
