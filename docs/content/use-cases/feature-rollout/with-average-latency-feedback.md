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
without compromising existing functionality. With Aperture, features can be
toggled on or off for specific user segments. The following policy enables you
to progressively introduce a new feature, all the while assessing its impact on
the application's latency. If the latency deteriorates beyond the configured
threshold, the rollout can be halted or reversed to ensure a seamless user
experience.

## Policy Key Concepts

For achieving a controlled rollout it is essential to monitor the service's
latency to be prepared to rollback the feature if the latency deteriorates
beyond the configured threshold.

The policy monitors the latency with the use of the following components:

- `average_latency_driver` persistently measures the application's latency.
- `criteria` determines the thresholds for both forward progression and rollback
  of the rollout.
- [`selectors`](../../concepts/flow-control/selector.md) define the rules that
  decide how components should select flows for requests processing.
- [`control point`](../../concepts/flow-control/selector.md) can be considered
  as a critical checkpoint in code or data plane, a strategically placed spot
  where flow control decisions are applied. Developers define these points
  during the integration of API Gateways or Service Meshes or by using Aperture
  SDKs.

The policy controls the rollout with the use of the following components:

- [`load_ramp`](/reference/policies/bundled-blueprints/policies/feature-rollout/base.md#load-ramp):
  controls the rollout's pace by incrementally increasing the percentage of
  requests served by the new feature.
- [`regulator`](../../concepts/flow-control/components/regulator.md): manages
  the flow of traffic to control points, facilitating either sticky or random
  sessions based on preference, thereby balancing the load and enabling
  controlled tests.
- [`steps`](/reference/policies/spec#load-ramp-parameters-step): incrementally
  increase the percentage of requests served by the new feature.

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
