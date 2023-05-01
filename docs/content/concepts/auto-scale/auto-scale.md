---
title: Auto Scale
sidebar_position: 3
keywords:
  - auto-scaling
  - Autoscaling
  - auto-scaler
  - Kubernetes
---

Auto-scaling is a crucial technique for effective load management of service
traffic. It enables service operators to automatically adjust the number of
instances or resources allocated to a service based on current or expected
demand and resource utilization. By doing so, auto-scaling ensures that a
service can handle incoming load while optimizing the cost of running the
service by allocating just the right number of resources.

Auto-scale is a core integration in Aperture that works hand in hand with the
flow control capabilities to provide a comprehensive load management platform.
Aperture policies allow defining auto-scaler(s) that consider the flow control
state for informed scaling decisions. For instance, during sudden traffic
spikes, if a concurrency limiter on a service sheds traffic, auto-scaler can
automatically add more instances to the service to handle the increased load.

With the auto-scaler capability in Aperture, service operators can configure
auto-scaling policies based on different service overload signals, such as load
shedding, in addition to resource utilization based on CPU, memory usage,
network I/O, and so on. This flexibility enables service operators to fine-tune
the auto-scaling behavior based on their specific service needs. Auto-scaling
policies can be set up to add or remove instances or resources based on these
signals, allowing for dynamic scaling in response to changing traffic patterns.

Auto-scaling is a powerful technique that enables service operators to maintain
service availability and performance while optimizing costs. In Aperture,
auto-scaling is an integral component of the load management platform, working
seamlessly with flow control to provide a comprehensive solution. These
capabilities allow services to dynamically adjust to incoming traffic patterns,
ensuring optimal performance while minimizing infrastructure costs.

## Insertion

Aperture Agents interface with cloud infrastructure APIs, such as Kubernetes
API, to discover, monitor, and scale infrastructure resources. The Aperture
Controller uses the information from the Agents to make informed auto-scaling
decisions that are then acted on by the Agents.

In an Agent group, the leader Agent is responsible for interfacing with the
cloud infrastructure APIs. For example, by maintaining a watch on scalable
Kubernetes resources, the Agent group leader can monitor changes to the resource
status, such as the number of replicas configured and currently deployed. The
up-to-date information is then used by the Aperture Controller to make informed
auto-scaling decisions.

<Zoom>

```mermaid
{@include: ./assets/insertion.mmd}
```

</Zoom>
