---
title: Auto Scale
sidebar_position: 10
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
{@include: ./assets/gen/auto-scale/insertion.mmd}
```

</Zoom>

## Auto Scaler

Components involved in the auto-scaling process:

:::info

See also [_Auto Scaler_ reference](/reference/policies/spec.md#auto-scaler)

:::

_Auto Scaler_ is a high-level component in Aperture that performs auto-scaling.
It can interface with infrastructure APIs such as Kubernetes to automatically
adjust the number of instances or resources allocated to a service to meet
changing workload demands. _Auto Scaler_ is designed to ensure that the service
is scaled out to meet demand and scaled in when demand is low. Scaling out is
done more aggressively than scaling in to ensure optimal performance and
availability.

- Controllers: _Auto Scaler_ leverages controllers (for example, Gradient
  Controller) to make scaling decisions. A Controller can be configured for
  either scaling in or out, and defines the criteria that determine when to
  scale. Controllers process one or more input signals to compute a desired
  scale value. By configuring Controllers, you can fine-tune the auto-scaling
  behavior to meet the specific needs of your service. See
  [Gradient Controller](#gradient-controller) for more details.
- A scale-in Controller is active only when its output is smaller than the
  actual scale value. A scale-out Controller is active only when its output is
  larger than the actual scale value. For example, the actual number of replicas
  of a Kubernetes Deployment. An inactive Controller does not contribute to the
  scaling decision.
- Scale decisions from multiple active Controllers are combined by the Auto
  Scaler by taking the largest scale value.
- Maximum scale-in and scale-out step sizes: The amount of scaling that happens
  at a time is limited by the maximum scale-in and scale-out step sizes. This is
  to prevent large-scale changes from happening at once.
- Cooldown periods: There are cooldown periods defined individually for
  scale-out and scale-in. The _Auto Scaler_ won't scale-out or scale-in again
  until the cooldown period has elapsed. The intention of cooldowns is to make
  the changes gradually and observe their effect to prevent overdoing either
  scale-in or scale-out.
  - Scale-in cooldown: The _Auto Scaler_ won't scale-in again until the cooldown
    period has elapsed. If there is a scale-out decision, the scale-in cooldown
    is reset. This immediately corrects excessive scale-in.
  - Scale-out cooldown: The _Auto Scaler_ won't scale-out again until the
    cooldown period has elapsed. If there is a scale-out decision which is much
    larger than the current scale value, the scale-out cooldown is reset. This
    is done to accommodate any urgent need for scale-out.

## Gradient Controller

The Gradient Controller computes a desired scale value based on a signal and
setpoint. The gradient controller tries to adjust the scale value proportionally
to the relative difference between setpoint and signal.

The `gradient` describes a corrective factor that should be applied to the scale
value to get the signal closer to the setpoint. It's computed as follows:

$$
\text{gradient} = \left(\frac{\text{signal}}{\text{setpoint}}\right)^{\text{slope}}
$$

`gradient` is then clamped to `[1.0, max_gradient]` range for the scale-out
controller and `[min_gradient, 1.0]` range for the scale-in controller.

The output of the gradient controller is computed as follows:

$$
\text{desired\_scale} = \text{gradient}_{\text{clamped}} \cdot \text{actual\_scale}.
$$

## Pod Scaler

:::info

See also [_Pod Scaler_ reference](/reference/policies/spec.md#pod-scaler)

:::

_Pod Scaler_ is a basic building block of an auto-scaling policy for Kubernetes.
It can scale out or scale in a Kubernetes Resources such as a Deployment. The
component takes the desired replicas as an input signal and scales the
underlying resources based on the value of the signal. To complete the feedback
loop, the component emits output signals for the number of configured replicas
and the actual number of replicas deployed.

A _Pod Scaler_ component is not required if you are using the high-level _Auto
Scaler_ component. It defines multiple scale in, scale out controllers and takes
care of instantiating the _Pod Scaler_ component.

## Kubernetes Object Selector

:::info

See also
[_Kubernetes Object Selector_ reference](/reference/policies/spec.md#kubernetes-object-selector)

:::

_Kubernetes Object Selectors_ are used by auto-scaling components in a policy,
such as [Pod Scaler](/reference/policies/spec.md#pod-scaler) or
[_Auto Scaler_](/reference/policies/spec.md#pod-scaler). A Kubernetes Object
Selector identifies a resource in the Kubernetes cluster.

A _Kubernetes Object Selector_ consists of:

- _Agent Group_: The Agent Group identifies Aperture Agents where the component
  gets applied.
- API Version: The Kubernetes API version of the resource.
- Kind: The Kind of the Kubernetes resource, such as Deployment, ReplicaSet,
  StatefulSet.
- Name: The name of the Kubernetes resource.
- Namespace: The Kubernetes namespace of the resource.

:::info Control Points Discovery

The Kubernetes resources identified by a _Kubernetes Object Selector_ are called
_Kubernetes Control Points_. These are a subset of resources in a Kubernetes
cluster resource that can be scaled in or out. Aperture Agents performs
automated discovery of Kubernetes Control Points in a cluster.

:::
