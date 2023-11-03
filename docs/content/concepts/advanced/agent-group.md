---
title: Agent Group
sidebar_label: Agent Group
sidebar_position: 1
keywords:
  - Agent
  - Configuration
---

## Agent Group {#agent-group}

The agent group is a versatile label that defines a collection of Agents
operating as peers. It serves as a means to organize and manage Agents within
Aperture. The agent group can be associated with different entities based on the
deployment mode of the Agents.

In sidecar mode, it is recommended to name the agent group based on the service.
This approach establishes a unified agent group for all pods within the service.
For example, all pods within the 'Checkout' service can be defined under the
same agent group.

In DaemonSet mode, the agent group name is typically based on the Kubernetes
cluster name. This ensures that all Agents deployed on each node of the cluster
belong to the same agent group.

:::note

The agent group can be configured during Agent installation. Refer to the
[agent config](/reference/configuration/agent.md#agent-info-config) for more
details.

:::

<!-- vale off -->

### Benefits of Agent Group

<!-- vale on -->

- **Efficient Management**: Agent groups facilitate the efficient management of
  multiple Agents within complex environments, such as
  [sidecar](/self-hosting/agent/kubernetes/operator/sidecar.md) or multi-cluster
  installations. They enable scaling of Aperture configuration to meet the
  demands of intricate setups.

- **State Synchronization**: An agent group defines the scope of agent-to-agent
  synchronization. Agents within the same group form a peer-to-peer network to
  synchronize fine-grained per label counters. These counters are crucial for
  [rate-limiting](/concepts/rate-limiter.md) and for implementing global token
  buckets used in [quota scheduling](/concepts/scheduler/quota-scheduler.md).
  Additionally, all Agents within an agent group instantiate the same set of
  flow control components as defined in the
  [policies](/concepts/advanced/policy.md) running at the Controller. This
  ensures consistent behavior and synchronization across the group.
