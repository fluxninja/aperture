---
title: Get Started
keywords:
  - setup
  - getting started
sidebar_position: 3
sidebar_label: Get Started
---

This guide will help you get started with Aperture in few steps. To begin with
you need to prepare your application to have Aperture integrated. Aperture can
be integrated in multiple ways. You can choose the one that best suits your
application.

## [Setting up your Application: Pick your Integration](./setting-up-application/setting-up-application.md)

1. [**Manually setting feature control points**](./setting-up-application/manual-control-points.md)

   Using Aperture SDKs, it is easier to manually set feature control points in
   your application. There are SDKs available for multiple languages. You can
   find the list of SDKs [here](../integrations/flow-control/sdk/sdk.md).

2. [**Middleware Insertions**](./setting-up-application/middleware-insertions.md)

   To make it easier to integrate Aperture in your application, we have created
   middleware for popular frameworks like Spring Boot, Netty, Armeria
   [see available middleware](../integrations/flow-control/sdk/java/java.md).
   With the help of middleware there isn't much code changes required in your
   application. Some middleware doesn't require any code change at all.

3. [**Service Meshes (Istio, Envoy) & API Gateways**](./setting-up-application/service-mesh-and-gateways.md)

   Aperture can be integrated with service meshes like Istio and Envoy. You can
   find the list of service meshes
   [here](../integrations/flow-control/envoy/envoy.md). With help of service
   meshes, you can control the flow of traffic in your application without any
   code change. It is recommended to use service meshes for Aperture integration
   as it is easier to get started with and doesn't require any code change. You
   can also integrate Aperture with API gateways, checkout the supported
   [API Gateways](../integrations/flow-control/gateway/gateway.md).

## [Installing Aperture](./installation/installation.md)

1. [**`aperturectl`**](./installation/aperture-cli/aperture-cli.md)

   Aperture includes its own CLI tool called `aperturectl`. You can use this
   tool to install Aperture in your Kubernetes cluster. Not just installation,
   it can help you do a many other things like creating policies, preview live
   traffic, and so on.

   Begin with Aperture installation by heading over to the
   [Installation](/get-started/installation/installation.md) section.

2. Helm

   Although there is a Helm chart available for installing Aperture, we
   recommend using `aperturectl` as it provides an easier and less cumbersome
   way to get started.

## [Your First Policy](./policies/policies.md)

As mentioned earlier, `aperturectl` is one of the powerful tools provided by
Aperture. It assists in creating policies, previewing live traffic, and more.
Let's explore how to create our first policy using `aperturectl` in the
[Generating and Applying Policies](/get-started/policies/policies.md) section.
