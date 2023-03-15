FROM docker.io/library/gradle:7.5.1-jdk18-alpine AS builder

COPY src /app/src
COPY build.gradle.kts /app/
WORKDIR /app
RUN gradle clean build

FROM openjdk:11-jre-slim

ENV FN_AGENT_HOST=localhost
ENV FN_AGENT_PORT=8089
ENV HOSTNAME=localhost
ENV CONCURRENCY=10
ENV LATENCY=50
ENV REJECT_RATIO=0.05

WORKDIR /JavaDemoApp

COPY build/libs/JavaDemoApp-0.0.1-SNAPSHOT.war demo.war

CMD java -jar demo.war --server.port=8080 --fn.agent.host=${FN_AGENT_HOST} --fn.agent.port=${FN_AGENT_PORT} --hostname=${HOSTNAME} --concurrency=${CONCURRENCY} --latency=${LATENCY} --reject-ratio=${REJECT_RATIO}
