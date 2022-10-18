---
title: Static Rate Limiting
keywords:
  - policies
  - ratelimit
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Regulating heavy-hitters

One of the simplest flow-control policies is static rate limiting. It uses
[rate limiting actuator](/concepts/flow-control/rate-limiter.md) to block
specific flow labels that exceed their quota within a certain timeframe (limit
reset interval).

### Example Policy

In this example, we will be rate limiting unique users based on the `User-Id`
header in the HTTP traffic. This header is provided by Envoy proxy under the
label key `http.request.header.user_id` (See
[Flow Labels](/concepts/flow-control/flow-label.md)).

We will be configuring our rate limiter to allow at most `120 requests` for each
user in the `60s` period.

In addition, we will be configuring our rate limiter to apply these limits to
`ingress` traffic on Kubernetes service
`service1-demo-app.demoapp.svc.cluster.local`.

#### Circuit Diagram

```mermaid
{@include: ./assets/static-rate-limiting/static-rate-limiting.mmd}
```

```mdx-code-block
<Tabs>
<TabItem value="Jsonnet">
```

```jsonnet
{@include: ./assets/static-rate-limiting/static-rate-limiting.jsonnet}
```

```mdx-code-block
</TabItem>
<TabItem value="YAML">
```

```yaml
{@include: ./assets/static-rate-limiting/static-rate-limiting.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```
