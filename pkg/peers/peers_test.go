package peers_test

import (
	"context"
	"errors"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ = Describe("Peers GetPeers", func() {
	ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
	When("client request comes in", func() {
		It("returns all the peer info that are added to peer discovery", func() {
			resp, err := svc.GetPeers(ctx, &emptypb.Empty{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(hardCodedPeers))
		})
	})

})

var _ = Describe("Peers GetPeer", func() {
	ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
	When("client request with peer address comes in", func() {
		It("returns the peer info that matches the provided peer address", func() {
			resp, err := svc.GetPeer(ctx, &peersv1.PeerRequest{Address: "1.2.3.4:54321"})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(hardCodedPeerInfo))
		})
	})

	When("empty client request comes in", func() {
		It("returns a peer not found error", func() {
			resp, err := svc.GetPeer(ctx, &peersv1.PeerRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("peer not found")))
			Expect(resp).To(BeNil())
		})
	})

	When("client request with non matching peer address comes in", func() {
		It("returns a peer not found error", func() {
			resp, err := svc.GetPeer(ctx, &peersv1.PeerRequest{Address: "1.2.3.4:12345"})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("peer not found")))
			Expect(resp).To(BeNil())
		})
	})
})

var (
	hardCodedIPAddress = "1.2.3.4:54321"
	hardCodedHostName  = "peers-hostname-info"
	hardCodedServices  = map[string]string{
		"service1": "peers1",
		"service2": "peers2",
		"service3": "peers3",
	}
	hardCodedPeerInfo = &peersv1.PeerInfo{
		Address:  hardCodedIPAddress,
		Hostname: hardCodedHostName,
		Services: hardCodedServices,
	}
	hardCodedPeers = &peersv1.Peers{
		PeerInfos: []*peersv1.PeerInfo{
			{
				Address:  hardCodedIPAddress,
				Hostname: hardCodedHostName,
				Services: hardCodedServices,
			},
			{
				Address:  "1.2.3.4:54322",
				Hostname: hardCodedHostName + "2",
				Services: hardCodedServices,
			},
			{
				Address:  "1.2.3.4:54323",
				Hostname: hardCodedHostName + "3",
				Services: hardCodedServices,
			},
		},
	}
)

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeRpcPeer() *peer.Peer {
	return &peer.Peer{Addr: fakeAddr("1.2.3.4:54321")}
}
