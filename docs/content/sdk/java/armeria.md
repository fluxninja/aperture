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

### Aperture Java Instrumentation Agent

All requests handled by an Armeria application can have Aperture SDK calls
automatically added into them using [Aperture Instrumentation Agent][javaagent].

:::info Agent API Key

You can create an Agent API key for your project in the Aperture Cloud UI. For
more information, refer to
[Define Control Points](/get-started/define-control-points.md).

:::

### Armeria Decorators

[Aperture Java SDK Armeria package](https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-armeria)
contains Armeria decorators that automatically set traffic control points for
decorated services:

```java
    public static HttpService createHelloHTTPService() {
        return new AbstractHttpService() {
            @Override
            protected HttpResponse doGet(ServiceRequestContext ctx, HttpRequest req) {
                return HttpResponse.of("Hello, world!");
            }
        };
    }

    ApertureHTTPService decoratedService = createHelloHTTPService()
        .decorate(ApertureHTTPService.newDecorator(apertureSDK, controlPointName));
    serverBuilder.service("/somePath", decoratedService);
```

You can instruct the decorator to exclude specific paths from being monitored by
the Aperture SDK. For example, you might want to exclude endpoints used for
health checks. To achieve this, you can add the path(s) you want to ignore to
the `ignoredPaths` field of the SDK, as shown in the following code:

```java
ApertureSDK sdk = ApertureSDK.builder()
        .setAddress("ORGANIZATION.app.fluxninja.com:443")
        .setAgentAPIKey("AGENT_API_KEY")
        ...
        .addIgnoredPaths("/healthz,/metrics")
        ...
        .build()
```

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
