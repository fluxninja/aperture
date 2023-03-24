FROM docker.io/library/gradle:7.5.1-jdk18-alpine AS builder

WORKDIR /app
COPY --link . .
RUN gradle clean build


FROM openjdk:18.0.2.1-slim-buster

COPY --link --from=builder /app/build/libs/JavaDemoApp-0.0.1-SNAPSHOT.war demo.war
CMD java -jar demo.war --server.port=${SIMPLE_SERVICE_PORT}  --debug
