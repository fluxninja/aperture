---
title: Get Started
keywords:
  - setup
  - getting started
sidebar_position: 2
sidebar_label: Get Started
---

This guide will help you get started with Aperture in few steps. To begin with
you need to prepare your application to have Aperture integrated. Aperture can
be integrated in multiple ways. You can choose the one that best suits your
application.

## [Install Aperture](./installation/installation.md)

1. [**`aperturectl`**](./installation/aperture-cli/aperture-cli.md)

   Aperture is accompanied by a tool called `aperturectl` that can be used to
   install Aperture in your Kubernetes cluster. Begin with Aperture installation
   by heading over to the
   [Installation](/get-started/installation/installation.md) section.

2. Helm

   Although there is a Helm chart available for installing Aperture, we
   recommend using `aperturectl` as it provides an easier and less cumbersome
   way to get started.

## [Set Up Your Application: Pick an integration](./set-up-application/set-up-application.md)

1. [**Manually setting feature control points**](./set-up-application/manual-control-points.md)

   Using Aperture SDKs, it is easier to manually set feature control points in
   your application. There are SDKs available for multiple languages. You can
   find the list of SDKs [here](../integrations/sdk/sdk.md).

2. [**Middleware Insertions**](./set-up-application/middleware-insertions.md)

   To make it easier to integrate Aperture in your application, we have created
   middleware for popular frameworks like Spring Boot, Netty, Armeria
   [see available middleware](../integrations/sdk/java/java.md). With the help
   of middleware, there isn't much code changes required in your application.
   Some middleware doesn't require any code change at all.

3. [**Service Meshes (Istio, Envoy) & API Gateways**](./set-up-application/service-mesh-and-gateways.md)

   Aperture can be integrated with service meshes like Istio and Envoy. You can
   find the list of supported service meshes
   [here](../integrations/istio/istio.md). With help of service meshes, you can
   control the flow of traffic in your application without any code change. It
   is recommended to use service meshes for Aperture integration as it is easier
   to get started with and doesn't require any code change. You can also
   integrate Aperture with API gateways, checkout the supported
   [API Gateways](../integrations/gateway/gateway.md).

## [Create Your First Policy](./policies/policies.md)

For creating policies, `aperturectl` can be of assistance. Apart from that, it
can help with listing policies, previewing live traffic, and many more things.
Let's explore how to create our first policy using `aperturectl` in the
[Generating and Applying Policies](/get-started/policies/policies.md) section.
