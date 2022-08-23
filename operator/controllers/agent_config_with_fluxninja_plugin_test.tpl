server:
  addr: ":80"

dist_cache:
  bind_addr: ":3320"
  memberlist_config_bind_addr: ":3322"

otel:
  addr: ":80"

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
