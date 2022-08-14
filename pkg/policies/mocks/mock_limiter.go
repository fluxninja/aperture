// Code generated by MockGen. DO NOT EDIT.
// Source: limiter.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	iface "github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	selectors "github.com/fluxninja/aperture/pkg/selectors"
	gomock "github.com/golang/mock/gomock"
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

// GetSelector mocks base method.
func (m *MockLimiter) GetSelector() *languagev1.Selector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelector")
	ret0, _ := ret[0].(*languagev1.Selector)
	return ret0
}

// GetSelector indicates an expected call of GetSelector.
func (mr *MockLimiterMockRecorder) GetSelector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelector", reflect.TypeOf((*MockLimiter)(nil).GetSelector))
}

// RunLimiter mocks base method.
func (m *MockLimiter) RunLimiter(labels selectors.Labels) *flowcontrolv1.LimiterDecision {
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
