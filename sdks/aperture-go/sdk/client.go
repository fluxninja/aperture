package aperture

import (
	"context"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	flowcontrol "github.com/fluxninja/aperture-go/gen/proto/flowcontrol/check/v1"
)

// Client is the interface that is provided to the user upon which they can perform Check calls for their service and eventually shut down in case of error.
type Client interface {
	StartFlow(ctx context.Context, feature string, labels map[string]string) (Flow, error)
	Shutdown(ctx context.Context) error
}

type apertureClient struct {
	flowControlClient flowcontrol.FlowControlServiceClient
	tracer            trace.Tracer
	timeout           time.Duration
	exporter          *otlptrace.Exporter
}

// Options that the user can pass to Aperture in order to receive a new Client.
// FlowControlClientConn and OTLPExporterClientConn are required.
type Options struct {
	ApertureAgentGRPCClientConn *grpc.ClientConn
	CheckTimeout                time.Duration
}

// NewClient returns a new Client that can be used to perform Check calls.
// The user will pass in options which will be used to create a connection with otel and a tracerProvider retrieved from such connection.
func NewClient(ctx context.Context, opts Options) (Client, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(opts.ApertureAgentGRPCClientConn),
	)
	if err != nil {
		return nil, err
	}

	res, err := newResource()
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer(libraryName)

	fcClient := flowcontrol.NewFlowControlServiceClient(opts.ApertureAgentGRPCClientConn)

	var timeout time.Duration
	if opts.CheckTimeout == 0 {
		timeout = defaultRPCTimeout
	} else {
		timeout = opts.CheckTimeout
	}

	c := &apertureClient{
		flowControlClient: fcClient,
		tracer:            tracer,
		timeout:           timeout,
		exporter:          exporter,
	}
	return c, nil
}

// StartFlow takes a feature name and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
// Return value is a Flow.
// The call returns immediately in case connection with Aperture Agent is not established.
// The default semantics are fail-to-wire. If StartFlow fails, calling Flow.Accepted() on returned Flow returns as true.
func (c *apertureClient) StartFlow(ctx context.Context, feature string, explicitLabels map[string]string) (Flow, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	labels := make(map[string]string)

	// Inherit labels from baggage
	baggageCtx := baggage.FromContext(ctx)
	for _, member := range baggageCtx.Members() {
		value, err := url.QueryUnescape(member.Value())
		if err != nil {
			continue
		}
		labels[member.Key()] = value
	}

	// Explicit labels override baggage
	for key, value := range explicitLabels {
		labels[key] = value
	}

	req := &flowcontrol.CheckRequest{
		Feature: feature,
		Labels:  labels,
	}

	_, span := c.tracer.Start(ctx, "Aperture Check", trace.WithAttributes(
		attribute.Int64(flowStartTimestampLabel, time.Now().UnixNano()),
		attribute.String(sourceLabel, "sdk"),
	))

	f := &flow{
		span: span,
	}

	res, err := c.flowControlClient.Check(ctx, req)
	if err != nil {
		f.checkResponse = nil
	}

	span.SetAttributes(
		attribute.Int64(workloadStartTimestampLabel, time.Now().UnixNano()),
	)

	f.checkResponse = res
	return f, nil
}

// Shutdown shuts down the aperture client.
func (c *apertureClient) Shutdown(ctx context.Context) error {
	return c.exporter.Shutdown(ctx)
}

// newResource returns a resource describing the running process, containing the library name and version.
func newResource() (*resource.Resource, error) {
	resourceDefault := resource.Default()
	r, err := resource.Merge(
		resourceDefault,
		resource.NewWithAttributes(
			resourceDefault.SchemaURL(),
			semconv.ServiceNameKey.String(libraryName),
			semconv.ServiceVersionKey.String(libraryVersion),
		),
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}
