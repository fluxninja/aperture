// Code generated by MockGen. DO NOT EDIT.
// Source: engine.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	iface "github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	selectors "github.com/fluxninja/aperture/pkg/selectors"
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
)

// MockEngineAPI is a mock of EngineAPI interface.
type MockEngineAPI struct {
	ctrl     *gomock.Controller
	recorder *MockEngineAPIMockRecorder
}

// MockEngineAPIMockRecorder is the mock recorder for MockEngineAPI.
type MockEngineAPIMockRecorder struct {
	mock *MockEngineAPI
}

// NewMockEngineAPI creates a new mock instance.
func NewMockEngineAPI(ctrl *gomock.Controller) *MockEngineAPI {
	mock := &MockEngineAPI{ctrl: ctrl}
	mock.recorder = &MockEngineAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEngineAPI) EXPECT() *MockEngineAPIMockRecorder {
	return m.recorder
}

// GetFluxMeterHist mocks base method.
func (m *MockEngineAPI) GetFluxMeterHist(metricID string) prometheus.Histogram {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFluxMeterHist", metricID)
	ret0, _ := ret[0].(prometheus.Histogram)
	return ret0
}

// GetFluxMeterHist indicates an expected call of GetFluxMeterHist.
func (mr *MockEngineAPIMockRecorder) GetFluxMeterHist(metricID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFluxMeterHist", reflect.TypeOf((*MockEngineAPI)(nil).GetFluxMeterHist), metricID)
}

// ProcessRequest mocks base method.
func (m *MockEngineAPI) ProcessRequest(controlPoint selectors.ControlPoint, svcs []string, labels selectors.Labels) *flowcontrolv1.CheckResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessRequest", controlPoint, svcs, labels)
	ret0, _ := ret[0].(*flowcontrolv1.CheckResponse)
	return ret0
}

// ProcessRequest indicates an expected call of ProcessRequest.
func (mr *MockEngineAPIMockRecorder) ProcessRequest(controlPoint, svcs, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessRequest", reflect.TypeOf((*MockEngineAPI)(nil).ProcessRequest), controlPoint, svcs, labels)
}

// RegisterConcurrencyLimiter mocks base method.
func (m *MockEngineAPI) RegisterConcurrencyLimiter(sa iface.Limiter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterConcurrencyLimiter", sa)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterConcurrencyLimiter indicates an expected call of RegisterConcurrencyLimiter.
func (mr *MockEngineAPIMockRecorder) RegisterConcurrencyLimiter(sa interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterConcurrencyLimiter", reflect.TypeOf((*MockEngineAPI)(nil).RegisterConcurrencyLimiter), sa)
}

// RegisterFluxMeter mocks base method.
func (m *MockEngineAPI) RegisterFluxMeter(fm iface.FluxMeter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterFluxMeter", fm)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterFluxMeter indicates an expected call of RegisterFluxMeter.
func (mr *MockEngineAPIMockRecorder) RegisterFluxMeter(fm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterFluxMeter", reflect.TypeOf((*MockEngineAPI)(nil).RegisterFluxMeter), fm)
}

// RegisterRateLimiter mocks base method.
func (m *MockEngineAPI) RegisterRateLimiter(l iface.RateLimiter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterRateLimiter", l)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterRateLimiter indicates an expected call of RegisterRateLimiter.
func (mr *MockEngineAPIMockRecorder) RegisterRateLimiter(l interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRateLimiter", reflect.TypeOf((*MockEngineAPI)(nil).RegisterRateLimiter), l)
}

// UnregisterConcurrencyLimiter mocks base method.
func (m *MockEngineAPI) UnregisterConcurrencyLimiter(sa iface.Limiter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterConcurrencyLimiter", sa)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterConcurrencyLimiter indicates an expected call of UnregisterConcurrencyLimiter.
func (mr *MockEngineAPIMockRecorder) UnregisterConcurrencyLimiter(sa interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterConcurrencyLimiter", reflect.TypeOf((*MockEngineAPI)(nil).UnregisterConcurrencyLimiter), sa)
}

// UnregisterFluxMeter mocks base method.
func (m *MockEngineAPI) UnregisterFluxMeter(fm iface.FluxMeter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterFluxMeter", fm)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterFluxMeter indicates an expected call of UnregisterFluxMeter.
func (mr *MockEngineAPIMockRecorder) UnregisterFluxMeter(fm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterFluxMeter", reflect.TypeOf((*MockEngineAPI)(nil).UnregisterFluxMeter), fm)
}

// UnregisterRateLimiter mocks base method.
func (m *MockEngineAPI) UnregisterRateLimiter(l iface.RateLimiter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterRateLimiter", l)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterRateLimiter indicates an expected call of UnregisterRateLimiter.
func (mr *MockEngineAPIMockRecorder) UnregisterRateLimiter(l interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterRateLimiter", reflect.TypeOf((*MockEngineAPI)(nil).UnregisterRateLimiter), l)
}
