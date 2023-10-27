---
title: Middleware Insertions
keywords:
  - Middleware Insertions
sidebar_position: 2
sidebar_label: Middleware Insertions
---

```mdx-code-block
import { Cards } from '@site/src/components/Cards';
```

Aperture supports inserting middleware into the request pipeline, which makes it
easy to integrate with less code changes.

<!-- vale off -->

## What is a Middleware Insertion?

<!-- vale on -->

It is a way to add a new layer of functionality to the request pipeline. It is a
piece of code that can be executed before or after the request is processed by
the application. Allowing you to add new functionality to the request pipeline
without changing the application code.

<!-- vale off -->

## How to add a Middleware?

<!-- vale of -->

There are multiple ways to add a middleware, and it depends on the language and
framework you are using. For example, in Spring Boot, you can register a Spring
Boot Filter, in Armeria you can register a decorator, in Netty, you can register
an Aperture Handler. Let's examine how to add a middleware insertion in Spring
Boot.

<!-- vale off -->

## How to add a Middleware in Spring Boot?

<!-- vale on -->

In Spring Boot, you can register an Aperture filter, which automatically sets
the feature control points. Here is an example:

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

:::info

Aperture provides different middlewares for different java frameworks, which you
can check out in the [Java section](/integrations/sdk/java/java.md). There are
examples available for each framework.

:::

<!-- vale off -->

## What's next?

<!-- vale on -->

Once the middleware insertion is done, head over to
[install Aperture](/get-started/installation/installation.md).
