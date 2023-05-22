---
title: CollectD write_http plugin JSON
description: Integrating CollectD write_http plugin JSON Metrics
keywords:
  - collectd
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [collectdreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `collectdreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/collectdreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
CollectD write_http plugin JSON as part of [Policy resources][policy-resources]
while [applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          collectd:
            per_agent_group: true
            pipeline:
              processors:
                - batch
              receivers:
                - collectd
            processors:
              batch:
                send_batch_size: 10
                timeout: 10s
            receivers:
              collectd: [collectdreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/collectdreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
