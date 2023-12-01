package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	aperture "github.com/fluxninja/aperture-go/v2/sdk"
	"github.com/fluxninja/aperture-go/v2/sdk/utils"
	checkhttpv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
)

// HTTPMiddleware is the interface for the HTTP middleware.
type HTTPMiddleware interface {
	Handle(http.Handler) http.Handler
}

type httpMiddleware struct {
	client           aperture.Client
	controlPoint     string
	middlewareParams aperture.MiddlewareParams
}

// NewHTTPMiddleware creates a new HTTPMiddleware struct.
func NewHTTPMiddleware(client aperture.Client, controlPoint string, middlewareParams aperture.MiddlewareParams) (HTTPMiddleware, error) {
	// Precompile the regex patterns for ignored paths
	if middlewareParams.IgnoredPaths != nil {
		compiledIgnoredPaths := make([]*regexp.Regexp, len(middlewareParams.IgnoredPaths))
		for i, pattern := range middlewareParams.IgnoredPaths {
			compiledPattern, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			} else {
				compiledIgnoredPaths[i] = compiledPattern
			}
		}
		middlewareParams.IgnoredPathsCompiled = compiledIgnoredPaths
	}

	return &httpMiddleware{
		client:           client,
		controlPoint:     controlPoint,
		middlewareParams: middlewareParams,
	}, nil
}

// Handle takes a http.Handler and returns a new http.Handler with the middleware applied.
func (m *httpMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the path is ignored, skip the middleware
		if m.middlewareParams.IgnoredPathsCompiled != nil {
			for _, compiledPattern := range m.middlewareParams.IgnoredPathsCompiled {
				if compiledPattern.MatchString(r.URL.Path) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		req := prepareCheckHTTPRequestForHTTP(r, m.client.GetLogger(), m.controlPoint, m.middlewareParams.FlowParams)

		flow := m.client.StartHTTPFlow(r.Context(), req, m.middlewareParams)
		if flow.Error() != nil {
			m.client.GetLogger().Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Error()", flow.Error().Error(), "flow.ShouldRun()", flow.ShouldRun())
		}

		defer func() {
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// SetStatus() method of Flow object can be used to capture whether the Flow was successful or resulted in an error.
			// If not set, status defaults to OK.
			err := flow.End()
			if err != nil {
				m.client.GetLogger().Info("Aperture flow control end got error.", "error", err)
			}
		}()

		if flow.ShouldRun() {
			next.ServeHTTP(w, r)
		} else {
			resp := flow.CheckResponse().GetDeniedResponse()
			// If there was connection error, the response will be nil.
			if resp == nil {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				w.WriteHeader(int(resp.GetStatus()))
				for key, value := range resp.GetHeaders() {
					w.Header().Set(key, value)
				}
				_, err := fmt.Fprint(w, resp.GetBody())
				if err != nil {
					m.client.GetLogger().Info("Aperture flow control respond body got an error.", "error", err)
				}
			}
		}
	})
}

func prepareCheckHTTPRequestForHTTP(req *http.Request, logger *slog.Logger, controlPoint string, flowParams aperture.FlowParams) *checkhttpv1.CheckHTTPRequest {
	labels := utils.LabelsFromCtx(req.Context())

	// override labels with explicit labels
	for key, value := range flowParams.Labels {
		labels[key] = value
	}

	// override labels with labels from headers
	for key, value := range req.Header {
		if strings.HasPrefix(key, ":") {
			continue
		}
		labels[key] = strings.Join(value, ",")
	}

	// We know that the protocol is TCP because Golang's http package doesn't support UDP
	// TODO: Should we support `httpu`?
	protocol := checkhttpv1.SocketAddress_TCP

	sourceHost, sourcePort, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		logger.Error("Failed to parse source address", "error", err)
	}

	sourcePortU32, err := strconv.ParseUint(sourcePort, 10, 32)
	if err != nil {
		logger.Error("Failed to parse source port", "error", err)
	}

	// TODO: Figure out if we can narrow down the port or figure out the host in a better way
	destinationPort := uint32(0)
	destinationHost := utils.GetLocalIP()

	var b bytes.Buffer
	req.Body = io.NopCloser(io.TeeReader(req.Body, &b))

	return &checkhttpv1.CheckHTTPRequest{
		Source: &checkhttpv1.SocketAddress{
			Address:  sourceHost,
			Protocol: protocol,
			Port:     uint32(sourcePortU32),
		},
		Destination: &checkhttpv1.SocketAddress{
			Address:  destinationHost,
			Protocol: protocol,
			Port:     destinationPort,
		},
		ControlPoint: controlPoint,
		RampMode:     flowParams.RampMode,
		Request: &checkhttpv1.CheckHTTPRequest_HttpRequest{
			Method:   req.Method,
			Path:     req.URL.Path,
			Host:     req.Host,
			Headers:  labels,
			Scheme:   req.URL.Scheme,
			Size:     req.ContentLength,
			Protocol: req.Proto,
			Body:     b.String(),
		},
	}
}
