FROM fabiocicerchia/nginx-lua:1.23.3-debian-compat

RUN apt update && apt-get install -y build-essential git

RUN git clone https://github.com/fluxninja/opentelemetry-lua.git && cd opentelemetry-lua && luarocks make

COPY --link gateways/lua/ .

COPY --link playground/resources/nginx/nginx_config.conf /etc/nginx/nginx.conf

RUN luarocks make aperture-nginx-plugin-0.1.0-1.rockspec

ENTRYPOINT [ "nginx", "-g", "daemon off;" ]
