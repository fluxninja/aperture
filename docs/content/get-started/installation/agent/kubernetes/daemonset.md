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

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion,apertureVersionWithOutV} from '../../../../apertureVersion.js';
```

The Aperture Agent can be installed as a
[Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
where it will get deployed on all the nodes of the cluster.

## Prerequisites

You can do the installation using the `aperturectl` CLI tool or using `Helm`.
Install the tool of your choice using the following links:

1. [aperturectl](/get-started/aperture-cli/aperture-cli.md)

   :::info Refer
   [aperturectl install agent](/reference/aperturectl/install/agent/agent.md) to
   see all the available command line arguments.

2. [Helm](https://helm.sh/docs/intro/install/)

   1. Once the Helm CLI is installed, add the
      [Aperture Agent Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-agent)
      repository in your environment for install/upgrade:

      ```bash
      helm repo add aperture https://fluxninja.github.io/aperture/
      helm repo update
      ```

<!-- vale off -->

## Installation {#agent-daemonset-installation}

<!-- vale on -->

The Aperture Agent in the DaemonSet mode will be installed using the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which will be managed by the Aperture Operator.

By following these instructions, you will have deployed the Aperture Agent into
your cluster.

:::info Refer

Kubernetes Objects which will be created by following steps are listed
[here](/reference/kubernetes-operator/agent.md).

:::

1. Configure the below parameters of etcd and Prometheus for the Agent Custom
   Resource by creating a `values.yaml` and passing it with `install` command:

   :::info

   The below parameters disable the FluxNinja ARC Extension for the Aperture
   Agent. If you want to keep it enabled, add parameters provided
   [here](/arc/extension.md#configuration) under the `agent.config` section.

   :::

   ```yaml
   agent:
     config:
       etcd:
         endpoints: ["ETCD_ENDPOINT_HERE"]
       prometheus:
         address: "PROMETHEUS_ADDRESS_HERE"
       agent_functions:
         endpoints: ["CONTROLLER_ENDPOINT_HERE"]
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of etcd and Prometheus, which is also being used by the
   Aperture Controller you want these Agents to connect with.
   `CONTROLLER_ENDPOINT_HERE` should point to the Aperture Controller. If you
   skip it, some sub-commands `aperturectl` commands won't work.

   If you have installed the
   [Aperture Controller](/get-started/installation/controller/controller.md) on
   the same cluster in `default` namespace, with etcd and Prometheus using
   `controller` as release name, the values for the values for
   `ETCD_ENDPOINT_HERE`, `PROMETHEUS_ADDRESS_HERE` and
   `CONTROLLER_ENDPOINT_HERE` would be as below:

   ```yaml
   agent:
     config:
       etcd:
         endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
       prometheus:
         address: "http://controller-prometheus-server.default.svc.cluster.local:80"
       agent_functions:
         endpoints: ["aperture-controller.default.svc.cluster.local:8080"]
   ```

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install agent --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install agent aperture/aperture-agent -f values.yaml`}
   </CodeBlock>
   </TabItem>
   </Tabs>

2. If you want to modify the default parameters or the Aperture Agent
   configuration, for example `log`, you can create or update the `values.yaml`
   file and pass it with `install` command:

   ```yaml
   agent:
     config:
       etcd:
         endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
       prometheus:
         address: "http://controller-prometheus-server.default.svc.cluster.local:80"
       log:
         level: debug
         pretty_console: true
         non_blocking: false
   ```

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install agent --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install agent aperture/aperture-agent -f values.yaml`}
   </CodeBlock>
   </TabItem>
   </Tabs>

   All the configuration parameters for the Aperture Agent are available
   [here](/reference/configuration/agent.md).

   A list of configurable parameters for the installation can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#parameters).

3. If you want to deploy the Aperture Agent into a namespace other than
   `default`, use the `--namespace` flag:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install agent --version ${apertureVersion} --values-file values.yaml --namespace aperture-agent`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install agent aperture/aperture-agent -f values.yaml --namespacen aperture-agent --create-namespace`}
   </CodeBlock>
   </TabItem>
   </Tabs>

