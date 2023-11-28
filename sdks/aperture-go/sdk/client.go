package aperture

import (
	"context"
	"log/slog"
	"regexp"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	checkv1 "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/check/v1"
	checkhttpv1 "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture-go/v2/sdk/utils"
)

// Options that the user can pass to Aperture in order to receive a new Client.
// FlowControlClientConn and OTLPExporterClientConn are required.
type Options struct {
	Logger      *slog.Logger
	Address     string
	APIKey      string
	DialOptions []grpc.DialOption
}

// MiddlewareParams is the interface for the middleware params.
type MiddlewareParams struct {
	IgnoredPaths         []string
	IgnoredPathsCompiled []*regexp.Regexp // New field for the compiled regex patterns
	FlowParams           FlowParams
	Timeout              time.Duration
}

// FlowParams is a struct that contains parameters for StartFlow call.
type FlowParams struct {
	// Labels are the labels that get passed to Aperture Agent via flowcontrolv1.Check call.
	Labels map[string]string
	// CallOptions are the grpc call options that get passed to Aperture Agent via flowcontrolv1.Check call.
	CallOptions []grpc.CallOption
	// If RampMode is set to true, then flow must be accepted by at least 1 LoadRamp component.
	RampMode bool
	// ResultCacheKey is key to the result cache entry which needs to be fetched at flow start.
	ResultCacheKey string
	// StateCacheKeys are keys to state cache entries that need to be fetched at flow start.
	StateCacheKeys []string
}

// Client is the interface that is provided to the user upon which they can perform Check calls for their service and eventually shut down in case of error.
type Client interface {
	StartFlow(ctx context.Context, controlPoint string, flowParams FlowParams) Flow
	StartHTTPFlow(ctx context.Context, request *checkhttpv1.CheckHTTPRequest, middlewareParams MiddlewareParams) HTTPFlow
	Shutdown(ctx context.Context) error
	GetLogger() *slog.Logger
	GetGRPClientConn() *grpc.ClientConn
}

type apertureClient struct {
	grpcClientConn        *grpc.ClientConn
	flowControlClient     checkv1.FlowControlServiceClient
	flowControlHTTPClient checkhttpv1.FlowControlServiceHTTPClient
	tracer                trace.Tracer
	exporter              *otlptrace.Exporter
	log                   *slog.Logger
}

// NewClient returns a new Client that can be used to perform Check calls.
// The user will pass in options which will be used to create a connection with otel and a tracerProvider retrieved from such connection.
func NewClient(ctx context.Context, opts Options) (Client, error) {
	if opts.DialOptions == nil {
		opts.DialOptions = []grpc.DialOption{}
	}

	if opts.APIKey != "" {
		dialOptions := grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
			md := metadata.Pairs("x-api-key", opts.APIKey)
			ctx = metadata.NewOutgoingContext(ctx, md)
			return invoker(ctx, method, req, reply, cc, callOpts...)
		})
		opts.DialOptions = append(opts.DialOptions, dialOptions)
	}

	conn, err := grpc.DialContext(ctx, opts.Address, opts.DialOptions...)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
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

	var logger *slog.Logger
	if opts.Logger != nil {
		logger = opts.Logger
	} else {
		logger = slog.Default().With("name", "aperture-go-sdk")
	}

	fcClient := checkv1.NewFlowControlServiceClient(conn)
	fcHTTPClient := checkhttpv1.NewFlowControlServiceHTTPClient(conn)

	c := &apertureClient{
		grpcClientConn:        conn,
		flowControlClient:     fcClient,
		flowControlHTTPClient: fcHTTPClient,
		tracer:                tracer,
		exporter:              exporter,
		log:                   logger,
	}
	return c, nil
}

// getSpan constructs new flow tracer span.
func (c *apertureClient) getSpan(ctx context.Context) trace.Span {
	_, span := c.tracer.Start(ctx, "Aperture Check", trace.WithAttributes(
		attribute.Int64(flowStartTimestampLabel, time.Now().UnixNano()),
		attribute.String(sourceLabel, "sdk"),
	))
	return span
}

// StartFlow takes a control point name and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
// Return value is a Flow.
// The call returns immediately in case connection with Aperture Agent is not established.
// The default semantics are fail-to-wire. If StartFlow fails, calling Flow.ShouldRun() on returned Flow returns as true.
func (c *apertureClient) StartFlow(ctx context.Context, controlPoint string, flowParams FlowParams) Flow {
	labels := utils.LabelsFromCtx(ctx)

	// Explicit labels override baggage
	for key, value := range flowParams.Labels {
		labels[key] = value
	}

	req := &checkv1.CheckRequest{
		ControlPoint: controlPoint,
		Labels:       labels,
		RampMode:     flowParams.RampMode,
		CacheLookupRequest: &checkv1.CacheLookupRequest{
			ResultCacheKey: flowParams.ResultCacheKey,
			StateCacheKeys: flowParams.StateCacheKeys,
		},
	}

	span := c.getSpan(ctx)

	f := newFlow(c.flowControlClient, span, flowParams.RampMode, flowParams.ResultCacheKey, flowParams.StateCacheKeys)

	defer f.Span().SetAttributes(
		attribute.Int64(workloadStartTimestampLabel, time.Now().UnixNano()),
	)

	res, err := c.flowControlClient.Check(ctx, req, flowParams.CallOptions...)
	if err != nil {
		f.err = err
	} else {
		f.checkResponse = res
	}

	return f
}

// StartHTTPFlow takes a control point name and labels that get passed to Aperture Agent via flowcontrolhttp.CheckHTTP call.
// Return value is a HTTPFlow.
// The default semantics are fail-to-wire. If StartHTTPFlow fails, calling HTTPFlow.ShouldRun() on returned HTTPFlow returns as true.
func (c *apertureClient) StartHTTPFlow(ctx context.Context, request *checkhttpv1.CheckHTTPRequest, middlewareParams MiddlewareParams) HTTPFlow {
	span := c.getSpan(ctx)

	f := newHTTPFlow(span, middlewareParams.FlowParams)

	defer f.Span().SetAttributes(
		attribute.Int64(workloadStartTimestampLabel, time.Now().UnixNano()),
	)

	// create a timeoutCtx if middlewareParams.Timeout is set
	if middlewareParams.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, middlewareParams.Timeout)
		defer cancel()
	}

	res, err := c.flowControlHTTPClient.CheckHTTP(ctx, request)
	if err != nil {
		f.err = err
	} else {
		f.checkResponse = res
	}

	return f
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

// GetLogger returns the logger used by the aperture client.
func (c *apertureClient) GetLogger() *slog.Logger {
	return c.log
}

// GetGRPClientConn returns the grpc client connection used by the aperture client.
func (c *apertureClient) GetGRPClientConn() *grpc.ClientConn {
	return c.grpcClientConn
}
