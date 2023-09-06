---
title: Using instrumentation agent to automatically set control points
sidebar_position: 2
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

:::info

The latest version of the Aperture Instrumentation Agent `.jar` file can be
found [here][aperture-javaagent].

:::

## Running the java agent

To statically load the java agent when running some application `.jar`, use the
following command:

`java -javaagent:path/to/javaagent.jar -jar path/to/application.jar`

### Configuring the java agent

Aperture Java Instrumentation Agent can be configured using a properties file,
system properties or environment variables:

<!-- vale off -->

| Property name                          | Environment variable name              | Default value | Description                                                                                                                                                                            |
| :------------------------------------- | :------------------------------------- | :------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| aperture.javaagent.config.file         | APERTURE_JAVAAGENT_CONFIG_FILE         |               | Path to a file containing configuration properties                                                                                                                                     |
| aperture.agent.hostname                | APERTURE_AGENT_HOSTNAME                | localhost     | Hostname of Aperture Agent to connect to                                                                                                                                               |
| aperture.agent.port                    | APERTURE_AGENT_PORT                    | 8089          | Port of Aperture Agent to connect to                                                                                                                                                   |
| aperture.control.point.name            | APERTURE_CONTROL_POINT_NAME            |               | (Required) Name of the control point this agent represents                                                                                                                             |
| aperture.javaagent.enable.fail.open    | APERTURE_JAVAAGENT_ENABLE_FAIL_OPEN    | true          | Sets the fail-open behavior for the client when the Aperture Agent is unreachable. <br /> If set to true, all traffic will pass through; if set to false, all traffic will be blocked. |
| aperture.javaagent.insecure.grpc       | APERTURE_JAVAAGENT_INSECURE_GRPC       | true          | Whether gRPC connection to Aperture Agent should be over plaintext                                                                                                                     |
| aperture.javaagent.root.certificate    | APERTURE_JAVAAGENT_ROOT_CERTIFICATE    |               | Path to a file containing root certificate to be used <br /> (insecure connection must be disabled)                                                                                    |
| aperture.connection.timeout.millis     | APERTURE_CONNECTION_TIMEOUT_MILLIS     | 1000          | Aperture Agent connection timeout in milliseconds                                                                                                                                      |
| aperture.javaagent.ignored.paths       | APERTURE_JAVAAGENT_IGNORED_PATHS       |               | Comma-separated list of paths that should not start a flow                                                                                                                             |
| aperture.javaagent.ignored.paths.regex | APERTURE_JAVAAGENT_IGNORED_PATHS_REGEX |               | Whether the configured ignored paths should be read as regular expressions                                                                                                             |

<!-- vale on -->

:::info

The priority order for configuration look is `system.property` >
`ENV_VARIABLE` > `properties file`.

:::

### Using command line properties

Example invocation with `commandline-set` properties:

```sh
java -javaagent:path/to/javaagent.jar \
-Daperture.agent.hostname="some_host" \
-Daperture.agent.port=12345 \
-Daperture.control.point.name="awesomeFeature" \
-Daperture.javaagent.ignored.paths="/health,/connected" \
-jar path/to/application.jar
```

### Using properties file

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
aperture.control.point.name=awesomeFeature
aperture.javaagent.ignored.paths=/health,/connected
```

[aperture-javaagent]:
  https://repo1.maven.org/maven2/com/fluxninja/aperture/aperture-javaagent
