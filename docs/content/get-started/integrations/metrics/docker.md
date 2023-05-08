---
title: Docker Stats
description: Integrating Docker Stats Metrics
keywords:
  - docker_stats
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [dockerstatsreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `dockerstatsreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/dockerstatsreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Docker Stats as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          docker_stats:
            per_agent_group: true
            pipeline:
              processors:
                - batch
              receivers:
                - docker_stats
            processors:
              batch:
                send_batch_size: 10
                timeout: 10s
            receivers:
              docker_stats: [dockerstatsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/dockerstatsreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
