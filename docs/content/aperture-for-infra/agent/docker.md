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
import {apertureVersion, apertureVersionWithOutV} from '../../apertureVersion.js';
```

Below are the instructions to install the Aperture Agent on Docker.

## Prerequisites

1. Install [Docker](https://docs.docker.com/get-docker/) on your system.

## Installation

1. Create a file named `agent.yaml` with the below content for passing the
   configuration to the Aperture Agent:

   ```yaml
   fluxninja:
     enable_cloud_controller: true
     endpoint: "ORGANIZATION_NAME.app.fluxninja.com:443"
     api_key: "API_KEY"
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

   Replace `ORGANIZATION_NAME` with the Aperture Cloud organization name and
   `API_KEY` with the API key linked to the project. Navigate to the
   **`Aperture`** tab in the sidebar menu and then select **`API Keys`** in the
   top bar. From there, you can either copy the existing key or create a new one
   by clicking on **`Create API Key`**.

   :::note

   If you are using a Self-Hosted Controller Aperture Controller, modify the
   above configuration as explained in
   [Self-hosted Agent Configuration](/aperture-for-infra/agent/agent.md#agent-self-hosted-controller).

   :::

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
