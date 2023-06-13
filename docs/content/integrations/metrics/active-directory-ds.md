---
title: Active Directory Domain Services
description: Integrating Active Directory Domain Services Metrics
keywords:
  - active_directory_ds
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [activedirectorydsreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `activedirectorydsreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/activedirectorydsreceiver` to the `bundled_extensions` list
to make [the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Active Directory Domain Services as part of [Policy resources][policy-resources]
while [applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          active_directory_ds:
            per_agent_group: true
            receivers:
              active_directory_ds:
                [activedirectorydsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/activedirectorydsreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/policies/spec.md#resources
