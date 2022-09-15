package peers_test

import (
	"context"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/peer"
)

var _ = Describe("Peers GetPeers", func() {
	ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
	When("empty client request comes in", func() {
		It("returns all the peer info that are added to peer discovery", func() {
			resp, err := svc.GetPeers(ctx, &peersv1.PeersRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(hardCodedPeers))
		})
	})
	When("client request with address that does not exist in peer discovery comes in", func() {
		It("returns all the peer info that are added to peer discovery", func() {
			resp, err := svc.GetPeers(ctx, &peersv1.PeersRequest{Address: []string{"1.2.3.4:55555"}})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(hardCodedPeers))
		})
	})
	When("client request with address that exists in peer discovery comes in", func() {
		It("returns specific peer info matches with 1.2.3.4:54321", func() {
			resp, err := svc.GetPeers(ctx, &peersv1.PeersRequest{Address: []string{"1.2.3.4:54321"}})
			Expect(err).NotTo(HaveOccurred())
			specificPeerInfo := hardCodedPeers
			specificPeerInfo.PeerInfo = hardCodedPeerInfo1
			Expect(resp).To(Equal(specificPeerInfo))
		})
		It("returns specific peer info matches with 3.4.5.6:54321", func() {
			resp, err := svc.GetPeers(ctx, &peersv1.PeersRequest{Address: []string{"3.4.5.6:54321"}})
			Expect(err).NotTo(HaveOccurred())
			specificPeerInfo := hardCodedPeers
			specificPeerInfo.PeerInfo = hardCodedPeerInfo3
			Expect(resp).To(Equal(specificPeerInfo))
		})
	})
})

var (
	// hardCodedPeerInfo
	hardCodedIPAddress = "1.2.3.4:54321"
	hardCodedHostName  = "peers-hostname-info"
	hardCodedServices  = map[string]string{
		"peers1": "service1",
		"peers2": "service2",
		"peers3": "service3",
	}
	hardCodedPeerInfo1 = &peersv1.PeerInfo{
		Address:  hardCodedIPAddress,
		Hostname: hardCodedHostName + "1",
		Services: hardCodedServices,
	}
	hardCodedPeerInfo2 = &peersv1.PeerInfo{
		Address:  "2.3.4.5:54321",
		Hostname: hardCodedHostName + "2",
		Services: hardCodedServices,
	}
	hardCodedPeerInfo3 = &peersv1.PeerInfo{
		Address:  "3.4.5.6:54321",
		Hostname: hardCodedHostName + "3",
		Services: hardCodedServices,
	}

	// hardCodedPeers
	hardCodedPeers = &peersv1.Peers{
		Peers: map[string]*peersv1.PeerInfo{
			"1.2.3.4:54321": hardCodedPeerInfo1,
			"2.3.4.5:54321": hardCodedPeerInfo2,
			"3.4.5.6:54321": hardCodedPeerInfo3,
		},
	}
)

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeRpcPeer() *peer.Peer {
	return &peer.Peer{Addr: fakeAddr("1.2.3.4:54321")}
}
