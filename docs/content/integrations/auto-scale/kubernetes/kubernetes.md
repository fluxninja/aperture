---
title: Kubernetes
sidebar_position: 1
keywords:
  - scaling
  - auto-scaler
  - Kubernetes
---

## Installation

To get started with auto-scaling in Kubernetes,
[install Aperture Agent](/get-started/installation/agent/kubernetes/operator/operator.md)
on your Kubernetes cluster. The Agent needs certain permissions to discover
scalable resources and perform auto-scale. The default installation mode using
the Aperture Operator should take care of creating the necessary roles and
permissions for Aperture Agent.

## Preview Discovered Control Points

Once the Aperture Agent is installed, it starts discovering control points,
which represent the Kubernetes Resources that can be scaled. This would include
Deployments, StatefulSets and any Custom Resources which are scalable.

The discovered control points can be viewed in the [FluxNinja ARC](/arc/arc.md)
UI. Navigate to the **Control Points** page and select the **Kubernetes** tab.
You should see a list of discovered control points. Alternatively, you can use
the
[introspection API](/reference/api/agent/flow-preview-service-preview-flow-labels.api.mdx)
or
[aperturectl](/reference/aperturectl/auto-scale/control-points/control-points.md)
to view this information.

## Configure Auto Scaling Policies

Auto-scaling policies are configured at the Aperture Controller. Refer to
[tutorials](/use-cases/auto-scale/auto-scale.md) for some example policies.
