---
title: Flow Control
sidebar_position: 1
keywords:
  - flows
  - tracing
  - opentracing
  - opentelemetry
---

# Flow Control

Reliable operations at web-scale are impossible without effective flow control.
Aperture provides sophisticated flow control capabilities by locating agents
next the services (sidecar).

## What is a flow? {#flow}

A flow is the fundamental unit of work from the perspective of an Aperture
Agent. It could be an API call, a feature, or even a database query. A flow in
Aperture is similar to
[OpenTelemetry Span](https://opentelemetry.io/docs/reference/specification/trace/api/#span).

Flow are observed and controlled by Aperture Agents by using data-plane
components such as classifiers, FluxMeters and actuators.
