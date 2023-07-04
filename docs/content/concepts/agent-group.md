---
title: Agent Group
sidebar_label: Agent Group
sidebar_position: 1
keywords:
  - Agent
  - Configuration
---

## Agent Group {#agent-group}

Agent group is a flexible label that defines a collection of agents that operate
as peers. For example, an agent group can be a Kubernetes cluster name in the
case of `DaemonSet` deployment of Agent, or it can be a service name for sidecar
deployments of Agent.

When employing sidecar mode, it's advisable to name the agent group based on the
respective service, therefore fostering a unified agent group for all pods
within a service. For instance, all pods within the 'Checkout' service can be
defined under the same agent group.

In DaemonSet mode, the Kubernetes cluster name typically becomes the agent group
name, which applies to all agents deployed on each node. This ensures all Agents
spanning the entire cluster comes under the same agent group.

:::note

Agent group can be configured at the agent during installation, refer to
[agent config](../reference/configuration/agent#agent-info-config)

:::

<!-- vale off -->

### Where does Agent Group help?

<!-- vale on -->

- **Complex Environments**: It helps manage multiple agents efficiently within
  intricate environments, like Kubernetes or multi-cluster installations.
  Basically, helping scale Aperture configuration.

- **State Synchronization**: an agent group defines the scope of agent-to-agent
  synchronization, with agents within the group forming a peer-to-peer network
  to synchronize fine-grained state per-label global counters that are used for
  rate-limiting purposes. Additionally, all agents within an agent group
  instantiate the same set of flow control components as published by the
  controller.
