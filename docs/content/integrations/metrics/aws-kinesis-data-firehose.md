---
title: AWS Kinesis Data Firehose
description: Integrating AWS Kinesis Data Firehose Metrics
keywords:
  - awsfirehose
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [awsfirehosereceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `awsfirehosereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/awsfirehosereceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for AWS
Kinesis Data Firehose as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          awsfirehose:
            per_agent_group: true
            receivers:
              awsfirehose: [awsfirehosereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsfirehosereceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
