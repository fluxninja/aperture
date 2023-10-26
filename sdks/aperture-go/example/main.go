package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/stdr"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	aperturego "github.com/fluxninja/aperture-go/v2/sdk"
	aperturegomiddleware "github.com/fluxninja/aperture-go/v2/sdk/middleware"
)

const (
	defaultAgentHost = "localhost"
	defaultAgentPort = "8089"
	defaultAppPort   = "8080"
)

// app struct contains the server and the Aperture client.
type app struct {
	server         *http.Server
	grpcClient     *grpc.ClientConn
	apertureClient aperturego.Client
}

// grpcClient creates a new gRPC client that will be passed in order to initialize the Aperture client.
func grpcClient(ctx context.Context, address string) (*grpc.ClientConn, error) {
	// creating a gRPC client connection is essential to allow the Aperture client to communicate with the Flow Control Service.
	var grpcDialOptions []grpc.DialOption
	grpcDialOptions = append(grpcDialOptions, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: time.Second * 10,
	}))
	grpcDialOptions = append(grpcDialOptions, grpc.WithUserAgent("aperture-go"))
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, //nolint:gosec // Requires enabling CLI option
	})
	grpcDialOptions = append(grpcDialOptions, grpc.WithTransportCredentials(cred))
	grpcDialOptions = append(grpcDialOptions, grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("Invoking method %s", method)
		md := metadata.Pairs("apikey", "7b712208201b4fccb9fa1ca46d2e63fa")
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}))

	return grpc.DialContext(ctx, address, grpcDialOptions...)
}

func main() {
	agentHost := getEnvOrDefault("APERTURE_AGENT_HOST", defaultAgentHost)
	agentPort := getEnvOrDefault("APERTURE_AGENT_PORT", defaultAgentPort)

	ctx := context.Background()

	apertureAgentGRPCClient, err := grpcClient(ctx, net.JoinHostPort(agentHost, agentPort))
	if err != nil {
		log.Fatalf("failed to create flow control client: %v", err)
	}

	stdr.SetVerbosity(2)

	opts := aperturego.Options{
		ApertureAgentGRPCClientConn: apertureAgentGRPCClient,
	}

	// initialize Aperture Client with the provided options.
	apertureClient, err := aperturego.NewClient(ctx, opts)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	appPort := getEnvOrDefault("FN_APP_PORT", defaultAppPort)
	// Create a server with passing it the Aperture client.
	mux := mux.NewRouter()
	a := &app{
		server: &http.Server{
			Addr:    net.JoinHostPort("localhost", appPort),
			Handler: mux,
		},
		apertureClient: apertureClient,
		grpcClient:     apertureAgentGRPCClient,
	}

	// Adding the http middleware to be executed before the actual business logic execution.
	superRouter := mux.PathPrefix("/super").Subrouter()
	superRouter.HandleFunc("", a.SuperHandler)
	superRouter.Use(aperturegomiddleware.NewHTTPMiddleware(apertureClient, "awesomeFeature", nil, nil, false, 200*time.Millisecond).Handle)

	mux.HandleFunc("/connected", a.ConnectedHandler)
	mux.HandleFunc("/health", a.HealthHandler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting example server")

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %+v", err)
		}
	}()

	<-done
	if err := apertureClient.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown aperture client: %+v", err)
	}
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %+v", err)
	}
}

// SuperHandler handles HTTP requests on /super endpoint.
func (a *app) SuperHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	// Simulate work being done
	time.Sleep(2 * time.Second)
}

// ConnectedHandler handles HTTP requests on /connected endpoint.
func (a *app) ConnectedHandler(w http.ResponseWriter, r *http.Request) {
	a.grpcClient.Connect()
	state := a.grpcClient.GetState()
	if state != connectivity.Ready {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_, _ = w.Write([]byte(state.String()))
}

// HealthHandler handles HTTP requests on /health endpoint.
func (a *app) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Healthy"))
}

func getEnvOrDefault(envName, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}
