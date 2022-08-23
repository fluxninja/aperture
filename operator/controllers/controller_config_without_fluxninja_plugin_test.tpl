server:
  addr: ":80"

otel:
  addr: ":80"

webhooks:
  addr: ":8086"
  tls:
    enable: true
    certs_path: "/etc/aperture/aperture-controller/certs"
    server_cert: "crt.pem"
    server_key: "key.pem"

log:
  pretty_console: false
  non_blocking:  true
  level: "info"
  file:  "stderr"

etcd:
  endpoints: [http://agent-etcd:2379]
  lease_ttl: 60s

prometheus:
  address: "http://agent-prometheus-server:80"

plugins:
  disable_plugins: false
  disabled_plugins:
    - aperture-plugin-fluxninja
