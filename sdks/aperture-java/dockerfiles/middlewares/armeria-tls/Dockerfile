# syntax=docker/dockerfile:1.4

ARG IMAGE_TAG

FROM aperture-java-all:$IMAGE_TAG

ARG ROOT_CERT
COPY --link ${ROOT_CERT} /root-cert

ENV APERTURE_ROOT_CERTIFICATE_FILE=/root-cert
ENV APERTURE_AGENT_INSECURE=false

CMD ["java", "-jar", "/armeria-example.jar"]
