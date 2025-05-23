// Code generated by MockGen. DO NOT EDIT.
// Source: go.etcd.io/etcd/client/v3 (interfaces: Lease)

// Package etcdmock is a generated GoMock package.
package etcdmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// MockLease is a mock of Lease interface.
type MockLease struct {
	ctrl     *gomock.Controller
	recorder *MockLeaseMockRecorder
}

// MockLeaseMockRecorder is the mock recorder for MockLease.
type MockLeaseMockRecorder struct {
	mock *MockLease
}

// NewMockLease creates a new mock instance.
func NewMockLease(ctrl *gomock.Controller) *MockLease {
	mock := &MockLease{ctrl: ctrl}
	mock.recorder = &MockLeaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLease) EXPECT() *MockLeaseMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockLease) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockLeaseMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLease)(nil).Close))
}

// Grant mocks base method.
func (m *MockLease) Grant(arg0 context.Context, arg1 int64) (*clientv3.LeaseGrantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Grant", arg0, arg1)
	ret0, _ := ret[0].(*clientv3.LeaseGrantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Grant indicates an expected call of Grant.
func (mr *MockLeaseMockRecorder) Grant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Grant", reflect.TypeOf((*MockLease)(nil).Grant), arg0, arg1)
}

// KeepAlive mocks base method.
func (m *MockLease) KeepAlive(arg0 context.Context, arg1 clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepAlive", arg0, arg1)
	ret0, _ := ret[0].(<-chan *clientv3.LeaseKeepAliveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeepAlive indicates an expected call of KeepAlive.
func (mr *MockLeaseMockRecorder) KeepAlive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepAlive", reflect.TypeOf((*MockLease)(nil).KeepAlive), arg0, arg1)
}

// KeepAliveOnce mocks base method.
func (m *MockLease) KeepAliveOnce(arg0 context.Context, arg1 clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepAliveOnce", arg0, arg1)
	ret0, _ := ret[0].(*clientv3.LeaseKeepAliveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeepAliveOnce indicates an expected call of KeepAliveOnce.
func (mr *MockLeaseMockRecorder) KeepAliveOnce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepAliveOnce", reflect.TypeOf((*MockLease)(nil).KeepAliveOnce), arg0, arg1)
}

// Leases mocks base method.
func (m *MockLease) Leases(arg0 context.Context) (*clientv3.LeaseLeasesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leases", arg0)
	ret0, _ := ret[0].(*clientv3.LeaseLeasesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Leases indicates an expected call of Leases.
func (mr *MockLeaseMockRecorder) Leases(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leases", reflect.TypeOf((*MockLease)(nil).Leases), arg0)
}

// Revoke mocks base method.
func (m *MockLease) Revoke(arg0 context.Context, arg1 clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Revoke", arg0, arg1)
	ret0, _ := ret[0].(*clientv3.LeaseRevokeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Revoke indicates an expected call of Revoke.
func (mr *MockLeaseMockRecorder) Revoke(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revoke", reflect.TypeOf((*MockLease)(nil).Revoke), arg0, arg1)
}

// TimeToLive mocks base method.
func (m *MockLease) TimeToLive(arg0 context.Context, arg1 clientv3.LeaseID, arg2 ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TimeToLive", varargs...)
	ret0, _ := ret[0].(*clientv3.LeaseTimeToLiveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TimeToLive indicates an expected call of TimeToLive.
func (mr *MockLeaseMockRecorder) TimeToLive(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeToLive", reflect.TypeOf((*MockLease)(nil).TimeToLive), varargs...)
}
