---
title: AWS Container Insights Receiver
description: Integrating AWS Container Insights Receiver Metrics
keywords:
  - awscontainerinsightreceiver
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [awscontainerinsightreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `awscontainerinsightreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/awscontainerinsightreceiver` to the `bundled_extensions` list
to make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for AWS
Container Insights Receiver as part of [Policy resources][policy-resources]
while [applying the policy][applying-policy]:

```yaml
policy:
  resources:
    infra_meters:
      awscontainerinsightreceiver:
        agent_group: default
        per_agent_group: true
        receivers:
          awscontainerinsightreceiver:
            [awscontainerinsightreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awscontainerinsightreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
