local access = require "kong.plugins.aperture-plugin.access"
local log = require "kong.plugins.aperture-plugin.log"
local json = require "cjson"

local ApertureHandler = {
  VERSION  = "0.1.0",
  PRIORITY = 10000,
}

function ApertureHandler:access(config)
  local authorized_status = access(kong.router.get_service().host, kong.router.get_service().port)
  if authorized_status ~= ngx.HTTP_OK then
    return ngx.exit(authorized_status)
  end
end

ApertureHandler.log = log

return ApertureHandler
