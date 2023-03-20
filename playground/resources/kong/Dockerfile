FROM kong:3.1.1-ubuntu

WORKDIR /usr/kong/aperture

COPY --link gateways/lua/ .

COPY --link playground/resources/kong/kong.yaml .

COPY --link playground/resources/kong/kong.conf .

USER root

RUN apt update && apt-get install -y build-essential git

RUN git clone https://github.com/fluxninja/opentelemetry-lua.git && cd opentelemetry-lua && luarocks make

RUN luarocks make aperture-kong-plugin-0.1.0-1.rockspec

USER kong

ENV KONG_DATABASE=off

ENV KONG_DECLARATIVE_CONFIG=kong.yaml

CMD [ "kong", "start", "-c", "kong.conf"]
