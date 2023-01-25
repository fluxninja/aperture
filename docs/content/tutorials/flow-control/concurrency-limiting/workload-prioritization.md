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

When services are resource constrained and concurrency limits are being applied,
it's often crucial to preserve key user-experience by gracefully degrading
application behavior. Graceful degradation allows prioritizing business critical
features while de-prioritizing background workloads and less critical features.
For instance, for an e-commerce application, ability to checkout a shopping cart
is more critical than personalized recommendations and should be prioritized
when resources are constrained.

Aperture's
[weighted fair queueing scheduler](/concepts/flow-control/components/concurrency-limiter.md#scheduler)
allows fairly prioriting certain flows over others based on their flow labels.
This enables graceful degradation in face of overloads and other failures, and
maximizes user-experience or revenue.

## Policy

We will be extending the policy we used in
[Basic Concurrency Limiting](../concurrency-limiting/basic-concurrency-limiting.md)
to classify requests into workloads and prioritize them.

In this example policy, we will be prioritizing traffic of different types of
users - `subscriber` users will get higher priority over `guest` users. That is,
under overload scenarios, subscribed users will get better quality of service
over guest users. We will be using 2 alternative ways to provide the `User-Type`
value to the scheduler:

- Subscribers: We will directly match the header value of `User-Type` to
  `subscriber` since all HTTP headers are directly available as flow labels
  within the scheduler.
- Guests: To identify guest users, we will first use a classification rule that
  uses an [extractor](concepts/flow-control/flow-classifier.md#extractors) to
  assign the header value to `user-type` flow label key. Ultimately, we will be
  using the `user_type` label key in the scheduler to match the request against
  `guest` value to identify the workload.

:::tip

You can quickly write classification rules on
[HTTP requests](concepts/flow-control/flow-classifier.md#live-previewing-requests)
and define scheduler priorities on
[Flow Labels](concepts/flow-control/flow-label.md#live-previewing-flow-labels)
by live previewing them first via introspection APIs.

:::

In addition, we will be configuring the scheduler to automatically assign the
tokens that need to be obtained in order to accept requests matching a given
workload. This continuous estimation (auto-tokens) helps with fair scheduling
and prioritization across workloads. This additional configuration is
highlighted in the Jsonnet spec below.

```mdx-code-block
<Tabs>
<TabItem value="YAML">
```

```yaml
{@include: ./assets/workload-prioritization/workload-prioritization.yaml}
```

```mdx-code-block
</TabItem>
<TabItem value="Jsonnet">
```

```jsonnet
{@include: ./assets/workload-prioritization/workload-prioritization.jsonnet}
```

```mdx-code-block
</TabItem>
</Tabs>
```

### Circuit Diagram

<Zoom>

```mermaid
{@include: ./assets/workload-prioritization/workload-prioritization.mmd}
```

</Zoom>

### Playground

The traffic generator in the [playground](/get-started/playground/playground.md)
is configured to generate similar traffic pattern (number of concurrent users)
for 2 types of users - subscribers and guests.

When we load the above policy in the playground, we will see the that during the
overload period, `subscriber` users have higher acceptance rate of their
requests than `guest` users.

<Zoom>

![Workload Prioritization](./assets/workload-prioritization/workload-prioritization-playground.png)

</Zoom>

### Demo Video

The below demo video shows the basic concurrency limiter and workload
prioritization policy in action within Aperture Playground.

[![Demo Video](https://img.youtube.com/vi/m070bAvrDHM/0.jpg)](https://www.youtube.com/watch?v=m070bAvrDHM)
