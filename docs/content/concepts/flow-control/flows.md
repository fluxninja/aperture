---
title: Flows
description: Concept of flows.
keywords:
  - flows
  - tracing
  - opentracing
  - opentelemetry
---

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Flows](#flows)

<!-- /code_chunk_output -->

# Flows

A flow is the fundamental unit of work from the perspective of an Aperture
Agent. It could be an API call, a feature, or even a database query. A flow in
Aperture is similar to
[OpenTelemetry Span](https://opentelemetry.io/docs/reference/specification/trace/api/#span).
