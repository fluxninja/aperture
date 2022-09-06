---
title: Flow Selector
sidebar_position: 1
keywords:
  - flows
  - services
  - discovery
  - labels
---

:::info See also
[Selector reference](/reference/configuration/policies#-v1selector) :::

Flow observability and control components are instantiated on Aperture Agents
and they select flows based on scoping rules defined in Selectors.

A Selector consists of following fields:

### Agent Group

Agent Group is a flexible label that defines a collection of agents that operate
as peers. For example, an Agent Group can be a Kubernetes cluster name in case
of DaemonSet deployment or it can be a service name for sidecar deployments.

Agents within an Agent Group form a peer-to-peer network. Agents synchronize
fine-grained state such as per label global counters that are used for rate
limiting purposes.

All the agents within a Agent Group instantiate the same set of policies,
published by Aperture Controller.

### Service

Service in Aperture is similar to services tracked in Kubernetes or Consul. A
Service is a collection of entities delivering a common functionality, such as
checkout, billing etc. Aperture maintains a mapping of entity IP addresses to
Service names. For each flow control decision request sent by an entity,
Aperture looks up the service name and then decides which flow control
components to execute.

Aperture Agents perform automated discovery of services and entities in
environments such as Kubernetes and watch for any changes. Service and entity
entries can also be created manually via configuration.

In addition, Aperture also has a concept of a `*` catch-all service. When the
Selector contains a catch-all service, it matches for all discovered entities
within a Agent Group.

### Control Point

A policy or rule is configured for a given Control Point within a service.
Control Point is either a library feature name or one of ingress/egress traffic
points.

### Label Matcher

Label Matcher is part of the classifier for whether a map of labels should be
considered a match or not. If multiple requirements are set, they are all ANDed.
An empty label matcher always matches.

This matcher allows matching the following labels:

- Flow labels - We can only match flow labels that were created at a previous
  control point

- Request labels - Request labels are always prefixed with request\_ and request
  headers are only available for traffic control points
