---
title: Scheduler
description: Scheduler
keywords:
  - scheduler
  - tokens
  - priority
  - fairness
  - queuing
---

:::info
See also what are [actuators](actuators.md) in general.
:::

:::info
See also [Concurrency Limiter reference](/reference/configuration/policies.md#-languagev1concurrencylimiter).
:::

# Scheduler

## Workload

Workloads are a group of flows based on common attributes. Workloads are
expressed by label matcher rules in Aperture. Aperture Agents schedule workloads
based on their priorities and by estimating their [tokens](#tokens).

## Tokens

Tokens represent the cost of admitting a flow in the system. Most commonly,
tokens are estimated based on milliseconds of response time that is observed
when a flow is processed. Token estimation of flows within a workload is crucial
when making flow control decisions. Concept of tokens is aligned with
[Little's Law](https://en.wikipedia.org/wiki/Little%27s_law), which defines a
relationship between response times, arrival rate and requests currently in the
system (concurrency).

In some cases, tokens can be represented as a number of requests instead of
response times, e.g. when performing flow control on external APIs that have
hard rate-limits.

## Token bucket

Aperture Agents use a variant of a
[token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) is used to
control the flows entering the system. Each flow has to acquire tokens from the
bucket within a deadline period in order to be admitted.

## Timeouts

The timeout parameter decides how long a request in the workload can wait for tokens.
This value impacts fairness because the larger the timeout the higher the chance a
request has to get scheduled.
