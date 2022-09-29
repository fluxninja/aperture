// Code generated by MockGen. DO NOT EDIT.
// Source: limiter.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	iface "github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
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

// GetCounter mocks base method.
func (m *MockLimiter) GetCounter() prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter")
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockLimiterMockRecorder) GetCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockLimiter)(nil).GetCounter))
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

// GetObserver mocks base method.
func (m *MockLimiter) GetObserver(labels map[string]string) prometheus.Observer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObserver", labels)
	ret0, _ := ret[0].(prometheus.Observer)
	return ret0
}

// GetObserver indicates an expected call of GetObserver.
func (mr *MockLimiterMockRecorder) GetObserver(labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObserver", reflect.TypeOf((*MockLimiter)(nil).GetObserver), labels)
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

// GetSelector mocks base method.
func (m *MockLimiter) GetSelector() *selectorv1.Selector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelector")
	ret0, _ := ret[0].(*selectorv1.Selector)
	return ret0
}

// GetSelector indicates an expected call of GetSelector.
func (mr *MockLimiterMockRecorder) GetSelector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelector", reflect.TypeOf((*MockLimiter)(nil).GetSelector))
}

// RunLimiter mocks base method.
func (m *MockLimiter) RunLimiter(labels map[string]string) *flowcontrolv1.LimiterDecision {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunLimiter", labels)
	ret0, _ := ret[0].(*flowcontrolv1.LimiterDecision)
	return ret0
}

// RunLimiter indicates an expected call of RunLimiter.
func (mr *MockLimiterMockRecorder) RunLimiter(labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunLimiter", reflect.TypeOf((*MockLimiter)(nil).RunLimiter), labels)
}
