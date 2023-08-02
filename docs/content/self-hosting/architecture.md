---
title: Architecture of Self-Hosted Aperture
sidebar_label: Architecture
sidebar_position: 0
keywords:
  - aperture
  - controller
  - self-hosted
  - open-source
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

Architecture of the self-hosted Aperture solution differs slightly from the
regular [FluxNinja Cloud + Aperture combination](/architecture/architecture.md).
The main difference is that the Aperture Controller is no longer part of
FluxNinja Cloud and is deployed separately. Aperture Controller also needs its
supporting databases.

Aperture uses two databases to store configuration, telemetry, and flow control
information: [Prometheus][] and [etcd][]. Prometheus enables Aperture to monitor
the system and detect deviations from the service-level objectives (SLOs)
defined in the declarative policies. Aperture Controller uses etcd (distributed
key-value store) to persist the declarative policies that define the control
circuits and their components, as well as the adjustments synchronized between
the Controller and Agents.

Existing etcd and
[scalable Prometheus](https://promlabs.com/blog/2021/10/14/promql-vendor-compatibility-round-three)
installations can be reused to minimize operational overhead and to integrate
into existing monitoring infrastructure.

<Zoom>

```mermaid
{@include: ../assets/diagrams/architecture/architecture_simple.mmd}
```

:::info

The roles of Aperture Agent and Aperture Controller are described on the
[Architecture][] page.

:::

</Zoom>

[Architecture]: /architecture/architecture.md
[Prometheus]: https://prometheus.io
[etcd]: https://etcd.io
