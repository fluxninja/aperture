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
[install Aperture Agent](/aperture-for-infra/agent/kubernetes/operator/operator.md)
on your Kubernetes cluster. The Agent needs certain permissions to discover
scalable resources and perform auto-scale. The default installation mode using
the Aperture Operator should take care of creating the necessary roles and
permissions for Aperture Agent.

## Preview Discovered Control Points

Once the Aperture Agent is installed, it starts discovering control points,
which represent the Kubernetes Resources that can be scaled. This would include
Deployments, StatefulSets and any Custom Resources which are scalable.

The discovered control points can be viewed in the
[Aperture Cloud](/introduction.md) UI. Navigate to the **Control Points** page
and select the **Kubernetes** tab. You should see a list of discovered control
points. Alternatively, you can use the
[introspection API](/reference/api/agent/auto-scale-kubernetes-control-points-service-get-control-points.api.mdx)
or
[aperturectl](/reference/aperture-cli/aperturectl/auto-scale/control-points/control-points.md)
to view this information.

For example:

```sh
aperturectl auto-scale control-points --kube
```

Returns:

```json
AGENT GROUP   NAME                                                NAMESPACE             KIND
default       coredns                                             kube-system           Deployment
default       coredns-5d78c9869d                                  kube-system           ReplicaSet
default       gateway                                             istio-system          Deployment
default       gateway-868c757988                                  istio-system          ReplicaSet
default       istiod                                              istio-system          Deployment
default       istiod-6d9df7fb7                                    istio-system          ReplicaSet
default       local-path-provisioner                              local-path-storage    Deployment
default       local-path-provisioner-6bc4bddd6b                   local-path-storage    ReplicaSet
default       service1-demo-app                                   demoapp               Deployment
default       service1-demo-app-7b4bc9bdcd                        demoapp               ReplicaSet
default       service2-demo-app                                   demoapp               Deployment
default       service2-demo-app-677bb57574                        demoapp               ReplicaSet
default       service3-demo-app                                   demoapp               Deployment
default       service3-demo-app-58656dcf95                        demoapp               ReplicaSet
default       wavepool-generator                                  demoapp               Deployment
default       wavepool-generator-5b4578bdd9                       demoapp               ReplicaSet
```
