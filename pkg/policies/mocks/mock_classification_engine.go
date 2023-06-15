// Code generated by MockGen. DO NOT EDIT.
// Source: classification-engine.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	iface "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	gomock "github.com/golang/mock/gomock"
)

// MockClassificationEngine is a mock of ClassificationEngine interface.
type MockClassificationEngine struct {
	ctrl     *gomock.Controller
	recorder *MockClassificationEngineMockRecorder
}

// MockClassificationEngineMockRecorder is the mock recorder for MockClassificationEngine.
type MockClassificationEngineMockRecorder struct {
	mock *MockClassificationEngine
}

// NewMockClassificationEngine creates a new mock instance.
func NewMockClassificationEngine(ctrl *gomock.Controller) *MockClassificationEngine {
	mock := &MockClassificationEngine{ctrl: ctrl}
	mock.recorder = &MockClassificationEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClassificationEngine) EXPECT() *MockClassificationEngineMockRecorder {
	return m.recorder
}

// GetClassifier mocks base method.
func (m *MockClassificationEngine) GetClassifier(classifierID iface.ClassifierID) iface.Classifier {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClassifier", classifierID)
	ret0, _ := ret[0].(iface.Classifier)
	return ret0
}

// GetClassifier indicates an expected call of GetClassifier.
func (mr *MockClassificationEngineMockRecorder) GetClassifier(classifierID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClassifier", reflect.TypeOf((*MockClassificationEngine)(nil).GetClassifier), classifierID)
}

// RegisterClassifier mocks base method.
func (m *MockClassificationEngine) RegisterClassifier(classifier iface.Classifier) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterClassifier", classifier)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterClassifier indicates an expected call of RegisterClassifier.
func (mr *MockClassificationEngineMockRecorder) RegisterClassifier(classifier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterClassifier", reflect.TypeOf((*MockClassificationEngine)(nil).RegisterClassifier), classifier)
}

// UnregisterClassifier mocks base method.
func (m *MockClassificationEngine) UnregisterClassifier(classifier iface.Classifier) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterClassifier", classifier)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterClassifier indicates an expected call of UnregisterClassifier.
func (mr *MockClassificationEngineMockRecorder) UnregisterClassifier(classifier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterClassifier", reflect.TypeOf((*MockClassificationEngine)(nil).UnregisterClassifier), classifier)
}
