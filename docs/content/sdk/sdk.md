---
title: SDKs
description: Setup Control Points using SDK libraries
keywords:
  - setup
  - flow
  - control
  - points
  - sdk
sidebar_position: 3
sidebar_label: SDKs
---

```mdx-code-block
import {apertureVersion} from '../apertureVersion.js';
import DocCardList from '@theme/DocCardList';
import Zoom from 'react-medium-image-zoom';
```

For services to control flows with Aperture Agent, [Control
Points][control-point] must be set within the service. Aperture SDKs can be used
to set feature or traffic (HTTP and gRPC) control points within the service
code. This approach allows for fine-grained flow control.

<a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/sdks/`}>Aperture
SDKs</a> are available for popular languages:

- [Golang][golang]
- [Java][java]
- [JavaScript][javascript]
- [Python][python]
- [.NET][dotnet]

Aperture SDK allows you to manually wrap any function call or code snippet
inside the service code as a feature control point. Every invocation of the
feature is a flow from the perspective of Aperture.

![SDK](./assets/sdks-dark.svg#gh-dark-mode-only)

![SDK](./assets/sdks-light.svg#gh-light-mode-only)

## 🧩 Middleware

Aperture includes middleware for the following frameworks, helping to set up
control points with less code changes:

- [Armeria][armeria]
- [Netty][netty]
- [Tomcat][tomcat]
- [Spring Boot][spring-boot]

<DocCardList />

[control-point]: /concepts/control-point.md
[golang]: ./go/manual.md
[java]: ./java/manual.md
[javascript]: ./javascript/manual.md
[python]: ./python/manual.md
[dotnet]: ./dotnet/manual.md
[netty]: ./java/netty.md#netty-handler
[tomcat]: ./java/tomcat.md#tomcat-filter
[spring-boot]: ./java/springboot.md#spring-boot-filter
[armeria]: ./java/armeria.md#armeria-decorators
