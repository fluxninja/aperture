---
title: Aperture Agent Operator
description: Kubernetes Objects installed with Aperture Agent Operator
keywords:
  - operator
  - kubernetes
  - agent
sidebar_position: 2
---

Aperture Operator for Aperture Agent manages the Kubernetes Objects required for
it.

## Kubernetes Objects for Operator

Below are the Kubernetes Objects which are created for the Operator, considering
`agent` is the release name and `default` is the namespace:

<!-- vale off -->

| API Version                     | Kind                         | Name                          | Namespace |
| ------------------------------- | ---------------------------- | ----------------------------- | --------- |
| apiextensions.k8s.io/v1         | CustomResourceDefinition     | agents.fluxninja.com          | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRole                  | agent-aperture-agent-operator | Global    |
| rbac.authorization.k8s.io/v1    | ClusterRoleBinding           | agent-aperture-agent-operator | Global    |
| admissionregistration.k8s.io/v1 | MutatingWebhookConfiguration | aperture-agent-defaulter      | Global    |
| v1                              | ServiceAccount               | agent-aperture-agent-operator | default   |
| v1                              | Service                      | agent-aperture-agent-manager  | default   |
| apps/v1                         | Deployment                   | agent-aperture-agent-manager  | default   |
| fluxninja.com/v1alpha1          | Agent                        | agent                         | default   |

<!-- vale on -->

## Kubernetes Objects by Operator

Below are the Kubernetes Objects which are created by the Operator, considering
`agent` is the Custom Resource name and `default` is the namespace:

### DaemonSet Mode

<!-- vale off -->

| API Version                  | Kind               | Name           | Namespace |
| ---------------------------- | ------------------ | -------------- | --------- |
| rbac.authorization.k8s.io/v1 | ClusterRole        | aperture-agent | Global    |
| rbac.authorization.k8s.io/v1 | ClusterRoleBinding | aperture-agent | Global    |
| v1                           | ConfigMap          | aperture-agent | default   |
| v1                           | Service            | aperture-agent | default   |
| v1                           | ServiceAccount     | aperture-agent | default   |
| apps/v1                      | DaemonSet          | aperture-agent | default   |

<!-- vale on -->

### Sidecar Mode

<!-- vale off -->

| API Version                     | Kind                         | Name              | Namespace                      |
| ------------------------------- | ---------------------------- | ----------------- | ------------------------------ |
| rbac.authorization.k8s.io/v1    | ClusterRole                  | aperture-agent    | Global                         |
| rbac.authorization.k8s.io/v1    | ClusterRoleBinding           | aperture-agent    | Global                         |
| admissionregistration.k8s.io/v1 | MutatingWebhookConfiguration | aperture-injector | Global                         |
| v1                              | ConfigMap                    | aperture-agent    | All Sidecar enabled namespaces |

<!-- vale on -->
