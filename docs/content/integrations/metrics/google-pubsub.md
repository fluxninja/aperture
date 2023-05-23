---
title: Google Cloud Pub/Sub
description: Integrating Google Cloud Pub/Sub Metrics
keywords:
  - googlecloudpubsub
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [googlecloudpubsubreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `googlecloudpubsubreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/googlecloudpubsubreceiver` to the `bundled_extensions` list
to make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Google Cloud Pub/Sub as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          googlecloudpubsub:
            per_agent_group: true
            receivers:
              googlecloudpubsub: [googlecloudpubsubreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/googlecloudpubsubreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /applying-policies/applying-policies.md
[policy-resources]: /reference/policies/spec.md#resources
