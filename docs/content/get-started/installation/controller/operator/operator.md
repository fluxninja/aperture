---
title: Install with Operator
description: Install Aperture Controller using Operator
keywords:
  - install
  - setup
  - controller
  - operator
sidebar_position: 1
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../../../apertureVersion.js';
```

## Controller Custom Resource Definition

The Aperture Controller is a Kubernetes-based application and is installed using
the
[Kubernetes Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
which is managed by the Aperture Operator.

The configuration for the Aperture Controller process is provided to the
Controller CRD under the `controller.config` section. All the configuration
parameters for the Aperture Controller are listed
[here](/reference/configuration/controller.md).

## Installation {#controller-installation}

By following these instructions, you will have deployed the Aperture Controller
into your cluster.

:::info Refer

Kubernetes Objects which will be created by the following steps are listed
[here](/reference/kubernetes-operator/controller.md).

:::

1. Run the following `install` command:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion}`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install controller aperture/aperture-controller`}
   </CodeBlock>
   </TabItem>
   </Tabs>

2. By default, Prometheus and etcd instances are installed. If you do not want
   to install and use your existing instances of Prometheus or etcd, configure
   the following values in the `values.yaml` file and pass it with the `install`
   command:

   ```yaml
   controller:
     config:
       etcd:
         endpoints: ["ETCD_ENDPOINT_HERE"]
       prometheus:
         address: "PROMETHEUS_ADDRESS_HERE"
   etcd:
     enabled: false

   prometheus:
     enabled: false
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of etcd and Prometheus, which will be used by the Aperture
   Controller.

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install controller aperture/aperture-controller -f values.yaml`}
   </CodeBlock>
   </TabItem>
   </Tabs>

   A list of all the configurable parameters for etcd are available
   [here](/reference/configuration/controller.md#etcd), and Prometheus are
   available [here](/reference/configuration/controller.md#prometheus).

   **Note**: Please ensure that the flag `web.enable-remote-write-receiver` is
   enabled on your existing Prometheus instance, as it is required by the
   Aperture Controller.

3. If you want to modify the default parameters or the Aperture Controller
   configuration, for example `log`, you can create or update the `values.yaml`
   file and pass it with `install` command:

   ```yaml
   controller:
     config:
       log:
         level: debug
         pretty_console: true
         non_blocking: false
   ```

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install controller aperture/aperture-controller -f values.yaml`}
   </CodeBlock>
   </TabItem>
   </Tabs>

   All the configuration parameters for the Aperture Controller are available
   [here](/reference/configuration/controller.md).

   A list of configurable parameters for the installation can be found in the
   [README](https://artifacthub.io/packages/helm/aperture/aperture-controller#parameters).

4. If you want to deploy the Aperture Controller into a namespace other than
   `default`, use the `--namespace` flag:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml --namespace aperture-controller`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm install controller aperture/aperture-controller -f values.yaml --namespace aperture-controller --create-namespace`}
   </CodeBlock>
   </TabItem>
   </Tabs>

