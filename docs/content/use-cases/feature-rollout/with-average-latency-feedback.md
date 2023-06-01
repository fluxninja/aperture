---
title: Feature Rollout with Average Latency Feedback
keywords:
  - policies
  - rollout
  - latency
  - feature-flags
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Policy Overview

Feature flags provide a mechanism for shipping new features to production
without compromising existing functionality. By utilizing Aperture, features can
be toggled on or off for specific user segments. The following policy enables
you to progressively introduce a new feature, all the while assessing its impact
on your application's latency. If the latency deteriorates beyond the configured
threshold, the rollout can be halted or reversed to ensure a seamless user
experience.

## Policy Key Concepts

At a high-level, this policy consists of:

- Latency monitoring: Continuously measure the application's latency using the
  [`average_latency_driver`](/reference/policies/bundled-blueprints/policies/feature-rollout/base.md#average-latency-driver).
- Rollout control: Gradually increase the percentage of requests that are to be
  served the new feature using
  [`steps`](/reference/policies/spec#load-ramp-parameters-step). Monitor the
  application's latency and roll back the feature if the latency deteriorates
  beyond the configured limit.

Some of the key concepts used in this policy are:

- [Selector](../../concepts/flow-control/selector.md): Selectors are the traffic
  signal managers for flow control and observability components in the Aperture
  Agents. They lay down the traffic rules determining how these components
  should select flows for their operations.
- [Control Point](../../concepts/flow-control/selector.md): Think of Control
  Points as designated checkpoints in your code or data plane. They're the
  strategic points where flow control decisions are applied. Developers define
  these using SDKs or during API Gateways or Service Meshes integration.
- [Regulator](../../concepts/flow-control/components/regulator.md): Picture the
  Regulator as a vigilant gatekeeper at a crowded event, letting in only a set
  number of guests at a time to maintain order. In Aperture's context, it
  controls the flow traffic to a Control Point, and can allow random or sticky
  sessions based on your preferences, helping balance load and enabling
  controlled tests.

## Policy Configuration

This policy uses the
[`Feature Rollout with Average Latency Feedback`](/reference/policies/bundled-blueprints/policies/feature-rollout/average-latency.md)
blueprint that enables incremental roll out of a new feature. In this example,
we will create a policy that slowly ramps up the percentage of requests that are
served with the new feature. We will continuously monitor the application's
latency and roll back the feature if the latency deteriorates beyond the
configured limit.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/with-average-latency-feedback/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>

```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/with-average-latency-feedback/policy.yaml}
```

</p>
</details>

## Playground

The above policy can be loaded using the `feature_rollout` scenario in
[Playground](https://github.com/fluxninja/aperture/blob/main/playground/README.md)

:::info

[Circuit Diagram](./assets/with-average-latency-feedback/graph.mmd.svg) for this
policy.

:::

<Zoom>

![Feature Rollout with Average Latency Feedback](./assets/with-average-latency-feedback/dashboard.png)

</Zoom>
