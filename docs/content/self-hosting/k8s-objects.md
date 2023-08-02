---
title: Managing Policies Using Kubernetes Objects
sidebar_position: 4
sidebar_label: Kubernetes Objects
keywords:
  - kubectl
  - k8s
  - crd
  - resource
  - custom
  - apply
---

With the self-hosted controller deployed on Kubernetes, you gain the possibility
to manage policies using Kubernetes Objects (in addition to the usual way via
`aperturectl`). Aperture Controller installation includes the Policy Custom
Resource Definition.

Policy objects can be created manually or prepared from blueprints via
[`aperturectl blueprints generate`][generate] command ([Generating Policies and
Dashboards][] contains an example how to run this command). [Take a look here,
how a Policy object could look like][Example] (look for "Generated Policy").
Such a Policy can be then applied with regular `kubectl apply`.

<!-- prettier-ignore-start -->

[generate]: /reference/aperturectl/blueprints/generate/generate.md
[Generating Policies and Dashboards]: /get-started/policies/policies.md#generating-policies-and-dashboards
[Example]: /use-cases/adaptive-service-protection/average-latency-feedback.md

<!-- prettier-ignore-end -->
