---
title: Aperture Agent
description: Install Aperture Agent
keywords:
  - install
  - setup
  - agent
  - sidecar
  - daemonset
sidebar_position: 3
---

## Overview

The Aperture Agent is the decision executor of the Aperture system. In addition
to gathering data, the Aperture Agent functions as a gatekeeper, acting on
traffic based on periodic adjustments made by the Aperture Controller.
Specifically, depending on feedback from the controller, the agent will
effectively allow or drop incoming requests. Further, supporting the controller,
the agent works to inject information into traffic, including the specific
traffic-shaping decisions made and classification labels which can later be used
for observability and closed loop feedback.

## Configuration

All the configuration parameters for the Aperture Agent are listed
[here](/reference/configuration/agent.md).

## Installation Modes {#agent-installation-modes}

The Aperture Agent can be installed in the following modes:

:::caution warning

Upgrading from one of the installation modes below to the other is discouraged
and can result in unpredictable behavior.

:::

1. **Kubernetes**

   1. [**Install with Operator**](kubernetes/operator/operator.md)

      The Aperture Agent can be installed using the Kubernetes Operator
      available for it. This method requires access to create cluster level
      resources like ClusterRole, ClusterRoleBinding, CustomResourceDefinition
      and so on.

      - [**DaemonSet**](kubernetes/operator/daemonset.md)

        The Aperture Agent can be installed as a
        [Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
        where it will get deployed on all the nodes of the cluster.

      - [**Sidecar**](kubernetes/operator/sidecar.md)

        The Aperture Agent can also be installed as a Sidecar. In this mode,
        whenever a new pod is started with required labels and annotations, the
        agent container will be attached with the pod.

   2. [**Namespace-Scoped Installation**](kubernetes/namespace-scoped/namespace-scoped.md)

      The Aperture Agent can also be installed with only namespace-scoped
      resources.

2. [**Bare Metal or VM**](bare_metal.md)

   The Aperture Agent can be installed as a system service on any Linux system
   that is [supported](../supported-platforms.md).

3. [**Docker**](docker.md)

   The Aperture Agent can also be installed on Docker as containers.

## Applying Policies

Once the
[application is set up](/get-started/setting-up-application/setting-up-application.md)
and both the Aperture Controller and Agents are installed, the next crucial step
is to create and apply policies.

[Your first policy](/get-started/policies/policies.md) section provides
step-by-step instructions on customizing, creating, and applying policies within
Aperture. Additionally, the [use-cases](/use-cases/use-cases.md) section serves
as a valuable resource for tailoring policies to meet specific requirements.
