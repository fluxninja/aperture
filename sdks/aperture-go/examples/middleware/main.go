package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	aperture "github.com/fluxninja/aperture-go/v2/sdk"
	"github.com/fluxninja/aperture-go/v2/sdk/middleware"
)

const (
	defaultAgentAddress = "localhost:8089"
	defaultAppPort      = "8080"
)

// app struct contains the server and the Aperture client.
type app struct {
	server         *http.Server
	apertureClient aperture.Client
}

// grpcOptions creates a new gRPC client that will be passed in order to initialize the Aperture client.
func grpcOptions(insecureMode, skipVerify bool) []grpc.DialOption {
	var grpcDialOptions []grpc.DialOption
	grpcDialOptions = append(grpcDialOptions, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: time.Second * 10,
	}))
	grpcDialOptions = append(grpcDialOptions, grpc.WithUserAgent("aperture-go"))
	if insecureMode {
		grpcDialOptions = append(grpcDialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else if skipVerify {
		grpcDialOptions = append(grpcDialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // For testing purposes only
		})))
	} else {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil
		}
		grpcDialOptions = append(grpcDialOptions, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, "")))
	}
	return grpcDialOptions
}

func main() {
	ctx := context.Background()

	apertureAgentAddr := getEnvOrDefault("APERTURE_AGENT_ADDRESS", defaultAgentAddress)
	apertureAgentInsecure := getEnvOrDefault("APERTURE_AGENT_INSECURE", "false")
	apertureAgentInsecureBool, _ := strconv.ParseBool(apertureAgentInsecure)
	apertureAgentSkipVerify := getEnvOrDefault("APERTURE_AGENT_SKIP_VERIFY", "false")
	apertureAgentSkipVerifyBool, _ := strconv.ParseBool(apertureAgentSkipVerify)

	opts := aperture.Options{
		Address:     apertureAgentAddr,
		DialOptions: grpcOptions(apertureAgentInsecureBool, apertureAgentSkipVerifyBool),
		APIKey:      getEnvOrDefault("APERTURE_API_KEY", ""),
	}

	// initialize Aperture Client with the provided options.
	apertureClient, err := aperture.NewClient(ctx, opts)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	appPort := getEnvOrDefault("APERTURE_APP_PORT", defaultAppPort)
	// Create a server with passing it the Aperture client.
	mux := mux.NewRouter()
	a := &app{
		server: &http.Server{
			Addr:    net.JoinHostPort("localhost", appPort),
			Handler: mux,
		},
		apertureClient: apertureClient,
	}

	// Adding the http middleware to be executed before the actual business logic execution.
	superRouter := mux.PathPrefix("/super").Subrouter()
	superRouter.HandleFunc("", a.SuperHandler)

	// START: middleware

	middlewareParams := aperture.MiddlewareParams{
		Timeout:              2000 * time.Millisecond,
		IgnoredPathsCompiled: []*regexp.Regexp{regexp.MustCompile("/health.*")},
		IgnoredPaths:         []string{"/connected"},
	}

	middleware, err := middleware.NewHTTPMiddleware(apertureClient, "awesomeFeature", middlewareParams)
	if err != nil {
		log.Fatalf("failed to create HTTP middleware: %v", err)
	}
	superRouter.Use(middleware.Handle)

	// END: middleware

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
	a.apertureClient.GetGRPClientConn().Connect()
	state := a.apertureClient.GetGRPClientConn().GetState()
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

func getEnvOrDefault(envName string, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}
