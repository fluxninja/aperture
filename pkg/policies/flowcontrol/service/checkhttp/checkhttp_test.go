package checkhttp_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/durationpb"

	entitiesv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/discovery/entities/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	flowcontrolhttpv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/discovery/entities"
	"github.com/fluxninja/aperture/v2/pkg/log"
	classification "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/checkhttp"
	"github.com/fluxninja/aperture/v2/pkg/policies/mocks"
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

func acceptedResponse() *flowcontrolv1.CheckResponse {
	return &flowcontrolv1.CheckResponse{
		DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
	}
}

func rejectedResponse() *flowcontrolv1.CheckResponse {
	return &flowcontrolv1.CheckResponse{
		DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED,
	}
}

var _ = Describe("CheckHTTP handler", func() {
	var (
		ctrl         *gomock.Controller
		checkHandler *mocks.MockHandlerWithValues
		handler      *checkhttp.Handler
	)

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
			entities.PutForTest(&entitiesv1.Entity{
				IpAddress: "1.2.3.4",
				Services:  []string{service1Selector.Service},
			})
			ctrl = gomock.NewController(GinkgoT())
			checkHandler = mocks.NewMockHandlerWithValues(ctrl)
			handler = checkhttp.NewHandler(
				classifier,
				servicegetter.FromEntities(entities),
				checkHandler,
			)
		})
		AfterEach(func() {
			ctrl.Finish()
		})
		It("returns ok response", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			// add "control-point" header to ctx
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs(),
			)
			checkHandler.EXPECT().CheckRequest(gomock.Any(), gomock.Any()).Return(acceptedResponse())
			resp, err := handler.CheckHTTP(ctxWithIp, &flowcontrolhttpv1.CheckHTTPRequest{ControlPoint: "ingress"})
			Expect(err).NotTo(HaveOccurred())
			Expect(code.Code(resp.GetStatus().GetCode())).To(Equal(code.Code_OK))
		})
		It("injects metadata", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs(),
			)
			checkHandler.EXPECT().CheckRequest(gomock.Any(), gomock.Any()).Return(acceptedResponse())
			resp, err := handler.CheckHTTP(ctxWithIp, &flowcontrolhttpv1.CheckHTTPRequest{ControlPoint: "ingress"})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetDynamicMetadata()).NotTo(BeNil())
		})
		It("sets retry-after header", func() {
			ctxWithIp := peer.NewContext(ctx, newFakeRpcPeer("1.2.3.4"))
			ctxWithIp = metadata.NewIncomingContext(
				ctxWithIp,
				metadata.Pairs(),
			)
			response := rejectedResponse()
			response.WaitTime = durationpb.New(10 * time.Second)
			checkHandler.EXPECT().CheckRequest(gomock.Any(), gomock.Any()).Return(response)
			resp, err := handler.CheckHTTP(ctxWithIp, &flowcontrolhttpv1.CheckHTTPRequest{ControlPoint: "ingress"})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetDynamicMetadata()).NotTo(BeNil())
			Expect(resp.GetDeniedResponse().Headers).To(HaveKeyWithValue("retry-after", "10"))
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
				package flowcontrol.checkhttp
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
