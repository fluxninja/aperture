---
title: Using instrumentation agent to automatically set control points
sidebar_position: 3
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

Aperture Agent host and port can be set using system properties:

| Property name           | Default value | Description                              |
| :---------------------- | :------------ | :--------------------------------------- |
| aperture.agent.hostname | localhost     | hostname of Aperture Agent to connect to |
| aperture.agent.port     | 8089          | port of Aperture Agent to connect to     |

Example invocation with properties:

```sh
java -javaagent:path/to/javaagent.jar \
-Daperture.agent.hostname="some_host" \
-Daperture.agent.port=12345 \
-jar path/to/application.jar
```

[download_link]: http://localhost:8080
