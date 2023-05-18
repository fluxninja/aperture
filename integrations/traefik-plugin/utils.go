package traefikplugin

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultAgentHost = "localhost"
	defaultAgentPort = "8080"
)

func grpcClient(ctx context.Context, address string) (*grpc.ClientConn, error) {
	// creating a gRPC client connection is essential to allow the Aperture client to communicate with the Flow Control Service.
	var grpcDialOptions []grpc.DialOption
	grpcDialOptions = append(grpcDialOptions, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: time.Second * 10,
	}))
	grpcDialOptions = append(grpcDialOptions, grpc.WithUserAgent("aperture-traefik-plugin"))
	grpcDialOptions = append(grpcDialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return grpc.DialContext(ctx, address, grpcDialOptions...)
}

func getEnvOrDefault(envName, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}
