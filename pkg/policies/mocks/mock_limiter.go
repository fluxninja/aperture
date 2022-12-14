// Code generated by MockGen. DO NOT EDIT.
// Source: limiter.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	checkv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	iface "github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
)

// MockLimiter is a mock of Limiter interface.
type MockLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockLimiterMockRecorder
}

// MockLimiterMockRecorder is the mock recorder for MockLimiter.
type MockLimiterMockRecorder struct {
	mock *MockLimiter
}

// NewMockLimiter creates a new mock instance.
func NewMockLimiter(ctrl *gomock.Controller) *MockLimiter {
	mock := &MockLimiter{ctrl: ctrl}
	mock.recorder = &MockLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLimiter) EXPECT() *MockLimiterMockRecorder {
	return m.recorder
}

// GetFlowSelector mocks base method.
func (m *MockLimiter) GetFlowSelector() *languagev1.FlowSelector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowSelector")
	ret0, _ := ret[0].(*languagev1.FlowSelector)
	return ret0
}

// GetFlowSelector indicates an expected call of GetFlowSelector.
func (mr *MockLimiterMockRecorder) GetFlowSelector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowSelector", reflect.TypeOf((*MockLimiter)(nil).GetFlowSelector))
}

// GetLimiterID mocks base method.
func (m *MockLimiter) GetLimiterID() iface.LimiterID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimiterID")
	ret0, _ := ret[0].(iface.LimiterID)
	return ret0
}

// GetLimiterID indicates an expected call of GetLimiterID.
func (mr *MockLimiterMockRecorder) GetLimiterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimiterID", reflect.TypeOf((*MockLimiter)(nil).GetLimiterID))
}

// GetPolicyName mocks base method.
func (m *MockLimiter) GetPolicyName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicyName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPolicyName indicates an expected call of GetPolicyName.
func (mr *MockLimiterMockRecorder) GetPolicyName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicyName", reflect.TypeOf((*MockLimiter)(nil).GetPolicyName))
}

// RunLimiter mocks base method.
func (m *MockLimiter) RunLimiter(ctx context.Context, labels map[string]string) *checkv1.LimiterDecision {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunLimiter", ctx, labels)
	ret0, _ := ret[0].(*checkv1.LimiterDecision)
	return ret0
}

// RunLimiter indicates an expected call of RunLimiter.
func (mr *MockLimiterMockRecorder) RunLimiter(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunLimiter", reflect.TypeOf((*MockLimiter)(nil).RunLimiter), ctx, labels)
}

// MockRateLimiter is a mock of RateLimiter interface.
type MockRateLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockRateLimiterMockRecorder
}

// MockRateLimiterMockRecorder is the mock recorder for MockRateLimiter.
type MockRateLimiterMockRecorder struct {
	mock *MockRateLimiter
}

// NewMockRateLimiter creates a new mock instance.
func NewMockRateLimiter(ctrl *gomock.Controller) *MockRateLimiter {
	mock := &MockRateLimiter{ctrl: ctrl}
	mock.recorder = &MockRateLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateLimiter) EXPECT() *MockRateLimiterMockRecorder {
	return m.recorder
}

// GetFlowSelector mocks base method.
func (m *MockRateLimiter) GetFlowSelector() *languagev1.FlowSelector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowSelector")
	ret0, _ := ret[0].(*languagev1.FlowSelector)
	return ret0
}

// GetFlowSelector indicates an expected call of GetFlowSelector.
func (mr *MockRateLimiterMockRecorder) GetFlowSelector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowSelector", reflect.TypeOf((*MockRateLimiter)(nil).GetFlowSelector))
}

// GetLimiterID mocks base method.
func (m *MockRateLimiter) GetLimiterID() iface.LimiterID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimiterID")
	ret0, _ := ret[0].(iface.LimiterID)
	return ret0
}

// GetLimiterID indicates an expected call of GetLimiterID.
func (mr *MockRateLimiterMockRecorder) GetLimiterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimiterID", reflect.TypeOf((*MockRateLimiter)(nil).GetLimiterID))
}

// GetPolicyName mocks base method.
func (m *MockRateLimiter) GetPolicyName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicyName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPolicyName indicates an expected call of GetPolicyName.
func (mr *MockRateLimiterMockRecorder) GetPolicyName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicyName", reflect.TypeOf((*MockRateLimiter)(nil).GetPolicyName))
}

// GetRequestCounter mocks base method.
func (m *MockRateLimiter) GetRequestCounter(labels map[string]string) prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequestCounter", labels)
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// GetRequestCounter indicates an expected call of GetRequestCounter.
func (mr *MockRateLimiterMockRecorder) GetRequestCounter(labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequestCounter", reflect.TypeOf((*MockRateLimiter)(nil).GetRequestCounter), labels)
}

