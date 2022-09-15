package peers

import (
	"context"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	"google.golang.org/grpc"
)

// PeerDiscoveryService is the implementation of peersv1.PeerDiscoveryServiceServer interface.
type PeerDiscoveryService struct {
	peersv1.UnimplementedPeerDiscoveryServiceServer
	peerDiscovery *PeerDiscovery
}

// RegisterPeerDiscoveryService registers a service for peer discovery.
func RegisterPeerDiscoveryService(server *grpc.Server, pd *PeerDiscovery) {
	svc := &PeerDiscoveryService{
		peerDiscovery: pd,
	}
	peersv1.RegisterPeerDiscoveryServiceServer(server, svc)
}

// GetPeers returns a matching peer info and peers if provided address matches peer;
// otherwise, it returns all the peer info that are added to PeerDiscovery.
func (svc *PeerDiscoveryService) GetPeers(ctx context.Context, req *peersv1.PeersRequest) (*peersv1.Peers, error) {
	pd := svc.peerDiscovery
	for _, address := range req.Address {
		if address == "" {
			continue
		}
		peerInfo, _ := pd.GetPeer(address)
		if peerInfo == nil {
			break
		} else {
			return &peersv1.Peers{
				PeerInfo: peerInfo,
				Peers:    pd.peers,
			}, nil
		}
	}
	return pd.GetPeers(), nil
}
