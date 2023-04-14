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
import {apertureVersion} from '../../../../apertureVersion.js';
import DocCardList from '@theme/DocCardList';
```

For services to communicate with Aperture Agent, [Control Points][flow-control]
must be set to describe where the Flows are happening.

This can be achieved in the following ways:

- [Istio/Envoy integration][istio]
- Use Aperture SDK to set feature or traffic Control Points within services.
  This approach allows fine-grained control.

<a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/sdks/`}>Aperture
SDKs</a> available for popular languages, such as :-

- [Golang][golang]
- [Java][java]
- [JavaScript][javascript]
- [Python][python]

Aperture SDK allows you to manually wrap any function call or code snippet
inside the Service code as a Feature Control Point. Every invocation of the
Feature is a Flow from the perspective of Aperture.

Middleware for popular frameworks is also provided, enabling simple
configuration of traffic Control Points within your service.

<DocCardList />

[flow-control]: /concepts/flow-control/flow-control.md
[istio]: /get-started/integrations/flow-control/envoy/istio.md
[golang]: ./go/manual.md
[java]: ./java/manual.md
[javascript]: ./javascript/manual.md
[python]: ./python/manual.md
