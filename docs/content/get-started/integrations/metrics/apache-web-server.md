---
title: Apache Web Server
description: Integrating Apache Web Server Metrics
keywords:
  - apache
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [apachereceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `apachereceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/apachereceiver` to the `bundled_extensions` list to make [the
receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Apache Web Server as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          apache:
            per_agent_group: true
            pipeline:
              processors:
                - batch
              receivers:
                - apache
            processors:
              batch:
                send_batch_size: 10
                timeout: 10s
            receivers:
              apache: [apachereceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/apachereceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