4. Alternatively, you can create the Agent Custom Resource directly on the
   Kubernetes cluster using the below steps:

   1. Create a `values.yaml` for just starting the operator and disabling the
      creation of Agent Custom Resource and pass it with `install` command:

      ```yaml
      agent:
        create: false
      ```

      <Tabs groupId="setup" queryString>
      <TabItem value="aperturectl" label="aperturectl">
      <CodeBlock language="bash">
      {`aperturectl install agent --version ${apertureVersion} --values-file values.yaml`}
      </CodeBlock>
      </TabItem>
      <TabItem value="Helm" label="Helm">
      <CodeBlock language="bash">
      {`helm install agent aperture/aperture-agent -f values.yaml`}
      </CodeBlock>
      </TabItem>
      </Tabs>

   2. Create a YAML file with the following specifications:

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
      ```

      Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE`
      with the actual values of etcd and Prometheus, which is also being used by
      the Aperture Controller you want these Agents to connect with.

      All the configuration parameters for the Agent Custom Resource are listed
      on the
      [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#agent-custom-resource-parameters)
      file of the Helm chart.

   3. Apply the YAML file to Kubernetes cluster using `kubectl`

      ```bash
      kubectl apply -f agent.yaml
      ```

5. Refer to steps on the
   [Istio Configuration](/get-started/integrations/flow-control/envoy/istio.md)
   if you don't have the
   [Envoy Filter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)
   configured on your cluster.

<!-- vale off -->

## Upgrade Procedure {#agent-daemonset-upgrade-procedure}

<!-- vale on -->

By following these instructions, you will have deployed the upgraded version of
Aperture Agent into your cluster.

1. Use the same `values.yaml` file created as part of the
   [Installation Steps](#agent-daemonset-installation) and pass it with below
   command:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install agent --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm template --include-crds --no-hooks agent aperture/aperture-agent -f values.yaml | kubectl apply -f -`}
   </CodeBlock>
   </TabItem>
   </Tabs>

2. If you have deployed the Aperture Agent into a namespace other than
   `default`, use the `--namespace` flag:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install agent --version ${apertureVersion} --values-file values.yaml --namespace aperture-agent`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm template --include-crds --no-hooks agent aperture/aperture-agent -f values.yaml --namespace aperture-agent | kubectl apply -f -`}
   </CodeBlock>
   </TabItem>
   </Tabs>

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

You can uninstall the Aperture Agent and its components installed above by
following these steps:

1. Uninstall the Aperture Agent:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl uninstall agent`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm uninstall agent`}
   </CodeBlock>
   </TabItem>
   </Tabs>

2. Alternatively, if you have installed the Aperture Agent Custom Resource as a
   standalone installation, follow the steps below:

   1. Delete the Aperture Agent Custom Resource:

      ```bash
      kubectl delete -f agent.yaml
      ```

   2. Delete the Aperture Agent to uninstall the Aperture Operator:

      <Tabs groupId="setup" queryString>
      <TabItem value="aperturectl" label="aperturectl">
      <CodeBlock language="bash">
      {`aperturectl uninstall agent`}
      </CodeBlock>
      </TabItem>
      <TabItem value="Helm" label="Helm">
      <CodeBlock language="bash">
      {`helm uninstall agent`}
      </CodeBlock>
      </TabItem>
      </Tabs>

3. If you have installed the chart in some other namespace than the `default`,
   execute the following commands:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl uninstall agent --namespace aperture-agent`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm uninstall agent --namespace aperture-agent
   kubectl delete namespace aperture-agent`}
   </CodeBlock>
   </TabItem>
   </Tabs>

4. **Optional**: Delete the CRD installed by the Helm chart:

:::note

Deleting a Helm chart via Helm doesn’t delete the Custom Resource Definitions
(CRD) installed via the chart.

:::

```bash
kubectl delete crd agents.fluxninja.com
```
