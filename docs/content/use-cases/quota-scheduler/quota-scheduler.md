---
title: Quota Scheduling
keywords:
  - tutorial
  - quota-scheduler
  - policies
  - quota
  - scheduling
  - external-api
  - prioritization
sidebar_position: 1
sidebar_label: Quota Scheduling
---

## Overview

Quota Scheduling is a method to maintain the balance of service-to-service
requests, ensuring that the request frequency stays within the given limit. More
than just limiting requests, Quota Scheduling allows for request prioritization
based on the workload, thereby offering a higher degree of control in specific
use cases.

<Zoom>

```mermaid
{@include: ../assets/quota-scheduler.mmd}
```

</Zoom>

The presented graph delineates the Token Bucket's operation, given a specified
bucket size and fill rate. The Token Bucket performs counting and distributes
tokens to all Agents. Inside each Agent, a scheduler organizes requests based on
priority (assigned through label matching) and the availability of tokens.

:::note

The Token Bucket is distributed across multiple Agents within the same cluster.

:::

## Real World Scenario

For Quota Scheduling, consider a cloud-based storage service managing requests
from numerous client applications. The service ensures fair usage by
implementing a quota scheduling policy. This way, it can prioritize critical
requests and manage resource allocation effectively, thereby preventing any
single client from monopolizing the service or exhausting the available quota.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
