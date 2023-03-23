# syntax=docker/dockerfile:1.4

ARG IMAGE_TAG

FROM aperture-java-all:$IMAGE_TAG

ARG SSL_CERT
COPY --link ${SSL_CERT} /ssl-cert

ENV FN_SSL_CERTIFICATE_FILE=/ssl-cert
ENV FN_INSECURE_GRPC=false

CMD ["java", "-jar", "/armeria-example.jar"]
