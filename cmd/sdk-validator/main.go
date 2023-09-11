package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	flowcontrolhttpv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/fluxninja/aperture/v2/cmd/sdk-validator/validator"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/checkhttp"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/envoy"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

var (
	logger      *log.Logger
	spanFailed  bool
	authzFailed bool
)

func init() {
	logLevel, logLevelSet := os.LookupEnv("LOG_LEVEL")
	if !logLevelSet {
		logLevel = log.TraceLevel.String()
	}
	logger = log.NewLogger(log.GetPrettyConsoleWriter(), logLevel)
	log.SetGlobalLogger(logger)
}

func main() {
	// setup flagset and flags
	fs := flag.NewFlagSet("sdk-validator", flag.ExitOnError)
	port := fs.String("port", "8089", "Port to start sdk-validator's gRPC server on.")
	requests := fs.Int("requests", 10, "Number of requests to make to SDK example server.")
	rejects := fs.Int64("rejects", 5, "Number of requests (out of 'requests') to reject.")
	sdkDockerImage := fs.String("sdk-docker-image", "", "Docker image of SDK example to run.")
	sdkPort := fs.String("sdk-port", "8080", "Port to expose on SDK's example container.")
	sslCertFilepath := fs.String("ssl-certificate", "", "Filepath of SSL certificate to configure server TLS.")
	sslKeyFilepath := fs.String("ssl-key", "", "Filepath of SSL key to configure server TLS.")

	// parse flags
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
		os.Exit(1)
	}

	id := ""
	if *sdkDockerImage != "" {
		log.Info().Msg("Starting Docker container")
		id, err = runDockerContainer(*sdkDockerImage, *sdkPort)
		if err != nil {
			log.Fatal().Err(err).Str("image", *sdkDockerImage).Msg("Failed to run Docker container")
		}
		log.Info().Str("image", *sdkDockerImage).Str("id", id).Msg("Container started")
	}

	sdkURL := fmt.Sprintf("http://localhost:%s", *sdkPort)

	// create listener for gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	// setup gRPC server and register various server instances to it
	var grpcServer *grpc.Server
	if *sslCertFilepath != "" && *sslKeyFilepath != "" {
		creds, err2 := credentials.NewServerTLSFromFile(*sslCertFilepath, *sslKeyFilepath)
		if err2 != nil {
			log.Fatal().Err(err2).Msg("Failed to create TLS creds from provided files")
		}
		log.Info().Msg("Starting TLS-secured gRPC server")
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor), grpc.Creds(creds))
	} else {
		log.Info().Msg("Starting insecure gRPC server")
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	}
	reflection.Register(grpcServer)

	commonHandler := &validator.CommonHandler{
		Rejects:  *rejects,
		Rejected: 0,
	}

	// instantiate and register flowcontrol handler
	flowcontrolHandler := &validator.FlowControlHandler{
		CommonHandler: commonHandler,
	}
	flowcontrolv1.RegisterFlowControlServiceServer(grpcServer, flowcontrolHandler)

	alerter := alerts.NewSimpleAlerter(100)
	reg := status.NewRegistry(log.GetGlobalLogger(), alerter)
	agentInfo := agentinfo.NewAgentInfo(metrics.DefaultAgentGroup)

	// instantiate and register flowcontrolhttp handler
	flowcontrolHTTPHandler := checkhttp.NewHandler(
		classifier.NewClassificationEngine(agentInfo, reg),
		servicegetter.NewEmpty(),
		commonHandler,
	)
	flowcontrolhttpv1.RegisterFlowControlServiceHTTPServer(grpcServer, flowcontrolHTTPHandler)

	authzHandler := envoy.NewHandler(
		classifier.NewClassificationEngine(agentInfo, reg),
		servicegetter.NewEmpty(),
		commonHandler,
	)
	authv3.RegisterAuthorizationServer(grpcServer, authzHandler)

	// initiate and register otel trace handler
	traceHandler := &validator.TraceHandler{}
	tracev1.RegisterTraceServiceServer(grpcServer, traceHandler)

	validation := 0

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		log.Info().Interface("signal", s).Msg("Got signal, attempting graceful shutdown")
		grpcServer.GracefulStop()

		log.Info().Msg("Validating fail-open behavior")
		rejected := startTraffic(sdkURL, *requests)
		l := log.With().Int("total requests", *requests).Int64("expected rejections", 0).Int("got rejections", rejected).Logger()
		if rejected != 0 {
			l.Error().Msg("Fail-open validation failed")
		} else {
			l.Info().Msg("Fail-open validation successful")
		}

		if *sdkDockerImage != "" {
			log.Info().Interface("id", id).Msg("Stopping Docker container")
			err = stopDockerContainer(id)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to stop Docker container")
			}
		}
		wg.Done()
	}()

	if *sdkDockerImage != "" {
		wg.Add(1)
		go func() {
			// give the gRPC server some time to initialize
			time.Sleep(2 * time.Second)

			rejected := confirmConnectedAndStartTraffic(sdkURL, *requests)
			l := log.With().Int("total requests", *requests).Int64("expected rejections", *rejects).Int("got rejections", rejected).Logger()
			if rejected != int(*rejects) {
				l.Error().Msg("FlowControl validation failed")
				validation = 1
			}

			if spanFailed {
				l.Error().Msg("Span attributes validation failed")
				validation = 1
			}
			if authzFailed {
				l.Error().Msg("Authz validation failed")
				validation = 1
			}

			if validation == 0 {
				l.Info().Msg("Validation successful")
				sigCh <- syscall.SIGTERM
			} else {
				l.Info().Msg("Validation failed")
				sigCh <- syscall.SIGTERM
			}
			wg.Done()
		}()
	}

	// start serving traffic on gRPC server
	log.Info().Str("add", lis.Addr().String()).Msg("Starting sdk-validator")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
	wg.Wait()
	log.Info().Msg("Successful graceful shutdown")
	os.Exit(validation)
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	log.Info().Str("method", info.FullMethod).Dur("latency", time.Since(start)).Msg("Request served")
	if err != nil {
		log.Error().Err(err).Msg("Handler returned error")
		if info.FullMethod == "/opentelemetry.proto.collector.trace.v1.TraceService/Export" {
			spanFailed = true
		} else if info.FullMethod == "/aperture.flowcontrol.checkhttp.v1.FlowControlServiceHTTP/CheckHTTP" || info.FullMethod == "/aperture.flowcontrol.check.v1.FlowControlService/Check" {
			authzFailed = true
		}
	}
	return h, err
}

