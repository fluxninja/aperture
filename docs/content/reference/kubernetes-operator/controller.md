---
title: Aperture Controller Operator
description: Kubernetes Objects installed with Aperture Controller Operator
keywords:
  - operator
  - kubernetes
  - controller
sidebar_position: 1
---

Aperture Operator for Aperture Controller manages the Kubernetes Objects
required for it.

## Kubernetes Objects for Operator

Below are the Kubernetes Objects which get created for the Operator considering
`controller` is the release name and `default` is the namespace:

<!-- vale off -->

| API Version                     | Kind                         | Name                                    | Namespace |
| ------------------------------- | ---------------------------- | --------------------------------------- | --------- |
| apiextensions.k8s.io/v1         | CustomResourceDefinition     | controllers.fluxninja.com               | Global    |
| apiextensions.k8s.io/v1         | CustomResourceDefinition     | policies.fluxninja.com                  | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRole                  | controller-kube-state-metrics           | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRole                  | controller-prometheus-server            | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRole                  | controller-aperture-controller-operator | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRoleBinding           | controller-kube-state-metrics           | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRoleBinding           | controller-prometheus-server            | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRoleBinding           | controller-aperture-controller-operator | Global    |
| admissionregistration.k8s.io/v1 | MutatingWebhookConfiguration | aperture-controller-defaulter           | Global    |
| policy/v1                       | PodDisruptionBudget          | controller-etcd                         | default   |
| v1                              | ServiceAccount               | controller-kube-state-metrics           | default   |
| v1                              | ServiceAccount               | controller-prometheus-server            | default   |
| v1                              | ServiceAccount               | controller-aperture-controller-operator | default   |
| v1                              | ConfigMap                    | controller-prometheus-server            | default   |
| v1                              | PersistentVolumeClaim        | controller-prometheus-server            | default   |
| v1                              | Service                      | controller-etcd-headless                | default   |
| v1                              | Service                      | controller-etcd                         | default   |
| v1                              | Service                      | controller-kube-state-metrics           | default   |
| v1                              | Service                      | controller-prometheus-server            | default   |
| v1                              | Service                      | controller-aperture-controller-manager  | default   |
| apps/v1                         | Deployment                   | controller-kube-state-metrics           | default   |
| apps/v1                         | Deployment                   | controller-aperture-controller-manager  | default   |
| apps/v1                         | StatefulSet                  | controller-etcd                         | default   |
| fluxninja.com/v1alpha1          | Controller                   | controller                              | default   |

<!-- vale on -->

## Kubernetes Objects by Operator

Below are the Kubernetes Objects which are created by the Operator, considering
`controller` is the Custom Resource name and `default` is the namespace:

<!-- vale off -->

| API Version                     | Kind                           | Name                              | Namespace |
| ------------------------------- | ------------------------------ | --------------------------------- | --------- |
| rbac.authorization.k8s.io/v1    | ClusterRole                    | aperture-controller               | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRoleBinding             | aperture-controller               | Global    |
| admissionregistration.k8s.io/v1 | ValidatingWebhookConfiguration | aperture-controller               | Global    |
| v1                              | ConfigMap                      | aperture-controller               | default   |
| v1                              | ConfigMap                      | controller-controller-client-cert | default   |
| v1                              | Service                        | aperture-controller               | default   |
| v1                              | Secret                         | controller-controller-cert        | default   |
| v1                              | ServiceAccount                 | aperture-controller               | default   |
| apps/v1                         | Deployment                     | aperture-controller               | default   |

<!-- vale on -->
