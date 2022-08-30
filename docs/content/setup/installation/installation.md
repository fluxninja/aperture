---
title: Installation of Aperture
description: Install Aperture Operator, Agent and Controller
keywords:
  - install
  - setup
  - agent
  - controller
  - operator
---

## Aperture Operator

### Overview

The Aperture Operator is a Kubernetes Operator used to deploy all the required resources
for Aperture Agent and Controller via
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).

### Configuration

The Aperture Operator related configurations can be passed via the `values.yaml` file
during the installation using Helm. All the configuration parameters
are listed on the
[README](https://artifacthub.io/packages/helm/aperture/aperture-operator#operator-parameters)
file of the Helm chart.

### Installation {#operator-installation}

(Consult [Supported Platforms](setup/supported-platforms.md) before installing.)

Below are the steps to install or upgrade the Aperture Operator into your setup using
the [Aperture Operator Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-operator).

By following these instructions, you will have deployed the Aperture Operator into your cluster.

1. Add the Helm chart repo in your environment:

   ```bash
   helm repo add aperture https://fluxninja.github.io/aperture/
   helm repo update
   ```

2. Install or upgrade the chart:

   ```bash
   helm upgrade --install aperture aperture/operator
   ```

3. If you want to install just the operator and not the Agent and Controller Custom Resources,
   create a `values.yaml` with below parameters and pass it with `helm upgrade`:

   ```yaml
   agent:
     create: false

   controller:
     create: false
   ```

   ```bash
   helm upgrade --install aperture aperture/operator -f values.yaml
   ```

4. The chart installs Prometheus and Etcd instances by default. If you
   don't want to install and use your existing instances of Prometheus or
   Etcd, configure below values in the `values.yaml` file and pass it with
   `helm upgrade`:

   ```yaml
   etcd:
     enabled: false

   prometheus:
     enabled: false
   ```

   ```bash
   helm upgrade --install aperture aperture/operator -f values.yaml
   ```

   A list of other configurable parameters for Etcd and Prometheus can be
   found in the [README](https://artifacthub.io/packages/helm/aperture/aperture-operator#istio).

   **Note**: Please make sure that the flag `web.enable-remote-write-receiver`
   is enabled on your existing Prometheus instance as it is required by the
   Agent.

5. If you want to modify the default parameters, you can create or update the
   `values.yaml` file and pass it with `helm upgrade`:

   ```bash
   helm upgrade --install aperture aperture/operator -f values.yaml
   ```

   A list of configurable parameters can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-operator#operator-parameters).

6. If you want to deploy the Aperture Operator into a namespace
   other than `default`, use the `-n` flag:

   ```bash
   NAMESPACE="aperture-system"; helm upgrade --install aperture aperture/operator -f values.yaml --set global.istioNamespace=$NAMESPACE -n $NAMESPACE --create-namespace
   ```

7. Once you have successfully deployed the Helm release, confirm that the
   Aperture Agent and Controller are up and running:

   ```bash
   kubectl get pod -A
   ```

   You should see pods for Prometheus, Etcd, and Aperture Operator in `RUNNING` state.

8. Refer steps on the [Istio Configuration](setup/istio.md) if you don't have the
   [Envoy Filter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)
   configured on your cluster.
