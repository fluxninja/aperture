server:
  addr: ":{{ .ServerPort }}"

dist_cache:
  bind_addr: ":{{ .DistributedCachePort }}"
  memberlist_config_bind_addr: ":{{ .MemberListPort }}"

otel:
  grpc_addr: "{{ .OtelConfig.GRPCAddr }}"
  http_addr: "{{ .OtelConfig.HTTPAddr }}"
  batch_prerollup:
    timeout: {{ .OtelConfig.BatchPrerollup.Timeout }}
    send_batch_size: {{ .OtelConfig.BatchPrerollup.SendBatchSize }}
  batch_postrollup:
    timeout: {{ .OtelConfig.BatchPostrollup.Timeout }}
    send_batch_size: {{ .OtelConfig.BatchPostrollup.SendBatchSize }}

log:
  pretty_console: {{ .Log.PrettyConsole }}
  non_blocking:  {{ .Log.NonBlocking }}
  level: "{{ .Log.Level }}"
  file:  "{{ .Log.File }}"

etcd:
  endpoints: {{ .Etcd.Endpoints }}
  lease_ttl: {{ .Etcd.LeaseTTL }}

prometheus:
  address: "{{ .PrometheusAddress }}"

plugins:
  disable_plugins: false
  {{- if not .FluxNinjaPlugin.Enabled }}
  disabled_plugins:
    - aperture-plugin-fluxninja
  {{- end }}

{{- if .FluxNinjaPlugin.Enabled }}

fluxninja_plugin:
  fluxninja_endpoint: "{{ .FluxNinjaPlugin.Endpoint }}"
  heartbeat_interval: "{{ .FluxNinjaPlugin.HeartbeatsInterval }}"
  client_grpc:
    insecure: {{ .FluxNinjaPlugin.TLS.Insecure }}
    tls:
      insecure_skip_verify: {{ .FluxNinjaPlugin.TLS.InsecureSkipVerify }}
      ca_file: {{ .FluxNinjaPlugin.TLS.CAFile }}
  client_http:
    tls:
      insecure_skip_verify: {{ .FluxNinjaPlugin.TLS.InsecureSkipVerify }}
      ca_file: {{ .FluxNinjaPlugin.TLS.CAFile }}
{{- end }}
