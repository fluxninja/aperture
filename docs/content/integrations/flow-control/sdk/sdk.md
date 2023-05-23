---
title: SDKs
description: Setup Control Points using SDK libraries
slug: setup-control-points-using-sdks-libraries
keywords:
  - setup
  - flow
  - control
  - points
  - sdk
sidebar_position: 3
---

```mdx-code-block
import {apertureVersion} from '../../../apertureVersion.js';
import DocCardList from '@theme/DocCardList';
```

For services to control flows with Aperture Agent, [Control
Points][flow-control] must be set within the service.

This can be achieved in the following ways:

- [Istio/Envoy integration][istio] for controlling HTTP or gRPC requests flowing
  through the service.
- Aperture SDKs can be used to set feature or traffic (HTTP and gRPC) control
  points within the service code. This approach allows for fine-grained flow
  control.

<a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/sdks/`}>Aperture
SDKs</a> available for popular languages, such as :-

- [Golang][golang]
- [Java][java]
- [JavaScript][javascript]
- [Python][python]

Aperture SDK allows you to manually wrap any function call or code snippet
inside the service code as a feature control point. Every invocation of the
feature is a flow from the perspective of Aperture.

Middleware for popular frameworks is also provided, enabling simple
configuration of traffic control points within your service.

<DocCardList />

[flow-control]: /concepts/flow-control/flow-control.md
[istio]: ../envoy/istio.md
[golang]: ./go/manual.md
[java]: ./java/manual.md
[javascript]: ./javascript/manual.md
[python]: ./python/manual.md
