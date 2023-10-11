// Server-side for handling agent functions
package agents

import (
	"context"
	"fmt"
	"regexp"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
)

// Module is fx module for controlling Agents on controller side.
var Module = fx.Provide(NewAgents)

// Agents wraps rpc.Clients where clients are agents and provides wrapper
//
// Agents wraps functions registered in agentfunctions, types should match.
type Agents struct {
	etcdTransport *transport.EtcdTransportServer
	etcdClient    *etcdclient.Client
}

// NewAgents wraps Clients with Agent-specific function wrappers.
func NewAgents(transport *transport.EtcdTransportServer, client *etcdclient.Client) Agents {
	return Agents{
		etcdTransport: transport,
		etcdClient:    client,
	}
}

// ListFlowControlPoints lists control points of all agents.
//
// Handled by agentfunctions.ControlPointsHandler.
func (a Agents) ListFlowControlPoints() ([]transport.Result[*cmdv1.ListFlowControlPointsAgentResponse], error) {
	var req cmdv1.ListFlowControlPointsRequest
	agents, err := a.GetAgents()
	if err != nil {
		return nil, err
	}
	return transport.SendRequests[cmdv1.ListFlowControlPointsAgentResponse](a.etcdTransport, agents, &req)
}

// ListAutoScaleControlPoints lists auto-scale control points of all agents.
func (a Agents) ListAutoScaleControlPoints() ([]transport.Result[*cmdv1.ListAutoScaleControlPointsAgentResponse], error) {
	var req cmdv1.ListAutoScaleControlPointsRequest
	agents, err := a.GetAgents()
	if err != nil {
		return nil, err
	}
	return transport.SendRequests[cmdv1.ListAutoScaleControlPointsAgentResponse](a.etcdTransport, agents, &req)
}

// ListDiscoveryEntities lists discovery entities.
func (a Agents) ListDiscoveryEntities() ([]transport.Result[*cmdv1.ListDiscoveryEntitiesAgentResponse], error) {
	var req cmdv1.ListDiscoveryEntitiesRequest
	agents, err := a.GetAgents()
	if err != nil {
		return nil, err
	}
	return transport.SendRequests[cmdv1.ListDiscoveryEntitiesAgentResponse](a.etcdTransport, agents, &req)
}

// ListDiscoveryEntity lists discovery entity by ip address or name.
func (a Agents) ListDiscoveryEntity(req *cmdv1.ListDiscoveryEntityRequest) (*cmdv1.ListDiscoveryEntityAgentResponse, error) {
	agents, err := a.GetAgents()
	if err != nil {
		return nil, err
	}

	return transport.SendRequest[cmdv1.ListDiscoveryEntityAgentResponse](a.etcdTransport, agents[0], req)
}

// PreviewFlowLabels previews flow labels on a given agent.
//
// Handled by agentfunctions.PreviewHandler.
func (a Agents) PreviewFlowLabels(
	agent string,
	req *previewv1.PreviewRequest,
) (*previewv1.PreviewFlowLabelsResponse, error) {
	return transport.SendRequest[previewv1.PreviewFlowLabelsResponse](a.etcdTransport, agent, &cmdv1.PreviewFlowLabelsRequest{Request: req})
}

// PreviewHTTPRequests previews flow labels on a given agent.
//
// Handled by agentfunctions.PreviewHandler.
func (a Agents) PreviewHTTPRequests(
	agent string,
	req *previewv1.PreviewRequest,
) (*previewv1.PreviewHTTPRequestsResponse, error) {
	return transport.SendRequest[previewv1.PreviewHTTPRequestsResponse](a.etcdTransport, agent, &cmdv1.PreviewHTTPRequestsRequest{Request: req})
}

// GetAgents lists the agents registered on etcd under /peers/aperture-agent.
func (a Agents) GetAgents() ([]string, error) {
	re := regexp.MustCompile(`/peers/aperture-agent/[^/]+/`)
	if re == nil {
		return nil, fmt.Errorf("failed to compile regular expression")
	}

	resp, err := a.etcdClient.Client.KV.Get(context.Background(), "/peers/aperture-agent/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	agents := []string{}
	for _, kv := range resp.Kvs {
		agents = append(agents, re.ReplaceAllString(string(kv.Key), ""))
	}
	return agents, nil
}
