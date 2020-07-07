// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudfoundry/bosh-cli/deployment/release (interfaces: JobResolver)

// Package mocks is a generated GoMock package.
package mocks

import (
	job "github.com/cloudfoundry/bosh-cli/release/job"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockJobResolver is a mock of JobResolver interface
type MockJobResolver struct {
	ctrl     *gomock.Controller
	recorder *MockJobResolverMockRecorder
}

// MockJobResolverMockRecorder is the mock recorder for MockJobResolver
type MockJobResolverMockRecorder struct {
	mock *MockJobResolver
}

// NewMockJobResolver creates a new mock instance
func NewMockJobResolver(ctrl *gomock.Controller) *MockJobResolver {
	mock := &MockJobResolver{ctrl: ctrl}
	mock.recorder = &MockJobResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobResolver) EXPECT() *MockJobResolverMockRecorder {
	return m.recorder
}

// Resolve mocks base method
func (m *MockJobResolver) Resolve(arg0, arg1 string) (job.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve", arg0, arg1)
	ret0, _ := ret[0].(job.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resolve indicates an expected call of Resolve
func (mr *MockJobResolverMockRecorder) Resolve(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockJobResolver)(nil).Resolve), arg0, arg1)
}
