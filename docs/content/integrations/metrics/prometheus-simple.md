---
title: Simple Prometheus
description: Integrating Simple Prometheus Metrics
keywords:
  - prometheus_simple
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [simpleprometheusreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `simpleprometheusreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/simpleprometheusreceiver` to the `bundled_extensions` list to
make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Simple Prometheus as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          prometheus_simple:
            per_agent_group: true
            receivers:
              prometheus_simple: [simpleprometheusreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/simpleprometheusreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
