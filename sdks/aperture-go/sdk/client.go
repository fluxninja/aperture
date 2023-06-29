package aperture

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	flowcontrol "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/check/v1"
	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
)

// Client is the interface that is provided to the user upon which they can perform Check calls for their service and eventually shut down in case of error.
type Client interface {
	StartFlow(ctx context.Context, controlPoint string, labels map[string]string) (Flow, error)
	StartHTTPFlow(ctx context.Context, request *flowcontrolhttp.CheckHTTPRequest) (HTTPFlow, error)
	Shutdown(ctx context.Context) error
	HTTPMiddleware(controlPoint string, labels map[string]string) mux.MiddlewareFunc
	GRPCUnaryInterceptor(controlPoint string, labels map[string]string) grpc.UnaryServerInterceptor
	GetLogger() logr.Logger
}

type apertureClient struct {
	flowControlClient     flowcontrol.FlowControlServiceClient
	flowControlHTTPClient flowcontrolhttp.FlowControlServiceHTTPClient
	tracer                trace.Tracer
	timeout               time.Duration
	exporter              *otlptrace.Exporter
	log                   logr.Logger
}

// Options that the user can pass to Aperture in order to receive a new Client.
// FlowControlClientConn and OTLPExporterClientConn are required.
type Options struct {
	ApertureAgentGRPCClientConn *grpc.ClientConn
	CheckTimeout                time.Duration
	Logger                      *logr.Logger
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

	var logger logr.Logger
	if opts.Logger != nil {
		logger = *opts.Logger
	} else {
		logger = stdr.New(log.Default()).WithName("aperture-go-sdk")
	}

	fcClient := flowcontrol.NewFlowControlServiceClient(opts.ApertureAgentGRPCClientConn)
	fcHTTPClient := flowcontrolhttp.NewFlowControlServiceHTTPClient(opts.ApertureAgentGRPCClientConn)

	c := &apertureClient{
		flowControlClient:     fcClient,
		flowControlHTTPClient: fcHTTPClient,
		tracer:                tracer,
		timeout:               opts.CheckTimeout,
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

// LabelsFromCtx extracts baggage labels from context.
func LabelsFromCtx(ctx context.Context) map[string]string {
	labels := make(map[string]string)
	baggageCtx := baggage.FromContext(ctx)
	for _, member := range baggageCtx.Members() {
		value, err := url.QueryUnescape(member.Value())
		if err != nil {
			continue
		}
		labels[member.Key()] = value
	}
	return labels
}

// StartFlow takes a control point name and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
// Return value is a Flow.
// The call returns immediately in case connection with Aperture Agent is not established.
// The default semantics are fail-to-wire. If StartFlow fails, calling Flow.ShouldRun() on returned Flow returns as true.
func (c *apertureClient) StartFlow(ctx context.Context, controlPoint string, explicitLabels map[string]string) (Flow, error) {
	// if c.timeout is not 0, then create a new context with timeout
	if c.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	labels := LabelsFromCtx(ctx)

	// Explicit labels override baggage
	for key, value := range explicitLabels {
		labels[key] = value
	}

	req := &flowcontrol.CheckRequest{
		ControlPoint: controlPoint,
		Labels:       labels,
	}

	span := c.getSpan(ctx)

	f := newFlow(span)

	res, err := c.flowControlClient.Check(ctx, req)
	if err != nil {
		f.checkResponse = nil
	} else {
		f.checkResponse = res
	}

	span.SetAttributes(
		attribute.Int64(workloadStartTimestampLabel, time.Now().UnixNano()),
	)

	return f, nil
}

// StartHTTPFlow takes a control point name and labels that get passed to Aperture Agent via flowcontrolhttp.CheckHTTP call.
// Return value is a HTTPFlow.
// The call returns immediately in case connection with Aperture Agent is not established.
// The default semantics are fail-to-wire. If StartHTTPFlow fails, calling HTTPFlow.ShouldRun() on returned HTTPFlow returns as true.
func (c *apertureClient) StartHTTPFlow(ctx context.Context, request *flowcontrolhttp.CheckHTTPRequest) (HTTPFlow, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	span := c.getSpan(ctx)

	f := newHTTPFlow(span)

	res, err := c.flowControlHTTPClient.CheckHTTP(ctx, request)
	if err != nil {
		f.checkResponse = nil
	} else {
		f.checkResponse = res
	}

	span.SetAttributes(
		attribute.Int64(workloadStartTimestampLabel, time.Now().UnixNano()),
	)

	return f, nil
}

// HTTPMiddleware takes a control point name, labels and timeout and creates a Middleware which can be used with HTTP server.
// Deprecated: 2.3.0 Use middlewares.HTTPMiddleware instead.
func (c *apertureClient) HTTPMiddleware(controlPoint string, labels map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newLabels := make(map[string]string, len(labels))
			maps.Copy(newLabels, labels)

			for key, value := range r.Header {
				newLabels[key] = strings.Join(value, ",")
			}

			flow, err := c.StartFlow(r.Context(), controlPoint, newLabels)
			if err != nil {
				c.log.Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.ShouldRun()", flow.ShouldRun())
			}

			if flow.ShouldRun() {
				// Simulate work being done
				next.ServeHTTP(w, r)
			} else {
				// TODO use HTTP Check and pull proper status
				w.WriteHeader(http.StatusServiceUnavailable)
				_, perr := fmt.Fprint(w, flow.CheckResponse().GetRejectReason().String())
				if perr != nil {
					c.log.Info("Aperture flow control end got error.", "error", perr)
				}
			}
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// SetStatus() method of Flow object can be used to capture whether the Flow was successful or resulted in an error.
			// If not set, status defaults to OK.
			err = flow.End()
			if err != nil {
				c.log.Info("Aperture flow control end got error.", "error", err)
			}
		})
	}
}

// GRPCUnaryInterceptor takes a control point name, labels and timeout and creates a UnaryInterceptor which can be used with gRPC server.
// Deprecated: 2.3.0 Use middlewares.GRPCUnaryInterceptor instead.
func (c *apertureClient) GRPCUnaryInterceptor(controlPoint string, labels map[string]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newLabels := make(map[string]string, len(labels))
		maps.Copy(newLabels, labels)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for key, value := range md {
				newLabels[key] = strings.Join(value, ",")
			}
		}

		flow, err := c.StartFlow(ctx, controlPoint, labels)
		if err != nil {
			c.log.Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.ShouldRun()", flow.ShouldRun())
		}

		if flow.ShouldRun() {
			// Simulate work being done
			resp, err := handler(ctx, req)
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// SetStatus() method of Flow object can be used to capture whether the Flow was successful or resulted in an error.
			// If not set, status defaults to OK.
			flowErr := flow.End()
			if flowErr != nil {
				c.log.Info("Aperture flow control end got error.", "error", err)
			}
			return resp, err
		} else {
			err := flow.End()
			if err != nil {
				c.log.Info("Aperture flow control end got error.", "error", err)
			}
			// TODO use HTTP Check and pull proper status
			return nil, status.Error(
				codes.Unavailable,
				fmt.Sprintf("Aperture rejected the request: %v", flow.CheckResponse().GetRejectReason().String()),
			)
		}
	}
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
func (c *apertureClient) GetLogger() logr.Logger {
	return c.log
}