// RunLimiter mocks base method.
func (m *MockRateLimiter) RunLimiter(ctx context.Context, labels map[string]string) *checkv1.LimiterDecision {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunLimiter", ctx, labels)
	ret0, _ := ret[0].(*checkv1.LimiterDecision)
	return ret0
}

// RunLimiter indicates an expected call of RunLimiter.
func (mr *MockRateLimiterMockRecorder) RunLimiter(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunLimiter", reflect.TypeOf((*MockRateLimiter)(nil).RunLimiter), ctx, labels)
}

// TakeN mocks base method.
func (m *MockRateLimiter) TakeN(labels map[string]string, count int) (string, bool, int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TakeN", labels, count)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(int)
	return ret0, ret1, ret2, ret3
}

// TakeN indicates an expected call of TakeN.
func (mr *MockRateLimiterMockRecorder) TakeN(labels, count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TakeN", reflect.TypeOf((*MockRateLimiter)(nil).TakeN), labels, count)
}

// MockConcurrencyLimiter is a mock of ConcurrencyLimiter interface.
type MockConcurrencyLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockConcurrencyLimiterMockRecorder
}

// MockConcurrencyLimiterMockRecorder is the mock recorder for MockConcurrencyLimiter.
type MockConcurrencyLimiterMockRecorder struct {
	mock *MockConcurrencyLimiter
}

// NewMockConcurrencyLimiter creates a new mock instance.
func NewMockConcurrencyLimiter(ctrl *gomock.Controller) *MockConcurrencyLimiter {
	mock := &MockConcurrencyLimiter{ctrl: ctrl}
	mock.recorder = &MockConcurrencyLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConcurrencyLimiter) EXPECT() *MockConcurrencyLimiterMockRecorder {
	return m.recorder
}

// GetFlowSelector mocks base method.
func (m *MockConcurrencyLimiter) GetFlowSelector() *languagev1.FlowSelector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowSelector")
	ret0, _ := ret[0].(*languagev1.FlowSelector)
	return ret0
}

// GetFlowSelector indicates an expected call of GetFlowSelector.
func (mr *MockConcurrencyLimiterMockRecorder) GetFlowSelector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowSelector", reflect.TypeOf((*MockConcurrencyLimiter)(nil).GetFlowSelector))
}

// GetLatencyObserver mocks base method.
func (m *MockConcurrencyLimiter) GetLatencyObserver(labels map[string]string) prometheus.Observer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatencyObserver", labels)
	ret0, _ := ret[0].(prometheus.Observer)
	return ret0
}

// GetLatencyObserver indicates an expected call of GetLatencyObserver.
func (mr *MockConcurrencyLimiterMockRecorder) GetLatencyObserver(labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatencyObserver", reflect.TypeOf((*MockConcurrencyLimiter)(nil).GetLatencyObserver), labels)
}

// GetLimiterID mocks base method.
func (m *MockConcurrencyLimiter) GetLimiterID() iface.LimiterID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimiterID")
	ret0, _ := ret[0].(iface.LimiterID)
	return ret0
}

// GetLimiterID indicates an expected call of GetLimiterID.
func (mr *MockConcurrencyLimiterMockRecorder) GetLimiterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimiterID", reflect.TypeOf((*MockConcurrencyLimiter)(nil).GetLimiterID))
}

// GetPolicyName mocks base method.
func (m *MockConcurrencyLimiter) GetPolicyName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicyName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPolicyName indicates an expected call of GetPolicyName.
func (mr *MockConcurrencyLimiterMockRecorder) GetPolicyName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicyName", reflect.TypeOf((*MockConcurrencyLimiter)(nil).GetPolicyName))
}

// GetRequestCounter mocks base method.
func (m *MockConcurrencyLimiter) GetRequestCounter(labels map[string]string) prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequestCounter", labels)
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// GetRequestCounter indicates an expected call of GetRequestCounter.
func (mr *MockConcurrencyLimiterMockRecorder) GetRequestCounter(labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequestCounter", reflect.TypeOf((*MockConcurrencyLimiter)(nil).GetRequestCounter), labels)
}

// RunLimiter mocks base method.
func (m *MockConcurrencyLimiter) RunLimiter(ctx context.Context, labels map[string]string) *checkv1.LimiterDecision {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunLimiter", ctx, labels)
	ret0, _ := ret[0].(*checkv1.LimiterDecision)
	return ret0
}

// RunLimiter indicates an expected call of RunLimiter.
func (mr *MockConcurrencyLimiterMockRecorder) RunLimiter(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunLimiter", reflect.TypeOf((*MockConcurrencyLimiter)(nil).RunLimiter), ctx, labels)
}
