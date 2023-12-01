package peers

import (
	"context"

	peersv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/peers/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PeerDiscoveryService is the implementation of peersv1.PeerDiscoveryServiceServer interface.
type PeerDiscoveryService struct {
	peersv1.UnimplementedPeerDiscoveryServiceServer
	peerDiscovery *PeerDiscovery
}

// RegisterPeersServiceIn bundles and annotates parameters.
type RegisterPeersServiceIn struct {
	fx.In
	Server        *grpc.Server `name:"default"`
	PeerDiscovery *PeerDiscovery
}

// RegisterPeerDiscoveryService registers a service for peer discovery.
func RegisterPeerDiscoveryService(in RegisterPeersServiceIn) {
	svc := &PeerDiscoveryService{
		peerDiscovery: in.PeerDiscovery,
	}
	peersv1.RegisterPeerDiscoveryServiceServer(in.Server, svc)
}

// GetPeers returns all peers.
func (svc *PeerDiscoveryService) GetPeers(ctx context.Context, _ *emptypb.Empty) (*peersv1.Peers, error) {
	pd := svc.peerDiscovery
	return pd.GetPeers(), nil
}

// GetPeer returns a peer.
func (svc *PeerDiscoveryService) GetPeer(ctx context.Context, req *peersv1.PeerRequest) (*peersv1.Peer, error) {
	pd := svc.peerDiscovery
	return pd.GetPeer(req.Address)
}
