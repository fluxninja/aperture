module github.com/fluxninja/aperture-go/v2

go 1.20

require (
	buf.build/gen/go/fluxninja/aperture/grpc/go v1.3.0-20230911085626-48dd42096c62.1
	buf.build/gen/go/fluxninja/aperture/protocolbuffers/go v1.31.0-20230911085626-48dd42096c62.1
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr v1.2.2
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.15.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.15.1
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

require (
	buf.build/gen/go/envoyproxy/protoc-gen-validate/protocolbuffers/go v1.31.0-20230814203303-eac44469a7af.1 // indirect
	buf.build/gen/go/grpc-ecosystem/grpc-gateway/protocolbuffers/go v1.31.0-20230822184712-f460f71081c1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20230717213848-3f92550aa753 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230726155614-23370e0ffb3e // indirect
)
