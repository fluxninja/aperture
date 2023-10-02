---
title: HAProxy
description: Integrating HAProxy Metrics
keywords:
  - haproxy
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [haproxyreceiver docs][receiver] in `opentelemetry-collector-contrib`
repository.

:::

:::note

The `haproxyreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/haproxyreceiver` to the `bundled_extensions` list to make the
[receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
HAProxy as part of [Policy resources][policy-resources] while applying the
policy:

```yaml
policy:
  resources:
    infra_meters:
      haproxy:
        agent_group: default
        per_agent_group: true
        receivers:
          haproxy: [haproxyreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/haproxyreceiver
[opentelemetry-collector]: /reference/configuration/spec.md#telemetry-collector
[policy-resources]: /reference/configuration/spec.md#resources
