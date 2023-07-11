---
title: Detecting Overload
keywords:
  - policies
  - signals
  - circuit
sidebar_position: 2
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Overview

Monitoring the health of a service is a critical aspect of ensuring reliable
operations. This policy provides a mechanism for detecting an overload state of
a service and sending alerts using Aperture's declarative policy language. The
policy creates a [circuit](/concepts/advanced/circuit.md) that models the
typical latency behavior of the service using an exponential moving average
(EMA). This automated learning of the normal latency threshold for each service
reduces the need for manual tuning of alert policies.

One reliable metric for detecting overload is the latency of service requests.
In Aperture, latency can be reported using a
[_Flux Meter_](/concepts/flux-meter.md).

:::tip

To prevent the mixing of latency measurements across different workloads, it's
recommended to apply the Flux Meter to a single type of workload. For instance,
if a service has both Select and Insert API calls, it is advised to measure the
latency of only one of these workloads using a Flux Meter. Refer to the
[_Selector_](/concepts/selector.md) documentation for guidance on applying the
Flux Meter to a subset of API calls for a service.

:::

## Configuration

In this example, the EMA of latency of `checkout-service.prod.svc.cluster.local`
is computed using metrics reported by the Flux Meter and obtained periodically
through a
[PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) query.
The EMA of latency is then multiplied by a tolerance factor to calculate the
setpoint latency, which serves as a threshold for detecting an overloaded
state - if the real-time latency of the service exceeds this setpoint (which is
based on the long-term EMA), the service is considered overloaded.

```mdx-code-block
<Tabs>
<TabItem value="YAML">
```

```yaml
{@include: ./assets/detecting-overload/detecting-overload.yaml}
```

```mdx-code-block
</TabItem>
<TabItem value="Jsonnet">
```

```jsonnet
{@include: ./assets/detecting-overload/detecting-overload.jsonnet}
```

```mdx-code-block
</TabItem>
</Tabs>
```

### Circuit Diagram

<Zoom>

```mermaid
{@include: ./assets/detecting-overload/detecting-overload.mmd}
```

</Zoom>

## Installation

Apply this custom policy to the `aperture-controller` namespace using
`aperturectl` or `kubectl`. The
[policy](/reference/aperturectl/apply/policy/policy.md) section within
`aperturectl` documentation provides additional information and examples related
the application of policies.

```mdx-code-block
<Tabs>
<TabItem value="aperturectl" label="aperturectl">
```

```bash
aperturectl apply policy --file=policy.yaml --kube
```

```mdx-code-block
</TabItem>
<TabItem value="kubectl" label="kubectl">
```

```bash
kubectl apply -f policy.yaml -n aperture-controller
```

```mdx-code-block
</TabItem>
</Tabs>
```

## Policy in Action

As the service processes traffic, various signal metrics collected from the
execution of the policy can be visualized:

<Zoom>

![LATENCY](./assets/detecting-overload/latency.png) `LATENCY`: Signal gathered
from the periodic execution of PromQL query on _Flux Meter_ metrics.

</Zoom>

<Zoom>

![LATENCY_EMA](./assets/detecting-overload/latency_ema.png) `LATENCY_EMA`:
Exponential Moving Average of `LATENCY` signal.

</Zoom>

<Zoom>

![LATENCY_SETPOINT](./assets/detecting-overload/latency_setpoint.png)
`LATENCY_SETPOINT`: Latency above which the service is considered to be
overloaded. This is calculated by multiplying the exponential moving average
with a tolerance factor (`LATENCY_EMA` \* `1.1`).

</Zoom>

<Zoom>

![IS_OVERLOAD_SWITCH](./assets/detecting-overload/is_overload_switch.png)
`IS_OVERLOAD_SWITCH` is a signal that represents whether the service is in an
overloaded state. This signal is derived by comparing `LATENCY` with
`LATENCY_SETPOINT`. A value of `0` indicates no overload, while a value of `1`
signals an overload.

</Zoom>
