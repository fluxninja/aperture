package traefikplugin

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	aperture "github.com/fluxninja/aperture-go/v2/sdk"

	"github.com/go-logr/stdr"
)

type Config struct {
	ControlPoint string
	//Labels       map[string]string
}

type TraefikPlugin struct {
	next         http.Handler
	ControlPoint string
	//Labels       map[string]string
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TraefikPlugin{
		next:         next,
		ControlPoint: config.ControlPoint,
		//Labels:       config.Labels,
	}, nil
}

func (a TraefikPlugin) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	agentHost := getEnvOrDefault("FN_AGENT_HOST", defaultAgentHost)
	agentPort := getEnvOrDefault("FN_AGENT_PORT", defaultAgentPort)

	ctx := context.Background()

	apertureAgentGRPCClient, err := grpcClient(ctx, net.JoinHostPort(agentHost, agentPort))
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
	a.next.ServeHTTP(rw, r)
}
