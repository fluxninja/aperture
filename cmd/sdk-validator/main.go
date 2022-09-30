package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

var (
	logger   *log.Logger
	accepted int32
)

func init() {
	logger = log.NewLogger(log.GetPrettyConsoleWriter(), "debug")
	log.SetGlobalLogger(logger)

	accepted = 0
}

type flowcontrol struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer

	accepts int
}

// Check is a dummy Check handler.
func (f *flowcontrol) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	logger.Info().Msg("Received Check request")

	services := []string{}
	rpcPeer, peerExists := peer.FromContext(ctx)
	if peerExists {
		services = append(services, rpcPeer.Addr.String())
	}

	start := time.Now()
	resp := f.check(ctx, req.Feature, req.Labels, services)
	end := time.Now()

	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)

	return resp, nil
}

func (f *flowcontrol) check(ctx context.Context, feature string, labels map[string]string, services []string) *flowcontrolv1.CheckResponse {
	resp := &flowcontrolv1.CheckResponse{
		DecisionType:  flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys: maps.Keys(labels),
		Services:      services,
		ControlPointInfo: &flowcontrolv1.ControlPointInfo{
			Feature: feature,
			Type:    flowcontrolv1.ControlPointInfo_TYPE_FEATURE,
		},
		RejectReason: flowcontrolv1.CheckResponse_REJECT_REASON_NONE,
	}

	if accepted >= int32(f.accepts) {
		resp.DecisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		atomic.StoreInt32(&accepted, 0)
	} else {
		atomic.AddInt32(&accepted, 1)
	}

	return resp
}

func main() {
	// setup flagset and flags
	fs := flag.NewFlagSet("sdk-validator", flag.ExitOnError)
	port := fs.String("port", "8080", "Port to start sdk-validator's grpc server on. Default is 8080.")
	accepts := fs.Int("accepts", 0, "Number of Check calls to accept. Default is 0 (accept all).")
	// parse flags
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
		os.Exit(1)
	}

	// create listener for grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// instantiate flowcontrol
	f := &flowcontrol{
		accepts: *accepts,
	}

	// setup grpc server and register FlowControlServiceServer instance to it
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	flowcontrolv1.RegisterFlowControlServiceServer(grpcServer, f)

	// start serving traffic on grpc server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
