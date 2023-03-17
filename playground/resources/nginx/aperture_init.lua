local json = require("cjson")
local socket = require("socket")
local http = require("socket.http")
http.TIMEOUT = 0.5
local ltn12 = require("ltn12")

local apertureAgentEndpoint = os.getenv("APERTURE_AGENT_ENDPOINT")
if apertureAgentEndpoint == nil or apertureAgentEndpoint == "" then
  error("Environment variable APERTURE_AGENT_ENDPOINT must be set")
end

local otlp_tracer_provider_new = require("opentelemetry.trace.tracer_provider").new
local otlp_simple_span_processor_new = require("opentelemetry.trace.simple_span_processor").new
local otlp_exporter_new = require("opentelemetry.trace.exporter.otlp").new
local otlp_resource_new = require("opentelemetry.resource").new
local otlp_exporter_client_new = require("opentelemetry.trace.exporter.http_client").new
local otlp_attr = require("opentelemetry.attribute")
local otlp_always_on_sampler = require("opentelemetry.trace.sampling.always_on_sampler").new()
local otlp_context = require("opentelemetry.context").new()
local otlp_trace_context_propagator = require("opentelemetry.trace.propagation.text_map.trace_context_propagator").new()
local otlp_span_kind = require("opentelemetry.trace.span_kind")

local exporter = otlp_exporter_new(otlp_exporter_client_new(apertureAgentEndpoint, 1))
local simple_span_processor = otlp_simple_span_processor_new(exporter)
local tp = otlp_tracer_provider_new(simple_span_processor,
        {sampler = otlp_always_on_sampler, resource = otlp_resource_new(otlp_attr.string("service.name", "aperture-nginx"), otlp_attr.string("service.version", "v0.1.0"))})

local otlp_tracer = tp:tracer("aperture-nginx")

function authorize_request(request)
  request.read_body()
  request_headers = request.get_headers()

  local server_addr = ""
  local server_port = ""

  if ngx.var.destination_hostname ~= nil then
    server_addr = socket.dns.toip(ngx.var.destination_hostname)
  end

  if not server_addr then
    server_addr = ngx.var.server_addr
  end

  if ngx.var.destination_port ~= nil then
    server_port = ngx.var.destination_port
  end
  if not server_port then
    server_port = ngx.var.server_port
  end

  local response_body = {}
  local request_body = {
    source = {
      address = ngx.var.remote_addr,
      protocol = ngx.var.protocol,
      port = ngx.var.remote_port
    },
    destination = {
      address = server_addr,
      protocol = ngx.var.protocol,
      port = server_port
    },
    request = {
      method = request.get_method(),
      headers = request_headers,
      path = ngx.var.uri,
      host = request_headers["Host"],
      scheme = ngx.var.scheme,
      size = request_headers["Content-Length"],
      protocol = ngx.var.http_version,
      body = request.get_body_data()
    }
  }

  request_body_json = json.encode(request_body)
  request_headers["Content-Type"] = "application/json"
  request_headers["Accept"] = "application/json"
  request_headers["control-point"] = "ingress"
  request_headers["content-length"] = string.len(request_body_json)

  local context, span = otlp_tracer:start(otlp_trace_context_propagator:extract(otlp_context, ngx.req), "Aperture CheckHTTP", {
    kind = otlp_span_kind.server,
    attributes = {otlp_attr.string("aperture.source", "lua")}
  })
  ngx.ctx.otlp_span = span

  local checkHTTPStart = ngx.now()
  local response, code, response_headers = http.request{
    url = apertureAgentEndpoint .. "/v1/flowcontrol/checkhttp",
    method = "POST",
    headers = request_headers,
    source = ltn12.source.string(request_body_json),
    ssl_verify = false,
    sink = ltn12.sink.table(response_body),
  }

  if response == nil then
    ngx.log(ngx.ERR, "failed to call Aperture CheckHTTP. Code: " .. code)
    code = 200
  end

  local response_json = json.decode(response_body[1])
  if code == 200 or code == 503 then
    ngx.ctx.aperture_check_reponse = response_json.dynamic_metadata["aperture.check_response"]
    if response_json.ok_response ~= nil then
      for header_name, header_value in pairs(response_json.ok_response.headers) do
        request.set_header(header_name, header_value)
      end
    elseif response_json.denied_response ~= nil then
      code = response_json.denied_response.status
    end
    local checkHTTPEnd = ngx.now()
    span:set_attributes(otlp_attr.string("checkhttp_duration", tostring((checkHTTPEnd - checkHTTPStart) * 1000)))
  else
    ngx.log(ngx.ERR, "failed to send Aperture CheckHTTP request. Code: " .. code .. ", Response: " .. response_body[1])
    code = 200
  end
  return code
end

function end_flow()
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
      span:set_attributes(otlp_attr.string("RESPONSE_TX_DURATION", tostring((tonumber(ngx.var.upstream_response_time) - ngx.var.upstream_header_time) * 1000)))
    end
    span:set_attributes(otlp_attr.string("aperture.check_response", ngx.ctx.aperture_check_reponse))
    span:finish()
  end
end
