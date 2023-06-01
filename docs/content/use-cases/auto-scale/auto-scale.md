---
title: Auto Scale
keywords:
  - tutorial
sidebar_position: 4
sidebar_label: Auto Scale
---

## Overview

Load-Based Auto-Scaling is a dynamic method designed to adjust the allocation of
instances or resources for a service, correlating with the workload demands. By
setting a scale-out Controller that responds to the load multiplier signal from
the Service Protection Policy, this technique accomplishes load-based
auto-scaling.

<Zoom>

```mermaid
{@include: ../assets/auto-scaling.mmd}
```

</Zoom>

The graph outlines the process of load-based auto-scaling. The Controller, on
receiving a scale in/out signal, dispatches a load-based auto-scaling signal to
the Agent, which interfaces with the infrastructure API, such as Kubernetes, to
perform the scaling.

## Real World Scenario

For Load-Based Auto-Scaling, consider a streaming service like Netflix during
peak viewing hours. To maintain high-quality streaming and minimize buffering,
the service employs load-based auto-scaling. This ensures the service resources
are dynamically adjusted according to the viewer demands, optimizing the
streaming experience for all users.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
