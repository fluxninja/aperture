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

The Aperture Controller functions as the brain of the Aperture system.
Leveraging an advanced control loop, the Controller routinely analyzes polled
metrics and indicators to determine how traffic should be shaped as defined by
set policies. Once determined, these decisions are then exported to all Aperture
Agents to effectively handle workloads. Only one Controller is needed to manage
each cluster.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (i.e. process variables) and compares them
against set points. The difference in the variable values against these points
is referred to as the error signal. The feedback loop then works to minimize
these error signals by determining and distributing control actions, that adjust
these process variables and maintain their values within the optimal range.

## Configuration

The Aperture Controller related configurations are stored in a configmap which
is created during the installation using Helm.

All the configuration parameters for Aperture Controller are listed
[here](/reference/configuration/controller.md).

## Installation {#controller-installation}

(Consult [Supported Platforms](/get-started/supported-platforms.md) before
installing.)

The Aperture Controller will be installed using the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which will be managed by the Aperture Operator.

Below are the steps to install or upgrade the Aperture Controller into your
setup using the
[Aperture Controller Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-controller).

By following these instructions, you will have deployed the Aperture Controller
into your cluster.

1. Add the Helm chart repo in your environment:

   ```bash
   helm repo add aperture https://fluxninja.github.io/aperture/
   helm repo update
   ```

2. The Aperture Controller can be installed using the default `values.yaml`:

   ```bash
   helm upgrade --install controller aperture/aperture-controller
   ```

3. Alternatively, you can create the Controller Custom Resource directly on the
   Kubernetes cluster using the below steps:

   1. Create a `values.yaml` for starting the operator and disabling the
      creation of Controller Custom Resource and pass it with `helm upgrade`:

      ```yaml
      controller:
        create: false
      ```

      ```bash
      helm upgrade --install controller aperture/aperture-controller -f values.yaml
      ```

   2. Create a YAML file with below specifications:

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

      All the configuration parameters for the Controller Custom Resource are
      listed on the
      [README](https://artifacthub.io/packages/helm/aperture/aperture-controller#controller-custom-resource-parameters)
      file of the Helm chart.

   3. Apply the YAML file to the Kubernetes cluster using `kubectl`

      ```bash
      kubectl apply -f controller.yaml
      ```

4. The chart installs Prometheus and Etcd instances by default. If you don't
   want to install and use your existing instances of Prometheus or Etcd,
   configure below values in the `values.yaml` file and pass it with
   `helm upgrade`:

   ```yaml
   controller:
     config:
       etcd:
         endpoints: ["ETCD_ENDPOINT_HERE"]
       prometheus:
         address: "PROMETHEUS_ADDRESS_HERE"

   etcd:
     enabled: false

   prometheus:
     enabled: false
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of Etcd and Prometheus, which will be used by the Aperture
   Controller.

   ```bash
   helm upgrade --install controller aperture/aperture-controller -f values.yaml
   ```

   A list of all the configurable parameters for Etcd are available
   [here](/reference/configuration/controller.md#etcd) and Prometheus are
   available [here](/reference/configuration/controller.md#prometheus).

   **Note**: Please make sure that the flag `web.enable-remote-write-receiver`
   is enabled on your existing Prometheus instance as it is required by the
   Controller.

5. If you want to modify the default parameters, you can create or update the
   `values.yaml` file and pass it with `helm upgrade`:

   ```bash
   helm upgrade --install controller aperture/aperture-controller -f values.yaml
   ```

   A list of configurable parameters can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-controller#parameters).

6. If you want to deploy the Aperture Controller into a namespace other than
   `default`, use the `-n` flag:

   ```bash
   helm upgrade --install controller aperture/aperture-controller -n aperture-controller --create-namespace
   ```

## Verifying the Installation

Once you have successfully deployed the resources, confirm that the Aperture
Controller is up and running:

```bash
kubectl get pod -A

kubectl get controller -A
```

You should see pods for Aperture Controller and Controller Manager in `RUNNING`
state and `Controller` Custom Resource in `created` state.

## Applying Policies

Follow the information on [Policy](/concepts/policy/circuit/sources.md) to
understand and design the policy circuits.

Once the design is ready, follow the steps on the
[Blueprints](/get-started/blueprints.md) to generate the Policy ConfigMap and
apply it on the Kubernetes.

## Uninstall

You can uninstall the Aperture Controller and it's components by uninstalling
the charts installed above:

1. Delete th Aperture Controller chart:

   ```bash
   helm uninstall controller
   ```

2. Alternativey, if you have installed the Aperture Controller Custom Resource
   separately, follow below steps:

   1. Delete the Aperture Controller Custom Resource:

      ```bash
      kubectl delete -f controller.yaml
      ```

   2. Delete the Aperture Controller chart to uninstall the Aperture Operator:

      ```bash
      helm uninstall controller
      ```

3. If you have installed the chart in some other namespace than `default`,
   execute below commands:

   ```bash
   helm uninstall controller -n aperture-controller
   kubectl delete namespace aperture-controller
   ```

   > Note: By design, deleting a chart via Helm doesn’t delete the Custom
   > Resource Definitions (CRDs) installed via the chart.

4. **Optional**: Delete the CRD installed by the chart:

   ```bash
   kubectl delete crd controllers.fluxninja.com
   ```
