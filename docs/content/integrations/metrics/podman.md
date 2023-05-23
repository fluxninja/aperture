---
title: Podman Stats
description: Integrating Podman Stats Metrics
keywords:
  - podman_stats
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [podmanreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `podmanreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/podmanreceiver` to the `bundled_extensions` list to make [the
receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Podman Stats as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          podman_stats:
            per_agent_group: true
            receivers:
              podman_stats: [podmanreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/podmanreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
