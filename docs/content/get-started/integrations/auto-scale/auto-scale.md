---
title: Auto Scale
sidebar_position: 2
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
---

# Auto Scale

uto-scaling is a powerful method for automatically adjusting the number of
instances or resources allocated to a service based on fluctuating workload
demands. Aperture's closed-loop control policies provide an ideal solution for
defining auto-scaling criteria. Aperture Agents interface with infrastructure
APIs in order to perform auto-scaling. For example, Aperture Agents can invoke
Kubernetes APIs to perform auto-scaling of any scalable resource in a Kubernetes
cluster.

## Setup

Aperture performs auto-scaling based on
[Signals](concepts/policy/circuit#signal) in an
[Aperture policy](concepts/policy/policy.md).

- [Kubernetes](./kubernetes/kubernetes.md): Auto-scaling any scalable resource
  in a Kubernetes cluster can be achieved using Aperture. In order to do this,
  you must first ensure that an Aperture Agent is installed on the cluster, and
  then configure
  [auto-scaling policies](tutorials/integrations/auto-scale/auto-scale.md) at
  the Aperture Controller.
