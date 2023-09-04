---
title: Docker
description: Install Aperture Controller on Docker
keywords:
  - install
  - setup
  - controller
  - docker
sidebar_position: 2
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import {apertureVersion, apertureVersionWithOutV} from '../../apertureVersion.js';
```

Below are the instructions to install the Aperture Controller on Docker.

## Prerequisites

1. Install [Docker](https://docs.docker.com/get-docker/) on your system.

2. Create a Docker network:

   ```bash
   docker network create aperture --driver bridge
   ```

## Installation of etcd

:::info Note

etcd is required for the Aperture Controller to function. If you already have an
etcd cluster running, you can skip these steps.

:::

1. Create a volume which will be used by etcd:

   ```bash
   docker volume create etcd-data
   ```

2. Create a file named `etcd.env` with below content for passing the environment
   variables to the etcd:

   ```bash
   ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
   ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379
   ETCD_DATA_DIR=/bitnami/etcd/data
   ETCDCTL_API=3
   ETCD_ON_K8S=no
   ETCD_START_FROM_SNAPSHOT=no
   ETCD_DISASTER_RECOVERY=no
   ALLOW_NONE_AUTHENTICATION=yes
   ETCD_AUTH_TOKEN=simple
   ETCD_AUTO_COMPACTION_MODE=periodic
   ETCD_AUTO_COMPACTION_RETENTION=24
   BITNAMI_DEBUG=false
   ```

3. Run the below command to start the etcd container:

   ```bash
   docker run -d \
   -p 2379:2379 \
   --env-file etcd.env \
   --name etcd \
   --network aperture \
   --volume etcd-data:/bitnami/etcd:rw \
   --health-cmd /opt/bitnami/scripts/etcd/healthcheck.sh \
   --health-interval 10s \
   --health-retries 5 \
   --health-start-period 30s \
   --health-timeout 5s \
   docker.io/bitnami/etcd:3.5.8-debian-11-r0
   ```

4. Verify that the etcd container is in the `healthy` state:

   ```bash
   CONTAINER_ID="etcd"; \
   while [[ "$(docker inspect -f '{{.State.Health.Status}}' "${CONTAINER_ID}")" != "healthy" ]]; \
   do echo "${CONTAINER_ID}" is starting; sleep 1; done; \
   echo "${CONTAINER_ID}" is now healthy!
   ```

## Installation of Prometheus

:::info Note

Prometheus is required for the Aperture Controller to function. If you already
have a Prometheus instance running, you can skip these steps.

:::

1. Create a volume which will be used by Prometheus:

   ```bash
   docker volume create prometheus-data
   ```

2. Create a file named `prometheus.yml` for passing the configuration to the
   Prometheus:

   ```yaml
   global:
     evaluation_interval: 1m
     scrape_interval: 1m
     scrape_timeout: 10s
   scrape_configs: []
   ```

3. Run the below command to set correct permissions for the Prometheus data
   volume:

   ```bash
   docker run --rm -v prometheus-data:/data:rw busybox chown -R 65534:65534 /data
   ```

4. Run the below command to start the Prometheus container:

   ```bash
   docker run -d \
   -p 9090:9090 \
   --name prometheus \
   --network aperture \
   --volume prometheus-data:/data:rw \
   --volume "$(pwd)"/prometheus.yaml:/etc/config/prometheus.yaml:ro \
   --health-cmd "wget --quiet --output-document=/dev/null --server-response --spider localhost:9090/-/ready 2>&1 | grep \"200 OK\" || exit 1" \
   --health-interval 10s \
   --health-retries 5 \
   --health-start-period 30s \
   --health-timeout 5s \
   quay.io/prometheus/prometheus:v2.33.5 \
   --config.file=/etc/config/prometheus.yaml \
   --storage.tsdb.path=/data \
   --storage.tsdb.retention.time=1d \
   --web.enable-remote-write-receiver \
   --web.console.libraries=/etc/prometheus/console_libraries \
   --web.console.templates=/etc/prometheus/consoles
   ```

5. Verify that the Prometheus container is in the `healthy` state:

   ```bash
   CONTAINER_ID="prometheus"; \
   while [[ "$(docker inspect -f '{{.State.Health.Status}}' "${CONTAINER_ID}")" != "healthy" ]]; \
   do echo "${CONTAINER_ID}" is starting; sleep 1; done; \
   echo "${CONTAINER_ID}" is now healthy!
   ```

## Installation of Aperture Controller

1. Create a file named `controller.yaml` with below content for passing the
   configuration to the Aperture Controller:

   ```yaml
   etcd:
     endpoints:
       - http://etcd:2379
   prometheus:
     address: "http://prometheus:9090"
   log:
     level: info
     pretty_console: true
     non_blocking: false
   ```

   All the configuration parameters for the Aperture Controller are available
   [here](/reference/configuration/controller.md).

2. Run the below command to start the Aperture Controller container:

   <CodeBlock language="bash">
   {`docker run -d \\
   -p 8080:8080 \\
   --name aperture-controller \\
   --network aperture \\
   -v "$(pwd)"/controller.yaml:/etc/aperture/aperture-controller/config/aperture-controller.yaml:ro \\
   docker.io/fluxninja/aperture-controller:${apertureVersionWithOutV}`}
   </CodeBlock>

## Installation using docker compose

1. Install [docker-compose](https://docs.docker.com/compose/install/).

2. Create `prometheus.yaml` and `controller.yaml` files with the same content as
   mentioned in the above steps.

3. Create a file named `docker-compose.yaml` with below content:

   <details><summary>docker-compose.yaml</summary>
   <p>
   <CodeBlock language="yaml">
   {`version: '3'
   services:
     etcd:
       image: docker.io/bitnami/etcd:3.5.8-debian-11-r0
       container_name: etcd
       healthcheck:
         test: ["CMD", "/opt/bitnami/scripts/etcd/healthcheck.sh"]
         interval: 10s
         timeout: 5s
         retries: 5
         start_period: 60s
       environment:
         - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
         - ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379
         - ETCD_DATA_DIR=/bitnami/etcd/data
         - ETCDCTL_API=3
         - ETCD_ON_K8S=no
         - ETCD_START_FROM_SNAPSHOT=no
         - ETCD_DISASTER_RECOVERY=no
         - ALLOW_NONE_AUTHENTICATION=yes
         - ETCD_AUTH_TOKEN=simple
         - ETCD_AUTO_COMPACTION_MODE=periodic
         - ETCD_AUTO_COMPACTION_RETENTION=24
         - BITNAMI_DEBUG=false
       ports:
         - 2379:2379
       volumes:
         - etcd-data:/bitnami/etcd:rw
       networks:
         - aperture
     prometheus-init:
       image: busybox
       user: root
       group_add:
         - '65534'
       volumes:
         - prometheus-data:/data:rw
       command: chown -R 65534:65534 /data
     prometheus:
       image: quay.io/prometheus/prometheus:v2.33.5
       container_name: prometheus
       command:
         - '--config.file=/etc/config/prometheus.yaml'
         - '--storage.tsdb.path=/data'
         - '--storage.tsdb.retention.time=1d'
         - '--web.enable-remote-write-receiver'
         - '--web.console.libraries=/etc/prometheus/console_libraries'
         - '--web.console.templates=/etc/prometheus/consoles'
       ports:
         - 9090:9090
       user: "65534:65534"
       depends_on:
         prometheus-init:
           condition: service_completed_successfully
       volumes:
         - prometheus.yaml:/etc/config/prometheus.yaml:ro
         - prometheus-data:/data:rw
       networks:
         - aperture
     controller:
       image: docker.io/fluxninja/aperture-controller:${apertureVersionWithOutV}
       container_name: aperture-controller
       ports:
         - 8080:8080
       volumes:
         - controller.yaml:/etc/aperture/aperture-controller/config/aperture-controller.yaml:ro
       networks:
         - aperture
       restart: on-failure
       depends_on:
         etcd:
           condition: service_healthy
   volumes:
     prometheus-data:
     etcd-data:
   networks:
     aperture:
       name: aperture
       external: true
   `}
   </CodeBlock>
   </p>
   </details>

4. Run the below command to start the etcd, Prometheus and Aperture Controller
   containers:

   ```bash
   docker compose up -d
   ```

## Verifying the Installation

1. Run the below command to verify that the Aperture Controller container is in
   the `healthy` state:

   ```bash
   docker run -it --rm \
   --network aperture curlimages/curl \
   sh -c \
   'while [[ \"$(curl -s -o /dev/null -w %{http_code} aperture-controller:8080/v1/status/system/readiness)\" != \"200\" ]]; \
   do echo "aperture-control is starting"; sleep 1; done && \
   echo "aperture-controller is now healty!"'
   ```

## Uninstall

1. If the Aperture Controller installation is done using `docker compose`, run
   the below command to stop and remove the Aperture Controller, Prometheus and
   etcd containers:

   ```bash
   docker compose down
   ```

2. Run the below command to stop and remove the Aperture Controller container:

   ```bash
   docker rm -f aperture-controller
   ```

3. Run the below command to stop and remove the Prometheus container:

   ```bash
    docker rm -f prometheus
   ```

4. Run the below command to stop and remove the etcd container:

   ```bash
   docker rm -f etcd
   ```

5. Run the below command to remove the Docker network:

   ```bash
   docker network rm aperture
   ```
