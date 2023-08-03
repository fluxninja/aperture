---
title: Adaptive Service Protection
keywords:
  - tutorial
sidebar_position: 1
sidebar_label: Adaptive Service Protection
---

```mdx-code-block
import Zoom from 'react-medium-image-zoom';
import DocCardList from '@theme/DocCardList';
```

## Overview

Adaptive service protection leverages closed-loop feedback of service health
telemetry to dynamically adjust the rate of requests processed by a service.
This adjustment is managed by Aperture Agents, which provide a virtual request
queue at the service's entry point.

The queue adjusts the rate of requests in real-time based on the service's
health, effectively mitigating potential service disruptions and maintaining
optimal performance under varying load conditions. This strategic management of
service load not only maximizes infrastructure utilization and service uptime,
but also ensures the fair admission of requests into the service based on the
priority and weight of each request.

Service overloads can result from a wide variety of failure scenarios, such as
cascading failures where a subset of service instances cause a wider outage, or
service slowdowns that result in failure at dependent services. Metastable
failures, where a system remains in a degraded state long after the original
failure condition has passed, can also lead to service overloads. In such
complex failure scenarios, Aperture's load scheduling feature offers a reliable
safeguard, ensuring that your system maintains optimal performance and uptime.

<Zoom>

```mermaid
{@include: ./assets/adaptive-service-protection/adaptive-service-protection.mmd}
```

</Zoom>

The diagram illustrates the working of a load scheduling policy. The policy is
evaluated at the Controller, which analyzes health signals in real-time. Based
on these metrics, it calculates a load multiplier, which is relayed to the
Agents. The Agents then adjust the rate of requests locally based on the load
multiplier applied to the recent rate of requests.

:::note

Aperture facilitates the observation of health signals from various services.
For instance, adaptive service protection can also be implemented based on the
health observation of an upstream service in relation to a downstream service.

:::

## Example Scenario

Consider the scenario of an e-commerce platform during a major sale event. To
handle the increased traffic, the load scheduling policy will monitor the health
of the service to detect overloads. The request rate will be adjusted in case
the service begins to deteriorate. This prevents service overloads, which reduce
good throughput and lead to cascading failures. Load scheduling policy ensures
smooth service operation and a consistent user experience for a successful sales
event.

<DocCardList />
