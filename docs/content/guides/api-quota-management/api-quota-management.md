---
title: API Quota Management
sidebar_position: 4
keywords:
  - api quota management
  - guides
---

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

## Overview

Quota scheduling is a sophisticated technique that enables effective management
of request quotas. This technique empowers services to enforce rate limits for
outbound or external API. This helps ensure that services stay within allocated
rate limits, therefore avoiding penalties, ensuring smooth and uninterrupted
operation.

Moreover, quota scheduling optimizes the utilization of request quotas by
prioritizing access based on business-critical workloads. This strategic
prioritization ensures that the most crucial requests receive their fair share
of request quotas, aligning API usage with business objectives and preventing
cost overages.

<Zoom>

```mermaid
{@include: ./assets/managing-quotas/managing-quotas.mmd}
```

</Zoom>

The diagram provides an overview of quota scheduling in action, including the
operation of the token bucket and its role in managing request admission. The
token bucket, specified by a given bucket size and fill rate, performs counting
and distributes tokens across Agents.

Requests coming into the system are categorized into different workloads, each
of which is defined by its priority and weight. This classification is crucial
for the scheduling process within each agent.

Inside every agent, there is a scheduler that priorities request admission based
on two factors: the priority and weight assigned to the corresponding workload,
and the availability of tokens from the global token bucket. This mechanism
ensures that high-priority requests are handled appropriately even under high
load or when the request rate is close to the rate limit.

## Example Scenario

Consider the scenario of a cloud-based database service handling requests from
several client applications with different priorities. By implementing a quota
scheduling policy using Aperture, service operators can ensure fair and
prioritized access. This policy enables prioritizing critical requests,
preventing any single client application from monopolizing the database service
or exhausting the available quota. Additionally, with Aperture's quota
scheduling, the system becomes expense-aware, ensuring it stays within quota
limits and avoids cost overages.

<DocCardList />
