---
title: SNMP
description: Integrating SNMP Metrics
keywords:
  - snmp
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [snmpreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `snmpreceiver` extension is available in the default agent image. If you're
[building][build] your own Aperture Agent, add `integrations/otel/snmpreceiver`
to the `bundled_extensions` list to make the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
SNMP as part of [Policy resources][policy-resources] while applying the policy:

```yaml
policy:
  resources:
    infra_meters:
      snmp:
        agent_group: default
        per_agent_group: true
        receivers:
          snmp: [snmpreceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/snmpreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
