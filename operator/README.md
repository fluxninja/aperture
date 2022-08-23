# Aperture Operator

This document is an overview of how the operator works from a user perspective.

## Introduction

The operator is used to deploy the Aperture Agent, Controller and its required resources on the Kubernetes.

## Custom Resources

The operator has below custome resources:

- Aperture

all custom resources use the api group `aperture.tech` and version `v1alpha1`.

## Deploying the operator using `make`

Follow the below steps to deploy the operator on the local cluster:

- Create Docker image for the operator:

  ```bash
  make docker-build
  ```

- [**Optional**] If you are using [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) cluster, upload the image to the cluster:

  ```bash
  kind load docker-image aperture-operator:latest
  ```

- Deploy the operator and its required resources:

  ```bash
  make deploy
  ```

- [**Optional**] Run the sample CR to deploy the Aperture Agent and Controller

  ```bash
  kubectl apply -f config/samples/v1alpha1_aperture.yaml
  ```

- To uninstall the operator, run below commands:

  ```bash
  kubectl delete -f config/samples/v1alpha1_aperture.yaml
  make undeploy
  ```

## Deploying the operator using `helm`

You can also deploy the operator and Aperture CR using the helm chart as well using below steps.

- Install the dependencies of the chart:

  ```bash
  helm dependency build manifests/charts/aperture-operator
  ```

- Install or upgrade the chart:

  ```bash
  helm upgrade --install aperture-operator manifests/charts/aperture-operator
  ```

- [**Optional**] If you want to install just the operator and not the Aperture CR, create a `values.yaml` with below parameters and pass it with `helm upgrade`:

  ```bash
  aperture:
    create: false
  ```

  ```bash
  helm upgrade --install aperture-operator manifests/charts/aperture-operator -f values.yaml
  ```

  All the configurable parameters for the Aperture Operator and CR can be found [here](./manifests/charts/aperture-operator/README.md).

- The chart installs Istio, Prometheus and Etcd instances by default. If you don't want to install and use your existing instances of Istio, Prometheus or Etcd, configure below values in the `values.yaml` file and pass it with `helm upgrade`:

  ```bash
  etcd:
    enabled: false
    endpoints: ["ETCD_INSTANCE_ENDPOINT"]

  prometheus:
    enabled: false
    address: "PROMETHEUS_INSTANCE_ADDRESS"

  istio:
    enabled: false
  ```

- To uninstall the operator, run below commands:

  ```bash
  helm uninstall aperture-operator
  ```
