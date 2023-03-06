package checkhttp_test

import (
	"context"

	flowcontrolhttpv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/discovery/entities"
	"github.com/fluxninja/aperture/pkg/log"
	classification "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/checkhttp"
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

var _ = Describe("CheckHTTP handler", func() {
	var handler *checkhttp.Handler

	When("it is queried with a request", func() {
		BeforeEach(func() {
			alerter := alerts.NewSimpleAlerter(100)
			classifier := classification.NewClassificationEngine(
				status.NewRegistry(log.GetGlobalLogger(), alerter),
			)
			_, err := classifier.AddRules(context.TODO(), "test", &hardcodedRegoRules)
			Expect(err).NotTo(HaveOccurred())
			entities := entities.NewEntities()
			entities.Put(&entitiesv1.Entity{
				IpAddress: "1.2.3.4",
				Services:  []string{service1FlowSelector.ServiceSelector.Service},
			})
			handler = checkhttp.NewHandler(
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
			resp, err := handler.CheckHTTP(ctxWithIp, &flowcontrolhttpv1.CheckHTTPRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(code.Code(resp.GetStatus().GetCode())).To(Equal(code.Code_OK))
		})
		It("injects metadata", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs("control-point", "ingress"),
			)
			resp, err := handler.CheckHTTP(ctxWithIp, &flowcontrolhttpv1.CheckHTTPRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetDynamicMetadata()).NotTo(BeNil())
		})
	})
})

var service1FlowSelector = policylangv1.FlowSelector{
	ServiceSelector: &policylangv1.ServiceSelector{
		Service: "service1-demo-app.demoapp.svc.cluster.local",
	},
	FlowMatcher: &policylangv1.FlowMatcher{
		ControlPoint: "ingress",
	},
}

var hardcodedRegoRules = policysyncv1.ClassifierWrapper{
	Classifier: &policylangv1.Classifier{
		FlowSelector: &service1FlowSelector,
		Rules: map[string]*policylangv1.Rule{
			"destination": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
						package flowcontrol.checkhttp
						destination := v {
							v := input.destination.address
						}
					`,
						Query: "data.flowcontrol.checkhttp.destination",
					},
				},
			},
			"source": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
						package flowcontrol.checkhttp
						source := v {
							v := input.source.address
						}
					`,
						Query: "data.flowcontrol.checkhttp.source",
					},
				},
			},
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
