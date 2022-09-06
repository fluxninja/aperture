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

Aperture provides sophisticated flow control capabilities by locating Aperture Agents next to your service instances as a sidecar. Aperture Agents monitor golden signals using an in-built telemetry system and a programmable, high-fidelity flow classifier used to label requests based on attributes such as customer tier or request type. These metrics are analyzed by the controller. 

The Aperture Controller is powered by always-on, dataflow-driven policies that continuously track deviations from service-level objectives (SLOs) and calculate recovery or escalation actions. The policies running in the controller are expressed as circuits, much like circuit networks in the game Factorio.


## What is a flow? {#flow}

A flow is the fundamental unit of work from the perspective of an Aperture
Agent. It could be an API call, a feature, or even a database query. A flow in
Aperture is similar to
[OpenTelemetry Span](https://opentelemetry.io/docs/reference/specification/trace/api/#span).

Flow are observed and controlled by Aperture Agents by using data-plane
components such as Classifiers, FluxMeters and Actuators.
