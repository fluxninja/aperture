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

[Aperture Java SDK servlet package](https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-servlet)
contains Aperture Filter that can be registered in Spring Boot application to
automatically set traffic control points for relevant services:

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

        registrationBean.addInitParameter("agent_address", "ORGANIZATION.app.fluxninja.com:443");
        registrationBean.addInitParameter("agent_api_key", "AGENT_API_KEY");

        return registrationBean;
    }
}
```

You can instruct the filter to exclude specific paths from being monitored by
the Aperture SDK. For example, you might want to exclude endpoints used for
health checks. To achieve this, you can add the path(s) you want to ignore to
the filter configuration, as shown in the following code:

```java
registrationBean.addInitParameter("ignored_paths", "/healthz,/metrics");
```

The paths should be specified as a comma-separated list. Note that the paths you
specify must match exactly. However, you can change this behavior to treat the
paths as regular expressions by setting the `ignored_paths_match_regex`
initialization parameter to true, like so:

```java
registrationBean.addInitParameter("ignored_paths_match_regex", "true");
```

For usage, you can view the [example app][spring-example] available in the
repository.

[spring-example]:
  https://github.com/fluxninja/aperture-java/blob/releases/aperture-java/v2.1.0/examples/spring-example/src/main/java/com/fluxninja/example/controller/AppController.java
