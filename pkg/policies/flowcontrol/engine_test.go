package flowcontrol

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	goprom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/durationpb"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/mocks"
)

var (
	lock          sync.Mutex
	limitersArray = []string{}
)

type MockTestLimiter struct {
	name         string
	shouldReject bool
}

// Decide implements iface.Limiter.
func (l *MockTestLimiter) Decide(ctx context.Context, flowLabels labels.Labels) *flowcontrolv1.LimiterDecision {
	lock.Lock()
	defer lock.Unlock()

	limitersArray = append(limitersArray, l.name)

	var deniedResponseStatusCode flowcontrolv1.StatusCode
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED
	if l.shouldReject {
		reason = flowcontrolv1.LimiterDecision_LIMITER_REASON_KEY_NOT_FOUND
		deniedResponseStatusCode = flowcontrolv1.StatusCode_BadRequest
	}

	return &flowcontrolv1.LimiterDecision{
		PolicyName:               l.name,
		PolicyHash:               "policy_hash",
		ComponentId:              "component_id",
		Dropped:                  l.shouldReject,
		DeniedResponseStatusCode: deniedResponseStatusCode,
		Reason:                   reason,
		WaitTime:                 &durationpb.Duration{Seconds: 1},
	}
}

// Revert implements iface.Limiter.
func (l *MockTestLimiter) Revert(context.Context, labels.Labels, *flowcontrolv1.LimiterDecision) {
	panic("unimplemented")
}

// GetLimiterID implements iface.Limiter.
func (l *MockTestLimiter) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  l.name,
		PolicyHash:  "policy_hash",
		ComponentID: "component_id",
	}
}

// GetRampMode implements iface.Limiter.
func (l *MockTestLimiter) GetRampMode() bool {
	panic("unimplemented")
}

// GetRequestCounter implements iface.Limiter.
func (l *MockTestLimiter) GetRequestCounter(labels map[string]string) goprom.Counter {
	panic("unimplemented")
}

// GetSelectors implements iface.Limiter.
func (l *MockTestLimiter) GetSelectors() []*policylangv1.Selector {
	return []*policylangv1.Selector{
		{
			ControlPoint: "ingress",
			Service:      "testService.testNamespace.svc.cluster.local",
			AgentGroup:   metrics.DefaultAgentGroup,
		},
	}
}

func (l *MockTestLimiter) GetPolicyName() string {
	return l.name
}

func (l *MockTestLimiter) Return(ctx context.Context, label string, tokens float64, requestID string) (bool, error) {
	panic("unimplemented")
}

func (l *MockTestLimiter) GetLatencyObserver(labels map[string]string) goprom.Observer {
	panic("unimplemented")
}

var _ iface.Limiter = &MockTestLimiter{}

