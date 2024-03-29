---
title: Armeria
sidebar_position: 3
slug: armeria
keywords:
  - java
  - sdk
  - control
  - points
  - middleware
  - armeria
---

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'
```

### Aperture Java Instrumentation Agent

All requests handled by an Armeria application can have Aperture SDK calls
automatically added into them using [Aperture Instrumentation Agent][javaagent].

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

### Armeria Decorators

[Aperture Java SDK Armeria package](https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-armeria)
contains Armeria decorators that automatically set traffic control points for
decorated services:

<CodeSnippet
lang="java"
snippetName="ArmeriaCreateHTTPService"/>

<CodeSnippet lang="java" snippetName="ArmeriadecorateService"/>

You can instruct the decorator to exclude specific paths from being monitored by
the Aperture SDK. For example, you might want to exclude endpoints used for
health checks. To achieve this, you can add the path(s) you want to ignore to
the `ignoredPaths` field of the SDK, as shown in the following code:

<CodeSnippet lang="java" snippetName="ArmeriaCreateApertureSDK"/>

The paths should be specified as a comma-separated list. Note that the paths you
specify must match exactly. However, you can change this behavior to treat the
paths as regular expressions by setting the `ignoredPathMatchRegex` flag to
true, like so:

```java
  builder
        .setIgnoredPathMatchRegex(true)
```

For more context on using Aperture Armeria Decorators to set control points,
refer to the [example app][armeria-example] available in the repository.

[armeria-example]:
  https://github.com/fluxninja/aperture-java/blob/releases/aperture-java/v2.1.0/examples/armeria-example/src/main/java/com/fluxninja/example/ArmeriaClient.java
[javaagent]: /sdk/java/auto-instrumentation.md
[api-keys]: /reference/cloud-ui/api-keys.md
