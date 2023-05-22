---
title: AWS ECS Container Metrics
description: Integrating AWS ECS Container Metrics Metrics
keywords:
  - awsecscontainermetrics
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [awsecscontainermetricsreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `awsecscontainermetricsreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/awsecscontainermetricsreceiver` to the `bundled_extensions`
list to make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for AWS
ECS Container Metrics as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          awsecscontainermetrics:
            per_agent_group: true
            receivers:
              awsecscontainermetrics:
                [awsecscontainermetricsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
