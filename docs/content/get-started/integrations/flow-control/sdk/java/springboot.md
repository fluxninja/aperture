---
title: Spring Boot
sidebar_position: 5
slug: spring-boot
keywords:
  - java
  - sdk
  - control
  - points
  - middleware
  - spring
  - boot
---

### Spring Boot Filter

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

[spring-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/spring-example
