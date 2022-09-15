package peers_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/cmd/aperture-agent/agent"
	"github.com/fluxninja/aperture/pkg/config"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/utils"
)

func TestPeers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Peers Suite")
}

var (
	l   *utils.GoLeakDetector
	app *fx.App
	svc peers.PeerDiscoveryService
)

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()

	pd, err := peers.NewPeerDiscovery("test", nil, nil)
	Expect(err).ToNot(HaveOccurred())
	for _, peerinfo := range hardCodedPeers.Peers {
		pd.AddPeer(peerinfo)
	}

	app = platform.New(
		config.ModuleConfig{
			MergeConfig: map[string]interface{}{
				"sentrywriter": map[string]interface{}{
					"disabled": true,
				},
			},
		}.Module(),
		fx.Supply(pd),
		fx.Provide(agent.ProvidePeersPrefix),
		fx.Provide(peers.ProvideDummyPeerDiscoveryService),
		grpcclient.ClientConstructor{Name: "peers-grpc-client", ConfigKey: "peer_discovery.client.grpc"}.Annotate(),
		fx.Populate(&svc),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = app.Start(ctx)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := app.Stop(ctx)
	Expect(err).NotTo(HaveOccurred())

	err = l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
