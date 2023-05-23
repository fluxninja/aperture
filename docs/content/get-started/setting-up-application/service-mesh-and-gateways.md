---
title: Service Mesh & API Getaways
keywords:
  - Service Meshes
  - API Getaways
sidebar_position: 3
sidebar_label: Service Meshes & API Getaways
---

```mdx-code-block
import { Cards } from '@site/src/components/Cards';
```

One of the most common and easy ways to integrate Aperture into your application
is to use a service mesh or API getaway. Aperture supports Envoy, Istio,
Linkerd, and Kong.

<!-- vale off -->

## How to integrate Aperture with Service Mesh?

<!-- vale on -->

Aperture, supports integration to EnvoyProxy and Istio. With the help
`aperturectl` which is a CLI tool provided by Aperture you can easily integrate
Aperture. Check out the complete installation guide in
[Service Mesh Integration](/integrations/flow-control/envoy/envoy.md).

<!-- vale off -->

## How to integrate Aperture with API Getaways?

<!-- vale on -->

With the help of Aperture Lua modules and Aperture Plugin for Nginx and Kong
respectively, it is easy to integrate Aperture with API Getaways. You can check
out the complete installation guide in
[API Getaway Integration](/integrations/flow-control/gateway/gateway.md).

<!-- vale off -->

## What's next?

<!-- vale on -->

Once you complete the Service Mesh or API Gateway Integration head over to
install Aperture.

```mdx-code-block

<Cards data={[
  {
    title: "Install Aperture",
    description: "Install Controller and Agent in your environment",
    url: "/get-started/installation/",
  },
  ]}/>
```
