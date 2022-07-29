package heartbeats

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/fluxninja/aperture/pkg/utils"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

const (
	heartbeatsGroup    = "heartbeats-job-group"
	jobName            = "aperture-heartbeats"
	jobNameHTTP        = "aperture-heartbeats-http"
	jobTimeoutDuration = time.Minute * 2
	healthKey          = "health"
	healthDetailsKey   = "health-details"
	healthyStatus      = "healthy"
	entityCacheKey     = "entity-cache"
	servicesKey        = entityCacheKey + ".services"
	overlappingKey     = entityCacheKey + ".overlapping-services"
	heartbeatsHTTPPath = "/plugins/fluxninja/v1/report"
)

type heartbeats struct {
	heartbeatsClient heartbeatv1.FluxNinjaServiceClient
	peersWatcher     *peers.PeerDiscovery
	clientHTTP       *http.Client
	agentInfo        *agentinfo.AgentInfo
	interval         config.Duration
	jobGroup         *jobs.JobGroup
	clientConn       *grpc.ClientConn
	statusRegistry   *status.Registry
	entityCache      *entitycache.EntityCache
	heartbeatsAddr   string
	APIKey           string
	jobName          string
}

func newHeartbeats(
	jobGroup *jobs.JobGroup,
	p pluginconfig.FluxNinjaPluginConfig,
	statusRegistry *status.Registry,
	entityCache *entitycache.EntityCache,
	agentInfo *agentinfo.AgentInfo,
	peersWatcher *peers.PeerDiscovery,
) *heartbeats {
	return &heartbeats{
		heartbeatsAddr: p.FluxNinjaEndpoint,
		interval:       p.HeartbeatInterval,
		APIKey:         p.APIKey,
		jobGroup:       jobGroup,
		statusRegistry: statusRegistry,
		entityCache:    entityCache,
		agentInfo:      agentInfo,
		peersWatcher:   peersWatcher,
	}
}

func (h *heartbeats) start(ctx context.Context, in *ConstructorIn) error {
	var job jobs.Job
	var err error

	if utils.IsHTTPUrl(h.heartbeatsAddr) {
		job, err = h.createHTTPJob(ctx, in.HTTPClient)
		if err != nil {
			return err
		}
	} else {
		job, err = h.createGRPCJob(ctx, in.GRPClientConnectionBuilder)
		if err != nil {
			return err
		}
	}
	h.registerHearbeatsJob(job)

	return nil
}

func (h *heartbeats) createGRPCJob(ctx context.Context, grpcClientConnBuilder grpcclient.ClientConnectionBuilder) (jobs.Job, error) {
	log.Debug().Str("heartbeatsAddr", h.heartbeatsAddr).Msg("Heartbeats service address")
	connWrapper := grpcClientConnBuilder.Build()
	conn, err := connWrapper.Dial(ctx, h.heartbeatsAddr)
	if err != nil {
		log.Warn().Err(err).Msg("Could not connect to heartbeat grpc server")
		return nil, err
	}

	h.clientConn = conn
	h.heartbeatsClient = heartbeatv1.NewFluxNinjaServiceClient(conn)

	job := jobs.BasicJob{
		JobFunc: h.sendSingleHeartbeat,
	}
	job.JobName = jobName

	return &job, nil
}

func (h *heartbeats) createHTTPJob(ctx context.Context, restapiClientConnection *http.Client) (jobs.Job, error) {
	h.heartbeatsAddr += heartbeatsHTTPPath

	h.clientHTTP = restapiClientConnection
	job := jobs.BasicJob{
		JobFunc: h.sendSingleHeartbeatByHTTP,
	}
	job.JobName = jobNameHTTP
	return &job, nil
}

func (h *heartbeats) registerHearbeatsJob(job jobs.Job) {
	executionTimeout := config.Duration{Duration: durationpb.New(jobTimeoutDuration)}
	jobConfig := jobs.JobConfig{
		InitiallyHealthy: true,
		ExecutionPeriod:  h.interval,
		ExecutionTimeout: executionTimeout,
	}
	// Setup with job registry
	err := h.jobGroup.RegisterJob(job, jobConfig)
	if err != nil {
		log.Error().Err(err).Str("group", heartbeatsGroup).Str("job", jobName).Msg("Error registering job")
	}
}

