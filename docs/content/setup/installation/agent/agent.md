---
title: Aperture Agent
description: Install Aperture Agent
keywords:
  - install
  - setup
  - agent
  - sidecar
  - daemonset
---

## Overview

The Aperture Agent is the decision executor of the Aperture system. In addition
to gathering data, the Aperture Agent functions as a gatekeeper, acting on
traffic based on decisions made by the Aperture Controller. Specifically,
depending on feedback from the Controller, the Agent will effectively allow or
drop incoming requests. Further supporting the Controller, the Agent works to
inject information into traffic, including the specific traffic-shaping
decisions made and classification labels which can later be used in policing.
One Agent is deployed per node.

## Configuration

The Aperture Agent related configurations are stored in a configmap which is
created during the installation using Helm. All the configuration parameters are
listed on the
[README](https://artifacthub.io/packages/helm/aperture/aperture-operator#aperture-custom-resource-parameters)
file of the Helm chart.

## Installation Modes {#agent-installation-modes}

The Aperture Agent can be installed in below listed modes:

1. **Kubernetes**

   1. **DaemonSet**

      The Aperture Agent can be installed as a
      [Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
      where it will get deployed on all the nodes of the cluster.

   2. **Sidecar**

      The Aperture Agent can also be installed as a Sidecar. In this mode,
      whenever a new pod is started with required labels and annotations, the
      Agent container will be attached with the pod.
