---
title: MongoDB Atlas
description: Integrating MongoDB Atlas Metrics
keywords:
  - mongodbatlas
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [mongodbatlasreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `mongodbatlasreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/mongodbatlasreceiver` to the `bundled_extensions` list to
make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
MongoDB Atlas as part of [Policy resources][policy-resources] while [applying
the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          mongodbatlas:
            per_agent_group: true
            receivers:
              mongodbatlas: [mongodbatlasreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/mongodbatlasreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
