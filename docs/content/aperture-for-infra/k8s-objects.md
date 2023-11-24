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

Policy objects can be created manually or prepared from blueprints via the
[`aperturectl blueprints generate`][generate] command ([Generating Policies and
Dashboards][generating-policies] contains an example how to run this command).
[Here](./guides/service-load-management/service-load-management.md) is an
example of how a Policy object could look such as (look for `Generated Policy`).
Such a Policy can be then applied with regular `kubectl apply`.

[generate]: /reference/aperture-cli/aperturectl/blueprints/generate/generate.md
[generating-policies]: /get-started/policies/policies.md#generating-policies
