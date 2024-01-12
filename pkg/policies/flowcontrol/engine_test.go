package flowcontrol

import (
	"context"
	"strings"
	"sync"
	"time"

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
	lock       sync.Mutex
	decideList = []string{}
	revertList = []string{}
	cacheMap   = map[string]string{}
)

type mockTestLimiter struct {
	name         string
	limiterType  string
	shouldReject bool
}

func noXAfterY(x string, y string) bool {
	lock.Lock()
	defer lock.Unlock()

	schedulerFound := false

	if len(decideList) > 0 {
		for i := 0; i < len(decideList)-1; i++ {
			limiterType := strings.Split(decideList[i], ".")[1]
			if limiterType == y {
				schedulerFound = true
				continue
			}
			if schedulerFound && limiterType == "x" {
				return false
			}
		}
	}

	return true
}

// Decide implements iface.Limiter.
func (l *mockTestLimiter) Decide(ctx context.Context, flowLabels labels.Labels) *flowcontrolv1.LimiterDecision {
	lock.Lock()
	defer lock.Unlock()

	decideList = append(decideList, l.name+"."+l.limiterType)

	var deniedResponseStatusCode flowcontrolv1.StatusCode
	reason := flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED
	if l.shouldReject {
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
func (l *mockTestLimiter) Revert(ctx context.Context, flowLabels labels.Labels, decision *flowcontrolv1.LimiterDecision) {
	lock.Lock()
	defer lock.Unlock()

	revertList = append(revertList, l.name)
}

// GetLimiterID implements iface.Limiter.
func (l *mockTestLimiter) GetLimiterID() iface.LimiterID {
	return iface.LimiterID{
		PolicyName:  l.name,
		PolicyHash:  "policy_hash",
		ComponentID: "component_id",
	}
}

// GetRampMode implements iface.Limiter.
func (l *mockTestLimiter) GetRampMode() bool {
	panic("unimplemented")
}

// GetRequestCounter implements iface.Limiter.
func (l *mockTestLimiter) GetRequestCounter(labels map[string]string) goprom.Counter {
	panic("unimplemented")
}

// GetSelectors implements iface.Limiter.
func (l *mockTestLimiter) GetSelectors() []*policylangv1.Selector {
	return []*policylangv1.Selector{
		{
			ControlPoint: "ingress",
			Service:      "testService.testNamespace.svc.cluster.local",
			AgentGroup:   metrics.DefaultAgentGroup,
		},
	}
}

func (l *mockTestLimiter) GetPolicyName() string {
	return l.name
}

func (l *mockTestLimiter) Return(ctx context.Context, label string, tokens float64, requestID string) (bool, error) {
	panic("unimplemented")
}

func (l *mockTestLimiter) GetLatencyObserver(labels map[string]string) goprom.Observer {
	panic("unimplemented")
}

var _ iface.Limiter = &mockTestLimiter{}

type mockTestCache struct{}

// Delete implements iface.Cache.
func (*mockTestCache) Delete(ctx context.Context, req *flowcontrolv1.CacheDeleteRequest) *flowcontrolv1.CacheDeleteResponse {
	panic("unimplemented")
}

// Lookup implements iface.Cache.
func (*mockTestCache) Lookup(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) *flowcontrolv1.CacheLookupResponse {
	panic("unimplemented")
}

// LookupGlobal implements iface.Cache.
func (*mockTestCache) LookupGlobal(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) map[string]*flowcontrolv1.KeyLookupResponse {
	panic("unimplemented")
}

// LookupGlobalNoWait implements iface.Cache.
func (*mockTestCache) LookupGlobalNoWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (map[string]*flowcontrolv1.KeyLookupResponse, *sync.WaitGroup) {
	panic("unimplemented")
}

// LookupNoWait implements iface.Cache.
func (*mockTestCache) LookupNoWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (*flowcontrolv1.CacheLookupResponse, *sync.WaitGroup, *sync.WaitGroup) {
	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}

	lock.Lock()
	defer lock.Unlock()

	key := request.ResultCacheKey
	value, ok := cacheMap[key]

	lookupStatus := flowcontrolv1.CacheLookupStatus_MISS
	if ok {
		lookupStatus = flowcontrolv1.CacheLookupStatus_HIT
	}

	return &flowcontrolv1.CacheLookupResponse{
		ResultCacheResponse: &flowcontrolv1.KeyLookupResponse{
			LookupStatus: lookupStatus,
			Value:        []byte(value),
		},
	}, &wg1, &wg2
}

// LookupResult implements iface.Cache.
func (*mockTestCache) LookupResult(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) *flowcontrolv1.KeyLookupResponse {
	panic("unimplemented")
}

// LookupResultNoWait implements iface.Cache.
func (*mockTestCache) LookupResultNoWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (*flowcontrolv1.KeyLookupResponse, *sync.WaitGroup) {
	panic("unimplemented")
}

// Upsert implements iface.Cache.
func (*mockTestCache) Upsert(ctx context.Context, req *flowcontrolv1.CacheUpsertRequest) *flowcontrolv1.CacheUpsertResponse {
	panic("unimplemented")
}

var _ iface.Cache = &mockTestCache{}

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

	Context("ProcessRequest", func() {
		BeforeEach(func() {
			decideList = make([]string, 0)
			revertList = make([]string, 0)
			cacheMap = make(map[string]string)

			cacheMap["key"] = "value"

			cl1 := &mockTestLimiter{
				name:         "concurrency-limiter1",
				limiterType:  "limiter",
				shouldReject: false,
			}
			err := engine.RegisterConcurrencyLimiter(cl1)
			Expect(err).NotTo(HaveOccurred())

			cs1 := &mockTestLimiter{
				name:         "concurrency-scheduler1",
				limiterType:  "scheduler",
				shouldReject: false,
			}
			err = engine.RegisterConcurrencyScheduler(cs1)
			Expect(err).NotTo(HaveOccurred())

			cl2 := &mockTestLimiter{
				name:         "concurrency-limiter2",
				limiterType:  "limiter",
				shouldReject: false,
			}
			err = engine.RegisterConcurrencyLimiter(cl2)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have 3 registered limiters and maintain order", func() {
			rl1 := &mockTestLimiter{
				name:         "rate-limiter1",
				limiterType:  "limiter",
				shouldReject: false,
			}
			err := engine.RegisterRateLimiter(rl1)
			Expect(err).NotTo(HaveOccurred())

			_ = engine.ProcessRequest(context.Background(), iface.RequestContext{
				FlowLabels:   make(labels.PlainMap),
				ControlPoint: "ingress",
				Services:     []string{"testService.testNamespace.svc.cluster.local"},
				RampMode:     false,
				ExpectEnd:    true,
			})

			lock.Lock()
			Expect(len(decideList)).To(Equal(4))
			lock.Unlock()
			Expect(noXAfterY("limiter", "scheduler")).To(BeTrue())

			err = engine.UnregisterRateLimiter(rl1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should revert other limiters if rate limiter rejects the request", func() {
			rl1 := &mockTestLimiter{
				name:         "rate-limiter1",
				limiterType:  "limiter",
				shouldReject: true,
			}
			err := engine.RegisterRateLimiter(rl1)
			Expect(err).NotTo(HaveOccurred())

			_ = engine.ProcessRequest(context.Background(), iface.RequestContext{
				FlowLabels:   make(labels.PlainMap),
				ControlPoint: "ingress",
				Services:     []string{"testService.testNamespace.svc.cluster.local"},
				RampMode:     false,
				ExpectEnd:    true,
			})
			time.Sleep(2 * time.Second)

			lock.Lock()
			Expect(len(decideList)).To(Equal(3))
			Expect(decideList).NotTo(ContainElement("concurrency-scheduler1"))
			lock.Unlock()

			lock.Lock()
			Expect(len(revertList)).To(Equal(2))
			Expect(revertList).NotTo(ContainElement("rate-limiter1"))
			lock.Unlock()

			err = engine.UnregisterRateLimiter(rl1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should retrieve cached data correctly", func() {
			cache := &mockTestCache{}
			engine.RegisterCache(cache)

			resp := engine.ProcessRequest(context.Background(), iface.RequestContext{
				FlowLabels:   make(labels.PlainMap),
				ControlPoint: "ingress",
				Services:     []string{"testService.testNamespace.svc.cluster.local"},
				RampMode:     false,
				ExpectEnd:    true,
				CacheLookupRequest: &flowcontrolv1.CacheLookupRequest{
					ControlPoint:   "ingress",
					ResultCacheKey: "key",
				},
			})
			Expect(resp).NotTo(BeNil())
			Expect(resp.CacheLookupResponse).NotTo(BeNil())
			Expect(resp.CacheLookupResponse.ResultCacheResponse.LookupStatus).To(Equal(flowcontrolv1.CacheLookupStatus_HIT))
			Expect(string(resp.CacheLookupResponse.ResultCacheResponse.Value)).To(Equal(cacheMap["key"]))
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
