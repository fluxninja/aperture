package peers

import (
	"context"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PeerDiscoveryService is the implementation of peersv1.PeerDiscoveryServiceServer interface.
type PeerDiscoveryService struct {
	peersv1.UnimplementedPeerDiscoveryServiceServer
	PeerDiscovery *PeerDiscovery
}

// RegisterPeerDiscoveryService registers a service for peer discovery.
func RegisterPeerDiscoveryService(server *grpc.Server, pd *PeerDiscovery) {
	svc := &PeerDiscoveryService{
		PeerDiscovery: pd,
	}
	peersv1.RegisterPeerDiscoveryServiceServer(server, svc)
}

// GetPeers returns all peers.
func (svc *PeerDiscoveryService) GetPeers(ctx context.Context, _ *emptypb.Empty) (*peersv1.Peers, error) {
	pd := svc.PeerDiscovery
	return pd.GetPeers(), nil
}

// GetPeer returns a peer.
func (svc *PeerDiscoveryService) GetPeer(ctx context.Context, req *peersv1.PeerRequest) (*peersv1.Peer, error) {
	pd := svc.PeerDiscovery
	return pd.GetPeer(req.Address)
}
