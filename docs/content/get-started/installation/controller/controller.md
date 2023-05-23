---
title: Aperture Controller
description: Install Aperture Controller
keywords:
  - install
  - setup
  - controller
sidebar_position: 2
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../../apertureVersion.js';
```

## Overview

The Aperture Controller functions as the brain of the Aperture system.
Leveraging an advanced control loop, the Controller routinely analyzes polled
metrics and indicators to determine how traffic should be shaped as defined by
set policies. Once determined, these decisions are then exported to all Aperture
Agents to effectively handle workloads.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (process variables) and compares them against
setpoints. The difference in the variable values against these points is
referred to as the error signal. The feedback loop then works to minimize these
error signals by determining and distributing control actions, that adjust these
process variables and maintain their values within the optimal range.

## Prerequisites

You can do the installation using the `aperturectl` CLI tool or using `Helm`.
Install the tool of your choice using the following links:

1. [aperturectl](/get-started/installation/aperture-cli/aperture-cli.md)

   :::info Refer

   [aperturectl install controller](/reference/aperturectl/install/controller/controller.md)
   to see all the available command line arguments.

   :::

2. [Helm](https://helm.sh/docs/intro/install/)

   1. Once the Helm CLI is installed, add the
      [Aperture Controller Helm chart](https://artifacthub.io/packages/helm/aperture/aperture-controller)
      repository in your environment for install or upgrade:

      ```bash
      helm repo add aperture https://fluxninja.github.io/aperture/
      helm repo update
      ```

## Installation

The Aperture Controller can be installed using the below options:

:::caution warning

Upgrading from one of the installation modes below to the other is discouraged
and can result in unpredictable behavior.

:::

1. [**Install with Operator**](operator/operator.md)

   The Aperture Controller can be installed using the Kubernetes Operator
   available for it. This method requires access to create cluster level
   resources like ClusterRole, ClusterRoleBinding, CustomResourceDefinition and
   so on.

2. [**Namespace-scoped Installation**](namespace-scoped/namespace-scoped.md)

   The Aperture Controller can also be installed with only namespace-scoped
   resources.

<!-- vale off -->

## Exposing etcd and Prometheus services {#expose-etcd-prometheus}

<!-- vale on -->

If the Aperture Controller is installed with the packaged etcd and Prometheus,
follow the following steps to expose them outside the Kubernetes cluster so that
the Aperture Agent running on Linux can access them.

:::info

[Contour](https://projectcontour.io/) is used as a
[Kubernetes Ingress Controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
in the following steps to expose the etcd and Prometheus services out of
Kubernetes cluster using Load Balancer.

Any other tools can also be used to expose the etcd and Prometheus services out
of the Kubernetes cluster based on your infrastructure.

:::

1. Add the Helm chart repository for Contour in your environment:

   ```bash
   helm repo add bitnami https://charts.bitnami.com/bitnami
   ```

2. Install the Contour chart by running the following command:

   ```bash
   helm install aperture bitnami/contour --namespace projectcontour --create-namespace
   ```

3. It might take a few minutes for the Contour Load Balancer IP to become
   available. You can watch the status by running:

   ```bash
   kubectl get svc aperture-contour-envoy --namespace projectcontour -w
   ```

4. Once `EXTERNAL-IP` is no longer `<pending>`, run the following command to get
   the External IP for the Load Balancer:

   ```bash
   kubectl describe svc aperture-contour-envoy --namespace projectcontour | grep Ingress | awk '{print $3}'
   ```

5. Add an entry for the above IP in the cloud provider's DNS configuration. For
   example, follow steps on
   [Cloud DNS on GKE](https://cloud.google.com/dns/docs/records) for Google
   Kubernetes Engine.

6. Configure the below parameters to install the
   [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
   with the Aperture Controller by updating the `values.yaml` created during
   installation and passing it with `install` command:

   ```yaml
   ingress:
     enabled: true
     domain_name: YOUR_DOMAIN_HERE

   etcd:
     service:
       annotations:
         projectcontour.io/upstream-protocol.h2c: "2379"
   ```

   Replace the values of `YOUR_DOMAIN_HERE` with the actual value the domain
   name under with the External IP is exposed.

   <Tabs groupId="setup" queryString>
   <TabItem value="aperturectl" label="aperturectl">
   <CodeBlock language="bash">
   {`aperturectl install controller --version ${apertureVersion} --values-file values.yaml`}
   </CodeBlock>
   </TabItem>
   <TabItem value="Helm" label="Helm">
   <CodeBlock language="bash">
   {`helm upgrade --install controller aperture/aperture-controller -f values.yaml`}
   </CodeBlock>
   </TabItem>
   </Tabs>

7. It might take a few minutes for the Ingress resource to get the `ADDRESS`.
   You can watch the status by running:

   ```bash
   kubectl get ingress controller-ingress -w
   ```

8. Once the `ADDRESS` matches the External IP, the etcd will be accessible on
   `http://etcd.YOUR_DOMAIN_HERE:80` and the Prometheus will be accessible on
   `http://prometheus.YOUR_DOMAIN_HERE:80`.

## Applying Policies

The process of creating policies for Aperture can be done either after the
installation of the controller or after the installation of the agent, depending
on your preference.
[Generating and applying policies](/get-started/policies/policies.md) guide
includes step-by-step instructions on how to create policies for Aperture in a
Kubernetes cluster, which you can follow to create policies according to your
needs.
