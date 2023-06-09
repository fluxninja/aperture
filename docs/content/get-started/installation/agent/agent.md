---
title: Aperture Agent
description: Install Aperture Agent
keywords:
  - install
  - setup
  - agent
  - sidecar
  - daemonset
sidebar_position: 4
---

## Overview

The Aperture Agent is the decision executor of the Aperture system. In addition
to gathering data, the Aperture Agent functions as a gatekeeper, acting on
traffic based on decisions made by the Aperture Controller. Specifically,
depending on feedback from the Controller, the Agent will effectively allow or
drop incoming requests. Further, supporting the Controller, the Agent works to
inject information into traffic, including the specific traffic-shaping
decisions made and classification labels which can later be used in policing.
One Agent is deployed per node.

## Configuration

All the configuration parameters for Aperture Agent are listed
[here][agent-configuration].

## Installation Modes {#agent-installation-modes}

The Aperture Agent can be installed in the following modes:

:::caution warning

Upgrading from one of the installation modes below to the other is discouraged
and can result in unpredictable behavior.

:::

1. **Kubernetes**

   1. [**Install with Operator**][agent-operator-installation]

      The Aperture Agent can be installed using the Kubernetes Operator
      available for it. This method requires access to create cluster level
      resources like ClusterRole, ClusterRoleBinding, CustomResourceDefinition
      and so on.

      - [**DaemonSet**][agent-daemonset]

      The Aperture Agent can be installed as a [Kubernetes
      DaemonSet][kubernetes-daemonset], where it will get deployed on all the
      nodes of the cluster.

      - [**Sidecar**][agent-sidecar]

      The Aperture Agent can also be installed as a Sidecar. In this mode,
      whenever a new pod is started with required labels and annotations, the
      Agent container will be attached with the pod.

   2. [**Namespace-Scoped Installation**][namespace-scoped-installation]

      The Aperture Agent can also be installed with only namespace-scoped
      resources.

2. [**Bare Metal or VM**][bare-metal]

   The Aperture Agent can be installed as a system service on any Linux system
   that is [supported][supported-platforms].

3. [**Docker**][docker]

   The Aperture Agent can also be installed on Docker as containers.

[docker]: ./docker.md
[bare-metal]: ./bare_metal.md
[agent-configuration]: /reference/configuration/agent.md
[namespace-scoped-installation]:
  ./kubernetes/namespace-scoped/namespace-scoped.md
[agent-sidecar]: ./kubernetes/operator/sidecar.md
[kubernetes-daemonset]:
  https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
[agent-daemonset]: ./kubernetes/operator/daemonset.md
[agent-operator-installation]: ./kubernetes/operator/operator.md
[supported-platforms]: ../supported-platforms.md
