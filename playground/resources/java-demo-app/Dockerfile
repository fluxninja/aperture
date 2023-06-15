# syntax=docker/dockerfile:1.4

FROM docker.io/library/gradle:7.5.1-jdk18-alpine AS builder

WORKDIR /app
COPY --link . .
RUN gradle clean build


FROM openjdk:22-jdk-slim-bullseye

RUN apt-get update && apt-get install -y wget
COPY --link --from=builder /app/build/libs/JavaDemoApp-0.0.1-SNAPSHOT.war demo.war
COPY --link ./config.yaml config.yaml
RUN wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.18.0/jmx_prometheus_javaagent-0.18.0.jar

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=20s \
   CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD java -javaagent:./jmx_prometheus_javaagent-0.18.0.jar=8087:config.yaml -jar demo.war --server.port=${SIMPLE_SERVICE_PORT}
