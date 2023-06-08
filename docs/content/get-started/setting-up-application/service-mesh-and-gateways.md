---
title: Service Mesh & API Gateways
keywords:
  - Service Meshes
  - API Gateways
sidebar_position: 3
sidebar_label: Service Meshes & API Gateways
---

```mdx-code-block
import { Cards } from '@site/src/components/Cards';
```

One of the most common and easy ways to integrate Aperture into your application
is to use a service mesh or API gateway. Aperture supports Envoy, Istio, Nginx
Gateway, and Kong Gateway.

<!-- vale off -->

## How to integrate Aperture with Service Mesh?

<!-- vale on -->

Aperture, supports integration with EnvoyProxy and Istio. With the help
`aperturectl` which is a CLI tool provided by Aperture, you can easily integrate
Aperture. Check out the complete installation guide in
[Service Mesh Integration](/integrations/istio/istio.md).

<!-- vale off -->

## How to integrate Aperture with API Gateways?

<!-- vale on -->

With the help of Aperture Lua modules and Aperture Plugin for Nginx and Kong
respectively, it is easy to integrate Aperture with API Gateway. You can check
out the complete installation guide in
[API Gateway Integration](/integrations/gateway/gateway.md).

<!-- vale off -->

## What's next?

<!-- vale on -->

Once you complete the Service Mesh or API Gateway integration, head over to
[install Aperture](/get-started/installation/installation.md).
