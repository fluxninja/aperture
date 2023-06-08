local otlp_attr = require("opentelemetry.attribute")
local socket = require("socket")

return function(config)
    local span = ngx.ctx.otlp_span
    if span ~= nil and ngx.ctx.aperture_check_reponse ~= nil then
        span:set_attributes(otlp_attr.string("http.status_code", ngx.var.status))

        local request_time = tonumber(ngx.var.request_time)
        local upstream_response_time = tonumber(ngx.var.upstream_response_time)
        local upstream_header_time = tonumber(ngx.var.upstream_header_time)
        local upstream_connect_time = tonumber(ngx.var.upstream_connect_time)
        if ngx.var.request_length ~= nil then
            span:set_attributes(otlp_attr.int("BYTES_RECEIVED", ngx.var.request_length))
        end
        if ngx.var.body_bytes_sent ~= nil then
            span:set_attributes(otlp_attr.int("BYTES_SENT", ngx.var.body_bytes_sent))
        end
        if request_time ~= nil then
            span:set_attributes(otlp_attr.int("DURATION", request_time * 1000))
            if upstream_response_time ~= nil then
                span:set_attributes(otlp_attr.int("REQUEST_DURATION",
                    math.floor((request_time - upstream_response_time) * 1000)))
                span:set_attributes(otlp_attr.int("RESPONSE_DURATION", upstream_response_time * 1000))
            end
        end
        if upstream_connect_time ~= nil then
            span:set_attributes(otlp_attr.int("REQUEST_TX_DURATION", upstream_connect_time * 1000))
        end
        if upstream_response_time ~= nil and upstream_header_time ~= nil then
            span:set_attributes(otlp_attr.int("RESPONSE_TX_DURATION",
                math.floor((upstream_response_time - upstream_header_time) * 1000)))
        end
        span:set_attributes(otlp_attr.int("aperture.flow_end_timestamp", math.floor(socket.gettime() * 1000)))
        span:set_attributes(otlp_attr.string("aperture.check_response", ngx.ctx.aperture_check_reponse))
        span:finish()
    end
end
