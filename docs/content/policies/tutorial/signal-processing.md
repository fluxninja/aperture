---
title: Signal Processing
keywords:
  - policies
  - signals
  - circuit
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

### Detecting Overload State

Aperture's control-loop policies are programmable "circuits" that are evaluated
periodically. One of the primary goals of these policies is to calculate the
deviation from objectives and apply counter-measures such as concurrency limits
to keep the system is safe operational zone. The policies are used to express
where the metrics are collected from and where the actuation happens, along with
signal processing needed to translate health metrics to corrective actions.

For instance, a policy can be written to detect overload build-up at an upstream
service and trigger load-shedding at a downstream service.

### Example Policy

One of the most reliable metrics to detect overload state is latency of the
service requests.

In this example, we will be computing exponential moving average (EMA) of
latency, gathered periodically from a
[PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) query.
Further, we will multiply EMA of latency with a tolerance factor to calculate
setpoint latency, which is a threshold to detect overloaded state. That is, if
the real-time latency of the service is more than this setpoint (which is based
on long-term EMA), then we can consider the service to be overloaded at that
time.

#### Circuit Diagram

<Zoom>

```mermaid
{@include: ./assets/signal-processing/signal-processing.mmd}
```

</Zoom>

```mdx-code-block
<Tabs>
<TabItem value="Jsonnet">
```

```jsonnet
{@include: ./assets/signal-processing/signal-processing.jsonnet}
```

```mdx-code-block
</TabItem>
<TabItem value="YAML">
```

```yaml
{@include: ./assets/signal-processing/signal-processing.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```

#### Monitoring the Policy

Signals flowing through policy's circuit are reported as Prometheus'
[Summaries](https://prometheus.io/docs/practices/histograms/). Therefore, they
can be monitored in real-time using tools such as
[Grafana](https://github.com/grafana/grafana).

Below is an example signal monitoring dashboard for the above policy that can be
imported into the Grafana instance.

```mdx-code-block
<Tabs>
<TabItem value="Jsonnet">
```

```jsonnet
{@include: ./assets/signal-processing/signals-dashboard.jsonnet}
```

```mdx-code-block
</TabItem>
<TabItem value="JSON">
```

```yaml
{@include: ./assets/signal-processing/signals-dashboard.json}
```

```mdx-code-block
</TabItem>
</Tabs>
```
