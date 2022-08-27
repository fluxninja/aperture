---
title: Selector
sidebar_position: 1
keywords:
  - flows
  - services
  - discovery
  - labels
---

# Selector

Flow observability and control components are instantiated on Aperture Agents
and select flows based on scoping rules defined in the Selectors.

A Selector consists of following fields:

## Agent Group

Agent Group is a flexible label that defines a collection of agents that operate
as peers. For example, an Agent Group can be a Kubernetes cluster name in case
of DaemonSet deployment or it can be a service name for sidecar deployments.

All the agents within a Agent Group instantiate the same set of policies,
published by Aperture Controller.

Agents within an Agent Group form a peer-to-peer network. Agents synchronize
fine-grained state such as per label global counters that are used for rate
limiting purposes.

## Service

Service in Aperture is similar to services tracked in Kubernetes, Consul etc. A
Service is a collection of entities delivering a common functionality, such as,
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

## Control Point

## Label Matcher
