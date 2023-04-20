---
title: Tomcat
sidebar_position: 6
slug: tomcat
keywords:
  - java
  - sdk
  - control
  - points
  - middleware
  - tomcat
---

### Tomcat Filter

[Aperture Java SDK servlet package](https://search.maven.org/artifact/com.fluxninja.aperture/aperture-java-servlet)
contains Aperture Filter that can be added to the web.xml file to automatically
set traffic control points for relevant services:

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

You can instruct the filter to exclude specific paths from being monitored by
the Aperture SDK. For example, you might want to exclude endpoints used for
health checks. To achieve this, you can add the path(s) you want to ignore to
the filter configuration, as shown in the following code:

```xml
<init-param>
  <param-name>ignored_paths</param-name>
  <param-value>/healthz,/metrics</param-value>
</init-param>
```

The paths should be specified as a comma-separated list. Note that the paths you
specify must match exactly. However, you can change this behavior to treat the
paths as regular expressions by setting the `ignored_paths_match_regex`
initialization parameter to true, like so:

```xml
<init-param>
  <param-name>ignored_paths_match_regex</param-name>
  <param-value>true</param-value>
</init-param>
```

For usage, you can view the [example app][tomcat-example] available in the
repository.

[tomcat-example]:
  https://github.com/fluxninja/aperture-java/tree/releases/aperture-java/v1.0.0/examples/tomcat-example
