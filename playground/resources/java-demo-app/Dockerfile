FROM docker.io/library/gradle:7.5.1-jdk18-alpine AS builder

COPY src /app/src
COPY build.gradle.kts /app/
WORKDIR /app
RUN gradle clean build

FROM openjdk:11-jre-slim

WORKDIR /java-demo-app

COPY build/libs/java-demo-app-0.0.1-SNAPSHOT.war demo.war

CMD java -jar demo.war --server.port=${SIMPLE_SERVICE_PORT}
