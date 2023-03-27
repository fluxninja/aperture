rockspec_format = "3.0"
package = "aperture-kong-plugin"
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
build = {
   type = "builtin",
   modules = {
      ["kong.plugins.aperture-plugin.access"] = "aperture-plugin/access.lua",
      ["kong.plugins.aperture-plugin.handler"] = "kong/handler.lua",
      ["kong.plugins.aperture-plugin.schema"] = "kong/schema.lua",
      ["kong.plugins.aperture-plugin.log"] = "aperture-plugin/log.lua",
   },
}
