---
title: Splunk HEC
description: Integrating Splunk HEC Metrics
keywords:
  - splunk_hec
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [splunkhecreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `splunkhecreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/splunkhecreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Splunk HEC as part of [Policy resources][policy-resources] while [applying the
policy][applying-policy]:

```yaml
policy:
  resources:
    infra_meters:
      splunk_hec:
        agent_group: default
        per_agent_group: true
        receivers:
          splunk_hec: [splunkhecreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/splunkhecreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/configuration/spec.md#resources
