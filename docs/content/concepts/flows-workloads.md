---
title: Flows and Workloads
slug: concepts/flows-workloads
description: Concept of flows and workloads.
keywords:
  - Flows
  - Workloads
---

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Flow](#flow)
- [Workload](#workload)

<!-- /code_chunk_output -->

# Flows and Workloads

## Flow

A flow is the fundamental unit of work from the perspective of an Aperture
Agent. It could be an API call, a feature, or even a database query. A flow in
Aperture is similar to
[OpenTelemetry Span](https://opentelemetry.io/docs/reference/specification/trace/api/#span).

## Workload

Workloads are a group of flows based on common attributes. Workloads are
expressed by label matcher rules in Aperture. Aperture Agents schedule workloads
based on their priorities and by estimating their [tokens](#tokens).
