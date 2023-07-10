---
title: Auto Scale
sidebar_position: 4
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
---

:::info See also

[Auto-scaling use-case](../../use-cases/auto-scaling/auto-scaling.md)

:::

Auto-scaling is a powerful method for automatically adjusting the number of
instances or resources allocated to a service based on fluctuating workload
demands. Aperture's closed-loop control policies provide an ideal solution for
defining auto-scaling criteria. Aperture Agents interface with infrastructure
APIs to perform auto-scaling. For example, Aperture Agents can invoke Kubernetes
APIs to perform auto-scaling of any scalable resource in a Kubernetes cluster.

![Auto Scaling](./assets/autoscale-dark.svg#gh-dark-mode-only)

![Auto Scaling](./assets/autoscale-light.svg#gh-light-mode-only)

## Setup

Aperture performs auto-scaling based on
[Signals](/concepts/advanced/circuit#signal) in an
[Aperture policy](/concepts/advanced/policy.md).

- [Kubernetes](./kubernetes/kubernetes.md): Auto-scaling any scalable resource
  in a Kubernetes cluster can be achieved using Aperture. To achieve this, you
  must first ensure that an Aperture Agent is installed on the cluster, and then
  configure auto-scaling policies at the Aperture Controller.
