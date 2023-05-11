package envoy_test

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	entitiesv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/discovery/entities/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/agentinfo"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/discovery/entities"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	classification "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/envoy"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/servicegetter"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

var _ = BeforeEach(func() {
	// Disable logs for cleaner tests output
	log.SetGlobalLevel(log.FatalLevel)
	ctx, cancel = context.WithCancel(context.Background())
})

var _ = AfterEach(func() {
	if cancel != nil {
		cancel()
	}
})

type AcceptingHandler struct{}

func (s *AcceptingHandler) CheckRequest(
	context.Context,
	iface.RequestContext,
) *flowcontrolv1.CheckResponse {
	resp := &flowcontrolv1.CheckResponse{
		DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
	}
	return resp
}

var _ = Describe("Authorization handler", func() {
	var handler *envoy.Handler

	When("it is queried with a request", func() {
		BeforeEach(func() {
			alerter := alerts.NewSimpleAlerter(100)
			classifier := classification.NewClassificationEngine(
				agentinfo.NewAgentInfo("testGroup"),
				status.NewRegistry(log.GetGlobalLogger(), alerter),
			)
			_, err := classifier.AddRules(context.TODO(), "test", &hardcodedRegoRules)
			Expect(err).NotTo(HaveOccurred())
			entities := entities.NewEntities()
			entities.Put(&entitiesv1.Entity{
				IpAddress: "1.2.3.4",
				Services:  []string{service1Selector.Service},
			})
			handler = envoy.NewHandler(
				classifier,
				servicegetter.FromEntities(entities),
				&AcceptingHandler{},
			)
		})
		It("returns ok response", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			// add "control-point" header to ctx
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("control-point", "ingress"),
			)
			resp, err := handler.Check(ctxWithIp, &authv3.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code.Code(resp.GetStatus().GetCode())).To(Equal(code.Code_OK))
		})
		It("injects metadata", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("control-point", "ingress"),
			)
			resp, err := handler.Check(ctxWithIp, &authv3.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetDynamicMetadata()).NotTo(BeNil())
		})
	})
})

var service1Selector = &policylangv1.Selector{
	ControlPoint: "ingress",
	Service:      "service1-demo-app.demoapp.svc.cluster.local",
	AgentGroup:   "testGroup",
}

var hardcodedRegoRules = policysyncv1.ClassifierWrapper{
	Classifier: &policylangv1.Classifier{
		Selectors: []*policylangv1.Selector{
			service1Selector,
		},
		Rego: &policylangv1.Rego{
			Labels: map[string]*policylangv1.Rego_LabelProperties{
				"destination": {
					Telemetry: true,
				},
				"source": {
					Telemetry: true,
				},
			},
			Module: `
				package envoy.authz
				destination := v {
					v := input.attributes.destination.address.socketAddress.address
				}
				source := v {
					v := input.attributes.source.address.socketAddress.address
				}
			`,
		},
	},
	ClassifierAttributes: &policysyncv1.ClassifierAttributes{
		PolicyName:      "test",
		PolicyHash:      "test",
		ClassifierIndex: 0,
	},
}

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeRpcPeer(ip string) *peer.Peer {
	return &peer.Peer{Addr: fakeAddr(ip + ":54321")}
}
