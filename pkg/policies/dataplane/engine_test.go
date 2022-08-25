package dataplane

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	goprom "github.com/prometheus/client_golang/prometheus"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/mocks"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

var _ = Describe("Dataplane Engine", func() {
	var (
		engine iface.Engine

		t             GinkgoTestReporter
		mockCtrl      *gomock.Controller
		mockLimiter   *mocks.MockLimiter
		mockFluxmeter *mocks.MockFluxMeter

		selector    *selectorv1.Selector
		histogram   goprom.Histogram
		fluxMeterID iface.FluxMeterID
		limiterID   iface.LimiterID
	)

	BeforeEach(func() {
		t = GinkgoTestReporter{}
		mockCtrl = gomock.NewController(t)
		mockLimiter = mocks.NewMockLimiter(mockCtrl)
		mockFluxmeter = mocks.NewMockFluxMeter(mockCtrl)

		engine = ProvideEngineAPI()
		selector = &selectorv1.Selector{
			AgentGroup: metrics.DefaultAgentGroup,
			Service:    "testService.testNamespace.svc.cluster.local",
			ControlPoint: &selectorv1.ControlPoint{
				Controlpoint: &selectorv1.ControlPoint_Traffic{Traffic: "ingress"},
			},
		}
		histogram = goprom.NewHistogram(goprom.HistogramOpts{
			Name: metrics.FluxMeterMetricName,
			ConstLabels: goprom.Labels{
				metrics.PolicyNameLabel:    "test",
				metrics.FluxMeterNameLabel: "test",
				metrics.PolicyHashLabel:    "test",
				metrics.DecisionTypeLabel:  flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED.String(),
			},
		})
		fluxMeterID = iface.FluxMeterID{
			PolicyName:    "test",
			FluxMeterName: "test",
			PolicyHash:    "test",
		}
		limiterID = iface.LimiterID{
			PolicyName:     "test",
			ComponentIndex: 0,
			PolicyHash:     "test",
		}
	})

	Context("Scheduler actuator", func() {
		BeforeEach(func() {
			mockLimiter.EXPECT().GetPolicyName().AnyTimes()
			mockLimiter.EXPECT().GetSelector().Return(selector).AnyTimes()
			mockLimiter.EXPECT().GetLimiterID().Return(limiterID).AnyTimes()
		})

		It("Registers scheduler actuator", func() {
			err := engine.RegisterConcurrencyLimiter(mockLimiter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Registers scheduler actuator second time", func() {
			err := engine.RegisterConcurrencyLimiter(mockLimiter)
			err2 := engine.RegisterConcurrencyLimiter(mockLimiter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("Unregisters not registered scheduler actuator", func() {
			err := engine.UnregisterConcurrencyLimiter(mockLimiter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Unregisters existing scheduler actuator", func() {
			err := engine.RegisterConcurrencyLimiter(mockLimiter)
			err2 := engine.UnregisterConcurrencyLimiter(mockLimiter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})
	})

	Context("Flux meter", func() {
		BeforeEach(func() {
			mockFluxmeter.EXPECT().GetPolicyName().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetFluxMeterName().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetPolicyHash().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetSelector().Return(selector).AnyTimes()
			mockFluxmeter.EXPECT().GetHistogram(flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED, "200").Return(histogram).AnyTimes()
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
			Expect(err).NotTo(HaveOccurred())
		})

		It("Unregisters existing Flux meter", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			err2 := engine.UnregisterFluxMeter(mockFluxmeter)
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("Tries to get unregistered fluxmeter hist", func() {
			hist := engine.GetFluxMeterHist("test", "test", "test", "200", flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED)
			Expect(hist).To(BeNil())
		})

		It("Returns registered fluxmeter hist", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			hist := engine.GetFluxMeterHist("test", "test", "test", "200", flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED)
			Expect(err).NotTo(HaveOccurred())
			Expect(hist).To(Equal(histogram))
		})
	})

	Context("Multimatch", func() {
		BeforeEach(func() {
			mockLimiter.EXPECT().GetPolicyName().AnyTimes()
			mockLimiter.EXPECT().GetSelector().Return(selector).AnyTimes()
			mockLimiter.EXPECT().GetLimiterID().Return(limiterID).AnyTimes()

			mockFluxmeter.EXPECT().GetPolicyName().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetFluxMeterName().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetPolicyHash().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetSelector().Return(selector).AnyTimes()
			mockFluxmeter.EXPECT().GetHistogram(flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED, "503").Return(histogram).AnyTimes()
			mockFluxmeter.EXPECT().GetFluxMeterID().Return(fluxMeterID).AnyTimes()
		})

		It("Return nothing for not compatible service", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterConcurrencyLimiter(mockLimiter)

			controlPoint := selectors.ControlPoint{
				Traffic: selectors.Ingress,
			}
			svcs := []services.ServiceID{{
				Service: "testService2.testNamespace2.svc.cluster.local",
			}}
			labels := selectors.NewLabels(selectors.LabelSources{
				Flow: map[string]string{"service": "whatever"},
			})

			mmr := engine.(*Engine).getMatches(controlPoint, svcs, labels)
			Expect(mmr.fluxMeters).To(BeEmpty())
			Expect(mmr.concurrencyLimiters).To(BeEmpty())
		})

		It("Return matched schedulers and fluxmeters", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterConcurrencyLimiter(mockLimiter)

			controlPoint := selectors.ControlPoint{
				Traffic: selectors.Ingress,
			}
			svcs := []services.ServiceID{{
				Service: "testService.testNamespace.svc.cluster.local",
			}}
			labels := selectors.NewLabels(selectors.LabelSources{
				Flow: map[string]string{"service": "testService.testNamespace.svc.cluster.local"},
			})

			mmr := engine.(*Engine).getMatches(controlPoint, svcs, labels)
			Expect(mmr.fluxMeters).NotTo(BeEmpty())
			Expect(mmr.concurrencyLimiters).NotTo(BeEmpty())
		})
	})
})
