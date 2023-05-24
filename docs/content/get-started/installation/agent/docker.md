---
title: Docker
description: Install Aperture Agent on Docker
keywords:
  - install
  - setup
  - agent
  - docker
sidebar_position: 4
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import {apertureVersion, apertureVersionWithOutV} from '../../../apertureVersion.js';
```

Below are the instructions to install the Aperture Agent on Docker.

## Prerequisites

1. Install [Docker](https://docs.docker.com/get-docker/) on your system.

## Installation

1. Create a file named `agent.yaml` with below content for passing the
   configuration to the Aperture Agent:

   ```yaml
   etcd:
     endpoints:
       - "ETCD_ENDPOINT_HERE"
   prometheus:
     address: "PROMETHEUS_ADDRESS_HERE"
   agent_functions:
     endpoints: ["CONTROLLER_ENDPOINT_HERE"]
   log:
     level: info
     pretty_console: true
     non_blocking: false
   otel:
     disable_kubernetes_scraper: true
     disable_kubelet_scraper: true
   auto_scale:
     kubernetes:
       enabled: false
   service_discovery:
     kubernetes:
       enabled: false
   ```

   Replace the values of `ETCD_ENDPOINT_HERE` and `PROMETHEUS_ADDRESS_HERE` with
   the actual values of etcd and Prometheus, which are also being used by the
   Aperture Controller you want these Agents to connect to.
   `CONTROLLER_ENDPOINT_HERE` should point to Aperture Controller. If you skip
   it, some sub-commands `aperturectl` commands won't work.

   If you have installed the
   [Aperture Controller](/get-started/installation/controller/docker.md) on the
   same Docker environment, with etcd and Prometheus, the values for the values
   for `ETCD_ENDPOINT_HERE`, `PROMETHEUS_ADDRESS_HERE` and
   `CONTROLLER_ENDPOINT_HERE` would be as below:

   ```yaml
   etcd:
     endpoints:
       - "http://etcd:2379"
   prometheus:
     address: "http://prometheus:9090"
   agent_functions:
     endpoints: ["aperture-controller:8080"]
     client:
       grpc:
         insecure: true
   log:
     level: info
     pretty_console: true
     non_blocking: false
   otel:
     disable_kubernetes_scraper: true
     disable_kubelet_scraper: true
   auto_scale:
     kubernetes:
       enabled: false
   service_discovery:
     kubernetes:
       enabled: false
   ```

   All the configuration parameters for the Aperture Agent are available
   [here](/reference/configuration/agent.md).

2. Run the below command to start the Aperture Agent container:

   <CodeBlock language="bash">
   {`docker run -d \\
   -p 8081:8080 \\
   --name aperture-agent \\
   --network aperture \\
   -v "$(pwd)"/agent.yaml:/etc/aperture/aperture-agent/config/aperture-agent.yaml:ro \\
   docker.io/fluxninja/aperture-agent:${apertureVersionWithOutV}`}
   </CodeBlock>

3. Verify that the Aperture Agent container is in the `healthy` state:

   ```bash
   docker run -it --rm \
   --network aperture curlimages/curl \
   sh -c \
   'while [[ \"$(curl -s -o /dev/null -w %{http_code} aperture-agent:8080/v1/status/system/readiness)\" != \"200\" ]]; \
   do echo "aperture-agent is starting"; sleep 1; done && \
   echo "aperture-agent is now healty!"'
   ```

## Uninstall

1. Run the below command to stop and remove the Aperture Agent container:

   ```bash
   docker rm -f aperture-agent
   ```
