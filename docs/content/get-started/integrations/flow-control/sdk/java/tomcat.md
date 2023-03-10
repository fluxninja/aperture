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

[tomcat-example]:
  https://github.com/fluxninja/aperture/tree/main/sdks/aperture-java/examples/tomcat-example
