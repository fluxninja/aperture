package peers

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	peersv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/peers/v1"
)

var (
	pd  *PeerDiscovery
	err error
)

var _ = Describe("Peers", func() {
	BeforeEach(func() {
		pd, err = NewPeerDiscovery("", nil, nil)
		Expect(err).ToNot(HaveOccurred())
		for _, peer := range hardCodedPeers.Peers {
			pd.addPeer(peer)
		}
	})

	Context("GetPeers", func() {
		It("returns all the peers that are added to peer discovery", func() {
			resp := pd.GetPeers()
			Expect(resp.SelfPeer).To(Equal(hardCodedPeers.SelfPeer))
			Expect(resp.Peers).To(Equal(hardCodedPeers.Peers))
			hardCodedPeersCopy := proto.Clone(hardCodedPeers).(*peersv1.Peers)
			Expect(resp).To(Equal(hardCodedPeersCopy))
		})
		It("returns all the peers except the removed peers", func() {
			pd.removePeer(hardCodedIPAddress1)
			pd.removePeer(hardCodedIPAddress2)
			resp := pd.GetPeers()
			hardCodedPeers3Copy := proto.Clone(createPeers(hardCodedIPAddress3, hardCodedPeer3)).(*peersv1.Peers)
			Expect(resp).To(Equal(hardCodedPeers3Copy))

			pd.addPeer(hardCodedPeer1)
			pd.removePeer(hardCodedIPAddress3)
			resp = pd.GetPeers()
			hardCodedPeers1Copy := proto.Clone(createPeers(hardCodedIPAddress1, hardCodedPeer1)).(*peersv1.Peers)
			Expect(resp).To(Equal(hardCodedPeers1Copy))
		})
		It("returns all the peer keys that are added to peer discovery", func() {
			resp := pd.GetPeerKeys()
			Expect(resp).To(ConsistOf([]string{hardCodedIPAddress1, hardCodedIPAddress2, hardCodedIPAddress3}))

			pd.removePeer(hardCodedIPAddress2)
			resp = pd.GetPeerKeys()
			Expect(resp).To(ConsistOf([]string{hardCodedIPAddress1, hardCodedIPAddress3}))
		})
		It("returns all the peer keys except the removed peers", func() {
			pd.removePeer(hardCodedIPAddress2)
			pd.removePeer(hardCodedIPAddress3)
			resp := pd.GetPeerKeys()
			Expect(resp).To(Equal([]string{hardCodedIPAddress1}))
		})
	})

	Context("GetPeer", func() {
		It("returns `peer not found` error", func() {
			resp, err := pd.GetPeer("1.2.3.4:55555")
			Expect(resp).To(BeNil())
			Expect(err).To(Equal(errors.New("peer not found")))
		})
		It("returns a peer with peer address 3.4.5.6:54321 then `peer not found` error after removing the peer", func() {
			resp, err := pd.GetPeer(hardCodedIPAddress3)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(hardCodedPeer3))

			pd.removePeer(hardCodedIPAddress3)
			_, err = pd.GetPeer(hardCodedIPAddress3)
			Expect(err).To(Equal(errors.New("peer not found")))
		})
	})
})

var (
	// hardCodedPeer
	hardCodedIPAddress1 = "1.2.3.4:54321"
	hardCodedIPAddress2 = "2.3.4.5:54321"
	hardCodedIPAddress3 = "3.4.5.6:54321"
	hardCodedHostName   = "peers-hostname-info"
	hardCodedPeer1      = &peersv1.Peer{
		Address:  hardCodedIPAddress1,
		Hostname: hardCodedHostName + "1",
	}
	hardCodedPeer2 = &peersv1.Peer{
		Address:  hardCodedIPAddress2,
		Hostname: hardCodedHostName + "2",
	}
	hardCodedPeer3 = &peersv1.Peer{
		Address:  hardCodedIPAddress3,
		Hostname: hardCodedHostName + "3",
	}

	// hardCodedPeers
	hardCodedPeers = &peersv1.Peers{
		SelfPeer: &peersv1.Peer{},
		Peers: map[string]*peersv1.Peer{
			hardCodedIPAddress1: hardCodedPeer1,
			hardCodedIPAddress2: hardCodedPeer2,
			hardCodedIPAddress3: hardCodedPeer3,
		},
	}
)

func createPeers(address string, peer *peersv1.Peer) *peersv1.Peers {
	return &peersv1.Peers{
		SelfPeer: &peersv1.Peer{},
		Peers: map[string]*peersv1.Peer{
			address: peer,
		},
	}
}
