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

Broadly, this policy revolves around two significant areas: Latency Monitoring
and Rollout Control.

- Latency Monitoring: By employing the
  [`average_latency_driver`](/reference/policies/bundled-blueprints/policies/feature-rollout/base.md#average-latency-driver),
  the policy persistently measures the application's latency. Within this, a
  criteria is set to establish the thresholds for both forward progression and
  rollback of the rollout. Furthermore, selectors play a crucial role in this
  policy, managing the traffic for flow control and observability components
  within the Aperture Agents. Selectors lay down the rules for traffic flow,
  determining how the components should operate. This capability allows
  developers to define [`control points`] within the code or data plane, which
  act as strategic locations where flow control decisions are applied.
  Developers set these control points using SDKs or during API Gateways or
  Service Meshes integration.

- Rollout Control: The [`load_ramp`] component is instrumental in performing a
  controlled rollout of the new feature. A
  [`regulator`](../../concepts/flow-control/components/regulator.md) within this
  component manages the flow of traffic to control points, facilitating either
  sticky or random sessions based on preference, thereby balancing the load and
  enabling controlled tests. The rollout's pace is managed using
  [`steps`](/reference/policies/spec#load-ramp-parameters-step) that
  incrementally increase the percentage of requests served by the new feature.
  This approach allows the policy to continuously monitor the application's
  latency and rollback the feature if the latency exceeds the configured
  threshold.

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
