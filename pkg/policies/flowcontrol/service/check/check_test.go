package check_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"google.golang.org/grpc/peer"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/entitycache"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
)

var (
	app *fx.App
	svc flowcontrolv1.FlowControlServiceServer
)

var _ = BeforeEach(func() {
	entities := entitycache.NewEntityCache()
	entities.Put(&entitycachev1.Entity{
		Prefix:    "",
		Uid:       "",
		IpAddress: hardCodedIPAddress,
		Name:      hardCodedEntityName,
		Services:  hardCodedServices,
	})
	app = platform.New(
		config.ModuleConfig{
			MergeConfig: map[string]interface{}{
				"flowcontrol": map[string]interface{}{
					"controller_addr": "",
					"policies_file":   "",
				},
				"sentrywriter": map[string]interface{}{
					"disabled": true,
				},
			},
		}.Module(),
		fx.Provide(agentinfo.ProvideAgentInfo),
		fx.Supply(entities),
		fx.Provide(servicegetter.FromEntityCache),
		fx.Provide(check.ProvideNopMetrics),
		fx.Provide(check.ProvideHandler),
		fx.Provide(flowcontrol.NewEngine),
		grpcclient.ClientConstructor{Name: "flowcontrol-grpc-client", ConfigKey: "flowcontrol.client.grpc"}.Annotate(),
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
	_ = app.Stop(ctx)
})

var _ = Describe("FlowControl Check", func() {
	When("client request comes in", func() {
		It("returns decision accepted response", func() {
			ctx := peer.NewContext(context.Background(), newFakeRpcPeer())
			resp, err := svc.Check(ctx, &flowcontrolv1.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect((resp.GetDecisionType())).To(Equal(flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED))
		})
	})
})

var (
	hardCodedIPAddress  = "1.2.3.4"
	hardCodedEntityName = "test-entity"
	hardCodedServices   = []string{"service1", "service2"}
)

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeRpcPeer() *peer.Peer {
	return &peer.Peer{Addr: fakeAddr("1.2.3.4:54321")}
}
