---
title: Sidecar
description: Install Aperture Agent in Sidecar mode
keywords:
  - install
  - setup
  - agent
  - sidecar
sidebar_position: 2
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion,apertureVersionWithOutV} from '../../../../apertureVersion.js';
```

The Aperture Agent can also be installed as a Sidecar. In this mode, whenever a
new pod is started with required labels and annotations, the agent container
will be attached with the pod.

## Controlling the Injection Policy

The Aperture Agent's Sidecar injection can be enabled by adding namespace labels
and pod annotations as below:

| Resource  | Label              | Annotation                      | Enabled Value | Disabled Value |
| --------- | ------------------ | ------------------------------- | ------------- | -------------- |
| Namespace | aperture-injection | -                               | enabled       | disabled       |
| Pod       | -                  | sidecar.fluxninja.com/injection | true          | false          |

The injector is configured with the following logic:

- If either label or annotation is disabled, the pod is not injected
- If pod annotation is enabled but the namespace label is not present, the pod
  is not injected
- If neither label nor annotation is set, the pod is injected if the namespace
  is listed under `.spec.sidecar.enableNamespacesByDefault`. This is not enabled
  by default, so generally this means the pod is not injected.

## Prerequisites

You can do the installation using `aperturectl` CLI tool or using `Helm`.
Install the tool of your choice using below links:

1.  [aperturectl](/get-started/aperture-cli/aperture-cli.md)

    :::info Refer
    [aperturectl install agent](/reference/aperturectl/install/agent/agent.md)
    to see all the available command line arguments. :::

2.  [Helm](https://helm.sh/docs/intro/install/)

    1. Once the Helm CLI is installed, add the
       [Aperture Agent Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-agent)
       repo in your environment for install/upgrade:

       ```bash
       helm repo add aperture https://fluxninja.github.io/aperture/
       helm repo update
       ```

## Installation {#agent-sidecar-installation}

The Aperture Agent in the Sidecar mode will be installed using the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which will be managed by the Aperture Operator.

By following these instructions, you will have deployed the Aperture Agent into
your cluster.

:::info Refer

Kubernetes Objects which will be created by following steps are listed
[here](/reference/kubernetes-operator/agent.md).

:::

1. Configure the below parameters of Plugins, Etcd and Prometheus for the Agent
   Custom Resource by creating a `values.yaml` with below parameters and pass it
   with `install` command:

   :::info

   The below parameters disable the FluxNinja ARC Plugin for the Aperture Agent.
   If you want to keep it enabled, add parameters provided
   [here](/arc/plugin.md#configuration) under the `agent.config` section.

   :::

   ```yaml
   agent:
     sidecar:
       enabled: true
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
   [Aperture Controller](/get-started/installation/controller/controller.md) on
   the same cluster in `default` namespace, with Etcd and Prometheus using
   `controller` as release name, the values for the values for
   `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` would be as below:

   ```yaml
   agent:
     sidecar:
       enabled: true
     config:
       etcd:
         endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
       prometheus:
         address: "http://controller-prometheus-server.default.svc.cluster.local:80"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
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

2. To enable the pod injection for a list of namespaces by default, you can
   create or update the `values.yaml` file and pass it with `install` command:

   ```yaml
   agent:
     sidecar:
       enabled: true
       enableNamespacesByDefault:
         - NAMESPACE1
         - NAMESPACE2
     config:
       etcd:
         endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
       prometheus:
         address: "http://controller-prometheus-server.default.svc.cluster.local:80"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
   ```

   Replace the `NAMESPACE1`, `NAMESPACE2` and so on, with the actual namespaces
   and add more if required.

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

3. If you want to modify the default parameters or the Aperture Agent config,
   for example `log`, you can create or update the `values.yaml` file and pass
   it with `install` command:

   ```yaml
   agent:
     sidecar:
       enabled: true
     config:
       etcd:
         endpoints: ["ETCD_ENDPOINT_HERE"]
       prometheus:
         address: "PROMETHEUS_ADDRESS_HERE"
       plugins:
         disabled_plugins:
           - aperture-plugin-fluxninja
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

   All the config parameters for the Aperture Agent are available
   [here](/reference/configuration/agent.md).

   A list of configurable parameters for the installation can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#parameters).

