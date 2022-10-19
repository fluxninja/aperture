client:
  proxy:
    http: ""
    https: ""
etcd:
  endpoints:
  - http://agent-etcd:2379
  lease_ttl: 60s
  password: ""
  tls:
    ca_file: ""
    cert_file: ""
    insecure_skip_verify: false
    key_file: ""
    key_log_file: ""
  username: ""
fluxninja_plugin:
  api_key: ""
  client:
    grpc:
      backoff:
        base_delay: 1s
        jitter: 0.2
        max_delay: 120s
        multiplier: 1.6
      insecure: false
      min_connection_timeout: 20s
      tls:
        ca_file: ""
        cert_file: ""
        insecure_skip_verify: false
        key_file: ""
        key_log_file: ""
      use_proxy: false
    http:
      disable_compression: false
      disable_keep_alives: false
      expect_continue_timeout: 1s
      idle_connection_timeout: 90s
      key_log_file: ""
      max_conns_per_host: 0
      max_idle_connections: 100
      max_idle_connections_per_host: 5
      max_response_header_bytes: 0
      network_keep_alive: 30s
      network_timeout: 30s
      read_buffer_size: 0
      response_header_timeout: 0s
      timeout: 60s
      tls:
        ca_file: ""
        cert_file: ""
        insecure_skip_verify: false
        key_file: ""
        key_log_file: ""
      tls_handshake_timeout: 10s
      use_proxy: false
      write_buffer_size: 0
  fluxninja_endpoint: ""
  heartbeat_interval: 5s
liveness:
  scheduler:
    max_concurrent_jobs: 0
  service:
    execution_period: 10s
    execution_timeout: 5s
    initial_delay: 0s
    initially_healthy: false
log:
  level: info
  non_blocking: true
  pretty_console: false
  writers:
  - compress: false
    file: stderr
    max_age: 7
    max_backups: 3
    max_size: 50
metrics:
  enable_go_metrics: false
  enable_process_collector: false
  pedantic: false
otel:
  batch_postrollup:
    send_batch_size: 10000
    send_batch_max_size: 20000
    timeout: 1s
  batch_prerollup:
    send_batch_size: 100
    send_batch_max_size: 200
    timeout: 1s
plugins:
  disable_plugins: false
  disabled_plugins:
  - aperture-plugin-fluxninja
  plugins_path: default
policies:
  promql_jobs_scheduler:
    max_concurrent_jobs: 0
profilers:
  cpu_profiler: false
  profiles_path: default
  register_http_routes: true
prometheus:
  address: http://aperture-prometheus-server:80
readiness:
  scheduler:
    max_concurrent_jobs: 0
  service:
    execution_period: 10s
    execution_timeout: 5s
    initial_delay: 0s
    initially_healthy: false
sentry_plugin:
  attach_stack_trace: true
  debug: true
  disabled: false
  dsn: https://6223f112b0ac4344aa67e94d1631eb85@o574197.ingest.sentry.io/6605877
  environment: production
  sample_rate: 1
  traces_sample_rate: 0.2
server:
  addr: :80
  grpc:
    connection_timeout: 120s
    enable_reflection: false
    latency_buckets_ms:
    - 10
    - 25
    - 100
    - 250
    - 1000
  grpc_gateway:
    grpc_server_address: 0.0.0.0:1
  http:
    disable_http_keep_alives: false
    idle_timeout: 30s
    latency_buckets_ms:
    - 10
    - 25
    - 100
    - 250
    - 1000
    max_header_bytes: 1048576
    read_header_timeout: 10s
    read_timeout: 10s
    write_timeout: 45s
  keep_alive: 180s
  network: tcp
  tls:
    allowed_cn: ""
    cert_file: /etc/aperture/aperture-controller/certs/crt.pem
    client_ca_file: ""
    enabled: true
    key_file: /etc/aperture/aperture-controller/certs/key.pem
watchdog:
  cgroup:
    adaptive_policy:
      enabled: false
      factor: 0.5
    watermarks_policy:
      enabled: false
      watermarks:
      - 0.5
      - 0.75
      - 0.8
      - 0.85
      - 0.9
      - 0.95
      - 0.99
  heap:
    adaptive_policy:
      enabled: false
      factor: 0.5
    limit: 268435456
    min_gogc: 25
    watermarks_policy:
      enabled: false
      watermarks:
      - 0.5
      - 0.75
      - 0.8
      - 0.85
      - 0.9
      - 0.95
      - 0.99
  job:
    execution_period: 10s
    execution_timeout: 5s
    initial_delay: 0s
    initially_healthy: false
  system:
    adaptive_policy:
      enabled: false
      factor: 0.5
    watermarks_policy:
      enabled: false
      watermarks:
      - 0.5
      - 0.75
      - 0.8
      - 0.85
      - 0.9
      - 0.95
      - 0.99
