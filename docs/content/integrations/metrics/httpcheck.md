---
title: HTTP Check
description: Integrating HTTP Check Metrics
keywords:
  - httpcheck
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [httpcheckreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `httpcheckreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/httpcheckreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
HTTP Check as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          httpcheck:
            per_agent_group: true
            receivers:
              httpcheck: [httpcheckreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/httpcheckreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
