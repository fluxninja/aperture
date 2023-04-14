---
title: Metrics
description: Integrating custom metrics pipelines
keywords:
  - setup
  - otel
  - opentelemetry
  - collector
sidebar_position: 3
---

Aperture allows feeding custom metrics to the controller Prometheus. This is
powered by adding custom metric receivers into OpenTelemetry Collector running
on Aperture Agent. [See configuration reference][config]

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />

[config]: /reference/configuration/agent.md#custom-metrics-config
