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

See also [awsfirehosereceiver docs][receiver] in opentelemetry-collect-contrib
repository.

:::

:::note

The `awsfirehosereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/awsfirehosereceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for AWS Kinesis Data Firehose
using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    awsfirehose:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - awsfirehose
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        awsfirehose: [awsfirehosereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsfirehosereceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