func runDockerContainer(image string, port string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	exposedPorts, portBindings, _ := nat.ParsePortSpecs([]string{
		fmt.Sprintf("0.0.0.0:%s:%s", port, port),
	})

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        image,
			Tty:          true,
			OpenStdin:    true,
			AttachStdout: true,
			AttachStderr: true,
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			PortBindings: portBindings,
			NetworkMode:  "host",
		},
		nil, nil, "")
	if err != nil {
		return "", err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	time.Sleep(time.Second * 2)

	for {
		containerJSON, err := cli.ContainerInspect(ctx, resp.ID)
		if err != nil {
			return "", err
		}
		if containerJSON.State != nil {
			if containerJSON.State.Status == "exited" {
				printContainerLogs(resp.ID)
				log.Fatal().Msg("Container exited")
			}
			if containerJSON.State.Health != nil {
				if containerJSON.State.Health.Status == "healthy" {
					return resp.ID, nil
				}
			}
		}
	}
}

func printContainerLogs(id string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}
	// print all container logs
	out, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return
	}
	_, _ = io.Copy(os.Stdout, out)
}

func stopDockerContainer(id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	printContainerLogs(id)

	err = cli.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}

func confirmConnectedAndStartTraffic(url string, requests int) int {
	for {
		req, err := http.NewRequest(http.MethodGet, url+"/connected", nil)
		if err != nil {
			log.Error().Err(err).Str("url", req.URL.String()).Msg("Failed to create http request")
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Err(err).Str("url", req.URL.String()).Msg("Failed to make http request")
		}
		res.Body.Close()
		if res.StatusCode == http.StatusOK {
			break
		}
	}
	log.Info().Msg("SDK example successfully connected to validator")

	rejected := startTraffic(url, requests)
	return rejected
}

func startTraffic(url string, requests int) int {
	rejected := 0
	for i := 0; i < requests; i++ {
		superReq, err := http.NewRequest(http.MethodGet, url+"/super", nil)
		if err != nil {
			log.Error().Err(err).Str("url", superReq.URL.String()).Msg("Failed to create http request")
		}
		res, err := http.DefaultClient.Do(superReq)
		if err != nil {
			log.Error().Err(err).Str("url", superReq.URL.String()).Msg("Failed to make http request")
		}
		if res.Body != nil {
			res.Body.Close()
		}
		if res.StatusCode >= 400 {
			log.Trace().Int("status code", res.StatusCode).Msg("Request rejected")
			rejected += 1
		}
	}
	return rejected
}
