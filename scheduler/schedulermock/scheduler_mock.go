// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/link/v3/scheduler (interfaces: Scheduler)

// Package schedulermock is a generated GoMock package.
package schedulermock

import (
	context "context"
	reflect "reflect"

	models "github.com/Scalingo/link/v3/models"
	scheduler "github.com/Scalingo/link/v3/scheduler"
	gomock "github.com/golang/mock/gomock"
)

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// ConfiguredEndpoints mocks base method.
func (m *MockScheduler) ConfiguredEndpoints(arg0 context.Context) scheduler.EndpointsWithStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfiguredEndpoints", arg0)
	ret0, _ := ret[0].(scheduler.EndpointsWithStatus)
	return ret0
}

// ConfiguredEndpoints indicates an expected call of ConfiguredEndpoints.
func (mr *MockSchedulerMockRecorder) ConfiguredEndpoints(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfiguredEndpoints", reflect.TypeOf((*MockScheduler)(nil).ConfiguredEndpoints), arg0)
}

// Failover mocks base method.
func (m *MockScheduler) Failover(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Failover", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Failover indicates an expected call of Failover.
func (mr *MockSchedulerMockRecorder) Failover(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Failover", reflect.TypeOf((*MockScheduler)(nil).Failover), arg0, arg1)
}

// GetEndpoint mocks base method.
func (m *MockScheduler) GetEndpoint(arg0 context.Context, arg1 string) *scheduler.EndpointWithStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEndpoint", arg0, arg1)
	ret0, _ := ret[0].(*scheduler.EndpointWithStatus)
	return ret0
}

// GetEndpoint indicates an expected call of GetEndpoint.
func (mr *MockSchedulerMockRecorder) GetEndpoint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEndpoint", reflect.TypeOf((*MockScheduler)(nil).GetEndpoint), arg0, arg1)
}

// Start mocks base method.
func (m *MockScheduler) Start(arg0 context.Context, arg1 models.Endpoint) (models.Endpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0, arg1)
	ret0, _ := ret[0].(models.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start.
func (mr *MockSchedulerMockRecorder) Start(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockScheduler)(nil).Start), arg0, arg1)
}

// Status mocks base method.
func (m *MockScheduler) Status(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockSchedulerMockRecorder) Status(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockScheduler)(nil).Status), arg0)
}

// Stop mocks base method.
func (m *MockScheduler) Stop(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockSchedulerMockRecorder) Stop(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockScheduler)(nil).Stop), arg0, arg1)
}

// UpdateEndpoint mocks base method.
func (m *MockScheduler) UpdateEndpoint(arg0 context.Context, arg1 models.Endpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEndpoint", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEndpoint indicates an expected call of UpdateEndpoint.
func (mr *MockSchedulerMockRecorder) UpdateEndpoint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEndpoint", reflect.TypeOf((*MockScheduler)(nil).UpdateEndpoint), arg0, arg1)
}
