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

Auto Scaling is the process of automatically adjusting the number of instances
or resources allocated to a service to meet changing workload demands.
Aperture's closed-loop control policies are a natural fit for defining
autoscaling criteria. Aperture Agents interface with Kubernetes APIs to perform
autoscaling on any scalable resource in a Kubernetes cluster. In this guide, we
will explore how to configure Autoscaling policies in Aperture and take
advantage of this powerful capability to optimize the performance and cost of
your services.

## Setup

Aperture can interface with Kubernetes to control and monitor all Kubernetes
resources that are scalable as Kubernetes Control Points.
