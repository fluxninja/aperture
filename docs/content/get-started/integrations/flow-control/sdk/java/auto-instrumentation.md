---
title: Using instrumentation agent to automatically set control points
sidebar_position: 2
slug: using-instrumentation-agent-to-automatically-set-control-points-using-java-sdk
keywords:
  - java
  - sdk
  - feature
  - flow
  - control
  - points
  - auto
  - instrumentation
---

Java application can be automatically instrumented using Aperture
Instrumentation Agent.

Supported technologies:

| Framework | Supported versions |
| :-------- | :----------------- |
| Armeria   | 1.15+              |
| Netty     | 4.1+               |

Latest version of the Aperture Instrumentation Agent jar file can be downloaded
[here][download_link].

## Running the java agent

To statically load the java agent when running some application jar, use the
following command:

`java -javaagent:path/to/javaagent.jar -jar path/to/application.jar`

### Configuring the java agent

Aperture Java Instrumentation Agent can be configured using a properties file,
system properties or environment variables:

| Property name                          | Environment variable name              | Default value | Description                                                                |
| :------------------------------------- | :------------------------------------- | :------------ | :------------------------------------------------------------------------- |
| aperture.javaagent.config.file         | APERTURE_JAVAAGENT_CONFIG_FILE         |               | Path to a file containing configuration properties                         |
| aperture.agent.hostname                | APERTURE_AGENT_HOSTNAME                | localhost     | Hostname of Aperture Agent to connect to                                   |
| aperture.agent.port                    | APERTURE_AGENT_PORT                    | 8089          | Port of Aperture Agent to connect to                                       |
| aperture.connection.timeout.millis     | APERTURE_CONNECTION_TIMEOUT_MILLIS     | 1000          | Aperture Agent connection timeout in milliseconds                          |
| aperture.javaagent.blocked.paths       | APERTURE_JAVAAGENT_BLOCKED_PATHS       |               | Comma-separated list of paths that should not start a flow                 |
| aperture.javaagent.blocked.paths.regex | APERTURE_JAVAAGENT_BLOCKED_PATHS_REGEX |               | Whether the configured blocked paths should be read as regular expressions |

The priority order is `system.property` > `ENV_VARIABLE` > `properties file`.

Example invocation with commandline-set properties:

```sh
java -javaagent:path/to/javaagent.jar \
-Daperture.agent.hostname="some_host" \
-Daperture.agent.port=12345 \
-Daperture.javaagent.blocked.paths="/health,/connected" \
-jar path/to/application.jar
```

Example invocation using a properties file:

```sh
java -javaagent:path/to/javaagent.jar \
-Daperture.javaagent.config.file="/config.properties" \
-jar path/to/application.jar
```

The `/config.properties` file:

```properties
aperture.agent.hostname=some_host
aperture.agent.port=12345
aperture.javaagent.blocked.paths=/health,/connected
```

[download_link]:
  https://repo1.maven.org/maven2/com/fluxninja/aperture/aperture-javaagent/0.26.0/aperture-javaagent-0.26.0.jar
