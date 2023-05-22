return function(config)
  local headers = ngx.ctx.denied_response_headers
  if headers ~= nil then
    for header_name, header_value in pairs(headers) do
      ngx.header[header_name] = header_value
    end
  end
end
