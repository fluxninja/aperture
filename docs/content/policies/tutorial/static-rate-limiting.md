---
title: Static Rate Limiting
keywords:
  - policies
  - ratelimit
sidebar_position: 1
---

```mdx-code-block
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Zoom from 'react-medium-image-zoom';
```

## Regulating heavy-hitters

One of the simplest flow-control policies is static rate limiting. It uses
[rate limiting actuator](/concepts/flow-control/rate-limiter.md) to block
specific flow labels that exceed their quota within a certain timeframe (limit
reset interval).

### Circuit Diagram

<Zoom>

```mermaid
{@include: ./assets/gen/static-rate-limiting/jsonnet/static-rate-limiting_0.mmd}
```

</Zoom>

### Policy

In this example, we will be rate limiting unique users based on the `User-Id`
header in the HTTP traffic. This header is provided by Envoy proxy under the
label key `http.request.header.user_id` (See
[Flow Labels](/concepts/flow-control/flow-label.md)).

We will be configuring our rate limiter to allow at most `120 requests` for each
user in the `60s` period.

In addition, we will be configuring our rate limiter to apply these limits to
`ingress` traffic on Kubernetes service
`service1-demo-app.demoapp.svc.cluster.local`.

```mdx-code-block
<Tabs>
<TabItem value="Jsonnet">
```

```jsonnet
local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local component = aperture.spec.v1.Component;
local rateLimiter = aperture.spec.v1.RateLimiter;
local selector = aperture.spec.v1.Selector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local circuit = aperture.spec.v1.Circuit;
local port = aperture.spec.v1.Port;

local rateLimitPort = port.new() + port.withSignalName('RATE_LIMIT');

local svcSelector =
  selector.new()
  + selector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + selector.withFlowSelector(
    flowSelector.new()
    + flowSelector.withControlPoint({ traffic: 'ingress' })
  );

local policyDef =
  policy.new()
  + policy.withCircuit(
    circuit.new()
    + circuit.withEvaluationInterval('300s')
    + circuit.withComponents([
      component.withRateLimiter(
        rateLimiter.new()
        + rateLimiter.withInPorts({ limit: port.withConstantValue(120) })
        + rateLimiter.withSelector(svcSelector)
        + rateLimiter.withLimitResetInterval('60s')
        + rateLimiter.withLabelKey('http.request.header.user_id')
      ),
    ]),
  );

local policyResource = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'service1-demo-app',
    labels: {
      'fluxninja.com/validate': 'true',
    },
  },
  spec: policyDef,
};

[
  policyResource,
]
```

```mdx-code-block
</TabItem>
<TabItem value="YAML">
```

```yaml
{@include: ./assets/gen/static-rate-limiting/jsonnet/static-rate-limiting_0.yaml}
```

```mdx-code-block
</TabItem>
</Tabs>
```
