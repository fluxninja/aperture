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
import {apertureVersion, apertureVersionWithOutV} from '../../../../../apertureVersion.js';
```

:::info

This method requires access to create cluster level resources like ClusterRole,
ClusterRoleBinding, CustomResourceDefinition and so on.

Use the
[Namespace-scoped Installation](/get-started/self-hosting/agent/kubernetes/namespace-scoped/namespace-scoped.md)
if you do not want to assign the cluster level permissions.

:::

The Aperture Agent can be installed as a
[Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
where it will get deployed on all the nodes of the cluster.

## Prerequisites

You can do the installation using the `aperturectl` CLI tool or using `Helm`.
Install the tool of your choice using the following links:

1. [Helm](https://helm.sh/docs/intro/install/)

   1. Once the Helm CLI is installed, add the
      [Aperture Agent Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-agent)
      repository in your environment for install or upgrade:

      ```bash
      helm repo add aperture https://fluxninja.github.io/aperture/
      helm repo update
      ```

2. [aperturectl](/get-started/installation/aperture-cli/aperture-cli.md)

   :::info Refer

   [aperturectl install agent](/reference/aperturectl/install/agent/agent.md) to
   see all the available command line arguments.

   :::

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

1. Configure the below parameters of Aperture Cloud endpoint and API key for the
   Agent Custom Resource by creating a `values.yaml` and passing it with
   `install` command:

   ```yaml
   agent:
     config:
       fluxninja:
         enable_cloud_controller: true
         endpoint: "ORGANIZATION_NAME.app.fluxninja.com:443"
     secrets:
       fluxNinjaExtension:
         create: true
         secretKeyRef:
           name: aperture-agent-apikey
           key: apiKey
         value: AGENT_API_KEY
   ```

   Replace `ORGANIZATION_NAME` with the Aperture Cloud organization name and
   `AGENT_API_KEY` with the API key linked to the project. If an API key has not
   been created, generate a new one through the Aperture Cloud UI. Refer to
   [Agent API Keys][agent-api-keys] for additional information.

   :::note

   If you are using a Self-Hosted Aperture Controller, modify the above
   configuration as explained in
   [Self-Hosting: Agent Configuration](/get-started/self-hosting/agent/agent.md#agent-self-hosted-controller).

   :::

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
       fluxninja:
         enable_cloud_controller: true
         endpoint: "ORGANIZATION_NAME.app.fluxninja.com:443"
       log:
         level: debug
         pretty_console: true
         non_blocking: false
     secrets:
       fluxNinjaExtension:
         create: true
         secretKeyRef:
           name: aperture-agent-apikey
           key: apiKey
         value: AGENT_API_KEY
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
   {`helm install agent aperture/aperture-agent -f values.yaml --namespace aperture-agent --create-namespace`}
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
          fluxninja:
            enable_cloud_controller: true
            endpoint: "ORGANIZATION_NAME.app.fluxninja.com:443"
        secrets:
          fluxNinjaExtension:
            create: true
            secretKeyRef:
              name: aperture-agent-apikey
              key: apiKey
            value: AGENT_API_KEY
      ```

      Replace `ORGANIZATION_NAME` with the Aperture Cloud organization name and
      `AGENT_API_KEY` with the API key linked to the project. If an API key has
      not been created, generate a new one through the Aperture Cloud UI. Refer
      to [Agent API Keys][agent-api-keys] for additional information.

      :::note

      If you are using a Self-Hosted Aperture Controller, modify the above
      configuration as explained in
      [Self-Hosting: Agent Configuration](/get-started/self-hosting/agent/agent.md#agent-self-hosted-controller).

      :::

      All the configuration parameters for the Agent Custom Resource are listed
      on the
      [README](https://artifacthub.io/packages/helm/aperture/aperture-agent#agent-custom-resource-parameters)
      file of the Helm chart.

   3. Apply the YAML file to Kubernetes cluster using `kubectl`

      ```bash
      kubectl apply -f agent.yaml
      ```

5. Refer to steps on the [Istio Configuration](/integrations/istio/istio.md) if
   you do not have the
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

   Once all the pods are in a running state after upgrade, run the below command
   to keep the Helm release updated:

   <CodeBlock language="bash">
   {`helm upgrade agent aperture/aperture-agent -f values.yaml`}
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

   Once all the pods are in a running state after upgrade, run the below command
   to keep the Helm release updated:

   <CodeBlock language="bash">
   {`helm upgrade agent aperture/aperture-agent -f values.yaml --namespace aperture-agent`}
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
   {`aperturectl uninstall agent --version ${apertureVersion}`}
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
      {`aperturectl uninstall agent --version ${apertureVersion}`}
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

4. By default, the Secret having SSL/TLS certificates generated by the
   Kubernetes Operator for itself is not deleted with above steps. If you want
   to delete it, run the following commands:

   ```bash
   kubectl delete secret -l app.kubernetes.io/instance=agent-aperture-agent-manager
   ```

5. **Optional**: Delete the CRD installed by the Helm chart:

   :::note

   Deleting a Helm chart using Helm does not delete the Custom Resource
   Definitions (CRD) installed from the chart.

   :::

   ```bash
   kubectl delete crd agents.fluxninja.com
   ```

[agent-api-keys]: /get-started/aperture-cloud/agent-api-keys.md
