---
title: Dynamic Configuration
description: Dynamically configuration of policy
keywords:
  - policy
  - configuration
sidebar_position: 3
---

Aperture's policies can be thought of as "Load Management Applications" running
on top of the Aperture platform. These applications (policies) are designed to
not only configure during startup, they can also be configured at runtime. This
helps preserve the runtime state of the policy while it receives updates to its
existing configuration.

:::note

Not all components support dynamic configuration updates. Look for fields such
as `*_config_key` in the components that support dynamic configuration updates.

:::

For instance, some actuation components can be switched between normal and dry
run modes at runtime through dynamic configuration. Certain blueprints use
[_Variable_](/reference/configuration/spec.md#variable) type components to
implement dynamic control. For instance, the
[feature rollout blueprint](/reference/blueprints/load-ramping/base#dynamic-configuration)
takes dynamic configuration to start or stop the rollout.

The dynamic configuration can be provided to an existing policy using the
[aperturectl CLI](/reference/aperture-cli/aperturectl/policy/apply/apply.md). To
learn more about its usage, see how the dynamic configuration is provided in the
service protection
[blueprint](/reference/blueprints/load-scheduling/average-latency#dynamic-configuration).
