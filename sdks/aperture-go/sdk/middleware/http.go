package middleware

import (
	"fmt"
	"net/http"

	aperture "github.com/fluxninja/aperture-go/v2/sdk"
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
func NewHTTPMiddleware(client aperture.Client, controlPoint string, middlewareParams aperture.MiddlewareParams) HTTPMiddleware {
	return &httpMiddleware{
		client:           client,
		controlPoint:     controlPoint,
		middlewareParams: middlewareParams,
	}
}

// Handle takes a http.Handler and returns a new http.Handler with the middleware applied.
func (m *httpMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the path is ignored, skip the middleware
		if m.middlewareParams.IgnoredPaths != nil {
			for _, ignoredPath := range m.middlewareParams.IgnoredPaths {
				if ignoredPath.MatchString(r.URL.Path) {
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
