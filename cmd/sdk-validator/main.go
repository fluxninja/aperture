package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/cmd/sdk-validator/validator"
	"github.com/fluxninja/aperture/pkg/log"
)

var logger *log.Logger

func init() {
	logger = log.NewLogger(log.GetPrettyConsoleWriter(), "debug")
	log.SetGlobalLogger(logger)
}

func main() {
	// setup flagset and flags
	fs := flag.NewFlagSet("sdk-validator", flag.ExitOnError)
	port := fs.String("port", "8080", "Port to start sdk-validator's grpc server on. Default is 8080.")
	rejectRatio := fs.Float64("reject-ratio", 0.5, "Ratio of calls to reject. Default is 0.5.")
	sdkDockerImage := fs.String("sdk-docker-image", "", "Location of SDK example to run. Default is ''.")
	sdkPort := fs.String("sdk-port", "8081", "Port to expose on SDK's example container. Default is 8081.")
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

	// create listener for grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	// instantiate flowcontrol
	f := &validator.FlowControlHandler{
		RejectRatio: *rejectRatio,
	}

	// setup grpc server and register FlowControlServiceServer instance to it
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	flowcontrolv1.RegisterFlowControlServiceServer(grpcServer, f)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		log.Info().Interface("signal", s).Msg("Got signal, attempting graceful shutdown")
		log.Info().Interface("id", id).Msg("Stopping Docker container")
		err = stopDockerContainer(id)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to stop Docker container")
		}
		grpcServer.GracefulStop()
		wg.Done()
	}()

	// start serving traffic on grpc server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
	wg.Wait()
	log.Info().Msg("Successful graceful shutdown")
}

func runDockerContainer(image string, port string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()
	_, _ = io.Copy(os.Stdout, reader)

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
		},
		nil, nil, "")
	if err != nil {
		return "", err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func stopDockerContainer(id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerStop(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}
