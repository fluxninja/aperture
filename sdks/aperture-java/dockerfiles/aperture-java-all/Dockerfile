# syntax=docker/dockerfile:1.4

FROM --platform=linux/amd64 gradle:7.5.1-jdk18-alpine AS builder

WORKDIR /

COPY --link . .

RUN ./gradlew clean \
    javaagent:agent:shadowJar \
    javaagent:test-services:armeria-test-service:shadowJar \
    javaagent:test-services:netty-test-service:shadowJar \
    examples:armeria-example:shadowJar \
    examples:netty-example:shadowJar \
    examples:standalone-example:shadowJar \
    examples:spring-example:assemble
# TODO: also add tomcat here

FROM --platform=linux/amd64 openjdk:18-jdk-alpine

COPY --link --from=builder /javaagent/agent/build/libs/agent*.jar /javaagent.jar
COPY --link --from=builder /javaagent/test-services/armeria-test-service/build/libs/*-all.jar /armeria.jar
COPY --link --from=builder /javaagent/test-services/netty-test-service/build/libs/*-all.jar /netty.jar

COPY --link --from=builder /examples/armeria-example/build/libs/*-all.jar /armeria-example.jar
COPY --link --from=builder /examples/netty-example/build/libs/*-all.jar /netty-example.jar
COPY --link --from=builder /examples/standalone-example/build/libs/*-all.jar /standalone-example.jar
COPY --link --from=builder /examples/spring-example/build/libs/*-SNAPSHOT.jar /spring-example.jar

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=20s \
   CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD ["echo", "This is a base testing image and should not be used by itself"]
