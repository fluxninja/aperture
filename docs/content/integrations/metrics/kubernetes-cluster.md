---
title: Kubernetes Cluster
description: Integrating Kubernetes Cluster Metrics
keywords:
  - k8s_cluster
  - otel
  - opentelemetry
  - collector
  - metrics
---

:::info

See also [k8sclusterreceiver docs][receiver] in
`opentelemetry-collector-contrib` repository.

:::

:::note

The `k8sclusterreceiver` extension is available in the default agent image. If
you're [building][build] your own Aperture Agent, add
`integrations/otel/k8sclusterreceiver` to the `bundled_extensions` list to make
[the receiver][receiver] available.

:::

You can configure the [OpenTelemetry Collector][opentelemetry-collector] for
Kubernetes Cluster as part of [Policy resources][policy-resources] while
[applying the policy][applying-policy]:

```yaml
policy:
  resources:
    telemetry_collectors:
      - agent_group: default
        infra_meters:
          k8s_cluster:
            per_agent_group: true
            receivers:
              k8s_cluster: [k8sclusterreceiver configuration here]
```

[build]: /reference/aperturectl/build/agent/agent.md
[receiver]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/k8sclusterreceiver
[opentelemetry-collector]: /reference/policies/spec.md#telemetry-collector
[applying-policy]: /use-cases/use-cases.md
[policy-resources]: /reference/policies/spec.md#resources