4. If you want to deploy the Aperture Agent into a namespace other than
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

5. Alternatively, you can create the Agent Custom Resource directly on the
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

   2. Create a YAML file with below specifications:

      ```yaml
      apiVersion: fluxninja.com/v1alpha1
      kind: Agent
      metadata:
        name: agent
      spec:
        sidecar:
          enabled: true
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

6. Refer steps on the
   [Istio Configuration](/get-started/integrations/flow-control/envoy/istio.md)
   if you don't have the
   [Envoy Filter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)
   configured on your cluster.

## Upgrade Procedure {#agent-sidecar-upgrade-procedure}

By following these instructions, you will have deployed the upgraded version of
Aperture Agent into your cluster.

1. Use the same `values.yaml` file created as part of
   [Installation Steps](#agent-sidecar-installation) and pass it with below
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

You should see pod for Aperture Agent Manager in `RUNNING` state and `Agent`
Custom Resource in `created` state.

Now, when you create a new pod in the above listed namespaces, you will be able
to see the Aperture Agent container attached with the existing pod containers.
Confirm that the container is injected:

```bash
kubectl describe po <POD_ID>
```

Replace the `POD_ID` with the actual pod ID and check the containers section in
the output. There should be a container with name `aperture-agent`.

## Customizing injection

Generally, the pod are injected based on the default and overridden parameters
provided in the Custom Resource.

Per-pod configuration is available to override these options on individual pods.
This is done by adding an `aperture-agent` container to your pod. The sidecar
injection will treat any configuration defined here as an override to the
default injection template.

Care should be taken when customizing these settings, as this allows complete
customization of the resulting Pod, including making changes that cause the
sidecar container to not function properly.

For example, the following configuration customizes a variety of settings,
including setting the CPU requests, adding a volume mount, and modifying
environment variables:

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: nginx
  name: nginx
spec:
  containers:
    - name: aperture-agent
      image: auto
      env:
        - name: "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP"
          value: "group1"
      resources:
        requests:
          cpu: 1
          memory: 256Mi
    - image: nginx
      name: nginx
      resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
```

In general, any field in a pod can be set. However, care must be taken for
certain fields:

- Kubernetes requires the image field to be set before the injection has run.
  While you can set a specific image to override the default one, it is
  recommended to set the image to `auto` which will cause the sidecar injector
  to automatically select the image to use.
- Some fields in Pod are dependent on related settings. For example, CPU request
  must be less than CPU limit. If both fields are not configured together, the
  pod may fail to start.

Additionally, `agent-group` field can be configured using the annotation like:

`sidecar.fluxninja.com/agent-group=group1`

## Uninstall

You can uninstall the Aperture Agent and it's components installed above by
following below steps:

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

3. If you have installed the chart in some other namespace than `default`,
   execute below commands:

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

4. Restarts all the Pods which were injected with the Sidecar:

5. If pods are running as part of a Kubernetes Deployment:

   ```bash
   kubectl rollout restart deployment <DEPLOYMENT_NAME> -n <NAMESPACE>
   ```

6. If pods are running as part of a Kubernetes DaemonSet:

   ```bash
   kubectl rollout restart daemonset <DAEMONSET_NAME> -n <NAMESPACE>
   ```

7. If pod is running standalone (not part of a deployment or replica set):

   ```bash
   kubectl delete pod <POD_ID> -n <NAMESPACE>
   k apply -f pod.yaml
   ```

8. **Optional**: Delete the CRD installed by the Helm chart:

   > Note: By design, deleting a Helm chart via Helm doesn’t delete the Custom
   > Resource Definitions (CRDs) installed via the chart.

   ```bash
   kubectl delete crd agents.fluxninja.com
   ```
