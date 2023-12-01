package flowcontrol

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	goprom "github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/mocks"
)

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
			Expect(mmr.schedulers).To(BeEmpty())
		})

		It("Return matched Load Schedulers and Flux Meters", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterScheduler(mockConLimiter)

			controlPoint := "ingress"
			svcs := []string{"testService.testNamespace.svc.cluster.local"}
			labels := labels.PlainMap{"service": "testService.testNamespace.svc.cluster.local"}

			mmr := engine.(*Engine).getMatches(controlPoint, svcs, labels)
			Expect(mmr.fluxMeters).NotTo(BeEmpty())
			Expect(mmr.schedulers).NotTo(BeEmpty())
		})
	})
})
