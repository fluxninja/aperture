---
title: Concurrency Limiting
sidebar_position: 1
keywords:
  - guides
  - concurrency limiting
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
import {apertureVersion} from '../apertureVersion.js';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from "@theme/TabItem";
import {BashTab, TabContent} from './blueprintsComponents.js';
import CodeSnippet from '../codeSnippet.js'

```

:::note

The following policy is based on the
[Concurrency Limiting](/reference/blueprints/concurrency-limiting/base.md)
blueprint.

:::

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

Concurrency limiting is a critical strategy for managing the load on an API.
