rockspec_format = "3.0"
package = "aperture-nginx-plugin"
version = "0.1.0-1"
source = {
   url = "git+https://github.com/fluxninja/aperture.git",
   tag = "v0.1.0",
}
description = {
   summary = "Integrate the FluxNinja Aperture with Kong API Gateway for Load management",
   homepage = "https://github.com/fluxninja/aperture/tree/master/gateways/lua",
   issues_url = "https://github.com/fluxninja/aperture/issues",
}
dependencies = {
   "lua-resty-http = 0.17.1-0",
   "lua-cjson = 2.1.0.10",
}
build = {
   type = "builtin",
   modules = {
      ["aperture-plugin.access"] = "aperture-plugin/access.lua",
      ["aperture-plugin.log"] = "aperture-plugin/log.lua",
      ["aperture-plugin.headers"] = "aperture-plugin/headers.lua",
   },
}
