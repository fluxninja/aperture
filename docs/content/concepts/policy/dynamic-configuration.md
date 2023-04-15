---
title: Dynamic Configuration
description: Dynamically configuration of policy
keywords:
  - policy
  - configuration
sidebar_position: 3
---

Aperture's policies can be thought of as "Reliability Applications" running on
top of the Aperture platform. And just like applications, the policies are
designed to not only configure during startup, they can also be configured at
runtime. This helps preserve the runtime state of the policy while it receives
updates to its existing configuration.

:::note

Not all components support dynamic configuration updates. Look for fields such
as `dynamic_config_key` in the components that support dynamic configuration
updates.

:::

If a policy contains a rate limiter, limit overrides for specific flow label
keys can be provided by dynamic configuration at runtime. This prevents the
resetting of distributed counters that would otherwise occur when a policy is
restarted.

The dynamic configuration can be provided to an existing policy using the
[aperturectl CLI](/reference/aperturectl/apply/apply.md). To learn more about
its usage, please see how the dynamic configuration is provided in the Latency
based AIMD (Additive Increase, Multiplicative Decrease) Concurrency Limiting
[Blueprint](/reference/policies/bundled-blueprints/policies/latency-aimd-concurrency-limiting.md).
