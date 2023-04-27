local json = require("cjson")
local socket = require("socket")
local http = require("resty.http")
local ltn12 = require("ltn12")

local apertureAgentEndpoint = os.getenv("APERTURE_AGENT_ENDPOINT")
if apertureAgentEndpoint == nil or apertureAgentEndpoint == "" then
  error("Environment variable APERTURE_AGENT_ENDPOINT must be set")
end

local apertureCheckTimeout = os.getenv("APERTURE_CHECK_TIMEOUT")
if apertureCheckTimeout == nil or apertureCheckTimeout == "" then
  apertureCheckTimeout = "250m"
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
local tp = otlp_tracer_provider_new(simple_span_processor, {
  sampler = otlp_always_on_sampler,
  resource = otlp_resource_new(otlp_attr.string("service.name", "aperture-lua"), otlp_attr.string("service.version", "v0.1.0"))
})
local otlp_tracer = tp:tracer("aperture-lua")

return function(destination_hostname, destination_port)
  local request = ngx.req
  request.read_body()
  local request_headers = ngx.req.get_headers()

  local server_addr = ""
  local server_port = ""

  if destination_hostname ~= nil and destination_hostname ~= "" then
    server_addr = socket.dns.toip(destination_hostname)
  end

  if server_addr == "" then
    server_addr = ngx.var.server_addr
  end

  if destination_port ~= nil and destination_port ~= "" then
    server_port = destination_port
  end

  if server_port == "" then
    server_port = ngx.var.server_port
  end

  local socket_type = ngx.var.server_protocol
  if socket_type ~= "UDP" then
    socket_type = "TCP"
  end

  local response_body = {}
  local request_body = {
    source = {
      address = ngx.var.remote_addr,
      protocol = socket_type,
      port = ngx.var.remote_port
    },
    destination = {
      address = server_addr,
      protocol = socket_type,
      port = server_port
    },
    request = {
      method = request.get_method(),
      headers = request_headers,
      path = ngx.var.uri,
      host = request_headers["Host"],
      scheme = ngx.var.scheme,
      size = request_headers["Content-Length"],
      protocol = ngx.var.server_protocol,
      body = request.get_body_data()
    },
    control_point = "ingress"
  }

  local request_body_json = json.encode(request_body)
  request_headers["Content-Type"] = "application/json"
  request_headers["Accept"] = "application/json"
  request_headers["content-length"] = string.len(request_body_json)
  request_headers["grpc-timeout"] = apertureCheckTimeout

  local context, span = otlp_tracer:start(otlp_trace_context_propagator:extract(otlp_context, ngx.req), "Aperture CheckHTTP", {
    kind = otlp_span_kind.server,
    attributes = {otlp_attr.string("aperture.source", "lua")}
  })
  ngx.ctx.otlp_span = span

  local httpc = http.new()
  local request_options = {
    method = "POST",
    body = request_body_json,
    headers = request_headers,
    ssl_verify = false,
  }
  local res, err = httpc:request_uri(apertureAgentEndpoint .. "/v1/flowcontrol/checkhttp", request_options)
  httpc:close()

  if not res or err ~= nil then
    ngx.log(ngx.ERR, "failed to call Aperture CheckHTTP. Error: ", err)
    return 200
  end

  local code = res.status
  if code == 200 then
    local response_json = json.decode(res.body)
    ngx.ctx.aperture_check_reponse = response_json.dynamic_metadata["aperture.check_response"]
    if response_json.ok_response ~= nil then
      for header_name, header_value in pairs(response_json.ok_response.headers) do
        request.set_header(header_name, header_value)
      end
    elseif response_json.denied_response ~= nil then
      code = response_json.denied_response.status
    else
      code = 200
    end
    span:set_attributes(otlp_attr.int("aperture.flow_start_timestamp", ngx.req.start_time() * 1000))
    span:set_attributes(otlp_attr.int("aperture.workload_start_timestamp", math.floor(socket.gettime() * 1000)))
  else
    ngx.log(ngx.ERR, "failed to send Aperture CheckHTTP request. Code: " .. code .. ", Response: " .. res.body)
    code = 200
  end
  return code
end
