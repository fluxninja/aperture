---
title: Static Rate Limiting
keywords:
  - policies
  - ratelimit
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Policy Overview

Regulating incoming traffic in the face of overwhelming requests is vital to
maintaining the health and availability of a service. A powerful tool to achieve
this is the implementation of static rate limiting, aimed at controlling the
intensity of 'heavy-hitters.' This policy utilizes the
[Rate-Limiting Actuator](/concepts/flow-control/components/rate-limiter.md) to
curtail specific flow labels that go beyond their allocated quota within a
defined time frame (limit reset interval). It is an efficient and
straightforward mechanism for mitigating traffic congestion and preventing
potential service degradation or downtime.

## Policy Key Concepts

This policy is centered around two fundamental aspects: the [`rate_limiter`] and
the [`selectors`].

- [Rate Limiter](../../concepts/flow-control/components/rate-limiter.md) is a
  component modeled on the token bucket algorithm. It serves as an effective
  bulwark against heavy and frequent traffic spikes, essentially operating like
  a 'traffic parking meter.' The rate limiter's flexibility lies in its ability
  to be configured according to various labels and limits, providing an
  adjustable mechanism to handle diverse traffic patterns and intensities.
  Selectors and Control Points: The Selectors and Control Points perform similar
  roles as described in the previous policies.

- [Selectors](../../concepts/flow-control/selector.md) act as traffic signal
  managers that lay down the rules for flow selection by the rate limiter. On
  the other hand, [Control Points](../../concepts/flow-control/selector.md)
  serve as strategic checkpoints where these flow control decisions are applied.
  These concepts remain central to the process of implementing effective rate
  limiting for a service.

## Policy Configuration

This example demonstrates rate limiting of unique users based on the `User-Id`
header in the HTTP traffic. Envoy proxy provides this header under the label key
`http.request.header.user_id` (see
[Flow Labels](/concepts/flow-control/flow-label.md) for more information).

This configuration limits each user to a burst of `40 requests` and `2 requests`
every `1s` period using the rate limiter. Additionally, the rate limiter applies
these limits to `ingress` traffic on the Kubernetes service
`service1-demo-app.demoapp.svc.cluster.local`.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/static-rate-limiting/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>

```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/static-rate-limiting/static-rate-limiting.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/static-rate-limiting/static-rate-limiting.mmd.svg)
for this policy.

:::

### Playground

When the policy above is loaded in the playground, no more than 2 requests per
second period (after an initial burst of 40 requests) are accepted.

<Zoom>

![Static Rate Limiting](./assets/static-rate-limiting/static-rate-limiting-02.png)

</Zoom>
