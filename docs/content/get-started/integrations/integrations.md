---
title: Integrations
sidebar_position: 3
description: Integrations provided by Aperture
image: ../assets/img/aperture_logo.png
keywords:
  - flow-control
  - integrations
  - sdk
  - enovy
  - go
  - java
  - cloud
  - flow control
---

## Overview

Aperture can be inserted into service instances with either _Service Meshes_ or
_SDKs_:

- Service Meshes: Aperture can be deployed with no changes to application code,
  using [Envoy](https://www.envoyproxy.io/). It latches onto Envoy’s
  [External Authorization API](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter)
  for control purposes and collects access logs for telemetry purposes. On each
  request, Envoy sends request metadata to the Aperture Agent for a flow control
  decision. Inside the Aperture Agent, the request traverses classifiers,
  rate-limiters, and schedulers, before the decision to accept or drop the
  request is sent back to Envoy. Aperture participates in the
  [OpenTelemetry](https://opentelemetry.io/) tracing protocol as it inserts flow
  classification labels into requests, enabling visualization in tracing tools
  such as [Jaeger](https://www.jaegertracing.io/).
- Aperture SDKs: In addition to service mesh insertion, Aperture provides SDKs
  that can be used by developers to achieve fine-grained flow control at the
  feature level inside service code. For example, an e-commerce app may
  prioritize users in the checkout flow over new sessions when the application
  is experiencing an overload. The Aperture Controller can be programmed to
  degrade features as an escalated recovery action when basic load shedding is
  triggered for several minutes.

## Integrations Modes

The Aperture can be integrated in following ways:

- **Flow Control**

  - [**Envoy**](./flow-control/envoy/istio.md)

    Aperture can be integrated with Envoy and Istio service mesh to provide flow
    control without changing the application code. Envoy's External
    Authorization API is used for control and access logs are collected for
    telemetry. On each request, Envoy sends request metadata to the Aperture
    Agent for flow control decision. For more details on how to set up Aperture
    with Envoy and Istio, refer to
    [our documentation](./flow-control/envoy/istio.md).

  - [**SDKs**](./flow-control/sdk/sdk.md)

    Aperture can also be integrated with service instances using SDKs, which
    allows for easy integration with your application. The SDKs provide a simple
    and flexible way to implement flow control, and can be easily customized to
    meet the specific needs of your application. For more details on how to set
    up Aperture with SDKs, refer to individual languages documentation listed
    below.

    - [GoLang](flow-control/sdk/go/manual.md)
    - [Java](flow-control/sdk/java/manual.md)
    - [JavaScript](./flow-control/sdk/javascript/manual.md)
