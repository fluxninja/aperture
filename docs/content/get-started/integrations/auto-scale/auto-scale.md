---
title: Auto Scale
sidebar_position: 2
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
---

# Auto Scale

Auto-scaling is the process of automatically adjusting the number of instances
or resources allocated to a service to meet changing workload demands.
Aperture's closed-loop control policies are a natural fit for defining
auto-scaling criteria. Aperture Agents interface with infrastructure APIs in
order to perform auto-scaling. For example, Aperture Agents can invoke
Kubernetes APIs to perform auto-scaling of any scalable resource in a Kubernetes
cluster. In this guide, we will explore how to configure auto-scaling policies
in Aperture and take advantage of this powerful capability to optimize the
performance and cost of your services.

## Setup

Aperture performs auto-scaling based on
[Signals](concepts/policy/circuit#signal) in an
[Aperture policy](concepts/policy/policy.md).

- [Kubernetes](./kubernetes/kubernetes.md): Any scalable resource in a
  Kubernetes cluster can be auto-scaled via Aperture. An Aperture Agent must be
  installed on the cluster. And auto-scale policies need to be configured at the
  Aperture Controller.
