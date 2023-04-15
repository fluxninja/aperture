---
title: Kubernetes Object Selector
sidebar_position: 1
keywords:
  - auto-scaling
  - Autoscaling
  - auto-scaler
  - Kubernetes
---

:::info

See also
[_Kubernetes Object Selector_ reference](/reference/policies/spec.md#kubernetes-object-selector)

:::

_Kubernetes Object Selectors_ are used by auto-scaling components in a policy,
such as [Pod Scaler](reference/policies/spec.md#pod-scaler) or
[_Auto Scaler_](reference/policies/spec.md#pod-scaler). A Kubernetes Object
Selector identifies a resource in the Kubernetes cluster.

A _Kubernetes Object Selector_ consists of:

- _Agent Group_: The Agent Group identifies Aperture Agents where the
  [component](components/components.md) gets applied.
- API Version: The Kubernetes API version of the resource.
- Kind: The Kind of the Kubernetes resource, such as Deployment, ReplicaSet,
  StatefulSet.
- Name: The name of the Kubernetes resource.
- Namespace: The Kubernetes namespace of the resource.

:::info Control Points Discovery

The Kubernetes resources identified by a _Kubernetes Object Selector_ are called
_Kubernetes Control Points_. These are a subset of resources in a Kubernetes
cluster resource that can be scaled in or out. Aperture Agents performs
automated discovery of Kubernetes Control Points in a cluster.

:::
