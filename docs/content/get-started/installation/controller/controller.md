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
Agents to effectively handle workloads.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (i.e. process variables) and compares them
against set points. The difference in the variable values against these points
is referred to as the error signal. The feedback loop then works to minimize
these error signals by determining and distributing control actions, that adjust
these process variables and maintain their values within the optimal range.

## Controller CRD

The Aperture Controller is a Kubernetes based application and is installed using
the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which is managed by the Aperture Operator.

The configuration for the Aperture Controller process is provided to the
Controller CRD under the `controller.config` section. All the configuration
parameters for Aperture Controller are listed
[here](/references/configuration/controller.md).

## Installation {#controller-installation}

Below are the steps to install the Aperture Controller into your setup using the
[Aperture Controller Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-controller).

By following these instructions, you will have deployed the Aperture Controller
into your cluster.

1. Add the Helm chart repo in your environment:

   ```bash
   helm repo add aperture https://fluxninja.github.io/aperture/
   helm repo update
   ```

2. Configure the below parameters for the Controller Custom Resource by creating
   a `values.yaml` with below parameters and pass it with `helm install`:

   :::info

   The below parameters disable the FluxNinja ARC Plugin for the Aperture
   Controller. If you want to keep it enabled, add parameters provided
   [here](/arc/plugin.md#configuration) under the `controller.config` section.

   :::

   ```yaml
   controller:
     config:
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
   ```

   ```bash
   helm install controller aperture/aperture-controller -f values.yaml
   ```

3. Alternatively, you can create the Controller Custom Resource directly on the
   Kubernetes cluster using the below steps:

   1. Create a `values.yaml` for starting the operator and disabling the
      creation of Controller Custom Resource and pass it with `helm install`:

      ```yaml
      controller:
        create: false
      ```

      ```bash
      helm install controller aperture/aperture-controller -f values.yaml
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
        config:
          etcd:
            endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
          prometheus:
            address: "http://controller-prometheus-server.default.svc.cluster.local:80"
          plugins:
            disabled_plugins:
              - aperture-plugin-fluxninja
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
   `helm install`:

   ```yaml
   controller:
     config:
       etcd:
         endpoints: ["ETCD_ENDPOINT_HERE"]
       prometheus:
         address: "PROMETHEUS_ADDRESS_HERE"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja

   etcd:
     enabled: false

   prometheus:
     enabled: false
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of Etcd and Prometheus, which will be used by the Aperture
   Controller.

   ```bash
   helm install controller aperture/aperture-controller -f values.yaml
   ```

   A list of all the configurable parameters for Etcd are available
   [here](/references/configuration/controller.md#etcd) and Prometheus are
   available [here](/references/configuration/controller.md#prometheus).

   **Note**: Please make sure that the flag `web.enable-remote-write-receiver`
   is enabled on your existing Prometheus instance as it is required by the
   Aperture Controller.

5. If you want to modify the default parameters or the Aperture Controller
   config, for example `log`, you can create or update the `values.yaml` file
   and pass it with `helm install`:

   ```yaml
   controller:
     config:
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
       log:
         level: debug
         pretty_console: true
         non_blocking: false
   ```

   ```bash
   helm install controller aperture/aperture-controller -f values.yaml
   ```

   All the config parameters for the Aperture Controller are available
   [here](/references/configuration/controller.md).

   A list of configurable parameters for the installation can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-controller#parameters).

6. If you want to deploy the Aperture Controller into a namespace other than
   `default`, use the `-n` flag:

   ```bash
   helm install controller aperture/aperture-controller -n aperture-controller --create-namespace
   ```

## Exposing Etcd and Prometheus services {#expose-etcd-prometheus}

If the Aperture Controller is installed with the packaged Etcd and Prometheus,
follow below steps to expose them outside of the Kubernetes cluster so that the
Aperture Agent running on Linux can access them.

:::info

[Contour](https://projectcontour.io/) is used as a
[Kubernetes Ingress Controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
in below steps to expose the Etcd and Prometheus services out of Kubernetes
cluster using Load Balancer.

Any other tools can also be used to expose the Etcd and Prometheus services out
of the Kubernetes cluster based on your infrastructure.

:::

1. Add the Helm chart repo for Contour in your environment:

   ```bash
   helm repo add bitnami https://charts.bitnami.com/bitnami
   ```

2. Install the Contour chart by running the following command:

   ```bash
   helm install aperture bitnami/contour --namespace projectcontour --create-namespace
   ```

3. It may take a few minutes for the Contour Load Balancer's IP to become
   available. You can watch the status by running:

   ```bash
   kubectl get svc aperture-contour-envoy --namespace projectcontour -w
   ```

4. Once `EXTERNAL-IP` is no longer `<pending>`, run below command to get the
   External IP for the Load Balancer:

   ```bash
   kubectl describe svc aperture-contour-envoy --namespace projectcontour | grep Ingress | awk '{print $3}'
   ```

5. Add an entry for the above IP in the Cloud provider's DNS configuration. For
   example, follow steps on
   [Cloud DNS on GKE](https://cloud.google.com/dns/docs/records) for Google
   Kubernetes Engine.

6. Configure the below parameters to install the
   [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
   with the Aperture Controller by updating the `values.yaml` created during
   installation and pass it with `helm upgrade`:

   ```yaml
   ingress:
     enabled: true
     domain_name: YOUR_DOMAIN_HERE

   etcd:
     service:
       annotations:
         projectcontour.io/upstream-protocol.h2c: "2379"
   ```

   Replace the values of `YOUR_DOMAIN_HERE` with the actual value the domain
   name under with the External IP is exposed.

   ```bash
   helm upgrade --install controller aperture/aperture-controller -f values.yaml
   ```

7. It may take a few minutes for the Ingress resource to get the `ADDRESS`. You
   can watch the status by running:

   ```bash
   kubectl get ingress controller-ingress -w
   ```

8. Once the `ADDRESS` matches the External IP, the Etcd will be accessible on
   `http://etcd.YOUR_DOMAIN_HERE:80` and the Prometheus will be accessible on
   `http://prometheus.YOUR_DOMAIN_HERE:80`.

## Upgrade Procedure {#controller-upgrade-procedure}

By following these instructions, you will have deployed the upgraded version of
Aperture Controller into your cluster.

1. Update the Helm chart repo in your environment:

   ```bash
   helm repo update
   ```

2. Use the same `values.yaml` file created as part of
   [Installation Steps](#controller-installation) and pass it with below
   command:

   ```bash
   helm template --include-crds --no-hooks controller aperture/aperture-controller -f values.yaml | kubectl apply -f -
   ```

3. If you have deployed the Aperture Controller into a namespace other than
   `default`, use the `-n` flag:

   ```bash
   helm template --include-crds --no-hooks controller aperture/aperture-controller -f values.yaml -n aperture-controller | kubectl apply -f -
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

The process of creating policies for Aperture can be done either after the
installation of the controller or after the installation of the agent, depending
on your preference.
[Create Policy](/get-started/policies/get-started-with-policy.md#create-policy)
guide includes step-by-step instructions on how to create policies for Aperture
in a Kubernetes cluster, which you can follow to create policies according to
your needs.

## Uninstall

You can uninstall the Aperture Controller and it's components by uninstalling
the chart installed above:

1. Delete the Aperture Controller chart:

   ```bash
   helm uninstall controller
   ```

2. Alternatively, if you have installed the Aperture Controller Custom Resource
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

4. If you have installed the Contour chart for exposing the Etcd and Prometheus
   service, execute the below command:

   ```bash
   helm uninstall aperture -n projectcontour
   kubectl delete namespace projectcontour
   ```

5. **Optional**: Delete the CRD installed by the chart:

   ```bash
   kubectl delete crd controllers.fluxninja.com
   kubectl delete crd policies.fluxninja.com
   ```
