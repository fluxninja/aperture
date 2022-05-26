package peers

import peersv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/peers/v1"

// PeerWatcher is used for tracking changes to peers.
type PeerWatcher interface {
	OnPeerAdded(peer *peersv1.PeerInfo)
	OnPeerRemoved(peer *peersv1.PeerInfo)
}

// PeerWatchers is a collection of PeerWatcher.
type PeerWatchers []PeerWatcher

// OnPeerAdded calls OnPeerAdded for each PeerWatcher in the collection.
func (pw PeerWatchers) OnPeerAdded(peer *peersv1.PeerInfo) {
	for _, watcher := range pw {
		watcher.OnPeerAdded(peer)
	}
}

// OnPeerRemoved calls OnPeerRemoved for each PeerWatcher in the collection.
func (pw PeerWatchers) OnPeerRemoved(peer *peersv1.PeerInfo) {
	for _, watcher := range pw {
		watcher.OnPeerRemoved(peer)
	}
}
