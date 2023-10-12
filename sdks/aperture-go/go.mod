module github.com/fluxninja/aperture-go/v2

go 1.20

require (
	buf.build/gen/go/fluxninja/aperture/grpc/go v1.3.0-20230922113103-98fe61af5f4b.1
	buf.build/gen/go/fluxninja/aperture/protocolbuffers/go v1.31.0-20230922113103-98fe61af5f4b.1
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr v1.2.2
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/otel v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.17.0
	go.opentelemetry.io/otel/sdk v1.17.0
	go.opentelemetry.io/otel/trace v1.17.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230911183012-2d3300fd4832
	google.golang.org/grpc v1.58.0
	google.golang.org/protobuf v1.31.0
)

require (
	buf.build/gen/go/envoyproxy/protoc-gen-validate/protocolbuffers/go v1.31.0-20230814203303-eac44469a7af.1 // indirect
	buf.build/gen/go/grpc-ecosystem/grpc-gateway/protocolbuffers/go v1.31.0-20230822184712-f460f71081c1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20230911183012-2d3300fd4832 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230911183012-2d3300fd4832 // indirect
)
