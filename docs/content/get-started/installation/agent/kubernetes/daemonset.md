---
title: DaemonSet
description: Install Aperture Agent in DaemonSet mode
keywords:
  - install
  - setup
  - agent
  - daemonset
sidebar_position: 1
---

The Aperture Agent can be installed as a
[Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
where it will get deployed on all the nodes of the cluster.

## Upgrade Procedure {#agent-daemonset-upgrade-procedure}

By following these instructions, you will have deployed the upgraded version of
Aperture Agent into your cluster.

1. Update the Helm chart repo in your environment:

   ```bash
   helm repo update
   ```

2. Use the same `values.yaml` file created as part of
   [Installation Steps](#agent-daemonset-installation) and pass it with below
   command:

   ```bash
   helm template --include-crds --no-hooks agent aperture/aperture-agent -f values.yaml | kubectl apply -f -
   ```

3. If you have deployed the Aperture Agent into a namespace other than
   `default`, use the `-n` flag:

   ```bash
   helm template --include-crds --no-hooks agent aperture/aperture-agent -f values.yaml -n aperture-agent | kubectl apply -f -
   ```

## Installation {#agent-daemonset-installation}

The Aperture Agent in the DaemonSet mode will be installed using the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which will be managed by the Aperture Operator.

Below are the steps to install the Aperture Agent into your setup using the
[Aperture Agent Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-agent).

By following these instructions, you will have deployed the Aperture Agent into
your cluster.

1. Add the Helm chart repo in your environment:

   ```bash
   helm repo add aperture https://fluxninja.github.io/aperture/
   helm repo update
   ```

2. Configure the below parameters of Plugins, Etcd and Prometheus for the Agent
   Custom Resource by creating a `values.yaml` with below parameters and pass it
   with `helm install`:

   :::info

   The below parameters disable the FluxNinja Cloud Plugin for the Aperture
   Agent. If you want to keep it enabled, add parameters provided
   [here](/cloud/plugin.md#configuration) under the `agent.config` section.

   :::

   ```yaml
   agent:
     config:
       etcd:
         endpoints: ["ETCD_ENDPOINT_HERE"]
       prometheus:
         address: "PROMETHEUS_ADDRESS_HERE"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of Etcd and Prometheus, which is also being used by the
   Aperture Controller you want these Agents to connect with.

   If you have installed the
   [Aperture Controller](/get-started/installation/controller.md) on the same
   cluster in `default` namespace, with Etcd and Prometheus using `controller`
   as release name, the values for the values for `ETCD_ENDPOINT_HERE` and
   `PROMETHEUS_ADDRESS_HERE` would be as below:

   ```yaml
   agent:
     config:
       etcd:
         endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
       prometheus:
         address: "http://controller-prometheus-server.default.svc.cluster.local:80"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
   ```

   ```bash
   helm install agent aperture/aperture-agent -f values.yaml
   ```

3. Alternatively, you can create the Agent Custom Resource directly on the
   Kubernetes cluster using the below steps:

   1. Create a `values.yaml` for just starting the operator and disabling the
      creation of Agent Custom Resource and pass it with `helm install`:

      ```yaml
      agent:
        create: false
      ```

      ```bash
      helm install agent aperture/aperture-agent -f values.yaml
      ```

   2. Create a YAML file with below specifications:

      ```yaml
      apiVersion: fluxninja.com/v1alpha1
      kind: Agent
      metadata:
        name: Agent
      spec:
        config:
          etcd:
            endpoints: ["ETCD_ENDPOINT_HERE"]
          prometheus:
            address: "PROMETHEUS_ADDRESS_HERE"
          plugins:
            disabled_plugins:
              - aperture-plugin-fluxninja
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

4. If you want to modify the default parameters or the Aperture Agent config,
   for example `log`, you can create or update the `values.yaml` file and pass
   it with `helm install`:

   ```yaml
   agent:
     config:
       etcd:
         endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
       prometheus:
         address: "http://controller-prometheus-server.default.svc.cluster.local:80"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
       log:
         level: debug
         pretty_console: true
         non_blocking: false
   ```

   ```bash
   helm install agent aperture/aperture-agent -f values.yaml
   ```

   All the config parameters for the Aperture Agent are available
   [here](/references/configuration/agent.md).

   A list of configurable parameters for the installation can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#parameters).

5. If you want to deploy the Aperture Agent into a namespace other than
   `default`, use the `-n` flag:

   ```bash
   helm install agent aperture/aperture-agent -n aperture-agent --create-namespace
   ```

6. Refer steps on the
   [Istio Configuration](/get-started/installation/agent/envoy/istio.md) if you
   don't have the
   [Envoy Filter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)
   configured on your cluster.

## Verifying the Installation

Once you have successfully deployed the resources, confirm that the Aperture
Agent is up and running:

```bash
kubectl get pod -A

kubectl get agent -A
```

You should see pods for Aperture Agent and Agent Manager in `RUNNING` state and
`Agent` Custom Resource in `created` state.

## Uninstall

You can uninstall the Aperture Agent and it's components by uninstalling the
chart installed above:

1. Delete the Aperture Agent chart:

   ```bash
   helm uninstall agent
   ```

2. Alternatively, if you have installed the Aperture Agent Custom Resource
   separately, follow below steps:

   1. Delete the Aperture Agent Custom Resource:

      ```bash
      kubectl delete -f agent.yaml
      ```

   2. Delete the Aperture Agent chart to uninstall the Aperture Operator:

      ```bash
      helm uninstall agent
      ```

3. If you have installed the chart in some other namespace than `default`,
   execute below commands:

   ```bash
   helm uninstall agent -n aperture-agent
   kubectl delete namespace aperture-agent
   ```

   > Note: By design, deleting a chart via Helm doesn’t delete the Custom
   > Resource Definitions (CRDs) installed via the chart.

4. **Optional**: Delete the CRD installed by the chart:

   ```bash
   kubectl delete crd agents.fluxninja.com
   ```
