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

## Policy Overview

Mitigating cascading failures is essential to maintain service stability, which
can be achieved effectively by matching a service's concurrency limit with its
processing capacity. However, determining the precise concurrency limit can be
challenging due to the evolving nature of service infrastructure. Factors such
as deployment of new versions, horizontal scaling, or fluctuating access
patterns can impact the concurrency limit. This policy is designed to address
this dynamic problem and offer reliable service protection.

## Policy Key Concepts

The `service_protection_core` incorporates the following components to ensure
that applications are protected from recurring overloads:

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

  For latency monitoring, the `latency_baseliner` encompasses the `flux_meter`
  which converts a flux of flows matching a flow selector into a Prometheus
  histogram, essentially measuring the scope of latency. By default, it tracks
  the workload duration of a flow, but it can flexibly track any metric from
  OpenTelemetry attributes depending on the insertion method.

## Policy Configuration

This policy learns the latency profile of a service using an exponential moving
average. Deviation of current latency from the historical latency indicates an
overload. In case of overload, the policy lowers the rate at which requests are
admitted into the service, making the excess requests wait in a queue. Once the
latency improves, the rate of requests is slowly increased to find the maximum
processing capacity of the service.

This policy uses the Service Protection with Average Latency Feedback
[Blueprint](/reference/policies/bundled-blueprints/policies/service-protection/average-latency.md).

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
