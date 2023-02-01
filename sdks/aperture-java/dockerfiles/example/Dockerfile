# syntax=docker/dockerfile:1.4
# This dockerfile defines java sdk example used in playground tilt setup
FROM --platform=linux/amd64 gradle:7.5.1-jdk18-alpine AS builder

WORKDIR /

COPY --link . .

RUN gradle examples:armeria-example:shadowJar

FROM --platform=linux/amd64 openjdk:18-jdk-alpine

COPY --link --from=builder /examples/armeria-example/build/libs/*-all.jar /armeria-example.jar

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=20s \
   CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD ["java", "-jar", "/armeria-example.jar"]
