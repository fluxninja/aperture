---
title: Alerting
keywords:
  - circuit
  - policy
sidebar_position: 5
sidebar_label: Alerting
---

## Overview

Aperture's declarative policies can be fine-tuned to identify overload states
and record them to databases. Beyond just configuring Aperture to send alerts to
an Alert Manager via database writing, it can also monitor alert signals
distinct from the standard four gold signals. This affords an extra layer of
flexibility and coverage, ensuring that potential issues are identified and
addressed promptly, maintaining the robustness and reliability of the service.

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
