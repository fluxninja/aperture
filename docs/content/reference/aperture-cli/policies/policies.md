---
title: Manage Policies
description: How to generate and apply policies in Aperture using aperturectl
keywords:
  - policy
  - jsonnet
  - grafana
  - policy
sidebar_position: 6
---

```mdx-code-block
import {apertureVersion} from '../../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

`aperturectl` is a powerful CLI that complements the Aperture Cloud UI. With
aperturectl, you can also manage blueprints and generate policies, dashboards,
and graphs. In this overview, you'll explore the various commands available in
aperturectl for managing and creating policies using blueprints.

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

## Generating Policies and Dashboards {#generating-policies}

Once the `values.yaml` file is ready, you can generate the blueprint using the
following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl blueprints generate --values-file=values.yaml --output-dir=policy-gen</CodeBlock>
<CodeBlock language="bash">aperturectl dashboard --policy-file=policy-gen/policies/rate-limiting-cr.yaml --output-dir=policy-gen</CodeBlock>
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
<TabItem value="aperture-cloud" label="Aperture Cloud">
```

You can pass the `--apply` flag with the `aperturectl cloud` to directly apply
the generated policies on the Aperture Cloud Controller.

:::info

See [Set up CLI (aperturectl)](/reference/aperture-cli/aperture-cli.md) for more
information on how to configure what aperturectl should connect to.

:::

```mdx-code-block
<CodeBlock language="bash">aperturectl cloud policy apply --file policy-gen/policies/rate-limiting.yaml</CodeBlock>
```

Run the following command to check if the policy was created.

```mdx-code-block
<CodeBlock language="bash">aperturectl cloud policies</CodeBlock>
```

```mdx-code-block
</TabItem>
<TabItem value="aperture-for-infra" label="Self Hosting">
```

```mdx-code-block
<Tabs>
<TabItem value="Kubernetes Operator" label="Kubernetes Operator">
```

If the Aperture Controller is deployed on
[Kubernetes using Operator](/aperture-for-infra/controller/kubernetes/operator/operator.md),
you can apply the policy using the following command:

```mdx-code-block
<CodeBlock language="bash">kubectl apply -f policy-gen/configuration/rate-limiting-cr.yaml -n aperture-controller</CodeBlock>
```

Run the following command to check if the policy was created.

```mdx-code-block
<CodeBlock language="bash">kubectl get policies.fluxninja.com -n aperture-controller</CodeBlock>
```

```mdx-code-block
</TabItem>
<TabItem value="Kubernetes Namespace-scoped" label="Kubernetes Namespace-scoped">
```

If the Aperture Controller is deployed on
[Kubernetes using Namespace-scoped](/aperture-for-infra/controller/kubernetes/namespace-scoped/namespace-scoped.md),
you can apply the policy using the following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl policy apply --file policy-gen/policies/rate-limiting.yaml --kube --controller-ns aperture-controller</CodeBlock>
```

Run the following command to check if the policy was created.

```mdx-code-block
<CodeBlock language="bash">aperturectl policies --controller-ns aperture-controller</CodeBlock>
```

```mdx-code-block
</TabItem>
<TabItem value="Docker" label="Docker">
```

If the Aperture Controller is deployed on
[Docker](/aperture-for-infra/controller/docker.md), you can apply the policy using the
following command:

```mdx-code-block
<CodeBlock language="bash">aperturectl policy apply --file policy-gen/policies/rate-limiting.yaml --controller localhost:8080 --insecure</CodeBlock>
```

Run the following command to check if the policy was created.

```mdx-code-block
<CodeBlock language="bash">aperturectl policies --controller localhost:8080 --insecure</CodeBlock>
```

```mdx-code-block
</TabItem>
</Tabs>
</TabItem>
</Tabs>
```

---

The policy runtime can be visualized in [Aperture Cloud][aperture-cloud],
Grafana or any other Prometheus compatible analytics tool. Refer to the
Prometheus compatible metrics available from the
[controller][controller-metrics] and [agent][agent-metrics].

## Deleting Policies

Run the following command to delete the above policy:

```mdx-code-block
<Tabs>
<TabItem value="aperturectl" label="aperturectl">
```

```bash
aperturectl policy delete --policy=rate-limiting
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
[aperture-cloud]: /introduction.md
