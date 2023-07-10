---
title: Service
sidebar_label: Service
sidebar_position: 3
keywords:
  - Service Discovery
---

In Aperture, a service represents a collection of entities that deliver a common
functionality, such as checkout or billing in an e-commerce example, or provide
a specific service, such as a database or a search service. It maps to the
concept of services in platforms like Kubernetes or Consul. Services in Aperture
are typically referred to by their fully qualified domain names (FQDN).

## Service Discovery

Aperture Agent maintains a mapping of entities' (Kubernetes pod or VM) IP
addresses to service names. Each agent is responsible for discovering the
entities in its vicinity and the services that they belong to. The discovery
database is local to each Aperture Agent.

In the case of a Kubernetes DaemonSet installation, each Aperture Agent
maintains a mapping of IP addresses to services for the pods running on its
worker node. Changes in services and entities are actively watched by the agent
to ensure that the discovery remains up to date. This allows for accurate and
reliable identification of services during flow control decision-making.

<!-- vale off -->

## How are Services used in flow control decisions?

<!-- vale on -->

Upon receiving a flow control decision request from an entity, Aperture uses the
IP address to service mapping to identify the service name(s). The service
name(s) along with flow labels determine which flow control components to
execute. Refer to [Selector concept](selector.md) for more details.

:::note

An entity might belong to multiple services.

:::

## Live Preview of Services

Use the
[`aperturectl discovery entities`](../reference/aperturectl/discovery/entities/)
CLI command to list the discovered entities (pods, VMs) and their mapping to
services.

For example:

```sh
aperturectl discovery entities --kube
```

Returns:

```json
{
  "entities": {
    "10.244.1.7": {
      "uid": "2cc868cd-7e1a-49e5-80b6-81bbeb506719",
      "ipAddress": "10.244.1.7",
      "name": "service1-demo-app-7b4bc9bdcd-2krh8",
      "namespace": "demoapp",
      "nodeName": "aperture-playground-worker2",
      "services": ["service1-demo-app.demoapp.svc.cluster.local"]
    },
    "10.96.11.97": {
      "uid": "10.96.11.97",
      "ipAddress": "10.96.11.97",
      "name": "ClusterIP-demoapp-service1-demo-app-10.96.11.97",
      "namespace": "demoapp",
      "services": ["service1-demo-app.demoapp.svc.cluster.local"]
    }
  }
}
```
