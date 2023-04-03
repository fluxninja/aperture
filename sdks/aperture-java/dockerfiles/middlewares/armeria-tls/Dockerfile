# syntax=docker/dockerfile:1.4

ARG IMAGE_TAG

FROM aperture-java-all:$IMAGE_TAG

ARG CA_CERT
COPY --link ${CA_CERT} /ca-cert

ENV FN_CA_CERTIFICATE_FILE=/ca-cert
ENV FN_INSECURE_GRPC=false

CMD ["java", "-jar", "/armeria-example.jar"]
