---
title: Auto Scale
sidebar_position: 2
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
  - HPA
---

# Auto Scaling

In today's digital age, enterprises must deal with a huge amount of data and
traffic that comes with the growth of online services. This creates a need for
auto-scaling, where an application can adjust its resources dynamically based on
demand. Autoscaling is crucial because it can help to optimize resource usage
and improve user experience, while also preventing system failures due to
overload. However, Auto-scaling in most of cloud environments is limited to a
few metrics such as CPU and memory utilization. This is where FluxNinja Aperture
Auto-scaling comes in where it address these issues by automatically adjusting
the resources allocated to an application using load-based autoscaling. This
helps manage load and ensure resources are used efficiently while keeping the
system stability is in check.

# Kubernetes Control Points

Kubernetes Control Points are a type of Kubernetes resource that can be scaled
up or down dynamically. These resources can be used to manage and optimize the
allocation of computing resources for containerized applications.

Some of Kubernetes Control Points include:

- Deployments
- ReplicaSets
- StatefulSets

## Kubernetes Resource Discovery

Aperture has a discover mechanism to discover all the Kubernetes resources that
are scalable as Kubernetes Control points.
