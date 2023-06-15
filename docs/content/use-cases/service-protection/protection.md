---
title: Service Protection
keywords:
  - policies
  - concurrency
  - service-protection
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

:::note

The following policy is based on the
[Service Protection with Average Latency Feedback](/reference/blueprints/policies/service-protection/average-latency.md)
blueprint.

:::

## Policy Overview

Mitigating cascading failures is essential to maintain service stability, which
can be achieved effectively by matching a service's concurrency limit with its
processing capacity. However, determining the precise concurrency limit can be
challenging due to the evolving nature of service infrastructure. Factors such
as deployment of new versions, horizontal scaling, or fluctuating access
patterns can impact the concurrency limit. This policy is designed to address
this dynamic problem and offer reliable service protection.

## Policy Configuration

In this policy, latency is of **`service1-demo-app.demoapp.svc.cluster.local`**
is monitored using an exponential moving average. Deviation of current latency
from the historical latency indicates an overload, which will lead to lower the
rate at which requests are admitted into the monitored service, making the
excess requests wait in a queue. Once the latency improves, the rate of requests
is slowly increased to the maximum processing capacity of the selected service.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/basic-service-protection/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/basic-service-protection/policy.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/basic-service-protection/graph.mmd.svg) for this
policy.

:::

### Playground

When the above policy is loaded in Aperture's
[Playground](https://github.com/fluxninja/aperture/blob/main/playground/README.md),
it demonstrates that when latency spikes due to high traffic at
`service1-demo-app.demoapp.svc.cluster.local`, the controller throttles the rate
of requests admitted into the service. This approach helps protect the service
from becoming unresponsive and maintains the current latency within the
tolerance limit (`1.1`) of historical latency.

<Zoom>

![Basic Service Protection](./assets/basic-service-protection/dashboard.png)

</Zoom>

### Dry Run Mode

You can run this policy in the `Dry Run` mode by setting the
`default_config.dry_run` option to `true`. In the `Dry Run` mode, the policy
does not throttle the request rate while still evaluating the decisions it would
take in each cycle. This is useful for evaluating the policy without impacting
the service.

:::note

The `Dry Run` mode can also be toggled dynamically at runtime, without reloading
the policy.

:::

### Demo Video

The below demo video shows the basic service protection and workload
prioritization policy in action within Aperture Playground.

[![Demo Video](https://img.youtube.com/vi/m070bAvrDHM/0.jpg)](https://www.youtube.com/watch?v=m070bAvrDHM)
