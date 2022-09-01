server:
  addr: ":80"

dist_cache:
  bind_addr: ":3320"
  memberlist_config_bind_addr: ":3322"

otel:
  grpc_addr: ":4317"
  http_addr: ":4318"
  batch_prerollup:
    timeout: 1s
    send_batch_size: 10000
  batch_postrollup:
    timeout: 1s
    send_batch_size: 10000

log:
  pretty_console: false
  non_blocking:  true
  level: "info"
  file:  "stderr"

etcd:
  endpoints: [http://agent-etcd:2379]
  lease_ttl: 60s

prometheus:
  address: "http://aperture-prometheus-server:80"

plugins:
  disable_plugins: false

fluxninja_plugin:
  fluxninja_endpoint: "test"
  heartbeat_interval: "10s"
  client_grpc:
    insecure: true
    tls:
      insecure_skip_verify: true
      ca_file: test
  client_http:
    tls:
      insecure_skip_verify: true
      ca_file: test
