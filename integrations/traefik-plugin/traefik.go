package aperturetraefikplugin

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	"github.com/go-logr/stdr"

	aperture "github.com/fluxninja/aperture-go/v2/sdk"
	aperturemiddlewares "github.com/fluxninja/aperture-go/v2/sdk/middlewares"
)

type Config struct {
	ControlPoint string
	AgentHost    string
	AgentPort    string
}

type TraefikPlugin struct {
	next         http.Handler
	ControlPoint string
	AgentHost    string
	AgentPort    string
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TraefikPlugin{
		next:         next,
		ControlPoint: config.ControlPoint,
		AgentHost:    config.AgentHost,
		AgentPort:    config.AgentPort,
	}, nil
}

func (a TraefikPlugin) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	apertureAgentGRPCClient, err := grpcClient(ctx, net.JoinHostPort(a.AgentHost, a.AgentPort))
	if err != nil {
		log.Fatalf("failed to create flow control client: %v", err)
	}

	// Initialize the logger
	logger := stdr.New(log.Default()).WithName("aperture-traefik-plugin")

	opts := aperture.Options{
		ApertureAgentGRPCClientConn: apertureAgentGRPCClient,
		CheckTimeout:                200 * time.Millisecond,
		Logger:                      &logger,
	}

	//initialize Aperture Client with the provided options.
	apertureClient, err := aperture.NewClient(ctx, opts)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	labels := aperture.LabelsFromCtx(r.Context())

	for key, value := range r.Header {
		if strings.HasPrefix(key, ":") {
			continue
		}
		labels[key] = strings.Join(value, ",")
	}

	protocol := flowcontrolhttp.SocketAddress_TCP

	sourceHost, sourcePort := aperturemiddlewares.SplitAddress(apertureClient.GetLogger(), r.RemoteAddr)

	destinationPort := uint32(0)
	destinationHost := aperturemiddlewares.GetLocalIP(apertureClient.GetLogger())

	bodyBytes, err := aperturemiddlewares.ReadClonedBody(r)
	if err != nil {
		apertureClient.GetLogger().V(2).Info("Error reading body", "error", err)
	}

	req := &flowcontrolhttp.CheckHTTPRequest{
		Source: &flowcontrolhttp.SocketAddress{
			Protocol: protocol,
			Address:  sourceHost,
			Port:     sourcePort,
		},
		Destination: &flowcontrolhttp.SocketAddress{
			Protocol: protocol,
			Address:  destinationHost,
			Port:     destinationPort,
		},
		ControlPoint: a.ControlPoint,
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

	flow, err := apertureClient.StartHTTPFlow(r.Context(), req)
	if err != nil {
		apertureClient.GetLogger().Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Accepted()", flow.Accepted())
	}

	if flow.Accepted() {
		a.next.ServeHTTP(rw, r)
		// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
		// The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
		// The second argument is error message for further diagnosis.
		err := flow.End(aperture.OK)
		if err != nil {
			apertureClient.GetLogger().Info("Aperture flow control end got error.", "error", err)
		}
	} else {
		resp := flow.CheckResponse().GetDeniedResponse()

		if resp == nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
		} else {
			rw.WriteHeader(int(resp.GetStatus()))
			for key, value := range resp.GetHeaders() {
				rw.Header().Set(key, value)
			}

			_, err := fmt.Fprint(rw, resp.GetBody())
			if err != nil {
				apertureClient.GetLogger().Info("Aperture flow control respond body got an error.", "error", err)
			}
		}

		err = flow.End(aperture.OK)
		if err != nil {
			apertureClient.GetLogger().Info("Aperture flow control end got error.", "error", err)
		}
	}
}