var _ = Describe("Dataplane Engine", func() {
	var (
		engine iface.Engine

		t              GinkgoTestReporter
		mockCtrl       *gomock.Controller
		mockConLimiter *mocks.MockScheduler
		mockFluxmeter  *mocks.MockFluxMeter

		selectors   []*policylangv1.Selector
		histogram   goprom.Histogram
		fluxMeterID iface.FluxMeterID
		limiterID   iface.LimiterID
	)

	BeforeEach(func() {
		t = GinkgoTestReporter{}
		mockCtrl = gomock.NewController(t)
		mockConLimiter = mocks.NewMockScheduler(mockCtrl)
		mockFluxmeter = mocks.NewMockFluxMeter(mockCtrl)

		limitersArray = make([]string, 0)

		engine = NewEngine(agentinfo.NewAgentInfo(metrics.DefaultAgentGroup))
		selectors = []*policylangv1.Selector{
			{
				ControlPoint: "ingress",
				Service:      "testService.testNamespace.svc.cluster.local",
				AgentGroup:   metrics.DefaultAgentGroup,
			},
		}

		histogram = goprom.NewHistogram(goprom.HistogramOpts{
			Name: metrics.FluxMeterMetricName,
			ConstLabels: goprom.Labels{
				metrics.PolicyNameLabel:    "test",
				metrics.FluxMeterNameLabel: "test",
				metrics.PolicyHashLabel:    "test",
				metrics.DecisionTypeLabel:  flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			},
		})
		fluxMeterID = iface.FluxMeterID{
			FluxMeterName: "test",
		}
		limiterID = iface.LimiterID{
			PolicyName:  "test",
			ComponentID: "0",
			PolicyHash:  "test",
		}
	})

	Context("ProcessRequest", func() {
		It("Should have 3 registered limiters", func() {
			cl1 := &MockTestLimiter{
				name:         "concurrency-limiter1",
				shouldReject: false,
			}
			err := engine.RegisterConcurrencyLimiter(cl1)
			Expect(err).NotTo(HaveOccurred())

			cs1 := &MockTestLimiter{
				name:         "concurrency-scheduler1",
				shouldReject: false,
			}
			err = engine.RegisterConcurrencyScheduler(cs1)
			Expect(err).NotTo(HaveOccurred())

			cl2 := &MockTestLimiter{
				name:         "concurrency-limiter2",
				shouldReject: false,
			}
			err = engine.RegisterConcurrencyLimiter(cl2)
			Expect(err).NotTo(HaveOccurred())

			resp := engine.ProcessRequest(context.Background(), iface.RequestContext{
				FlowLabels:   make(labels.PlainMap),
				ControlPoint: "ingress",
				Services:     []string{"testService.testNamespace.svc.cluster.local"},
				RampMode:     false,
				ExpectEnd:    true,
			})

			fmt.Printf("\nprocess request response: %+v\n", resp)
			fmt.Printf("\n limiters array: %+v\n", limitersArray)

			Expect(len(limitersArray)).To(Equal(3))
			Expect(limitersArray[0]).To(Equal("concurrency-limiter1"))
			Expect(limitersArray[1]).To(Equal("concurrency-limiter2"))
			Expect(limitersArray[2]).To(Equal("concurrency-scheduler1"))
		})
	})

	Context("Load Scheduler", func() {
		BeforeEach(func() {
			mockConLimiter.EXPECT().GetPolicyName().AnyTimes()
			mockConLimiter.EXPECT().GetSelectors().Return(selectors).AnyTimes()
			mockConLimiter.EXPECT().GetLimiterID().Return(limiterID).AnyTimes()
		})

		It("Registers Load Scheduler", func() {
			err := engine.RegisterScheduler(mockConLimiter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Registers Load Scheduler second time", func() {
			err := engine.RegisterScheduler(mockConLimiter)
			err2 := engine.RegisterScheduler(mockConLimiter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).To(HaveOccurred())
		})

		It("Unregisters not registered Load Scheduler", func() {
			err := engine.UnregisterScheduler(mockConLimiter)
			Expect(err).To(HaveOccurred())
		})

		It("Unregisters existing Load Scheduler", func() {
			err := engine.RegisterScheduler(mockConLimiter)
			err2 := engine.UnregisterScheduler(mockConLimiter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})
	})

	Context("Flux meter", func() {
		var labels map[string]string

		BeforeEach(func() {
			mockFluxmeter.EXPECT().GetFluxMeterName().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetSelectors().Return(selectors).AnyTimes()
			labels = map[string]string{
				metrics.FlowStatusLabel:   metrics.FlowStatusOK,
				metrics.StatusCodeLabel:   "200",
				metrics.DecisionTypeLabel: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			}
			mockFluxmeter.EXPECT().GetHistogram(labels).Return(histogram).AnyTimes()
			mockFluxmeter.EXPECT().GetFluxMeterID().Return(fluxMeterID).AnyTimes()
		})

		It("Registers Flux meter", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Registers Flux meter second time", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			err2 := engine.RegisterFluxMeter(mockFluxmeter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).To(HaveOccurred())
		})

		It("Unregisters not registered Flux meter", func() {
			err := engine.UnregisterFluxMeter(mockFluxmeter)
			Expect(err).To(HaveOccurred())
		})

		It("Unregisters existing Flux meter", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			err2 := engine.UnregisterFluxMeter(mockFluxmeter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("Tries to get unregistered fluxmeter histogram", func() {
			fluxMeter := engine.GetFluxMeter("test")
			Expect(fluxMeter).To(BeNil())
		})

		It("Returns registered fluxmeter histogram", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			Expect(err).NotTo(HaveOccurred())
			fluxMeter := engine.GetFluxMeter("test")
			h := fluxMeter.GetHistogram(labels)
			Expect(h).To(Equal(histogram))
		})
	})

	Context("Multimatch", func() {
		BeforeEach(func() {
			mockConLimiter.EXPECT().GetPolicyName().AnyTimes()
			mockConLimiter.EXPECT().GetSelectors().Return(selectors).AnyTimes()
			mockConLimiter.EXPECT().GetLimiterID().Return(limiterID).AnyTimes()

			mockFluxmeter.EXPECT().GetFluxMeterName().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetSelectors().Return(selectors).AnyTimes()
			mockFluxmeter.EXPECT().GetFluxMeterID().Return(fluxMeterID).AnyTimes()
		})

		It("Return nothing for not compatible service", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterScheduler(mockConLimiter)

			controlPoint := "ingress"
			svcs := []string{"testService2.testNamespace2.svc.cluster.local"}
			labels := labels.PlainMap{"service": "whatever"}

			mmr := engine.(*Engine).getMatches(controlPoint, svcs, labels)
			Expect(mmr.fluxMeters).To(BeEmpty())
			Expect(mmr.quotaAndLoadSchedulers).To(BeEmpty())
		})

		It("Return matched Load Schedulers and Flux Meters", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterScheduler(mockConLimiter)

			controlPoint := "ingress"
			svcs := []string{"testService.testNamespace.svc.cluster.local"}
			labels := labels.PlainMap{"service": "testService.testNamespace.svc.cluster.local"}

			mmr := engine.(*Engine).getMatches(controlPoint, svcs, labels)
			Expect(mmr.fluxMeters).NotTo(BeEmpty())
			Expect(mmr.quotaAndLoadSchedulers).NotTo(BeEmpty())
		})
	})
})
