---
title: Concurrency Limiting
keywords:
  - policies
  - concurrency
sidebar_position: 2
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

One of the basic techniques to protect services from cascading failures is to
limit the concurrency on the service to match the processing capacity of the
service.
