---
title: Kubernetes Object Selector
sidebar_position: 2
keywords:
  - auto-scaling
  - Autoscaling
  - auto-scaler
  - Kubernetes
---

:::info

See also
[Kubernetes Object Selector reference](/reference/policies/spec.md#kubernetes-object-selector)

:::

_Kubernetes Object Selectors_ are used by auto-scaling components in a policy
such as [PodScaler](reference/policies/spec.md#pod-scaler) or
[AutoScaler](reference/policies/spec.md#pod-scaler). A _Kubernetes Object
Selector_ identifies a resource in the Kubernetes cluster.

A Kubernetes Object Selector consists of:

- agent group name: The Agent Group identifies Aperture Agents that the
  component apples to.
- api version: The Kubernetes API version of the resource.
- kind: The Kubernetes Kind of the resource. E.g. Deployment, ReplicaSet,
  StatefulSet.
- name: The name of the resource.
- namespace: The Kubernetes namespace of the resource.

:::info Control Points Discovery

The Kubernetes resources identified by a Kubernetes Object Selector are called
Kubernetes Control Points. These are a subset of resources in a Kubernetes
cluster resource that can be scaled in or out. Aperture Agents perform automated
discovery of Kubernetes Control Points in a cluster.

:::
