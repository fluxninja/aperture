package authz_test

import (
	"context"

	ext_authz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/peer"

	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/classification/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/authz"
	"github.com/fluxninja/aperture/pkg/classification"
	"github.com/fluxninja/aperture/pkg/flowcontrol"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

var (
	ctx        context.Context
	cancel     context.CancelFunc
	classifier *classification.Classifier
	handler    *authz.Handler
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

type AcceptingHandler struct {
	flowcontrol.HandlerWithValues
}

func (s *AcceptingHandler) CheckWithValues(
	context.Context,
	selectors.ControlPoint,
	[]services.ServiceID,
	selectors.Labels,
) *flowcontrolv1.CheckResponse {
	resp := &flowcontrolv1.CheckResponse{
		DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
	}
	return resp
}

var _ = Describe("Authorization handler", func() {
	When("it is queried with a request", func() {
		BeforeEach(func() {
			classifier = classification.New()
			_, err := classifier.AddRules(context.TODO(), "test", &hardcodedRegoRules)
			Expect(err).NotTo(HaveOccurred())
			handler = authz.NewHandler(classifier, nil, &AcceptingHandler{})
		})
		It("returns ok response", func() {
			Eventually(func(g Gomega) {
				ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer())
				resp, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(code.Code(resp.GetStatus().GetCode())).To(Equal(code.Code_OK))
			}).Should(Succeed())
		})
		It("injects metadata", func() {
			Eventually(func(g Gomega) {
				ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer())
				resp, err := handler.Check(ctxWithIp, &ext_authz.CheckRequest{})
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(resp.GetDynamicMetadata()).ShouldNot(BeNil())
			}).Should(Succeed())
		})
	})
})

var demo1Selector = policylangv1.Selector{
	Service: "demo1-demo-app.demoapp.svc.cluster.local",
	ControlPoint: &policylangv1.ControlPoint{
		Controlpoint: &policylangv1.ControlPoint_Traffic{
			Traffic: "ingress",
		},
	},
}

var hardcodedRegoRules = classificationv1.Classifier{
	Selector: &demo1Selector,
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
}

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

func newFakeRpcPeer() *peer.Peer {
	return &peer.Peer{Addr: fakeAddr("1.2.3.4:54321")}
}
