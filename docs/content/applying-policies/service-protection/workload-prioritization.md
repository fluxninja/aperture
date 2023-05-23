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

When services are resource constrained, it is often crucial to preserve key
user-experience. Business critical features need to be prioritized while
throttling background workloads and less critical features. For instance, for an
e-commerce application, the ability to check out a shopping cart is more
critical than personalized recommendations and should be prioritized when
resources are constrained.

Aperture's
[weighted fair queuing scheduler](/concepts/flow-control/components/load-scheduler.md#scheduler)
enables prioritization of certain flows over others based on their flow labels,
ensuring that the user experience or revenue is maximized in the face of
overloads and other failures.

## Policy

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

The traffic generator in the [playground](/playground/playground.md) is
configured to generate similar traffic pattern (number of concurrent users) for
2 types of users - subscribers and guests.

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
