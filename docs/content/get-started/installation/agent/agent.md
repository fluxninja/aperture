---
title: Aperture Agent
description: Install Aperture Agent
keywords:
  - install
  - setup
  - agent
  - sidecar
  - daemonset
sidebar_position: 2
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
[here](/reference/configuration/agent.md).

## Installation Modes {#agent-installation-modes}

The Aperture Agent can be installed in the following modes:

1. **Kubernetes**

   1. [**DaemonSet**](kubernetes/daemonset.md)

      The Aperture Agent can be installed as a
      [Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
      where it will get deployed on all the nodes of the cluster.

   2. [**Sidecar**](kubernetes/sidecar.md)

      The Aperture Agent can also be installed as a Sidecar. In this mode,
      whenever a new pod is started with required labels and annotations, the
      Agent container will be attached with the pod.

2. [**Bare Metal/VM**](bare_metal.md)

   The Aperture Agent can be installed as a system service on any Linux system
   that is [supported](supported-platforms.md).
