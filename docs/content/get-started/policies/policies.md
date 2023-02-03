---
title: Generating and Applying Policies
description: How to generate and apply policies in Aperture
keywords:
  - policy
  - jsonnet
  - grafana
  - policy
sidebar_position: 4
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
```

## Introduction

Aperture comes with a pre-packaged [Blueprints][blueprints] that can be used to
generate [Policies][policies] and Grafana Dashboards. Blueprints can be used
both as a guide for creating new policies, or used as-is by providing required
parameters or customizations.

In order to install Aperture Blueprints and generate policies, you can use
`aperturectl` CLI tool.

## Manage Aperture Blueprints

Follow the installation steps for `aperturectl`
[here](/get-started/aperture-cli/aperture-cli.md#installation)

Now that we have `aperturectl` installed, we can use it to manage Aperture
Blueprints.

The below command will pull the Aperture Blueprints from the Aperture repository
and store it locally:

<CodeBlock language="bash">
aperturectl blueprints pull --version {apertureVersion}
</CodeBlock>

Run the following command to list the policies and dashboards associated with
the installed Blueprints:

<CodeBlock language="bash">
aperturectl blueprints list --version {apertureVersion}
</CodeBlock>

To learn more about `aperturectl` and the other commands it supports, visit
[aperturectl](/reference/aperture-cli/aperturectl.md).

## Generate Aperture Policies and Grafana Dashboards {#generating-aperture-policies-and-grafana-dashboards}

Once you have the Blueprints installed locally, you can use `aperturectl` to
generate Aperture Policies and Grafana Dashboards.

Suppose you want to generate the `policies/static-rate-limiting` policy with the
following values file `rate-limiting-values.yaml`:

```yaml
common:
  policy_name: static-rate-limiting
dashboard:
  datasource:
    filter_regex: ""
    name: $datasource
  refresh_interval: 10s
policy:
  classifiers: []
  evaluation_interval: 300s
  rate_limiter:
    dynamic_config:
      overrides:
        - label_value: gold
          limit_scale_factor: 1
    flow_selector:
      flow_matcher:
        control_point: ingress
      service_selector:
        agent_group: default
        service: service1-demo-app.demoapp.svc.cluster.local
    parameters:
      label_key: http.request.header.user_type
      lazy_sync:
        enabled: true
        num_sync: 5
      limit_reset_interval: 1s
    rate_limit: "50.0"
```

You can run the following command to generate the policy:

```bash
# Generate the policy
aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting-values.yaml

# See what was generated
tree ${apertureVersion}/policies/static-rate-limiting
```

The output would look something like below:

```bash
${apertureVersion}/policies/static-rate-limiting
├── dashboards
│   └── static-rate-limiting.json
├── graphs
│   ├── static-rate-limiting.dot
│   └── static-rate-limiting.svg
└── policies
    └── static-rate-limiting.yaml

3 directories, 4 files
```

## Apply Policy

<Tabs>
<TabItem value="aperturectl" label="aperturectl">

You can pass `--apply` flag with the `apertuectl` to directly apply the
generated policies on the Kubernetes cluster in the namespace where the Aperture
Controller is installed.

<CodeBlock language="bash">
aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting-values.yaml --version {apertureVersion} --apply
</CodeBlock>

It uses the default configuration for Kubernetes cluster under `~/.kube/config`.
You can pass the `--kube-config` flag to pass any other path.

<CodeBlock language="bash">
aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting-values.yaml --version {apertureVersion} --kube-config=/path/to/config --apply
</CodeBlock>

</TabItem>
<TabItem value="kubectl" label="kubectl">

The policy YAML generated using the above example can also be applied using
`kubectl`.

<CodeBlock language="bash">
kubectl apply -f ${apertureVersion}/policies/static-rate-limiting/policies/static-rate-limiting.yaml -n aperture-controller
</CodeBlock>

</TabItem>
</Tabs>

Run the following command to check if the Aperture Policy was created.

```bash
kubectl get policies.fluxninja.com -n aperture-controller
```

The Aperture Policy runtime can be visualized in Grafana or any other Prometheus
compatible analytics tool. Refer to the Prometheus compatible metrics available
from [Controller][controller-metrics] and [Agent][agent-metrics]. Some of the
Policy [Blueprints][blueprints] come with recommended Grafana dashboards.

## Delete Policy

Run the following command to delete the above Aperture Policy:

```bash
kubectl delete policies.fluxninja.com static-rate-limiting -n aperture-controller
```

[controller-metrics]: /reference/observability/prometheus-metrics/controller.md
[agent-metrics]: /reference/observability/prometheus-metrics/agent.md
[blueprints]: /reference/policies/bundled-blueprints/bundled-blueprints.md
[policies]: /concepts/policy/policy.md
[service]: /concepts/integrations/flow-control/service.md
