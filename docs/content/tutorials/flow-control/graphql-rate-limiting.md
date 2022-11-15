---
title: GraphQL Query Static Rate Limiting
keywords:
  - policies
  - ratelimit
  - graphql
sidebar_position: 5
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

In this tutorial, we will use [Flow Classifier Rego Rules][rego-rules] to rate
limit a GraphQL query. We will build upon what we've seen in
[Static Rate Limiting](static-rate-limiting.md) and
[Workload Prioritization](workload-prioritization.md) tutorials.

[rego-rules]: ../../concepts/flow-control/flow-classifier#rego
