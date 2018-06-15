// Code generated by MockGen. DO NOT EDIT.
// Source: gachamachine/storage/iface (interfaces: FolderOperator)

// Package mock_s3_sdk is a generated GoMock package.
package mock_s3_sdk

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFolderOperator is a mock of FolderOperator interface
type MockFolderOperator struct {
	ctrl     *gomock.Controller
	recorder *MockFolderOperatorMockRecorder
}

// MockFolderOperatorMockRecorder is the mock recorder for MockFolderOperator
type MockFolderOperatorMockRecorder struct {
	mock *MockFolderOperator
}

// NewMockFolderOperator creates a new mock instance
func NewMockFolderOperator(ctrl *gomock.Controller) *MockFolderOperator {
	mock := &MockFolderOperator{ctrl: ctrl}
	mock.recorder = &MockFolderOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFolderOperator) EXPECT() *MockFolderOperatorMockRecorder {
	return m.recorder
}

// Readdirnames mocks base method
func (m *MockFolderOperator) Readdirnames(arg0 int) ([]string, error) {
	ret := m.ctrl.Call(m, "Readdirnames", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Readdirnames indicates an expected call of Readdirnames
func (mr *MockFolderOperatorMockRecorder) Readdirnames(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readdirnames", reflect.TypeOf((*MockFolderOperator)(nil).Readdirnames), arg0)
}
