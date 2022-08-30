---
title: Aperture Controller
description: Install Aperture Controller
keywords:
  - install
  - setup
  - controller
sidebar_position: 1
---

## Overview

The Aperture Controller functions as the primary decision maker of the system.
Leveraging our advanced control loop, the controller routinely analyzes polled
metrics and indicators to determine how traffic should be shaped as defined by
set policies. Once determined, these decisions are then exported to all agents
in order to effectively handle workloads. Only one controller is needed to
effectively manage each cluster.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (i.e. process variables) and compares them
against set points. The difference in the variable values against these points
is referred to as the error signal. The feedback loop then works to minimize
these error signals by determining and distributing specific decisions, or
control actions, that adjust these process variables and maintain their values
within the optimal range.

## Configuration

The Aperture Controller related configurations are stored in a configmap which
is created during the installation using Helm. All the configuration parameters
are listed on the
[README](https://artifacthub.io/packages/helm/aperture/aperture-operator#aperture-custom-resource-parameters)
file of the Helm chart.

## Installation {#controller-installation}

(Consult [Supported Platforms](setup/supported-platforms.md) before installing.)

**Note**: Make sure you have the [Aperture Operator](setup/installation.md#installation-operator-installation) in running state before following the below steps.

1. The Aperture Controller can be installed using the Helm chart of Aperture Operator
   by using the default `values.yaml` or create a `values.yaml` with below parameters and pass it with `helm upgrade`:

   ```yaml
   controller:
     create: true
   ```

   ```bash
   helm upgrade --install aperture aperture/operator -f values.yaml
   ```

2. Alternatively, you can create the Controller Custom Resource directly on the Kubernetes cluster using the below YAML:

   1. Create a YAML file with below specifications:

   ```yaml
   apiVersion: fluxninja.com/v1alpha1
   kind: Controller
   metadata:
     name: controller
   spec:
     image:
       registry: docker.io/fluxninja
       repository: aperture-controller
       tag: latest
   ```

   All the configuration parameters for the Controller Custom Resource are listed on the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-operator#aperture-custom-resource-parameters)
   file of the Helm chart.

   2. Apply the YAML file to Kubernetes cluster using `kubectl`

   ```bash
   kubectl apply -f controller.yaml
   ```

3. Once you have successfully deployed the Custom Resource, confirm that the
   Aperture Controller is up and running:

   ```bash
   kubectl get pod -A

   kubectl get controller -A
   ```

You should see pods for Aperture Controller in `RUNNING` state and `Controller` Custom Resource in `created` state.
