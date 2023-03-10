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

### Armeria Decorators

<a
href={`https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-armeria`}>Aperture
Java SDK Armeria package</a> contains Armeria decorators that automatically set
traffic Control Points for decorated services:

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
        .decorate(ApertureHTTPService.newDecorator(apertureSDK));
    serverBuilder.service("/somePath", decoratedService);
```

For more context on how to use Aperture Armeria Decorators to set Control
Points, you can take a look at the [example app][armeria-example] available in
our repository.

[armeria-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/armeria-example
[javaagent]:
  /get-started/integrations/flow-control/sdk/java/using-instrumentation-agent-to-automatically-set-control-points-using-java-sdk
