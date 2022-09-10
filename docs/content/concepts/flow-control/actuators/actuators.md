---
title: Actuators
sidebar_position: 4
---

An actuator is a [dataplane component][components] that is responsible for
controlling [flows at specified control point][selector], for example by load
shedding certain types of traffic at a given control point.

Actuators are configured as a part of a [policy][policies]. Actuators are
instantiated in Aperture Agents, but to actually affect the flows, they need
[the integration][integrations] to communicate with Agent on every flow.

Aperture support two kinds of acturators: Scheduler-powered Concurrency Limiter
and a Rate-limiter.

## [Concurrency Limiter][concurrency-limiter]

Scheduler-powered Concurrency Limiter is about protecting your services from
overload. Its goal is to limit number of concurrent requests to service to a
level the service can handle. It's implemented by configurable load-shedding.
Thanks to the ability to define workloads of different priorities and weights,
it allows to shed some “less useful” flows, while not affecting the more
important ones.

## [Rate Limiter][rate-limiter]

Rate Limiter is a versatile tool that could be used to protect your services,
but also can be used for different purposes, like giving fairness to your users.
Some potential usecases are:

- protecting the service from bot traffic,
- ratelimiting users,
- preventing from reaching external API quotas.

Thanks to lazy-syncing, rate-limiter can offer low latency overhead.

## Dynamic configuration

Because actuators are defined as circuit components, they can be dynamically
reconfigured based on signals. This is helpful when eg. you want to enable
rate-limiter only on overload.

[concurrency-limiter]: ./concurrency-limiter.md
[rate-limiter]: ./rate-limiter.md
[policies]: /concepts/policies/policies.md
[integrations]: /concepts/flow-control/flow-control.md#integrations
[components]: /concepts/flow-control/flow-control.md#components
[selector]: /concepts/flow-control/selector.md
