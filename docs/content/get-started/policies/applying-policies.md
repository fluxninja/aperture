---
title: Applying Policies
description: How to apply Policies in Aperture
keywords:
  - policy
sidebar_position: 1
---

```mdx-code-block
import {apertureVersion} from '../../introduction.md';
import VersionCode from '../../assets/scripts/version-code';
```

Aperture policies are applied at the Kubernetes cluster where the Aperture
controller is running. Policies are represented as Kubernetes objects of kind
Policy, which is a custom resource definition. Once an Aperture Policy spec is
defined, it can be applied like any other Kubernetes resource.:

```bash
kubectl apply -f <aperture-policy-file>
```

Here's an example of a policy configuration file:

```mdx-code-block
<CodeBlock language="yaml">
```

```yaml
{@include: ../../tutorials/flow-control/assets/static-rate-limiting/static-rate-limiting.yaml}
```

```mdx-code-block
</CodeBlock>
```

Follow the steps given below to create the above Policy:

1. Create the Policy by running the following command:

<VersionCode version="${apertureVersion}">
kubectl apply -f https://raw.githubusercontent.com/fluxninja/aperture/v{apertureVersion}/docs/content/tutorials/flow-control/assets/static-rate-limiting/static-rate-limiting.yaml
</VersionCode>

2. Run the following to check if the Policy was created.

```bash
kubectl get policies
```
