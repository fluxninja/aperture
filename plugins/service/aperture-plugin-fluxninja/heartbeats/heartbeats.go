package heartbeats

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	guuid "github.com/google/uuid"
	"github.com/technosophos/moniker"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	controlpointcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/controlpointcache/v1"
	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/peers/v1"
	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/kubernetes"
	"github.com/fluxninja/aperture/pkg/entitycache"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
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
	heartbeatv1.UnimplementedControllerInfoServiceServer
	heartbeatsClient            heartbeatv1.FluxNinjaServiceClient
	statusRegistry              status.Registry
	agentInfo                   *agentinfo.AgentInfo
	clientHTTP                  *http.Client
	interval                    config.Duration
	jobGroup                    *jobs.JobGroup
	clientConn                  *grpc.ClientConn
	peersWatcher                *peers.PeerDiscovery
	entityCache                 *entitycache.EntityCache
	policyFactory               *controlplane.PolicyFactory
	ControllerInfo              *heartbeatv1.ControllerInfo
	serviceControlPointCache    *cache.Cache[selectors.ControlPointID]
	kubernetesControlPointCache kubernetes.ControlPointCache
	heartbeatsAddr              string
	APIKey                      string
	jobName                     string
	installationMode            string
}

func newHeartbeats(
	jobGroup *jobs.JobGroup,
	p pluginconfig.FluxNinjaPluginConfig,
	statusRegistry status.Registry,
	entityCache *entitycache.EntityCache,
	agentInfo *agentinfo.AgentInfo,
	peersWatcher *peers.PeerDiscovery,
	policyFactory *controlplane.PolicyFactory,
	serviceControlPointCache *cache.Cache[selectors.ControlPointID],
	kubernetesControlPointCache kubernetes.ControlPointCache,
	installationMode string,
) *Heartbeats {
	return &Heartbeats{
		heartbeatsAddr:              p.FluxNinjaEndpoint,
		interval:                    p.HeartbeatInterval,
		APIKey:                      p.APIKey,
		jobGroup:                    jobGroup,
		statusRegistry:              statusRegistry,
		entityCache:                 entityCache,
		agentInfo:                   agentInfo,
		peersWatcher:                peersWatcher,
		policyFactory:               policyFactory,
		serviceControlPointCache:    serviceControlPointCache,
		kubernetesControlPointCache: kubernetesControlPointCache,
		installationMode:            installationMode,
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
	newID := guuid.NewString()
	parts := strings.Split(newID, "-")
	moniker := strings.Replace(moniker.New().Name(), " ", "-", 1)
	controllerID := fmt.Sprintf("%s-%s", moniker, parts[0])

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

	h.ControllerInfo = &heartbeatv1.ControllerInfo{
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

	job := jobs.NewBasicJob(jobName, h.sendSingleHeartbeat)
	return job, nil
}

func (h *Heartbeats) createHTTPJob(ctx context.Context, restapiClientConnection *http.Client) (jobs.Job, error) {
	h.heartbeatsAddr += heartbeatsHTTPPath

	h.clientHTTP = restapiClientConnection
	job := jobs.NewBasicJob(jobNameHTTP, h.sendSingleHeartbeatByHTTP)
	return job, nil
}

func (h *Heartbeats) registerHearbeatsJob(job jobs.Job) {
	executionTimeout := config.MakeDuration(jobTimeoutDuration)
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
	var servicesList *heartbeatv1.ServicesList
	if h.entityCache != nil {
		servicesList = populateServicesList(h.entityCache)
	}

	var agentGroup string
	if h.agentInfo != nil {
		agentGroup = h.agentInfo.GetAgentGroup()
	}

	var peers *peersv1.Peers
	if h.peersWatcher != nil {
		peers = h.peersWatcher.GetPeers()
	}

	policies := &policysyncv1.PolicyWrappers{}
	if h.policyFactory != nil {
		policies.PolicyWrappers = h.policyFactory.GetPolicyWrappers()
	}

	serviceControlPointObjects := make(map[selectors.ControlPointID]struct{})
	if h.serviceControlPointCache != nil {
		serviceControlPointObjects = h.serviceControlPointCache.GetAllAndClear()
	}

	serviceControlPoints := make([]*heartbeatv1.ServiceControlPoint, 0, len(serviceControlPointObjects))
	for cp := range serviceControlPointObjects {
		serviceControlPoints = append(serviceControlPoints, &heartbeatv1.ServiceControlPoint{
			Name:        cp.ControlPoint,
			ServiceName: cp.Service,
		})
	}

	var kubernetesControlPoints []*controlpointcachev1.KubernetesControlPoint
	if h.kubernetesControlPointCache != nil {
		kubernetesControlPointObjects := h.kubernetesControlPointCache.Keys()
		kubernetesControlPoints = make([]*controlpointcachev1.KubernetesControlPoint, 0, len(kubernetesControlPointObjects))
		for _, cp := range kubernetesControlPointObjects {
			kubernetesControlPoints = append(kubernetesControlPoints, cp.ToProto())
		}
	}

	return &heartbeatv1.ReportRequest{
		VersionInfo:             info.GetVersionInfo(),
		ProcessInfo:             info.GetProcessInfo(),
		HostInfo:                info.GetHostInfo(),
		AgentGroup:              agentGroup,
		ControllerInfo:          h.ControllerInfo,
		Peers:                   peers,
		ServicesList:            servicesList,
		AllStatuses:             h.statusRegistry.GetGroupStatus(),
		Policies:                policies,
		ServiceControlPoints:    serviceControlPoints,
		KubernetesControlPoints: kubernetesControlPoints,
		InstallationMode:        h.installationMode,
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

func (h *Heartbeats) GetControllerInfo(context.Context, *emptypb.Empty) (*heartbeatv1.ControllerInfo, error) {
	return h.ControllerInfo, nil
}

func RegisterControllerInfoService(grpc *grpc.Server, handler *Heartbeats) error {
	heartbeatv1.RegisterControllerInfoServiceServer(grpc, handler)
	return nil
}

func RegisterControllerInfoServiceHTTP() fx.Option {
	return grpcgateway.RegisterHandler{Handler: heartbeatv1.RegisterControllerInfoServiceHandlerFromEndpoint}.Annotate()
}
