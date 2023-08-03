---
title: Alerting
keywords:
  - circuit
  - policy
sidebar_position: 5
sidebar_label: Alerting
---

```mdx-code-block
import DocCardList from '@theme/DocCardList';
import Zoom from 'react-medium-image-zoom';
```

## Overview

Aperture provides a versatile solution for defining complex alert criteria using
circuit-based policies. These policies enable intricate signal processing on one
or more metrics, aiding in the precise detection of alert conditions.

Policies are evaluated periodically and query metrics from Prometheus using
PromQL, which are subsequently processed in a circuit to derive intermediate
signals. A [_Decider_](/reference/configuration/spec.md#decider) component can
be used to check if a specific signal surpasses a predetermined threshold and
alert on it.

The alert events are forwarded to a Prometheus Alert Manager endpoint. Operators
can then consume notifications from these alerts through their preferred
channels, ensuring prompt response to any potential issues.

<Zoom>

```mermaid
{@include: ./assets/alerting/alerting.mmd}
```

</Zoom>

The diagram depicts Agents collecting metrics into Prometheus (the metrics can
be collected by any other mechanism). These metrics are queried by the
Controller, processed in the circuit-based policies to compute Alert events
which are relayed to the Alert Manager.

## Example Scenario

Examples of complex alerting scenarios include:

1. Auto-learning the alerting threshold based on past trends.
2. Threshold based alert conditions combining multiple metrics.

<DocCardList />
