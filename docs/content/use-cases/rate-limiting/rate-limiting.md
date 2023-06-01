---
title: Rate Limiting
keywords:
  - tutorial
sidebar_position: 3
sidebar_label: Rate Limiting
---

## Overview

Rate Limiting is a preventive mechanism designed to regulate the volume of
requests dispatched to a service. By capping the quantity of requests that can
be transmitted to a service within a specified interval, this technique
effectively curbs the potential for service abuse.

<Zoom>

```mermaid
{@include: ../assets/rate-limiting.mmd}
```

</Zoom>

This graph demonstrates the functionality of the Token Bucket in distributing
tokens across all Agents, based on a pre-set bucket size and fill rate. The
Agents, in response to the count, determine when to limit requests and when to
allow them.

## Real World Scenario

Considering Rate Limiting, imagine an online voting system where votes are
submitted via API calls. To prevent misuse or abuse of the service (like vote
spamming), a rate-limiting policy would be implemented. This ensures that each
user can only send a certain number of requests per set interval, maintaining
the integrity of the system.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
