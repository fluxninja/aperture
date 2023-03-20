local otlp_attr = require("opentelemetry.attribute")

return function(config)
  local span = ngx.ctx.otlp_span
  if span ~= nil and ngx.ctx.aperture_check_reponse ~= nil then
    span:set_attributes(otlp_attr.string("http.status_code", ngx.var.status))

    if ngx.var.request_length ~= nil then
      span:set_attributes(otlp_attr.string("BYTES_RECEIVED", tostring(ngx.var.request_length)))
    end
    if ngx.var.body_bytes_sent ~= nil then
      span:set_attributes(otlp_attr.string("BYTES_SENT", tostring(ngx.var.body_bytes_sent)))
    end
    if ngx.var.request_time ~= nil then
      span:set_attributes(otlp_attr.string("DURATION", tostring(ngx.var.request_time * 1000)))
    end
    if ngx.var.request_time ~= nil and ngx.var.upstream_response_time ~= nil then
      span:set_attributes(otlp_attr.string("REQUEST_DURATION", tostring((ngx.var.request_time - tonumber(ngx.var.upstream_response_time)) * 1000)))
    end
    if ngx.var.upstream_connect_time ~= nil then
      span:set_attributes(otlp_attr.string("REQUEST_TX_DURATION", tostring(tonumber(ngx.var.upstream_response_time) * 1000)))
    end
    if ngx.var.upstream_response_time ~= nil then
      span:set_attributes(otlp_attr.string("RESPONSE_DURATION", tostring(tonumber(ngx.var.upstream_response_time) * 1000)))
    end
    if ngx.var.upstream_response_time ~= nil and ngx.var.upstream_header_time ~= nil then
      span:set_attributes(otlp_attr.string("RESPONSE_TX_DURATION", tostring((tonumber(ngx.var.upstream_response_time) - tonumber(ngx.var.upstream_header_time)) * 1000)))
    end
    span:set_attributes(otlp_attr.string("aperture.check_response", ngx.ctx.aperture_check_reponse))
    span:finish()
  end
end
