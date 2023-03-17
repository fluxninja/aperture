FROM fabiocicerchia/nginx-lua:1.23.3-debian-compat

RUN apt update && apt-get install -y wget build-essential git

RUN luarocks install lua-cjson

RUN git clone https://github.com/fluxninja/opentelemetry-lua.git && cd opentelemetry-lua && luarocks make

COPY aperture_init.lua /aperture_init.lua

COPY nginx_config.conf /etc/nginx/nginx.conf
