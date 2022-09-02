dist_cache:
  bind_addr: :3320
  memberlist_config_bind_addr: :8000
etcd:
  endpoints:
  - http://agent-etcd:2379
  lease_ttl: 60s
log:
  file: stderr
  level: info
  non_blocking: true
  pretty_console: false
otel:
  batch_postrollup:
    send_batch_size: 10000
    timeout: 1s
  batch_prerollup:
    send_batch_size: 500
    timeout: 1s
  grpc_addr: :4317
  http_addr: :4318
plugins:
  disable_plugins: false
  disabled_plugins:
  - aperture-plugin-fluxninja
prometheus:
  address: http://aperture-prometheus-server:80
server:
  addr: :80
