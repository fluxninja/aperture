package peers_test

import (
	"context"
	"errors"
	"time"

	peersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/peers/v1"
	"github.com/fluxninja/aperture/cmd/aperture-agent/agent"
	"github.com/fluxninja/aperture/pkg/config"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	app *fx.App
	svc peersv1.PeerDiscoveryServiceServer
)

var _ = BeforeEach(func() {
	app = platform.New(
		config.ModuleConfig{
			MergeConfig: map[string]interface{}{
				"sentrywriter": map[string]interface{}{
					"disabled": true,
				},
			},
		}.Module(),
		fx.Provide(agent.ProvidePeersPrefix),
		fx.Provide(peers.ProvideDummyPeerDiscoveryService),
		grpcclient.ClientConstructor{Name: "peers-grpc-client", ConfigKey: "peer_discovery.client.grpc"}.Annotate(),
		fx.Populate(&svc),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := app.Start(ctx)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := app.Stop(ctx)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Peers GetPeers", func() {
	pd, err := peers.NewPeerDiscovery("peers", nil, peers.PeerWatchers{})
	Expect(err).NotTo(HaveOccurred())
	for _, peerinfo := range hardCodedPeers.PeerInfos {
		pd.AddPeer(peerinfo)
	}

	When("client request comes in", func() {
		It("returns all the peer info that are added to peer discovery", func() {
			ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
			resp, err := svc.GetPeers(ctx, &emptypb.Empty{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(hardCodedPeers))
		})
	})
})

var _ = Describe("Peers GetPeer", func() {
	pd, err := peers.NewPeerDiscovery("peers", nil, peers.PeerWatchers{})
	Expect(err).NotTo(HaveOccurred())
	peerInfo := &peersv1.PeerInfo{
		Address:  hardCodedIPAddress,
		Hostname: hardCodedHostName,
		Services: hardCodedServices,
	}
	pd.AddPeer(peerInfo)

	When("client request with peer address comes in", func() {
		It("returns the peer info that matches the provided peer address", func() {
			ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
			resp, err := svc.GetPeer(ctx, &peersv1.PeerRequest{Address: "1.2.3.4:54321"})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(peerInfo))
		})
	})

	When("empty client request comes in", func() {
		It("returns a peer not found error", func() {
			ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
			resp, err := svc.GetPeer(ctx, &peersv1.PeerRequest{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("peer not found")))
			Expect(resp).To(BeNil())
		})
	})

	When("client request with non matching peer address comes in", func() {
		It("returns a peer not found error", func() {
			ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
			resp, err := svc.GetPeer(ctx, &peersv1.PeerRequest{Address: "1.2.3.4:50000"})
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
	hardCodedPeers = &peersv1.Peers{
		PeerInfos: []*peersv1.PeerInfo{
			{
				Address:  "1.2.3.4:54321",
				Hostname: hardCodedHostName + "1",
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
