---
title: Self-Hosting Aperture
sidebar_position: 6.5
keywords:
  - self-hosted
---

The easiest way to install Aperture is to install just Aperture Agents and point
them to [Aperture Cloud][aperture-cloud], which also [provides the Aperture
Controller][cloud-controller].

If you want to have control over the infrastructure and data, it's also possible
to self-host your own Aperture Controller. In such a setup, Agents and
Controller comprise a fully functional Aperture system.

Note that [Aperture Cloud can integrate][extension-config] with Self-Hosted
Controller too, providing an easy way to manage policies and a holistic view of
the infrastructure, along with tools for OLAP analysis of traffic.

[aperture-cloud]: /introduction.md
[cloud-controller]: /reference/fluxninja.md#cloud-controller
[extension-config]: /reference/fluxninja.md#configuration
