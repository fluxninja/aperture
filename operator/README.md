# Aperture Operator

This document is an overview of how the operator works from a user perspective.

## Introduction

The operator is used to deploy the Aperture Agent, Controller and its required resources on the Kubernetes.

## Custom Resources

The operator has below custome resources:

- Agent
- Controller

all custom resources use the api group `fluxninja.com` and version `v1alpha1`.

## Deploying the operator using `make`

Follow the below steps to deploy the operator on the local cluster:

- Create Docker image for the operator:

  ```bash
  make operator-docker-build
  ```

- [**Optional**] If you are using [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) cluster, upload the image to the cluster:

  ```bash
  kind load docker-image aperture-operator:latest
  ```

- Deploy the operator and its required resources:

  ```bash
  make operator-deploy
  ```

- [**Optional**] Run the sample CR to deploy the Aperture Agent and Controller

  ```bash
  kubectl apply -f config/samples/fluxninja.com_v1alpha1_agent.yaml
  kubectl apply -f config/samples/fluxninja.com_v1alpha1_controller.yaml
  ```

- To uninstall the operator, run below commands:

  ```bash
  kubectl delete -f config/samples/fluxninja.com_v1alpha1_agent.yaml
  kubectl delete -f config/samples/fluxninja.com_v1alpha1_controller.yaml
  make undeploy
  ```

## Deploying the operator using `helm`

You can also deploy the operator and Agent CR using the helm chart as well using below steps.

- Install the dependencies of the chart:

  ```bash
  helm dependency build manifests/charts/aperture-controller
  helm dependency build manifests/charts/aperture-agent
  ```

- Configure Etcd configuration for Agent by modifying the Aperture Agent chart as below if you are installing the Aperture Controller chart as well:

  ```yaml
  agent:
    etcd:
      endpoints: ["http://controller-etcd:2379"]
    prometheus:
      address: "http://controller-prometheus-server:80"
  ```

- Install or upgrade the chart:

  ```bash
  helm upgrade --install controller manifests/charts/aperture-controller
  helm upgrade --install agent manifests/charts/aperture-agent
  ```

- [**Optional**] If you want to install just the operator and not the Aperture Agent and Controller CR, create a `values.yaml` with below parameters and pass it with `helm upgrade`:

  ```yaml
  agent:
    create: false
  ```

  ```yaml
  controller:
    create: false
  ```

  ```bash
  helm upgrade --install controller manifests/charts/aperture-controller -f controller-values.yaml
  helm upgrade --install agent manifests/charts/aperture-agent -f agent-values.yaml
  ```

  All the configurable parameters for the Aperture Operator and CR can be found at
  [Agent](./manifests/charts/aperture-agent/README.md).
  and [Controller](./manifests/charts/aperture-controller/README.md)

- The Controller chart installs Prometheus and Etcd instances by default. If you don't want to install and use your existing instances of Prometheus or Etcd, configure below values in the `values.yaml` file and pass it with `helm upgrade`:

  ```yaml
  controller:
    etcd:
      endpoints: ["ETCD_ENDPOINT"]
    prometheus:
      address: "PROMETHEUS_ENDPOINT"

  etcd:
    enabled: false

  prometheus:
    enabled: false
  ```

  ```yaml
  agent:
    etcd:
      endpoints: ["ETCD_ENDPOINT"]
    prometheus:
      address: "PROMETHEUS_ENDPOINT"
  ```

- To uninstall the operator, run below commands:

  ```bash
  helm uninstall controller
  helm uninstall agent
  ```
