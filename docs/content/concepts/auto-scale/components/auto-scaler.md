---
title: Auto Scaler
keywords:
  - Autoscaling
  - auto-scaler
sidebar_position: 1
---

:::info

See also [_Auto Scaler_ reference](/reference/configuration/spec.md#auto-scaler)

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
