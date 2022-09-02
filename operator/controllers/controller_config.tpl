etcd:
  endpoints: {{ .Etcd.Endpoints }}
  lease_ttl: {{ .Etcd.LeaseTTL }}
{{- if .FluxNinjaPlugin.Enabled }}
fluxninja_plugin:
  client_grpc:
    insecure: {{ .FluxNinjaPlugin.TLS.Insecure }}
    tls:
      ca_file: {{ .FluxNinjaPlugin.TLS.CAFile }}
      insecure_skip_verify: {{ .FluxNinjaPlugin.TLS.InsecureSkipVerify }}
  client_http:
    tls:
      ca_file: {{ .FluxNinjaPlugin.TLS.CAFile }}
      insecure_skip_verify: {{ .FluxNinjaPlugin.TLS.InsecureSkipVerify }}
  fluxninja_endpoint: "{{ .FluxNinjaPlugin.Endpoint }}"
  heartbeat_interval: "{{ .FluxNinjaPlugin.HeartbeatsInterval }}"
{{- end }}
log:
  file: "{{ .Log.File }}"
  level: "{{ .Log.Level }}"
  non_blocking: {{ .Log.NonBlocking }}
  pretty_console: {{ .Log.PrettyConsole }}
otel:
  batch_postrollup:
    send_batch_size: {{ .OtelConfig.BatchPostrollup.SendBatchSize }}
    timeout: {{ .OtelConfig.BatchPostrollup.Timeout }}
  batch_prerollup:
    send_batch_size: {{ .OtelConfig.BatchPrerollup.SendBatchSize }}
    timeout: {{ .OtelConfig.BatchPrerollup.Timeout }}
  grpc_addr: "{{ .OtelConfig.GRPCAddr }}"
  http_addr: "{{ .OtelConfig.HTTPAddr }}"
plugins:
  disable_plugins: false
  {{- if not .FluxNinjaPlugin.Enabled }}
  disabled_plugins:
  - aperture-plugin-fluxninja
  {{- end }}
prometheus:
  address: "{{ .PrometheusAddress }}"
server:
  addr: ":{{ .ServerPort }}"
webhooks:
  addr: ":8086"
  tls:
    certs_path: "{{ .CertPath }}"
    enable: true
    server_cert: "{{ .CertName }}"
    server_key: "{{ .CertKey }}"
