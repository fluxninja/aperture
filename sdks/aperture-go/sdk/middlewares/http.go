package middlewares

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	"github.com/gorilla/mux"
)

// TODO: Create struct for middleware

// HTTPMiddleware takes a control point name and paths to ignore, and creates a Middleware which can be used with HTTP server.
func HTTPMiddleware(client *aperture.Client, controlPoint string, ignoredPaths *[]regexp.Regexp) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the path is ignored, skip the middleware
			if ignoredPaths != nil {
				for _, ignoredPath := range *ignoredPaths {
					if ignoredPath.MatchString(r.URL.Path) {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			labels := labelsFromCtx(r.Context())

			for key, value := range r.Header {
				if strings.HasPrefix(key, ":") {
					continue
				}
				labels[key] = strings.Join(value, ",")
			}

			// We know that the protocol is TCP because Golang's http package doesn't support UDP
			// TODO: Should we support `httpu`?
			protocol := flowcontrolhttp.SocketAddress_TCP

			sourceHost, sourcePort := c.splitAddress(r.RemoteAddr)
			destinationHost, destinationPort := c.splitAddress(r.URL.Host)

			bodyBytes, err := readClonedBody(r)
			if err != nil {
				c.log.V(2).Info("Error reading body", "error", err)
			}

			req := &flowcontrolhttp.CheckHTTPRequest{
				Source: &flowcontrolhttp.SocketAddress{
					Address:  sourceHost,
					Protocol: protocol,
					Port:     sourcePort,
				},
				Destination: &flowcontrolhttp.SocketAddress{
					Address:  destinationHost,
					Protocol: protocol,
					Port:     destinationPort,
				},
				ControlPoint: controlPoint,
				Request: &flowcontrolhttp.CheckHTTPRequest_HttpRequest{
					Method:   r.Method,
					Path:     r.URL.Path,
					Host:     r.Host,
					Headers:  labels,
					Scheme:   r.URL.Scheme,
					Size:     r.ContentLength,
					Protocol: r.Proto,
					Body:     string(bodyBytes),
				},
			}

			flow, err := c.StartHTTPFlow(r.Context(), req)
			if err != nil {
				c.log.Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Accepted()", flow.Accepted())
			}

			if flow.Accepted() {
				// Simulate work being done
				next.ServeHTTP(w, r)
				// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
				// The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
				// The second argument is error message for further diagnosis.
				err := flow.End(OK)
				if err != nil {
					c.log.Info("Aperture flow control end got error.", "error", err)
				}
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
						c.log.Info("Aperture flow control end got error.", "error", err)
					}
				}
				err = flow.End(OK)
				if err != nil {
					c.log.Info("Aperture flow control end got error.", "error", err)
				}
			}
		})
	}
}
