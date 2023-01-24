# syntax=docker/dockerfile:1.4

ARG IMAGE_TAG
FROM aperture-java-all:$IMAGE_TAG

CMD ["java", "-javaagent:/javaagent.jar", "-jar", "/armeria.jar"]
