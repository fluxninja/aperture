---
title: Auto Scaling
keywords:
  - tutorial
sidebar_position: 4
sidebar_label: Auto Scaling
---

## Overview

Auto-scaling is a vital pillar of load management. It empowers service operators
to adjust the number of instances or resources allocated to a service
automatically, based on current or anticipated demand and resource utilization.
This way, auto-scaling ensures a service can handle incoming load while
optimizing operational costs by allocating the appropriate number of resources.

In Aperture, service operators can configure auto-scaling policies based on
different overload signals such as load throttling, in addition to resource
utilization based on CPU, memory usage, network I/O, and more. This versatility
enables service operators to fine-tune auto-scaling behavior according to their
specific needs. Auto-scaling policies can be set up to add or remove instances
or resources based on these signals, enabling dynamic scaling in response to
changing traffic patterns.

<Zoom>

```mermaid
{@include: ./assets/auto-scaling/auto-scaling.mmd}
```

</Zoom>

The diagram outlines the process of auto-scaling. The controller, on receiving a
scale in or out signal, dispatches an auto-scaling signal to the agent, which
interfaces with the infrastructure APIs, such as Kubernetes, to perform the
scaling.

## Example Scenario

Imagine a task management application whose usage varies significantly between
working days and holidays, and even fluctuates during different hours of a
standard working day. To maintain responsive APIs while managing infrastructure
costs effectively, the service can use load-based auto-scaling. This strategy
ensures that service resources are dynamically adjusted in line with usage
patterns, therefore optimizing both user experience and infrastructure
expenditure.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
