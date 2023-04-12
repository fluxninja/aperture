---
title: CouchDB
description: Integrating CouchDB Metrics
keywords:
  - couchdb
  - otel
  - opentelemetry
  - collector
  - metrics
---

Before proceeding, ensure that you have [built][build] Aperture Agent with the
`couchdbreceiver` extension enabled, so that [couchdbreceiver][receiver] is
available.

You can configure [Custom metrics][custom-metrics] for CouchDB using the
following configuration in the [Aperture Agent's config][agent-config]:

```yaml
otel:
  custom_metrics:
    couchdb:
      per_agent_group: true
      pipeline:
        processors:
          - batch
        receivers:
          - couchdb
      processors:
        batch:
          send_batch_size: 10
          timeout: 10s
      receivers:
        couchdb: [couchdbreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/couchdbreceiver
[custom-metrics]: /reference/configuration/agent.md#custom-metrics-config
[agent-config]: /reference/configuration/agent.md#agent-o-t-e-l-config
