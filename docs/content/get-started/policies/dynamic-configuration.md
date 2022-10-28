---
title: Dynamic Configuration
description: Dynamically configuration of policy
keywords:
  - policy
  - configuration
sidebar_position: 2
---

Aperture's policies can be thought of as "Reliability Applications" running on
top of the Aperture platform. And just like applications, the policies are
designed to not only configure during startup, they can also be configured at
runtime. This helps preserve runtime state of the policy while it receives
updates to it's existing configuration.

:::note

Not all components support dynamic configuration updates. Look for fields such
as `dynamic_config_key` in the components that support dynamic configuration
updates.

:::

For instance, if a policy contains a rate limiter, we can provide limit
overrides for specific flow label keys via dynamic configuration. This prevents
resetting of distributed counters that would otherwise happen when a policy is
restarted.

The dynamic configuration can be provided in the Policy Custom Resource using
the `dynamicConfig` key. To learn more about it's usage, please see how the
dynamic configuration is provided in the
[Latency Gradient](https://github.com/fluxninja/aperture/blob/main/blueprints/lib/1.0/blueprints/latency-gradient/policy.libsonnet)
Blueprint.
