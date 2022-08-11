package dataplane

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	goprom "github.com/prometheus/client_golang/prometheus"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/mocks"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

var _ = Describe("Dataplane Engine", func() {
	var (
		engine iface.EngineAPI

		t             GinkgoTestReporter
		mockCtrl      *gomock.Controller
		mockLimiter   *mocks.MockLimiter
		mockFluxmeter *mocks.MockFluxMeter

		selector  *policylangv1.Selector
		histogram goprom.Histogram
	)

	agentGroup := ""

	BeforeEach(func() {
		t = GinkgoTestReporter{}
		mockCtrl = gomock.NewController(t)
		mockLimiter = mocks.NewMockLimiter(mockCtrl)
		mockFluxmeter = mocks.NewMockFluxMeter(mockCtrl)

		engine = ProvideEngineAPI()
		selector = &policylangv1.Selector{
			AgentGroup: "default",
			Service:    "testService.testNamespace.svc.cluster.local",
			ControlPoint: &policylangv1.ControlPoint{
				Controlpoint: &policylangv1.ControlPoint_Traffic{Traffic: "ingress"},
			},
		}
		histogram = goprom.NewHistogram(goprom.HistogramOpts{
			Name:        "test",
			ConstLabels: goprom.Labels{"metric_id": "test"},
		})
	})

	Context("Scheduler actuator", func() {
		BeforeEach(func() {
			mockLimiter.EXPECT().GetPolicyName().AnyTimes()
			mockLimiter.EXPECT().GetSelector().Return(selector).AnyTimes()
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
			mockFluxmeter.EXPECT().GetMetricID().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetSelector().Return(selector).AnyTimes()
			mockFluxmeter.EXPECT().GetHistogram().Return(histogram).AnyTimes()
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
			hist := engine.GetFluxMeterHist("test")
			Expect(hist).To(BeNil())
		})

		It("Returns registered fluxmeter hist", func() {
			err := engine.RegisterFluxMeter(mockFluxmeter)
			hist := engine.GetFluxMeterHist("test")
			Expect(err).NotTo(HaveOccurred())
			Expect(hist).To(Equal(histogram))
		})
	})

	Context("Multimatch", func() {
		BeforeEach(func() {
			mockLimiter.EXPECT().GetPolicyName().AnyTimes()
			mockLimiter.EXPECT().GetSelector().Return(selector).AnyTimes()

			mockFluxmeter.EXPECT().GetMetricID().Return("test").AnyTimes()
			mockFluxmeter.EXPECT().GetSelector().Return(selector).AnyTimes()
			mockFluxmeter.EXPECT().GetHistogram().Return(histogram).AnyTimes()
		})

		It("Return nothing for not compatible service", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterConcurrencyLimiter(mockLimiter)

			controlPoint := selectors.ControlPoint{
				Traffic: selectors.Ingress,
			}
			svcs := []services.ServiceID{{
				AgentGroup: "default",
				Service:    "testService2.testNamespace2.svc.cluster.local",
			}}
			labels := selectors.NewLabels(selectors.LabelSources{
				Flow: map[string]string{"service": "whatever"},
			})

			mmr := engine.(*Engine).getMatches(agentGroup, controlPoint, svcs, labels)
			Expect(mmr.FluxMeters).To(BeEmpty())
			Expect(mmr.ConcurrencyLimiters).To(BeEmpty())
		})

		It("Return matched schedulers and fluxmeters", func() {
			_ = engine.RegisterFluxMeter(mockFluxmeter)
			_ = engine.RegisterConcurrencyLimiter(mockLimiter)

			controlPoint := selectors.ControlPoint{
				Traffic: selectors.Ingress,
			}
			svcs := []services.ServiceID{{
				AgentGroup: "default",
				Service:    "testService.testNamespace.svc.cluster.local",
			}}
			labels := selectors.NewLabels(selectors.LabelSources{
				Flow: map[string]string{"service": "testService.testNamespace.svc.cluster.local"},
			})

			mmr := engine.(*Engine).getMatches(agentGroup, controlPoint, svcs, labels)
			Expect(mmr.FluxMeters).NotTo(BeEmpty())
			Expect(mmr.ConcurrencyLimiters).NotTo(BeEmpty())
		})
	})
})