func (h *heartbeats) stop() {
	log.Info().Msg("Stopping gRPC heartbeats")

	if h.clientConn != nil {
		_ = h.clientConn.Close()
	}

	_ = h.jobGroup.DeregisterJob(h.jobName)
}

func (h *heartbeats) newHeartbeat(
	jobCtxt context.Context,
	health,
	healthDetails string,
) *heartbeatv1.ReportRequest {
	allStasuses := make(map[string]*anypb.Any)

	flatResults, err := h.statusRegistry.GetAllFlat()
	if err != nil {
		log.Error().Err(err).Msg("could not get results from status registry")
	} else {
		for flatKey, flatValue := range flatResults {
			allStasuses[flatKey] = flatValue.Status.GetMessage()
		}
	}

	if h.entityCache != nil {
		serviceList, overlappingList := h.entityCache.Services()
		packedSvcs := &heartbeatv1.Services{Services: serviceList}
		anySvcs, err := anypb.New(packedSvcs)
		if err != nil {
			log.Error().Err(err).Msg("Cannot cast packed services to Any")
		}

		packedOverlaps := &heartbeatv1.OverlappingServices{OverlappingServices: overlappingList}
		anyOverlaps, err := anypb.New(packedOverlaps)
		if err != nil {
			log.Error().Err(err).Msg("Cannot cast packed overlapping services to Any")
		}
		allStasuses[servicesKey] = anySvcs
		allStasuses[overlappingKey] = anyOverlaps
	}

	anyHealth, _ := anypb.New(wrapperspb.String(health))
	anyHealthDetails, _ := anypb.New(wrapperspb.String(healthDetails))
	allStasuses[healthKey] = anyHealth
	allStasuses[healthDetailsKey] = anyHealthDetails

	var agentGroup string
	if h.agentInfo != nil {
		agentGroup = h.agentInfo.GetAgentGroup()
	}

	var peerInfos []*peersv1.PeerInfo
	if h.peersWatcher != nil {
		peerInfos = h.peersWatcher.GetPeers()
	}

	return &heartbeatv1.ReportRequest{
		VersionInfo: info.GetVersionInfo(),
		ProcessInfo: info.GetProcessInfo(),
		HostInfo:    info.GetHostInfo(),
		AgentGroup:  agentGroup,
		PeerInfos:   peerInfos,
		AllStatuses: allStasuses,
	}
}

func (h *heartbeats) sendSingleHeartbeat(jobCtxt context.Context) (proto.Message, error) {
	report := h.newHeartbeat(jobCtxt, healthyStatus, "")

	// Add api key value to metadata
	md := metadata.Pairs("apiKey", h.APIKey)
	ctx := metadata.NewOutgoingContext(jobCtxt, md)
	_, err := h.heartbeatsClient.Report(ctx, report)
	if err != nil {
		log.Warn().Err(err).Msg("could not send heartbeat report")
	}
	return &emptypb.Empty{}, nil
}

func (h *heartbeats) sendSingleHeartbeatByHTTP(jobCtxt context.Context) (proto.Message, error) {
	report := h.newHeartbeat(jobCtxt, healthyStatus, "")
	reqBody, err := protojson.MarshalOptions{
		UseProtoNames: true,
	}.Marshal(report)
	if err != nil {
		log.Warn().Err(err).Msg("could not marshal report")
		return &emptypb.Empty{}, err
	}

	req, err := http.NewRequestWithContext(jobCtxt, "POST", h.heartbeatsAddr, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Warn().Err(err).Msg("could not create request")
		return &emptypb.Empty{}, err
	}
	req.Header.Add("apiKey", h.APIKey)
	cli := http.DefaultClient
	cli.Transport = &http.Transport{}
	_, err = cli.Do(req)
	if err != nil {
		log.Warn().Err(err).Msg("could not send heartbeat report")
	}
	return &emptypb.Empty{}, nil
}
