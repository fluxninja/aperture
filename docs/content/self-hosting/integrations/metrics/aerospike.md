---
title: Aerospike
description: Integrating Aerospike Metrics
keywords:
  - aerospike
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [aerospikereceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `aerospikereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/aerospikereceiver` to the `bundled_extensions` list to make
the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Aerospike as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      aerospike:
        agent_group: default
        per_agent_group: true
        receivers:
          aerospike: [aerospikereceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/aerospikereceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
