---
title: Pod Scaler
keywords:
  - Autoscaling
  - auto-scaler
  - Kubernetes
sidebar_position: 1
---

:::info

See also [_Pod Scaler_ reference](/reference/policies/spec.md#pod-scaler)

:::

_Pod Scaler_ is a basic building block of an auto-scaling policy for Kubernetes.
It can scale out or scale in a Kubernetes Resources such as a Deployment. The
component takes the desired replicas as an input signal and scales the
underlying resources based on the value of the signal. To complete the feedback
loop, the component emits output signals for the number of configured replicas
and the actual number of replicas deployed.

A _Pod Scaler_ component is not required if you are using the high-level
[_Auto Scaler_](auto-scaler.md) component. It defines multiple scale in, scale
out controllers and takes care of instantiating the _Pod Scaler_ component.
