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
	peerDiscovery *PeerDiscovery
}

// RegisterPeerDiscoveryService registers a service for peer discovery.
func RegisterPeerDiscoveryService(server *grpc.Server, pd *PeerDiscovery) {
	svc := &PeerDiscoveryService{
		peerDiscovery: pd,
	}
	peersv1.RegisterPeerDiscoveryServiceServer(server, svc)
}

// GetPeers returns all the peer info that are added to PeerDiscovery.
func (pd *PeerDiscoveryService) GetPeers(ctx context.Context, _ *emptypb.Empty) (*peersv1.Peers, error) {
	return pd.peerDiscovery.Peers(), nil
}

// GetPeer returns the peer info in the PeerDiscovery with the given address.
func (pd *PeerDiscoveryService) GetPeer(ctx context.Context, req *peersv1.PeerRequest) (*peersv1.PeerInfo, error) {
	return pd.peerDiscovery.Peer(req.Address)
}

// GetPeerKeys returns all the peer keys that are added to PeerDiscovery.
func (pd *PeerDiscoveryService) GetPeerKeys(ctx context.Context, _ *emptypb.Empty) (*peersv1.PeerKeysResponse, error) {
	keys := pd.peerDiscovery.PeerKeys()
	return &peersv1.PeerKeysResponse{
		Keys: keys,
	}, nil
}

// AddPeer adds given peer to the PeerDiscovery.
func (pd *PeerDiscoveryService) AddPeer(ctx context.Context, req *peersv1.PeerInfo) (*emptypb.Empty, error) {
	pd.peerDiscovery.addPeer(req)
	return nil, nil
}

// RemovePeer removes the peer with the given address in the PeerDiscovery.
func (pd *PeerDiscoveryService) RemovePeer(ctx context.Context, req *peersv1.PeerRequest) (*emptypb.Empty, error) {
	pd.peerDiscovery.removePeer(req.Address)
	return nil, nil
}
