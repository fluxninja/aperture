package peers

import peersv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/peers/v1"

// PeerWatcher is used for tracking changes to peers.
type PeerWatcher interface {
	OnPeerAdded(peer *peersv1.Peer)
	OnPeerRemoved(peer *peersv1.Peer)
}

// PeerWatchers is a collection of PeerWatcher.
type PeerWatchers []PeerWatcher

// OnPeerAdded calls OnPeerAdded for each PeerWatcher in the collection.
func (pw PeerWatchers) OnPeerAdded(peer *peersv1.Peer) {
	for _, watcher := range pw {
		watcher.OnPeerAdded(peer)
	}
}

// OnPeerRemoved calls OnPeerRemoved for each PeerWatcher in the collection.
func (pw PeerWatchers) OnPeerRemoved(peer *peersv1.Peer) {
	for _, watcher := range pw {
		watcher.OnPeerRemoved(peer)
	}
}
