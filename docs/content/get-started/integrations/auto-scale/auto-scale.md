---
title: Auto Scale
sidebar_position: 2
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
  - HPA
---

# Auto Scale

Auto Scale in Aperture refers to the mechanism of automatically scaling the
number of Kubernetes pods based on the load on the pods using metrics from flow
control or k8s resources ultization for a any Kubernetes Resource such as
Deployment, ReplicaSet (CPU, memory, etc.).

## Setup

Aperture can interface with Kubernetes to control and monitor all Kubernetes
resources that are scalable as Kubernetes Control Points. This is accomplished
through the use of the Aperture Controller, which watches a specific node in the
Kubernetes cluster and discovers all Kubernetes resources that are scalable as
Kubernetes Control Points.