5. Alternatively, you can create the Controller Custom Resource directly on the
   Kubernetes cluster using the below steps:

   1. Create a `values.yaml` for starting the operator and disabling the
      creation of Controller Custom Resource and pass it with `install` command:

      ```yaml
      controller:
        create: false
      ```

      <Tabs groupId="setup" queryString>
      <TabItem value="aperturectl" label="aperturectl">
      <CodeBlock language="bash">
      {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml`}
      </CodeBlock>
      </TabItem>
      <TabItem value="Helm" label="Helm">
      <CodeBlock language="bash">
      {`helm install controller aperture/aperture-controller -f values.yaml`}
      </CodeBlock>
      </TabItem>
      </Tabs>

   2. Create a YAML file with the below specifications:

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
      ```

      All the configuration parameters for the Controller Custom Resource are
      listed on the
      [README](https://artifacthub.io/packages/helm/aperture/aperture-controller#controller-custom-resource-parameters)
      file of the Helm chart.

   3. Apply the YAML file to the Kubernetes cluster using `kubectl`

      ```bash
      kubectl apply -f controller.yaml
      ```

## Upgrade Procedure {#controller-upgrade-procedure}

By following these instructions, you will have deployed the upgraded version of
Aperture Controller into your cluster.

1. Use the same `values.yaml` file created as part of
   [Installation Steps](#controller-installation) and pass it with below
   command:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm template --include-crds --no-hooks controller aperture/aperture-controller -f values.yaml | kubectl apply -f -`}
   </CodeBlock>

   Once all the pods are in a running state after upgrade, run the below command
   to keep the Helm release updated:

   <CodeBlock language="bash">
   {`helm upgrade controller aperture/aperture-controller -f values.yaml`}
   </CodeBlock>
   </TabItem>
   </Tabs>

2. If you have deployed the Aperture Controller into a namespace other than
   `default`, use the `--namespace` flag:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml --namespace aperture-controller`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm template --include-crds --no-hooks controller aperture/aperture-controller -f values.yaml --namespace aperture-controller | kubectl apply -f -`}
   </CodeBlock>

   Once all the pods are in a running state after upgrade, run the below command
   to keep the Helm release updated:

   <CodeBlock language="bash">
   {`helm upgrade controller aperture/aperture-controller -f values.yaml --namespace aperture-controller`}
   </CodeBlock>
   </TabItem>
   </Tabs>

## Verifying the Installation

Once you have successfully deployed the resources, confirm that the Aperture
Controller is up and running:

```bash
kubectl get pod -A

kubectl get controller -A
```

You should see pods for Aperture Controller and Controller Manager in `RUNNING`
state and `Controller` Custom Resource in `created` state.

## Uninstall

You can uninstall the Aperture Controller and its components installed above by
following the below steps:

1. Uninstall the Aperture Controller:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl uninstall controller`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm uninstall controller`}
   </CodeBlock>
   </TabItem>
   </Tabs>

2. Alternatively, if you have installed the Aperture Controller Custom Resource
   separately, follow the below steps:

   1. Delete the Aperture Controller Custom Resource:

      ```bash
      kubectl delete -f controller.yaml
      ```

   2. Delete the Aperture Controller to uninstall the Aperture Operator:

      <Tabs groupId="setup" queryString>
      <TabItem value="aperturectl" label="aperturectl">
      <CodeBlock language="bash">
      {`aperturectl uninstall controller`}
      </CodeBlock>
      </TabItem>
      <TabItem value="Helm" label="Helm">
      <CodeBlock language="bash">
      {`helm uninstall controller`}
      </CodeBlock>
      </TabItem>
      </Tabs>

3. If you have installed the chart in some other namespace than `default`,
   execute the below commands:

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl uninstall controller --namespace aperture-controller`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm uninstall controller --namespace aperture-controller
   kubectl delete namespace aperture-controller`}
   </CodeBlock>
   </TabItem>
   </Tabs>

4. If you have installed the Contour chart for exposing the etcd and Prometheus
   service, execute the below command:

   ```bash
   helm uninstall aperture -n projectcontour
   kubectl delete namespace projectcontour
   ```

5. By default, the Secret and ConfigMap having SSL/TLS certificates generated by
   the Kubernetes Operator for itself and Aperture Controller are not deleted
   with above steps. If you want to delete them, run the below commands:

   ```bash
   kubectl delete secret -l app.kubernetes.io/instance=controller-aperture-controller-manager
   kubectl delete secret -l app.kubernetes.io/instance=controller
   kubectl delete configmap -l app.kubernetes.io/instance=controller
   ```

6. **Optional**: Delete the CRD installed by the Helm chart:

   > Note: Intentionally, deleting a chart by using Helm does not delete the
   > Custom Resource Definitions installed by using the Helm chart.

   ```bash
   kubectl delete crd controllers.fluxninja.com
   kubectl delete crd policies.fluxninja.com
   ```
