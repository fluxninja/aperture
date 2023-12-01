// Server-side for handling agent functions
package agents

import (
	"context"
	"fmt"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	proto "google.golang.org/protobuf/proto"

	cmdv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/preview/v1"
	peersv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/peers/v1"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

// Module is fx module for controlling Agents on controller side.
var Module = fx.Provide(NewAgents)

// Agents wraps rpc.Clients where clients are agents and provides wrapper
//
// Agents wraps functions registered in agentfunctions, types should match.
type Agents struct {
	etcdTransport *transport.EtcdTransportClient
	etcdClient    *etcdclient.Client
}

// NewAgents wraps Clients with Agent-specific function wrappers.
func NewAgents(transport *transport.EtcdTransportClient, client *etcdclient.Client) Agents {
	return Agents{
		etcdTransport: transport,
		etcdClient:    client,
	}
}

// ListFlowControlPoints lists control points of all agents.
//
// Handled by agentfunctions.ControlPointsHandler.
func (a Agents) ListFlowControlPoints(ctx context.Context) ([]transport.Result[*cmdv1.ListFlowControlPointsAgentResponse], error) {
	var req cmdv1.ListFlowControlPointsRequest
	agents, err := a.GetAgents(ctx)
	if err != nil {
		return nil, err
	}
	return transport.SendRequests[cmdv1.ListFlowControlPointsAgentResponse](ctx, a.etcdTransport, flattenAgents(agents), &req)
}

// ListAutoScaleControlPoints lists auto-scale control points of all agents.
func (a Agents) ListAutoScaleControlPoints(ctx context.Context) ([]transport.Result[*cmdv1.ListAutoScaleControlPointsAgentResponse], error) {
	var req cmdv1.ListAutoScaleControlPointsRequest
	agents, err := a.GetAgents(ctx)
	if err != nil {
		return nil, err
	}
	return transport.SendRequests[cmdv1.ListAutoScaleControlPointsAgentResponse](ctx, a.etcdTransport, flattenAgents(agents), &req)
}

// ListDiscoveryEntities lists discovery entities.
func (a Agents) ListDiscoveryEntities(ctx context.Context, agentGroup string) ([]transport.Result[*cmdv1.ListDiscoveryEntitiesAgentResponse], error) {
	var req cmdv1.ListDiscoveryEntitiesRequest
	agents, err := a.GetAgentsForGroup(ctx, agentGroup)
	if err != nil {
		return nil, err
	}
	return transport.SendRequests[cmdv1.ListDiscoveryEntitiesAgentResponse](ctx, a.etcdTransport, agents, &req)
}

// ListDiscoveryEntity lists discovery entity by ip address or name.
func (a Agents) ListDiscoveryEntity(ctx context.Context, req *cmdv1.ListDiscoveryEntityRequest) (*cmdv1.ListDiscoveryEntityAgentResponse, error) {
	agents, err := a.GetAgents(ctx)
	if err != nil {
		return nil, err
	}

	return transport.SendRequest[cmdv1.ListDiscoveryEntityAgentResponse](ctx, a.etcdTransport, flattenAgents(agents)[0], req)
}

// PreviewFlowLabels previews flow labels on a given agent.
//
// Handled by agentfunctions.PreviewHandler.
func (a Agents) PreviewFlowLabels(
	ctx context.Context,
	agent string,
	req *previewv1.PreviewRequest,
) (*previewv1.PreviewFlowLabelsResponse, error) {
	return transport.SendRequest[previewv1.PreviewFlowLabelsResponse](ctx, a.etcdTransport, agent, &cmdv1.PreviewFlowLabelsRequest{Request: req})
}

// PreviewHTTPRequests previews flow labels on a given agent.
//
// Handled by agentfunctions.PreviewHandler.
func (a Agents) PreviewHTTPRequests(
	ctx context.Context,
	agent string,
	req *previewv1.PreviewRequest,
) (*previewv1.PreviewHTTPRequestsResponse, error) {
	return transport.SendRequest[previewv1.PreviewHTTPRequestsResponse](ctx, a.etcdTransport, agent, &cmdv1.PreviewHTTPRequestsRequest{Request: req})
}

// GetAgents lists the agents registered on etcd under /peers/aperture-agent.
func (a Agents) GetAgents(ctx context.Context) (map[string][]string, error) {
	resp, err := a.etcdClient.Get(ctx, paths.AgentPeerPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	agents := map[string][]string{}
	for _, kv := range resp.Kvs {
		// Extract agent from kv
		var peer peersv1.Peer
		err = proto.Unmarshal(kv.Value, &peer)
		if err != nil {
			return nil, err
		}

		// Extract agent-group from kv
		// The etcd value currently doesn't have the agent-group in it, so need to extract it from the key
		parts := strings.SplitN(strings.TrimPrefix(string(kv.Key), "/peers/aperture-agent/"), "/", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("error fetching agent-group from etcd key %s", kv.Key)
		}
		agentGroup := parts[0]

		if _, ok := agents[agentGroup]; !ok {
			agents[agentGroup] = []string{peer.Hostname}
		} else {
			agents[agentGroup] = append(agents[agentGroup], string(peer.Hostname))
		}
	}
	return agents, nil
}

// GetAgentsForGroup lists the agents under an agent-group registered on etcd under /peers/aperture-agent.
func (a Agents) GetAgentsForGroup(ctx context.Context, agentGroup string) ([]string, error) {
	agents, err := a.GetAgents(ctx)
	if err != nil {
		return nil, err
	}

	if agentGroup == "" {
		agentGroup = "default"
	}

	if _, ok := agents[agentGroup]; !ok {
		return nil, fmt.Errorf("no agents found for agent-group %s", agentGroup)
	}
	return agents[agentGroup], nil
}

func flattenAgents(agents map[string][]string) []string {
	var flattened []string
	for _, agent := range agents {
		flattened = append(flattened, agent...)
	}
	return flattened
}
