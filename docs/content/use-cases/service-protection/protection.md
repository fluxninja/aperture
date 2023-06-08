---
title: Protection
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
[Service Protection with Average Latency Feedback](/reference/policies/bundled-blueprints/policies/service-protection/average-latency.md)
blueprint.

:::

## Overview

The response times of a service can start to deteriorate when the service's
concurrency limit is surpassed. Consequently, a degradation in response latency
can serve as a reliable signal for identifying service overload. This policy is
designed to detect overload situations based on latency deterioration and
adjusts the request rate accordingly to restore latency to its acceptable range.

## Policy Configuration

This policy monitors the latency of requests processed by the
**`service1-demo-app.demoapp.svc.cluster.local`** service. It calculates the
deviations in current latency from a baseline historical latency, which serves
as an indicator of service overload. A deviation of **`1.1`** from the baseline
is considered as a signal for service overload.

To mitigate service overload, the requests to
**`service1-demo-app.demoapp.svc.cluster.local`** service are passed through a
load scheduler. The load scheduler reduces the request rate in overload
scenarios, temporarily placing excess requests in a queue.

As service latency improves, indicating a return to normal operational state,
the request rate is incrementally increased until it matches the incoming
request rate. This responsive mechanism helps ensure that service performance is
optimized while mitigating the risk of service disruptions due to overload.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl values.yaml">
```

```yaml
{@include: ./assets/protection/values.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

<details><summary>Generated Policy</summary>
<p>

```yaml
{@include: ./assets/protection/policy.yaml}
```

</p>
</details>

:::info

[Circuit Diagram](./assets/protection/graph.mmd.svg) for this policy.

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

![Basic Service Protection](./assets/protection/dashboard.png)

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