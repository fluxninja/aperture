---
title: Rate Limiting
keywords:
  - tutorial
sidebar_position: 3
sidebar_label: Rate Limiting
---

## Overview

Rate limiting is a critical strategy for managing the load on an API. By
imposing restrictions on the number of requests a unique consumer can make
within a specific time frame, rate limiting prevents a small set of users from
monopolizing the majority of resources on a service, ensuring fair access for
all API consumers.

Aperture implements this strategy through its high-performance, distributed rate
limiter. This system enforces per-key limits based on fine-grained labels,
thereby offering precise control over API usage. For each unique key, Aperture
maintains a token bucket of a specified bucket capacity and fill rate. The fill
rate dictates the sustained requests per second (RPS) permitted for a key, while
transient overages over the fill rate are accommodated for brief periods, as
determined by the bucket capacity.

This intricate system of rate-limiting plays a pivotal role in maintaining the
integrity of a service. It effectively safeguards against excessive usage that
could potentially result in API abuse, while simultaneously ensuring optimal
performance and resource allocation.

<Zoom>

```mermaid
{@include: ./assets/rate-limiting/rate-limiting.mmd}
```

</Zoom>

The diagram depicts the distribution of tokens across agents through a global
token bucket. Each incoming request prompts the agents to decrement tokens from
the bucket. If the bucket has run out of tokens, indicating that the rate limit
has been reached, the incoming request is rejected. Conversely, if tokens are
available in the bucket, the request is accepted. The token bucket is
continually replenished at a predefined fill rate, up to the maximum number of
tokens specified by the bucket capacity.

## Example Scenario

Consider a social media platform that implements rate limits to prevent abuse of
its APIs. Each user of the platform gets identified by a unique key and assigned
a specific rate limit, controlled by the Aperture's distributed rate limiter.
For instance, the platform might allow a user to make a certain number of
requests per minute to post content, retrieve posts, or send messages.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
