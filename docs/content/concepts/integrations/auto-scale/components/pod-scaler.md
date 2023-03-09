---
title: Pod Scaler
keywords:
  - Autoscaling
  - auto-scaler
  - Kubernetes
sidebar_position: 1
---

:::info

See also [Pod Scaler reference](/reference/policies/spec.md#pod-scaler)

:::

Pod Scaler is a basic building block of an auto-scaling policy for Kubernetes.
It can scale out or scale in a Kubernetes Resources such as a Deployment. The
component takes the desired replicas as an input signal and scales the
underlying resources based on the value of the signal. To complete the feedback
loop the component emits output signals for the number of configured replicas
and the actual number of replicas deployed.

A Pod Scaler component is not required if you are using the high-level
[Auto Scaler](auto-scaler.md) component. It defines multiple scale in, scale out
controllers and takes care of instantiating the Pod Scaler component.
