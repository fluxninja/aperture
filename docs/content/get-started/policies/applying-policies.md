---
title: Get Started with First Policy
description: How to apply Policies in Aperture
keywords:
  - policy
sidebar_position: 1
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
```

Aperture [Policies][policy-concept] are applied at the Kubernetes cluster where
the Aperture Controller is running. Aperture Policies are represented as
Kubernetes objects of kind Policy, which is a
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).

Once the Aperture Policy spec is defined, it can be applied like any other
Kubernetes resource:

```bash
kubectl --namespace=aperture-controller apply -f <aperture-policy-file>
```

Here's an example of the Aperture Policy configuration file:

```yaml
{@include: ../../tutorials/flow-control/assets/static-rate-limiting/static-rate-limiting.yaml}
```

## Create Policy

Follow the steps given below to create the above Aperture Policy:

:::info

The Aperture Policy has to be created in the same Kubernetes namespace where the
Aperture Controller is running.

:::

1. Create the Aperture Policy by running the following command:

<CodeBlock language="bash">
{`kubectl --namespace=aperture-controller apply -f https://raw.githubusercontent.com/fluxninja/aperture/${apertureVersion}/docs/content/tutorials/flow-control/assets/static-rate-limiting/static-rate-limiting.yaml`}
</CodeBlock>

2. Run the following command to check if the Aperture Policy was created.

```bash
kubectl --namespace=aperture-controller get policies
```

3. The Aperture Policy runtime can be visualized in Grafana or any other
   Prometheus compatible analytics tool. Refer to the Prometheus compatible
   metrics available from [Controller][controller-metrics] and
   [Agent][agent-metrics]. Some of the Policy [Blueprints][blueprints] come with
   recommended Grafana dashboards.

## Delete Policy

Run the following command to delete the above Aperture Policy:

```bash
kubectl --namespace=aperture-controller delete policy static-rate-limiting
```

[controller-metrics]: /references/observability/prometheus-metrics/controller.md
[agent-metrics]: /references/observability/prometheus-metrics/agent.md
[policy-concept]: /concepts/policy/policy.md
[blueprints]: /get-started/policies/blueprints/blueprints.md
