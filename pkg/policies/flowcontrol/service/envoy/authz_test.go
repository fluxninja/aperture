package envoy_test

import (
	"context"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	grpc_status "google.golang.org/grpc/status"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/entitycache/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy"
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
	selectors.ControlPoint,
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
			handler = envoy.NewHandler(classifier, entities, &AcceptingHandler{})
		})
		It("returns ok response", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			// add "traffic-direction" header to ctx
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("traffic-direction", "INBOUND"),
			)
			resp, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code.Code(resp.GetStatus().GetCode())).To(Equal(code.Code_OK))
		})
		It("injects metadata", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("traffic-direction", "INBOUND"),
			)
			resp, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetDynamicMetadata()).ShouldNot(BeNil())
		})
		It("handles entity cache lookup failure", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("9.9.9.9"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("traffic-direction", "INBOUND"),
			)
			_, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
			Expect(err).To(HaveOccurred())
			status, _ := grpc_status.FromError(err)
			Expect(status.Code()).To(Equal(codes.NotFound))
		})
	})
})

var service1Selector = policylangv1.Selector{
	ServiceSelector: &policylangv1.ServiceSelector{
		Service: "service1-demo-app.demoapp.svc.cluster.local",
	},
	FlowSelector: &policylangv1.FlowSelector{
		ControlPoint: &policylangv1.ControlPoint{
			Controlpoint: &policylangv1.ControlPoint_Traffic{
				Traffic: "ingress",
			},
		},
	},
}

var hardcodedRegoRules = policysyncv1.ClassifierWrapper{
	Classifier: &classificationv1.Classifier{
		Selector: &service1Selector,
		Rules: map[string]*classificationv1.Rule{
			"destination": {
				Source: &classificationv1.Rule_Rego_{
					Rego: &classificationv1.Rule_Rego{
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
				Source: &classificationv1.Rule_Rego_{
					Rego: &classificationv1.Rule_Rego{
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
