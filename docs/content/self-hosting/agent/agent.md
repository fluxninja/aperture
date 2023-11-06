---
title: Agents
description: Aperture Agents
keywords:
  - install
  - setup
  - agent
  - sidecar
  - daemonset
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import CodeBlock from '@theme/CodeBlock';
```

## Overview

The Aperture Agent is the decision executor of the Aperture system. In addition
to gathering data, the Aperture Agent functions as a gatekeeper, acting on
traffic based on periodic adjustments made by the Aperture Controller.
Specifically, depending on feedback from the Controller, the agent will
effectively allow or drop incoming requests. Further, supporting the Controller,
the agent works to inject information into traffic, including the specific
traffic-shaping decisions made and classification labels which can later be used
for observability and closed loop feedback.

## Configuration

All the configuration parameters for the Aperture Agent are listed
[here](/reference/configuration/agent.md).

## Installation Modes {#agent-installation-modes}

The Aperture Agent can be installed in the following modes:

:::caution warning

Upgrading from one of the installation modes below to the other is discouraged
and can result in unpredictable behavior.

:::

1. **Kubernetes**

   1. [**Namespace-scoped Installation**][namespace-scoped-installation]

      The Aperture Agent can also be installed with only namespace-scoped
      resources.

   2. [**Install with Operator**](kubernetes/operator/operator.md)

      The Aperture Agent can be installed using the Kubernetes Operator
      available for it.

      :::info

      This method requires access to create cluster level resources like
      ClusterRole, ClusterRoleBinding, CustomResourceDefinition and so on.

      Use the [Namespace-scoped Installation][namespace-scoped-installation] if
      you do not want to assign the cluster level permissions.

      :::

      - [**DaemonSet**](kubernetes/operator/daemonset.md)

        The Aperture Agent can be installed as a
        [Kubernetes DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
        where it will get deployed on all the nodes of the cluster.

      - [**Sidecar**](kubernetes/operator/sidecar.md)

        The Aperture Agent can also be installed as a Sidecar. In this mode,
        whenever a new pod is started with required labels and annotations, the
        agent container will be attached with the pod.

2. [**Bare Metal or VM**](bare-metal.md)

   The Aperture Agent can be installed as a system service on any Linux system
   that is [supported](../supported-platforms.md).

3. [**Docker**](docker.md)

   The Aperture Agent can also be installed on Docker as containers.

## Self-Hosted Aperture Controller {#agent-self-hosted-controller}

When using the self-hosted Aperture Controller instead of the Aperture Cloud
Controller, you need to turn off the `enable_cloud_controller` flag and
configure Controller, etcd and Prometheus endpoints directly, for example:

```mdx-code-block
<Tabs>
  <TabItem value="aperturectl or helm">
```

```yaml
agent:
  config:
    fluxninja:
      enable_cloud_controller: false
      endpoint: "ORGANIZATION_NAME.app.fluxninja.com:443"
    etcd:
      endpoints: ["http://controller-etcd.default.svc.cluster.local:2379"]
    prometheus:
      address: "http://controller-prometheus-server.default.svc.cluster.local:80"
    agent_functions:
      endpoints: ["aperture-controller.default.svc.cluster.local:8080"]
  secrets:
    fluxNinjaExtension:
      create: true
      secretKeyRef:
        name: aperture-agent-apikey
        key: apiKey
      value: API_KEY
```

The values above assume that you have installed the
[Aperture Controller](/self-hosting/controller/controller.md) on the same
cluster in `default` namespace, with etcd and Prometheus using `controller` as
release name. If your setup is different, adjust these endpoints accordingly.

```mdx-code-block
  </TabItem>

  <TabItem value="Docker or Bare Metal">
```

```yaml
fluxninja:
  enable_cloud_controller: false
  endpoint: "ORGANIZATION_NAME.app.fluxninja.com:443"
  api_key: API_KEY
etcd:
  endpoints: ["http://etcd:2379"]
prometheus:
  address: "http://prometheus:9090"
agent_functions:
  endpoints: ["aperture-controller:8080"]
otel:
  disable_kubernetes_scraper: true
  disable_kubelet_scraper: true
auto_scale:
  kubernetes:
    enabled: false
service_discovery:
  kubernetes:
    enabled: false
log:
  level: info
  pretty_console: true
  non_blocking: false
```

You might need to adjust the endpoints, depending on your exact setup.

```mdx-code-block
  </TabItem>
</Tabs>
```

:::info

If you're not using [Aperture Cloud](/reference/fluxninja.md), simply remove the
`fluxninja` and `secrets` sections.

:::

## Applying Policies

Once the
[application is set up](/get-started/set-up-application/set-up-application.md)
and Agents are installed, you can start creating and applying policies.

[Your first policy](/get-started/policies/policies.md) section provides
step-by-step instructions on customizing, creating, and applying policies within
Aperture. Additionally, the [guides](/guides/guides.md) section serves as a
valuable resource for tailoring policies to meet specific requirements.

[namespace-scoped-installation]: kubernetes/namespace-scoped/namespace-scoped.md
