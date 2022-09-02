etcd:
  endpoints: [http://controller-etcd:2379]
  lease_ttl: 60s
fluxninja_plugin:
  client_grpc:
    insecure: true
    tls:
      ca_file: test
      insecure_skip_verify: true
  client_http:
    tls:
      ca_file: test
      insecure_skip_verify: true
  fluxninja_endpoint: "test"
  heartbeat_interval: "10s"
log:
  file: "stderr"
  level: "info"
  non_blocking: true
  pretty_console: false
otel:
  batch_postrollup:
    send_batch_size: 10000
    timeout: 1s
  batch_prerollup:
    send_batch_size: 10000
    timeout: 1s
  grpc_addr: ":4317"
  http_addr: ":4318"
plugins:
  disable_plugins: false
prometheus:
  address: "http://aperture-prometheus-server:80"
server:
  addr: ":80"
webhooks:
  addr: ":8086"
  tls:
    certs_path: "/etc/aperture/aperture-controller/certs"
    enable: true
    server_cert: "crt.pem"
    server_key: "key.pem"
