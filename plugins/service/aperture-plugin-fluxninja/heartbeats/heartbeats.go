package heartbeats

import (
	"bytes"
	"context"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/entitycache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
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
	entityCacheKey     = "entity-cache"
	servicesKey        = entityCacheKey + ".services"
	overlappingKey     = entityCacheKey + ".overlapping-services"
	heartbeatsHTTPPath = "/plugins/fluxninja/v1/report"
)

type Heartbeats struct {
	heartbeatsClient heartbeatv1.FluxNinjaServiceClient
	peersWatcher     *peers.PeerDiscovery
	clientHTTP       *http.Client
	agentInfo        *agentinfo.AgentInfo
	interval         config.Duration
	jobGroup         *jobs.JobGroup
	clientConn       *grpc.ClientConn
	statusRegistry   status.Registry
	entityCache      *entitycache.EntityCache
	controllerInfo   *heartbeatv1.ControllerInfo
	heartbeatsAddr   string
	APIKey           string
	jobName          string
}

func (h *Heartbeats) GetControllerInfo() *heartbeatv1.ControllerInfo {
	return h.controllerInfo
}

func newHeartbeats(
	jobGroup *jobs.JobGroup,
	p pluginconfig.FluxNinjaPluginConfig,
	statusRegistry status.Registry,
	entityCache *entitycache.EntityCache,
	agentInfo *agentinfo.AgentInfo,
	peersWatcher *peers.PeerDiscovery,
) *Heartbeats {
	return &Heartbeats{
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

func (h *Heartbeats) start(ctx context.Context, in *ConstructorIn) error {
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

func (h *Heartbeats) setupControllerInfo(ctx context.Context, etcdClient *etcdclient.Client) error {
	etcdPath := "/fluxninja/controllerid"
	controllerID := guuid.NewString()

	txn := etcdClient.Client.Txn(etcdClient.Client.Ctx())
	resp, err := txn.If(clientv3.Compare(clientv3.CreateRevision(etcdPath), "=", 0)).
		Then(clientv3.OpPut(etcdPath, controllerID)).
		Else(clientv3.OpGet(etcdPath)).Commit()
	if err != nil {
		log.Error().Err(err).Msg("Could not read/write controller id to etcd")
		return err
	}

	// Succeeded is true if the If condition above is true - meaning there were no controller Id in etcd
	if !resp.Succeeded {
		for _, res := range resp.Responses {
			controllerID = string(res.GetResponseRange().Kvs[0].Value)
			break
		}
	}

	h.controllerInfo = &heartbeatv1.ControllerInfo{
		Id: controllerID,
	}

	return nil
}

func (h *Heartbeats) createGRPCJob(ctx context.Context, grpcClientConnBuilder grpcclient.ClientConnectionBuilder) (jobs.Job, error) {
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

func (h *Heartbeats) createHTTPJob(ctx context.Context, restapiClientConnection *http.Client) (jobs.Job, error) {
	h.heartbeatsAddr += heartbeatsHTTPPath

	h.clientHTTP = restapiClientConnection
	job := jobs.BasicJob{
		JobFunc: h.sendSingleHeartbeatByHTTP,
	}
	job.JobName = jobNameHTTP
	return &job, nil
}

func (h *Heartbeats) registerHearbeatsJob(job jobs.Job) {
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

func (h *Heartbeats) stop() {
	log.Info().Msg("Stopping gRPC heartbeats")

	if h.clientConn != nil {
		_ = h.clientConn.Close()
	}

	_ = h.jobGroup.DeregisterJob(h.jobName)
}

func (h *Heartbeats) newHeartbeat(
	jobCtxt context.Context,
) *heartbeatv1.ReportRequest {
	var entityCache *entitycachev1.EntityCache
	if h.entityCache != nil {
		entityCache = h.entityCache.Services()
	}

	var agentGroup string
	if h.agentInfo != nil {
		agentGroup = h.agentInfo.GetAgentGroup()
	}

	var peers *peersv1.Peers
	if h.peersWatcher != nil {
		peers = h.peersWatcher.GetPeers()
	}

	return &heartbeatv1.ReportRequest{
		VersionInfo:    info.GetVersionInfo(),
		ProcessInfo:    info.GetProcessInfo(),
		HostInfo:       info.GetHostInfo(),
		AgentGroup:     agentGroup,
		ControllerInfo: h.controllerInfo,
		Peers:          peers,
		EntityCache:    entityCache,
		AllStatuses:    h.statusRegistry.GetGroupStatus(),
	}
}

func (h *Heartbeats) sendSingleHeartbeat(jobCtxt context.Context) (proto.Message, error) {
	report := h.newHeartbeat(jobCtxt)

	// Add api key value to metadata
	md := metadata.Pairs("apiKey", h.APIKey)
	ctx := metadata.NewOutgoingContext(jobCtxt, md)
	_, err := h.heartbeatsClient.Report(ctx, report)
	if err != nil {
		log.Warn().Err(err).Msg("could not send heartbeat report")
	}
	return &emptypb.Empty{}, nil
}

func (h *Heartbeats) sendSingleHeartbeatByHTTP(jobCtxt context.Context) (proto.Message, error) {
	report := h.newHeartbeat(jobCtxt)
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
