---
title: Applying Policies
description: How to apply Policies in Aperture
keywords:
  - policy
sidebar_position: 1
---

```mdx-code-block
import {apertureVersion} from '../../introduction.md';
import CodeBlock from '@theme/CodeBlock';
```

Aperture Policies are applied at the Kubernetes cluster where the Aperture
Controller is running. Aperture Policies are represented as Kubernetes objects
of kind Policy, which is a
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).

Once the Aperture Policy spec is defined, it can be applied like any other
Kubernetes resource:

```bash
kubectl apply -f <aperture-policy-file>
```

Here's an example of the Aperture Policy configuration file:

```yaml
{@include: ../../tutorials/flow-control/assets/static-rate-limiting/static-rate-limiting.yaml}
```

Follow the steps given below to create the above Aperture Policy:

:::info

The Aperture Policy has to be created in the same Kubernetes namespace where the
Aperture Controller is running.

:::

1. Create the Aperture Policy by running the following command:

<CodeBlock language="bash">
{`kubectl apply -f https://raw.githubusercontent.com/fluxninja/aperture/v${apertureVersion}/docs/content/tutorials/flow-control/assets/static-rate-limiting/static-rate-limiting.yaml`}
</CodeBlock>

2. Run the following command to check if the Aperture Policy was created.

```bash
kubectl get policies
```
