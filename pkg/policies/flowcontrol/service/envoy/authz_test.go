package envoy_test

import (
	"context"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/entitycache/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/servicegetter"
	"github.com/fluxninja/aperture/pkg/status"
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

func (s *AcceptingHandler) CheckWithValues(
	context.Context,
	[]string,
	string,
	map[string]string,
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
			classifier := classification.NewClassificationEngine(
				status.NewRegistry(log.GetGlobalLogger()),
			)
			_, err := classifier.AddRules(context.TODO(), "test", &hardcodedRegoRules)
			Expect(err).NotTo(HaveOccurred())
			entities := entitycache.NewEntityCache()
			entities.Put(&entitycachev1.Entity{
				IpAddress: "1.2.3.4",
				Services:  []string{service1Selector.ServiceSelector.Service},
			})
			handler = envoy.NewHandler(
				classifier,
				servicegetter.FromEntityCache(entities),
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
			resp, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code.Code(resp.GetStatus().GetCode())).To(Equal(code.Code_OK))
		})
		It("injects metadata", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("control-point", "ingress"),
			)
			resp, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetDynamicMetadata()).ShouldNot(BeNil())
		})
	})
})

var service1Selector = policylangv1.Selector{
	ServiceSelector: &policylangv1.ServiceSelector{
		Service: "service1-demo-app.demoapp.svc.cluster.local",
	},
	FlowSelector: &policylangv1.FlowSelector{
		ControlPoint: "ingress",
	},
}

var hardcodedRegoRules = policysyncv1.ClassifierWrapper{
	Classifier: &policylangv1.Classifier{
		Selector: &service1Selector,
		Rules: map[string]*policylangv1.Rule{
			"destination": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
						package envoy.authz
						destination := v {
							v := input.attributes.destination.address.socketAddress.address
						}
					`,
						Query: "data.envoy.authz.destination",
					},
				},
			},
			"source": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
						package envoy.authz
						source := v {
							v := input.attributes.destination.address.socketAddress.address
						}
					`,
						Query: "data.envoy.authz.source",
					},
				},
			},
		},
	},
	CommonAttributes: &policysyncv1.CommonAttributes{
		PolicyName:     "test",
		PolicyHash:     "test",
		ComponentIndex: 0,
	},
}

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeRpcPeer(ip string) *peer.Peer {
	return &peer.Peer{Addr: fakeAddr(ip + ":54321")}
}
