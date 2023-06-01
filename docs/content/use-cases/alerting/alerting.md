---
title: Alerting
keywords:
  - circuit
  - policy
sidebar_position: 5
sidebar_label: Alerting
---

## Overview

Monitoring the health of a service is a critical aspect of ensuring reliable
operations. In this example, we will demonstrate how to detect an overload state
of a service and send an alert using Aperture's declarative policy language. The
policy will create a circuit that models the normal latency behavior of the
service using an exponential moving average (EMA). This enables the alerting
policy to automatically learn the normal latency threshold of each service,
reducing the need for manual tuning of alert policies for individual services.

<Zoom>

```mermaid
{@include: ../assets/alerting.mmd}
```

The graph depicts the process of the Agent writing any overload detection to
Prometheus. These metrics are queried by the Controller and relayed to the Alert
Manager. The Controller includes a signal processing mechanism that can
distinguish false positives by comparing current latency with the setpoint.

:::note

Aperture facilitates the observation of health signals from various services.
For instance, service protection can also be implemented based on the health
observation of an upstream service in relation to a downstream service.

:::

</Zoom>

## Real World Scenario

Regarding monitoring service health, envision a data center managing multiple
servers. Using Aperture's declarative policy language, they can monitor each
server's health. In case a server enters an overload state, the policy triggers
an alert, allowing immediate corrective action. This proactive approach
minimizes downtime and ensures uninterrupted service delivery.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
