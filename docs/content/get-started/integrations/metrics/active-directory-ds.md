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
opentelemetry-collect-contrib repository.

:::

:::note

The `activedirectorydsreceiver` extension is available in the default agent
image. If you're [building][build] your own Aperture Agent, add
`integrations/otel/activedirectorydsreceiver` to the `bundled_extensions` list
to make [the receiver][receiver] available.

:::

You can configure [Custom metrics][custom-metrics] for Active Directory Domain
Services using the following configuration in the [Aperture Agent's
config][agent-config]:

```yaml
otel:
  custom_metrics:
    active_directory_ds:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - active_directory_ds
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        active_directory_ds: [activedirectorydsreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/activedirectorydsreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
