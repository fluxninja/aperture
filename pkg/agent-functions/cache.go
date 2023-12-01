package agentfunctions

import (
	"context"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CacheHandler is a handler for cache-family of functions.
type CacheHandler struct {
	cache      iface.Cache
	agentGroup string
}

// NewCacheHandler creates a new CacheHandler.
func NewCacheHandler(
	cache iface.Cache,
	agentInfo *agentinfo.AgentInfo,
) (*CacheHandler, error) {
	return &CacheHandler{
		cache:      cache,
		agentGroup: agentInfo.GetAgentGroup(),
	}, nil
}

// CacheLookup looks up given keys in the cache.
func (h *CacheHandler) CacheLookup(ctx context.Context, req *cmdv1.GlobalCacheLookupRequest) (*flowcontrolv1.CacheLookupResponse, error) {
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	return h.cache.Lookup(ctx, req.Request), nil
}

// CacheUpsert inserts or updates given keys into the cache.
func (h *CacheHandler) CacheUpsert(ctx context.Context, req *cmdv1.GlobalCacheUpsertRequest) (*flowcontrolv1.CacheUpsertResponse, error) {
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	return h.cache.Upsert(ctx, req.Request), nil
}

// CacheDelete deletes given keys from the cache.
func (h *CacheHandler) CacheDelete(ctx context.Context, req *cmdv1.GlobalCacheDeleteRequest) (*flowcontrolv1.CacheDeleteResponse, error) {
	if req.Request == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	return h.cache.Delete(ctx, req.Request), nil
}

// RegisterCacheHandlers registers cache handler functions in handler registry.
func RegisterCacheHandlers(handler *CacheHandler, t *transport.EtcdTransportServer) error {
	err := transport.RegisterFunction(t, handler.CacheUpsert)
	if err != nil {
		return err
	}
	err = transport.RegisterFunction(t, handler.CacheLookup)
	if err != nil {
		return err
	}
	err = transport.RegisterFunction(t, handler.CacheDelete)
	if err != nil {
		return err
	}
	return nil
}
