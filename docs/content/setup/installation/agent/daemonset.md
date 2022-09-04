---
title: DaemonSet
description: Install Aperture Agent in DaemonSet mode
keywords:
  - install
  - setup
  - agent
  - daemonset
---

The Aperture Agent can be installed as a
[Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
where it will get deployed on all the nodes of the cluster.

## Installation {#agent-daemonset-installation}

(Consult [Supported Platforms](setup/supported-platforms.md) before installing.)

The Aperture Agent in the DaemonSet mode will be installed using the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which will be managed by the Aperture Operator.

Below are the steps to install or upgrade the Aperture Agent into your setup
using the
[Aperture Agent Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-agent).

By following these instructions, you will have deployed the Aperture Agent into
your cluster.

1. Add the Helm chart repo in your environment:

   ```bash
   helm repo add aperture https://fluxninja.github.io/aperture/
   helm repo update
   ```

2. Configure the required parameters of Etcd and Prometheus for the Agent Custom
   Resource by creating a `values.yaml` with below parameters and pass it with
   `helm upgrade`:

   ```yaml
   agent:
     etcd:
       endpoints: ["ETCD_ENDPOINT_HERE"]
     prometheus:
       address: "PROMETHEUS_ADDRESS_HERE"
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of Etcd and Prometheus, which is also being used by the
   Aperture Controller you want these Agents to connect with.

   ```bash
   helm upgrade --install agent aperture/aperture-agent -f values.yaml
   ```

3. Alternatively, you can create the Agent Custom Resource directly on the
   Kubernetes cluster using the below steps:

   1. Create a `values.yaml` for just starting the operator and disabling the
      creation of Agent Custom Resource and pass it with `helm upgrade`:

      ```yaml
      agent:
        create: false
      ```

      ```bash
      helm upgrade --install agent aperture/aperture-agent -f values.yaml
      ```

   2. Create a YAML file with below specifications:

      ```yaml
      apiVersion: fluxninja.com/v1alpha1
      kind: Agent
      metadata:
        name: Agent
      spec:
        etcd:
          endpoints: ["ETCD_ENDPOINT_HERE"]
        prometheus:
          address: "PROMETHEUS_ADDRESS_HERE"
      ```

      Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE`
      with the actual values of Etcd and Prometheus, which is also being used by
      the Aperture Controller you want these Agents to connect with.

      All the configuration parameters for the Agent Custom Resource are listed
      on the
      [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#agent-custom-resource-parameters)
      file of the Helm chart.

   3. Apply the YAML file to Kubernetes cluster using `kubectl`

      ```bash
      kubectl apply -f agent.yaml
      ```

4. If you want to modify the default parameters, you can create or update the
   `values.yaml` file and pass it with `helm upgrade`:

   ```bash
   helm upgrade --install agent aperture/aperture-agent -f values.yaml
   ```

   A list of configurable parameters can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#parameters).

5. If you want to deploy the Aperture Agent into a namespace other than
   `default`, use the `-n` flag:

   ```bash
   helm upgrade --install agent aperture/aperture-agent -n "aperture-controller" --create-namespace
   ```

6. Once you have successfully deployed the resources, confirm that the Aperture
   Agent is up and running:

   ```bash
   kubectl get pod -A

   kubectl get agent -A
   ```

You should see pods for Aperture Agent and Agent Manager in `RUNNING` state and
`Agent` Custom Resource in `created` state.

7. Refer steps on the [Istio Configuration](setup/istio.md) if you don't have
   the
   [Envoy Filter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)
   configured on your cluster.
