---
title: Monitoring
description: Monitoring signals in the policy
keywords:
  - jsonnet
  - grafana
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Signals flowing through a policy's circuit are reported as Prometheus'
[Summaries](https://prometheus.io/docs/practices/histograms/). Therefore, they
can be monitored in real-time using tools such as
[Grafana](https://github.com/grafana/grafana).

Below is an example signal monitoring dashboard that can be imported into the
Grafana instance.

```mdx-code-block
<Tabs>
<TabItem value="Jsonnet">
```

```jsonnet
{@include: ./assets/monitoring/signals-dashboard.jsonnet}
```

```mdx-code-block
</TabItem>
<TabItem value="JSON">
```

```yaml
{@include: ./assets/monitoring/signals-dashboard.json}
```

```mdx-code-block
</TabItem>
</Tabs>
```
