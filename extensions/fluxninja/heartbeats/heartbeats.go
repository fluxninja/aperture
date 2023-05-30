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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	heartbeatv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/fluxninja/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/extconfig"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/cache"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/discovery/entities"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/etcd/election"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	grpcclient "github.com/fluxninja/aperture/v2/pkg/net/grpc"
	"github.com/fluxninja/aperture/v2/pkg/peers"
	autoscalek8sdiscovery "github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/discovery"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	flowcontrolpoints "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/controlpoints"
	"github.com/fluxninja/aperture/v2/pkg/status"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

const (
	heartbeatsGroup    = "heartbeats-job-group"
	jobName            = "aperture-heartbeats"
	jobNameHTTP        = "aperture-heartbeats-http"
	jobTimeoutDuration = time.Minute * 2
	entitiesKey        = "entity-cache"
	servicesKey        = entitiesKey + ".services"
	overlappingKey     = entitiesKey + ".overlapping-services"
	heartbeatsHTTPPath = "/fluxninja/v1/report"
)

// Heartbeats is the struct that holds information about heartbeats.
type Heartbeats struct {
	heartbeatv1.UnimplementedControllerInfoServiceServer
	heartbeatsClient          heartbeatv1.FluxNinjaServiceClient
	statusRegistry            status.Registry
	autoscalek8sControlPoints autoscalek8sdiscovery.AutoScaleControlPoints
	policyFactory             *controlplane.PolicyFactory
	ControllerInfo            *heartbeatv1.ControllerInfo // set in OnStart
	jobGroup                  *jobs.JobGroup
	clientConn                *grpc.ClientConn // set in OnStart
	peersWatcher              *peers.PeerDiscovery
	entities                  *entities.Entities
	clientHTTP                *http.Client // set in OnStart
	interval                  config.Duration
	flowControlPoints         *cache.Cache[selectors.TypedControlPointID]
	agentInfo                 *agentinfo.AgentInfo
	election                  *election.Election
	APIKey                    string
	jobName                   string
	installationMode          string
	heartbeatsAddr            string
}

func newHeartbeats(
	jobGroup *jobs.JobGroup,
	extensionConfig *extconfig.FluxNinjaExtensionConfig,
	statusRegistry status.Registry,
	entities *entities.Entities,
	agentInfo *agentinfo.AgentInfo,
	peersWatcher *peers.PeerDiscovery,
	policyFactory *controlplane.PolicyFactory,
	election *election.Election,
	flowControlPoints *cache.Cache[selectors.TypedControlPointID],
	autoscalek8sControlPoints autoscalek8sdiscovery.AutoScaleControlPoints,
) *Heartbeats {
	return &Heartbeats{
		heartbeatsAddr:            extensionConfig.Endpoint,
		interval:                  extensionConfig.HeartbeatInterval,
		APIKey:                    extensionConfig.APIKey,
		jobGroup:                  jobGroup,
		statusRegistry:            statusRegistry,
		entities:                  entities,
		agentInfo:                 agentInfo,
		peersWatcher:              peersWatcher,
		policyFactory:             policyFactory,
		election:                  election,
		flowControlPoints:         flowControlPoints,
		autoscalek8sControlPoints: autoscalek8sControlPoints,
		installationMode:          extensionConfig.InstallationMode,
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
		log.Warn().Err(err).Msg("Could not connect to heartbeat gRPC server")
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
	report := &heartbeatv1.ReportRequest{
		VersionInfo:      info.GetVersionInfo(),
		ProcessInfo:      info.GetProcessInfo(),
		HostInfo:         info.GetHostInfo(),
		ControllerInfo:   h.ControllerInfo,
		AllStatuses:      h.statusRegistry.GetGroupStatus(),
		InstallationMode: h.installationMode,
	}

	var agentGroup string
	if h.agentInfo != nil {
		agentGroup = h.agentInfo.GetAgentGroup()
		report.AgentGroup = agentGroup
	}

	policies := &policysyncv1.PolicyWrappers{}
	if h.policyFactory != nil {
		policies.PolicyWrappers = h.policyFactory.GetPolicyWrappers()
		report.Policies = policies
	}

	if h.flowControlPoints != nil {
		report.FlowControlPoints = flowcontrolpoints.ToProto(h.flowControlPoints)
	}

	if h.election != nil && h.election.IsLeader() {
		var servicesList *heartbeatv1.ServicesList
		if h.entities != nil {
			servicesList = populateServicesList(h.entities)
			report.ServicesList = servicesList
		}

		if h.peersWatcher != nil {
			peers := h.peersWatcher.GetPeers()
			report.Peers = peers
		}

		if h.autoscalek8sControlPoints != nil {
			report.AutoScaleKubernetesControlPoints = h.autoscalek8sControlPoints.ToProto()
		}
	}
	return report
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

// GetControllerInfo returns the controller info.
func (h *Heartbeats) GetControllerInfo(context.Context, *emptypb.Empty) (*heartbeatv1.ControllerInfo, error) {
	return h.ControllerInfo, nil
}
