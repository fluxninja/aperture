# syntax=docker/dockerfile:1.4

FROM docker.io/library/gradle:7.5.1-jdk18-alpine AS builder

WORKDIR /app
COPY --link . .
RUN gradle clean build


FROM openjdk:18-jdk-alpine

COPY --link --from=builder /app/build/libs/JavaDemoApp-0.0.1-SNAPSHOT.war demo.war

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=20s \
   CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD java -jar demo.war --server.port=${SIMPLE_SERVICE_PORT}
