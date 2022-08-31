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

1. The Aperture Agent can be installed in the DaemonSet mode using the Helm chart of Aperture Operator
   by using the default `values.yaml` or create a `values.yaml` with below parameters and pass it with `helm upgrade`:

   ```yaml
   agent:
     create: true
   ```

   ```bash
   helm upgrade --install aperture aperture/operator -f values.yaml
   ```

2. Alternatively, you can create the Agent Custom Resource directly on the Kubernetes cluster using the below YAML:

   1. Create a YAML file with below specifications:

      ```yaml
      apiVersion: fluxninja.com/v1alpha1
      kind: Agent
      metadata:
        name: agent
      spec:
        image:
          registry: docker.io/fluxninja
          repository: aperture-agent
          tag: latest
      ```

   All the configuration parameters for the Agent Custom Resource are listed on the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-operator#aperture-custom-resource-parameters)
   file of the Helm chart.

   2. Apply the YAML file to Kubernetes cluster using `kubectl`

      ```bash
      kubectl apply -f agent.yaml
      ```

3. Once you have successfully deployed the Helm release, confirm that the
   Aperture Agent is up and running:

   ```bash
   kubectl get pod -A
   ```

   You should see pods for Aperture Agent in `RUNNING` state.
