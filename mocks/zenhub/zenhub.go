// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/44smkn/zenhub_exporter/pkg/zenhub (interfaces: Client)

// Package mocks_zenhub is a generated GoMock package.
package mocks_zenhub

import (
	context "context"
	reflect "reflect"

	model "github.com/44smkn/zenhub_exporter/pkg/model"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// FetchWorkspaceIssues mocks base method.
func (m *MockClient) FetchWorkspaceIssues(arg0 context.Context) ([]model.Issue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchWorkspaceIssues", arg0)
	ret0, _ := ret[0].([]model.Issue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchWorkspaceIssues indicates an expected call of FetchWorkspaceIssues.
func (mr *MockClientMockRecorder) FetchWorkspaceIssues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchWorkspaceIssues", reflect.TypeOf((*MockClient)(nil).FetchWorkspaceIssues), arg0)
}