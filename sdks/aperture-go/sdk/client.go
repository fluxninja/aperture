package aperture

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"

	flowcontrol "github.com/fluxninja/aperture-go/gen/proto/flowcontrol/check/v1"
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Client is the interface that is provided to the user upon which they can perform Check calls for their service and eventually shut down in case of error.
type Client interface {
	StartFlow(ctx context.Context, controlPoint string, labels map[string]string) (Flow, error)
	Shutdown(ctx context.Context) error
	HTTPMiddleware(controlPoint string, labels map[string]string, timeout time.Duration) func(http.Handler) http.Handler
	GRPCUnaryInterceptor(controlPoint string, labels map[string]string, timeout time.Duration) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type apertureClient struct {
	flowControlClient flowcontrol.FlowControlServiceClient
	tracer            trace.Tracer
	timeout           time.Duration
	exporter          *otlptrace.Exporter
	log               logr.Logger
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

	fcClient := flowcontrol.NewFlowControlServiceClient(opts.ApertureAgentGRPCClientConn)

	var timeout time.Duration
	if opts.CheckTimeout == 0 {
		timeout = defaultRPCTimeout
	} else {
		timeout = opts.CheckTimeout
	}

	var logger logr.Logger
	if opts.Logger != nil {
		logger = *opts.Logger
	} else {
		logger = stdr.New(log.Default()).WithName("aperture-go-sdk")
	}

	c := &apertureClient{
		flowControlClient: fcClient,
		tracer:            tracer,
		timeout:           timeout,
		exporter:          exporter,
		log:               logger,
	}
	return c, nil
}

// StartFlow takes a control point name and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
// Return value is a Flow.
// The call returns immediately in case connection with Aperture Agent is not established.
// The default semantics are fail-to-wire. If StartFlow fails, calling Flow.Accepted() on returned Flow returns as true.
func (c *apertureClient) StartFlow(ctx context.Context, controlPoint string, explicitLabels map[string]string) (Flow, error) {
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
		ControlPoint: controlPoint,
		Labels:       labels,
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
	} else {
		f.checkResponse = res
	}

	span.SetAttributes(
		attribute.Int64(workloadStartTimestampLabel, time.Now().UnixNano()),
	)

	return f, nil
}

// HTTPMiddleware takes a control point name, labels and timeout and creates a Middleware which can be used with HTTP server.
func (client *apertureClient) HTTPMiddleware(controlPoint string, labels map[string]string, timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newLabels := make(map[string]string, len(labels))
			maps.Copy(newLabels, labels)

			for key, value := range r.Header {
				newLabels[key] = strings.Join(value, ",")
			}

			flow := client.executeFlow(controlPoint, newLabels, timeout)

			if flow.Accepted() {
				// Simulate work being done
				next.ServeHTTP(w, r)
				// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
				// The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
				// The second argument is error message for further diagnosis.
				err := flow.End(OK)
				if err != nil {
					client.log.Info("Aperture flow control end got error.", "error", err)
				}
			} else {
				// TODO use HTTP Check and pull proper status
				w.WriteHeader(http.StatusServiceUnavailable)
				_, err := fmt.Fprint(w, flow.CheckResponse().GetRejectReason().String())
				if err != nil {
					client.log.Info("Aperture flow control end got error.", "error", err)
				}
				err = flow.End(OK)
				if err != nil {
					client.log.Info("Aperture flow control end got error.", "error", err)
				}
			}
		})
	}
}

// GRPCUnaryInterceptor takes a control point name, labels and timeout and creates a UnaryInterceptor which can be used with gRPC server.
func (client *apertureClient) GRPCUnaryInterceptor(controlPoint string, labels map[string]string, timeout time.Duration) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newLabels := make(map[string]string, len(labels))
		maps.Copy(newLabels, labels)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for key, value := range md {
				newLabels[key] = strings.Join(value, ",")
			}
		}

		flow := client.executeFlow(controlPoint, newLabels, timeout)

		if flow.Accepted() {
			// Simulate work being done
			resp, err := handler(ctx, req)
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
			// The second argument is error message for further diagnosis.
			flowErr := flow.End(OK)
			if flowErr != nil {
				client.log.Info("Aperture flow control end got error.", "error", err)
			}
			return resp, err
		} else {
			err := flow.End(OK)
			if err != nil {
				client.log.Info("Aperture flow control end got error.", "error", err)
			}
			// TODO use HTTP Check and pull proper status
			return nil, status.Error(
				codes.Unavailable,
				fmt.Sprintf("Aperture Rejected the Request: %v", flow.CheckResponse().GetRejectReason().String()),
			)
		}
	}
}

func (client *apertureClient) executeFlow(controlPoint string, labels map[string]string, timeout time.Duration) Flow {
	context, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
	flow, err := client.StartFlow(context, controlPoint, labels)
	if err != nil {
		client.log.Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Accepted()", flow.Accepted())
	}

	return flow
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
