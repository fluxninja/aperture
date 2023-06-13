---
title: OTLP JSON File
description: Integrating OTLP JSON File Metrics
keywords:
  - otlpjsonfile
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [otlpjsonfilereceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `otlpjsonfilereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/otlpjsonfilereceiver` to the `bundled_extensions` list to
make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
OTLP JSON File as part of [Policy resources][policy-resources] while [applying
the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          otlpjsonfile:
            per_agent_group: true
            receivers:
              otlpjsonfile: [otlpjsonfilereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/otlpjsonfilereceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/policies/spec.md#resources
