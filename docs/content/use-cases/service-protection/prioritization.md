---
title: Workload Prioritization
keywords:
  - policies
  - scheduler
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Policy Overview

When dealing with services in resource-limited scenarios, it becomes paramount
to prioritize key user experiences and business-critical features over less
crucial tasks or background workloads. For instance, in an e-commerce platform,
the checkout process must take precedence over functionalities like personalized
recommendations, especially during resource shortage or high traffic. Aperture's
[Weighted Fair Queuing Scheduler (WFQ)](/concepts/flow-control/components/load-scheduler.md#scheduler)
enables such prioritization of flows over others based on their labels, ensuring
user experience or revenue is maximized during overloads or other failure
scenarios.

## Policy Key Concepts

The `service_protection_core` incorporates the following components to ensure
that applications are protected from overload:

- [`adaptive_load_scheduler`](../../concepts/flow-control/components/load-scheduler.md),
  it manages incoming request traffic to prevent service overload by throttling
  concurrent requests to a service.
- [`selectors`](../../concepts/flow-control/selector.md) define the rules that
  decide how components should select flows for requests processing.
- [`control point`](../../concepts/flow-control/selector.md) can be considered
  as a critical checkpoint in code or data plane, a strategically placed spot
  where flow control decisions are applied. Developers define these points
  during the integration of API Gateways or Service Meshes or by using Aperture
  SDKs.
- [`scheduler`](../../concepts/flow-control/components/load-scheduler.md)
  ensures that requests are serviced based on their priority and the number of
  tokens available, which are assigned by the token bucket.

## Policy Configuration

In this example policy, traffic of different types of users will be prioritized,
with `subscriber` users receiving higher priority over `guest` users. This means
that under overload scenarios, subscribed users will receive better quality of
service than guest users. Two alternative methods will be used to provide the
`User-Type` value to the scheduler:

- Subscribers: The header value of `User-Type` will be directly matched to
  `subscriber`, since all HTTP headers are directly available as flow labels
  within the scheduler.
- Guests: To identify guest users, a classification rule will be used that
  utilizes an
  [extractor](/concepts/flow-control/resources/classifier.md#extractors) to
  assign the header value to the `user-type` flow label key. The `user_type`
  label key will then be used in the scheduler to match the request against the
  `guest` value to identify the workload.

:::tip

Classification rules can be written for
[HTTP requests](/concepts/flow-control/resources/classifier.md#live-previewing-requests),
and scheduler priorities can be defined for
[Flow Labels](/concepts/flow-control/flow-label.md#live-previewing-flow-labels)
by live previewing them first using introspection APIs.

:::

To improve fairness and prioritization across workloads, the scheduler will be
configured to automatically assign tokens for accepting requests that match a
given workload. This is achieved through continuous estimation of tokens (auto
tokens) performed by the scheduler itself.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/workload-prioritization/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/workload-prioritization/policy.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/workload-prioritization/graph.mmd.svg) for this
policy.

:::

### Playground

The traffic generator in the
[playground](https://github.com/fluxninja/aperture/blob/main/playground/README.md)
is configured to generate similar traffic pattern (number of concurrent users)
for 2 types of users - subscribers and guests.

Loading the policy highlighted above in the playground will reveal that, during
overload periods, requests from `subscriber` users have a higher acceptance rate
than those from `guest` users.

<Zoom>

![Workload Prioritization](./assets/workload-prioritization/dashboard.png)

</Zoom>

### Demo Video

The below demo video shows the basic concurrency limiter and workload
prioritization policy in action within Aperture Playground.

[![Demo Video](https://img.youtube.com/vi/m070bAvrDHM/0.jpg)](https://www.youtube.com/watch?v=m070bAvrDHM)
