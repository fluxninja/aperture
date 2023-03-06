---
title: Setting traffic control points using middleware
sidebar_position: 2
slug: setting-traffic-control-points-using-middleware-using-java-sdk
keywords:
  - java
  - sdk
  - control
  - points
  - middleware
---

Aperture Java SDK provides several middlewares for popular frameworks that
provide a simple way to define traffic Control Points.

### Armeria decorators

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

### Spring Boot filter

<a
href={`https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-servlet`}>Aperture
Java SDK servlet package</a> contains Aperture Filter that can be registered in
Spring Boot application to automatically set traffic Control Points for relevant
services:

```java

import com.fluxninja.aperture.servlet.jakarta.ApertureFilter;

...

@RestController
public class AppController {

    ...

    @RequestMapping(value = "/super", method = RequestMethod.GET)
    public String hello() {
        return "Hello World";
    }

    ...

    @Bean
    public FilterRegistrationBean<ApertureFilter> apertureFilter(Environment env){
        FilterRegistrationBean<ApertureFilter> registrationBean = new FilterRegistrationBean<>();

        registrationBean.setFilter(new ApertureFilter());
        registrationBean.addUrlPatterns("/super");

        registrationBean.addInitParameter("agent_host", "localhost");
        registrationBean.addInitParameter("agent_port", "8089");

        return registrationBean;
    }
}
```

For example usage, you can view the [example app][spring-example] available in
our repository.

### Tomcat filter

<a
href={`https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-servlet`}>Aperture
Java SDK servlet package</a> contains Aperture Filter that can be added to the
web.xml file to automatically set traffic Control Points for relevant services:

```xml
    <filter>
        <filter-name>ApertureFilter</filter-name>
        <filter-class>com.fluxninja.aperture.servlet.javax.ApertureFilter</filter-class>
        <init-param>
            <param-name>agent_host</param-name>
            <param-value>localhost</param-value>
        </init-param>
        <init-param>
            <param-name>agent_port</param-name>
            <param-value>8089</param-value>
        </init-param>
    </filter>
```

For example usage, you can view the [example app][tomcat-example] available in
our repository.

[armeria-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/armeria-example
[tomcat-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/tomcat-example
[spring-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/spring-example
