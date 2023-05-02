---
title: Basic Service Protection
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

The most effective technique to protect services from cascading failures is to
limit the concurrency of the service to match the processing capacity of the
service. However, figuring out the concurrency limit of a service is a hard
problem in the face of continuously changing service infrastructure. Each new
version deployed, horizontal scaling, or a change in access patterns can change
the concurrency limit of a service.

This policy learns the latency profile of a service using an exponential moving
average. Deviation of current latency from the historical latency indicates an
overload. In case of overload, the policy lowers the rate at which requests are
admitted into the service, making the excess requests wait in a queue. Once the
latency improves, the rate of requests is slowly increased to find the maximum
processing capacity of the service.

## Policy

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
[Playground](/get-started/playground/playground.md), it demonstrates that when
latency spikes due to high traffic at
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
