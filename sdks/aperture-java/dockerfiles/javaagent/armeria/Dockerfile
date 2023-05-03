# syntax=docker/dockerfile:1.4

ARG IMAGE_TAG
FROM aperture-java-all:$IMAGE_TAG

ENV APERTURE_JAVAAGENT_INSECURE_GRPC=true
ENV APERTURE_CONTROL_POINT_NAME=awesomeFeature

CMD ["java", "-javaagent:/javaagent.jar", "-jar", "/armeria.jar"]
