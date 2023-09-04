---
title: F5 Big-IP
description: Integrating F5 Big-IP Metrics
keywords:
  - bigip
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [bigipreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `bigipreceiver` extension is available in the default agent image. If you're
[building][build] your own Aperture Agent, add `integrations/otel/bigipreceiver`
to the `bundled_extensions` list to make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for F5
Big-IP as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    infra_meters:
      bigip:
        agent_group: default
        per_agent_group: true
        receivers:
          bigip: [bigipreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/bigipreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
