package servicegetter

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc/peer"

	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
)

// ServiceGetter can be used to query services based on client context.
type ServiceGetter interface {
	ServicesFromContext(ctx context.Context) []string
}

// FromEntityCache creates a new EntityCache-powered ServiceGetter.
func FromEntityCache(ec *entitycache.EntityCache) ServiceGetter {
	return &ecServiceGetter{entityCache: ec}
}

// NewEmpty creates a new ServiceGetter that always returns nil.
func NewEmpty() ServiceGetter { return emptyServiceGetter{} }

type ecServiceGetter struct {
	entityCache    *entitycache.EntityCache
	ecHasDiscovery bool
	metrics        *Metrics
}

// ServicesFromContext returns list of services associated with IP extracted from context
//
// The returned list of services depends only on state of entityCache.
// However, emitted warnings will depend on whether service discovery is enabled or not.
func (sg *ecServiceGetter) ServicesFromContext(ctx context.Context) []string {
	svcs, ok := sg.servicesFromContext(ctx)
	sg.metrics.inc(ok)
	return svcs
}

func (sg *ecServiceGetter) servicesFromContext(ctx context.Context) (svcs []string, ok bool) {
	rpcPeer, peerExists := peer.FromContext(ctx)
	if !peerExists {
		if sg.ecHasDiscovery {
			log.Bug().Msg("cannot get client info from context")
		}
		return nil, false
	}

	tcpAddr, isTCPAddr := rpcPeer.Addr.(*net.TCPAddr)
	if !isTCPAddr {
		if sg.ecHasDiscovery {
			log.Bug().Msg("client addr is not TCP")
		}
		return nil, false

	}

	clientIP := tcpAddr.IP.String()
	entity, err := sg.entityCache.GetByIP(clientIP)
	if err != nil {
		if sg.ecHasDiscovery {
			log.Sample(noEntitySampler).Warn().Err(err).Str("clientIP", clientIP).
				Msg("cannot get services")
		}
		return nil, false
	}

	return entity.Services, true
}

var noEntitySampler = log.NewRatelimitingSampler()

// FxIn are FX arguments to ProvideFromEntityCache.
type FxIn struct {
	fx.In
	Lifecycle      fx.Lifecycle
	EntityCache    *entitycache.EntityCache
	EntityTrackers *entitycache.EntityTrackers
	Metrics        *Metrics `optional:"true"`
}

// ProvideFromEntityCache provides an EntityCache-powered ServiceGetter.
func ProvideFromEntityCache(in FxIn) ServiceGetter {
	sg := &ecServiceGetter{
		entityCache: in.EntityCache,
		metrics:     in.Metrics,
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// Checking this flag on OnStart so that all registrations done in
			// provide/invoke stage would be visible.
			sg.ecHasDiscovery = in.EntityTrackers.HasDiscovery()
			return nil
		},
		OnStop: func(context.Context) error { return nil },
	})

	return sg
}

type emptyServiceGetter struct{}

// ServicesFromContext implements ServiceGetter interface.
func (sg emptyServiceGetter) ServicesFromContext(ctx context.Context) []string { return nil }
