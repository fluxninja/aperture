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

// GetPeers returns all the peer info that are added to PeerDiscovery.
func (pd *PeerDiscoveryService) GetPeers(ctx context.Context, _ *emptypb.Empty) (*peersv1.Peers, error) {
	return pd.PeerDiscovery.GetPeers(), nil
}

// GetPeer returns the peer info in the PeerDiscovery with the given address.
func (pd *PeerDiscoveryService) GetPeer(ctx context.Context, req *peersv1.PeerRequest) (*peersv1.PeerInfo, error) {
	return pd.PeerDiscovery.GetPeer(req.Address)
}
