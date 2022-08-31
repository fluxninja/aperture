---
title: Sidecar
description: Install Aperture Agent in Sidecar mode
keywords:
  - install
  - setup
  - agent
  - daemonset
---

The Aperture Agent can also be installed as a Sidecar. In this mode, whenever a new pod is started
with required labels and annotations, the agent container will be attached with the pod.

## Controlling the Injection Policy

The Aperture Agent's Sidecar injection can be enabled by adding namespace labels and pod annotations as below:

| Resource  | Label              | Annotation                      | Enabled Value | Disabled Value |
| --------- | ------------------ | ------------------------------- | ------------- | -------------- |
| Namespace | aperture-injection | -                               | enabled       | disabled       |
| Pod       | -                  | sidecar.fluxninja.com/injection | true          | false          |

The injector is configured with the following logic:

- If either label or annotation is disabled, the pod is not injected
- If pod annotation is enabled but the namespace label is not present, the pod is not injected
- If neither label nor annotation is set, the pod is injected if the namespace is listed under
  `.spec.sidecar.enableNamespacesByDefault`. This is not enabled by default, so generally this
  means the pod is not injected.

## Installation {#agent-sidecar-installation}

1. The Aperture Agent can be installed in the Sidecar mode using the Helm chart of Aperture Operator
   by using the default `values.yaml` or create a `values.yaml` with below parameters and pass it with `helm upgrade`:

   ```yaml
   agent:
     create: true
     sidecar:
       enabled: true
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
        sidecar:
          enabled: true
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

3. To enable the pod injection for a list of namespaces by default for injection, add the below parameters in YAML specification:

   ```yaml
   apiVersion: fluxninja.com/v1alpha1
   kind: Agent
   metadata:
     name: agent
   spec:
     sidecar:
       enabled: true
       enableNamespacesByDefault:
         - NAMESPACE1
         - NAMESPACE2
     image:
       registry: docker.io/fluxninja
       repository: aperture-agent
       tag: latest
   ```

   Replace the NAMESPACE1 and NAMESPACE2 with the actual namespaces and add more if required.

   The similar parameters can be added under the `.values.agent.sidecar.enableNamespacesByDefault` to deploy with the Helm chart.

4. Once you have successfully deployed the Custom Resource, confirm that the
   Aperture Agent is up and running:

   ```bash
   kubectl get agent -A
   ```

   You should see the `Agent` Custom Resource in `created` state.

5. Now, when you create a new pod in the above listed namespaces, you will be able to see the Aperture Agent container attached
   with the existing pod containers. Confirm that the container is injected:

   ```bash
   kubectl describe po <POD_ID>
   ```

   Replace the `POD_ID` with the actual pod ID and check the containers section in the output. There should be a container
   with name `aperture-agent`.

## Customizing injection

Generally, the pod are injected based on the default and overridden parameters provided in the Custom Resource.

Per-pod configuration is available to override these options on individual pods. This is done by adding an `aperture-agent` container
to your pod. The sidecar injection will treat any configuration defined here as an override to the default injection template.

Care should be taken when customizing these settings, as this allows complete customization of the resulting Pod,
including making changes that cause the sidecar container to not function properly.

For example, the following configuration customizes a variety of settings, including setting the CPU requests, adding a volume mount, and
modifying environment variables:

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

In general, any field in a pod can be set. However, care must be taken for certain fields:

- Kubernetes requires the image field to be set before the injection has run. While you can set a specific
  image to override the default one, it is recommended to set the image to `auto` which will cause the
  sidecar injector to automatically select the image to use.
- Some fields in Pod are dependent on related settings. For example, CPU request must be less than CPU limit.
  If both fields are not configured together, the pod may fail to start.

Additionally, `agent-group` field can be configured using the annotation like:

`sidecar.fluxninja.com/agent-group=group1`
