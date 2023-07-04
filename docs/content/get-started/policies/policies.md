---
title: Your First Policy
description: How to generate and apply policies in Aperture
keywords:
  - policy
  - jsonnet
  - grafana
  - policy
sidebar_position: 3
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Introduction

The easiest way to get started with policies in Aperture is to use the built-in
blueprint system. The Aperture repository contains several
[blueprints][blueprints] that can be used to generate [policies][policies] and
[Grafana dashboards][grafana]. Blueprints can be used both as a guide for
creating new policies, or used as-is by providing required parameters or
customizations.

To manage blueprints and generate policies, you can use the `aperturectl`
[CLI tool](/reference/aperturectl/aperturectl.md), by following the
[installation steps](../../get-started/installation/aperture-cli/aperture-cli.md#installation)
first.

Advanced users can learn about designing new policies by following the example
circuit created in

[Detecting overload](../../use-cases/alerting/detecting-overload.md) use-case.

<Zoom>

```mermaid
{@include: ./assets/blueprints.mmd}
```

</Zoom>

## Listing Available Blueprints

The following command can be used to list available blueprints:

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints list --version={apertureVersion}</CodeBlock>
```

Which will output the following:

```bash
auto-scaling/pod-auto-scaler
load-ramping/base
load-scheduling/average-latency
load-scheduling/postgresql
load-scheduling/promql
quota-scheduling/base
rate-limiting/base
```

## Customizing Blueprints

Blueprints use a configuration file to provide required fields and to customize
the generated policy and dashboard files.

For example, to generate a `policies/rate-limiting` policy, you can first
generate a `values.yaml` file using the following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints values --name=policies/rate-limiting --version={apertureVersion} --output-file=values.yaml</CodeBlock>
```

You can then edit the `values.yaml` to provide the required fields
(`__REQUIRED_FIELD__` placeholder) as follows:

<Tabs>
<TabItem value="Final/Edited Values">

```yaml
{@include: ./assets/values.yaml}
```

You can then run the following command to generate the blueprint:

</TabItem>
<TabItem value="Placeholder Values">

```yaml
{@include: ./assets/raw_values.yaml}
```

</TabItem>
</Tabs>

## Generating Policies and Dashboards

Once the `values.yaml` file is ready, you can generate the blueprint using the
following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints generate --name=policies/rate-limiting
--values-file=values.yaml --output-dir=policy-gen --version={apertureVersion}</CodeBlock>
```

The following directory structure will be generated:

```bash
policy-gen
├── dashboards
│   └── rate-limiting.json
├── graphs
│   ├── rate-limiting.dot
│   └── rate-limiting.mmd
└── policies
│   ├── rate-limiting-cr.yaml
│   └── rate-limiting.yaml
```

## Applying Policies

The generated policies can be applied using `aperturectl` or `kubectl`.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl" label="aperturectl">
```

You can pass the `--apply` flag with the `aperturectl` to directly apply the
generated policies on a Kubernetes cluster in the namespace where the Aperture
Controller is installed.

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints generate --name=policies/rate-limiting
--values-file=values.yaml --apply --version={apertureVersion}</CodeBlock>
```

It uses the default configuration for Kubernetes cluster under `~/.kube/config`.
You can pass the `--kube-config` flag to pass any other path.

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints generate --name=policies/rate-limiting
--values-file=values.yaml --kube-config=/path/to/config --apply --version={apertureVersion}</CodeBlock>
```

```mdx-code-block
</TabItem>
<TabItem value="kubectl" label="kubectl">
```

The policy YAML generated (Kubernetes Custom Resource) using the above example
can also be applied using `kubectl`.

```bash
kubectl apply -f policy-gen/configuration/rate-limiting-cr.yaml -n aperture-controller
```

```mdx-code-block
</TabItem>
</Tabs>
```

Run the following command to check if the policy was created.

```bash
kubectl get policies.fluxninja.com -n aperture-controller
```

The policy runtime can be visualized in Grafana or any other Prometheus
compatible analytics tool. Refer to the Prometheus compatible metrics available
from the [controller][controller-metrics] and [agent][agent-metrics]. Some
policy [blueprints][blueprints] come with recommended Grafana dashboards.

## Deleting Policies

Run the following command to delete the above policy:

```bash
kubectl delete policies.fluxninja.com rate-limiting -n aperture-controller
```

[controller-metrics]: /reference/observability/prometheus-metrics/controller.md
[agent-metrics]: /reference/observability/prometheus-metrics/agent.md
[blueprints]: /reference/blueprints/blueprints.md
[policies]: /concepts/advanced/policy.md
[grafana]: https://grafana.com/docs/grafana/latest/dashboards/
