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

To simplify the process of creating policies in Aperture, the built-in blueprint
system can be utilized. The Aperture repository contains several
[blueprints][blueprints] that can generate [policies][policies], and [Grafana
dashboards][grafana]. These blueprints serve as starting points for creating new
policies, or can be used as-is by providing the required parameters or
customizations. The [use-cases](/use-cases/use-cases.md) section showcases
practical examples of blueprints in action.

To manage blueprints and generate policies, use the
[aperturectl](/reference/aperturectl/aperturectl.md) CLI.

For advanced users interested in designing new policies, explore the example
circuit created in the
[detecting overload](../../use-cases/alerting/detecting-overload.md) use-case.
This example serves as a valuable reference for understanding the process of
creating custom policies in Aperture.

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
<CodeBlock language="bash">aperturectl blueprints values --name=rate-limiting/base --version={apertureVersion} --output-file=values.yaml</CodeBlock>
```

You can then edit the `values.yaml` to provide the required fields
(`__REQUIRED_FIELD__` placeholder) as follows:

<Tabs>
<TabItem value="Final/Edited Values">

```yaml
{@include: ./assets/values.yaml}
```

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
<CodeBlock language="bash">aperturectl blueprints generate --name=rate-limiting/base
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
<CodeBlock language="bash">aperturectl blueprints generate --name=rate-limiting/base
--values-file=values.yaml --apply --version={apertureVersion}</CodeBlock>
```

:::info

See [aperturectl configuration](/get-started/installation/configure-cli.md) on
how to configure what aperturectl should connect to.

:::

```mdx-code-block
</TabItem>
<TabItem value="kubectl" label="kubectl">
```

:::caution

You can only apply policies with kubectl on [Self-Hosted][self-hosted] Aperture
Controller.

:::

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

```mdx-code-block
<Tabs>
<TabItem value="aperturectl" label="aperturectl">
```

```bash
aperturectl policies
```

```mdx-code-block
</TabItem>
<TabItem value="kubectl" label="kubectl">
```

```bash
kubectl get policies.fluxninja.com -n aperture-controller
```

```mdx-code-block
</TabItem>
</Tabs>
```

The policy runtime can be visualized in [FluxNinja Cloud][], Grafana or any
other Prometheus compatible analytics tool. Refer to the Prometheus compatible
metrics available from the [controller][controller-metrics] and
[agent][agent-metrics]. Some policy [blueprints][blueprints] come with
recommended Grafana dashboards.

## Deleting Policies

Run the following command to delete the above policy:

```mdx-code-block
<Tabs>
<TabItem value="aperturectl" label="aperturectl">
```

```bash
aperturectl delete policy --policy=rate-limiting
```

```mdx-code-block
</TabItem>
<TabItem value="kubectl" label="kubectl">
```

```bash
kubectl delete policies.fluxninja.com rate-limiting -n aperture-controller
```

```mdx-code-block
</TabItem>
</Tabs>
```

[controller-metrics]: /reference/observability/prometheus-metrics/controller.md
[agent-metrics]: /reference/observability/prometheus-metrics/agent.md
[blueprints]: /reference/blueprints/blueprints.md
[policies]: /concepts/advanced/policy.md
[grafana]: https://grafana.com/docs/grafana/latest/dashboards/
[self-hosted]: /self-hosting/self-hosting.md
[FluxNinja Cloud]: /introduction.md
