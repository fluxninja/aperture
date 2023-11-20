---
title: JMX
description: Integrating JMX Metrics
keywords:
  - jmx
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [jmxreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `jmxreceiver` extension is available in the default agent image. If you're
[building][build] your own Aperture Agent, add `integrations/otel/jmxreceiver`
to the `bundled_extensions` list to make the [receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for JMX
as part of [Policy resources][policy-resources] while applying the policy:

```yaml
policy:
  resources:
    infra_meters:
      jmx:
        agent_group: default
        per_agent_group: true
        receivers:
          jmx: [jmxreceiver configuration here]
```

[build]: /reference/aperture-cli/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/jmxreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
