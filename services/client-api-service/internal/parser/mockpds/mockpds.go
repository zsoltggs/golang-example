// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zsoltggs/golang-example/pkg/pds (interfaces: ServiceClient)

// Package mockpds is a generated GoMock package.
package mockpds

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pds "github.com/zsoltggs/golang-example/pkg/pds"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockServiceClient is a mock of ServiceClient interface.
type MockServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockServiceClientMockRecorder
}

// MockServiceClientMockRecorder is the mock recorder for MockServiceClient.
type MockServiceClientMockRecorder struct {
	mock *MockServiceClient
}

// NewMockServiceClient creates a new mock instance.
func NewMockServiceClient(ctrl *gomock.Controller) *MockServiceClient {
	mock := &MockServiceClient{ctrl: ctrl}
	mock.recorder = &MockServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceClient) EXPECT() *MockServiceClientMockRecorder {
	return m.recorder
}

// GetPortByID mocks base method.
func (m *MockServiceClient) GetPortByID(arg0 context.Context, arg1 *pds.GetPortByIDRequest, arg2 ...grpc.CallOption) (*pds.GetPortByIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPortByID", varargs...)
	ret0, _ := ret[0].(*pds.GetPortByIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPortByID indicates an expected call of GetPortByID.
func (mr *MockServiceClientMockRecorder) GetPortByID(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPortByID", reflect.TypeOf((*MockServiceClient)(nil).GetPortByID), varargs...)
}

// GetPortsPaginated mocks base method.
func (m *MockServiceClient) GetPortsPaginated(arg0 context.Context, arg1 *pds.GetPortsPaginatedRequest, arg2 ...grpc.CallOption) (*pds.GetPortsPaginatedResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPortsPaginated", varargs...)
	ret0, _ := ret[0].(*pds.GetPortsPaginatedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPortsPaginated indicates an expected call of GetPortsPaginated.
func (mr *MockServiceClientMockRecorder) GetPortsPaginated(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPortsPaginated", reflect.TypeOf((*MockServiceClient)(nil).GetPortsPaginated), varargs...)
}

// UpsertPort mocks base method.
func (m *MockServiceClient) UpsertPort(arg0 context.Context, arg1 *pds.UpsertPortRequest, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertPort", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertPort indicates an expected call of UpsertPort.
func (mr *MockServiceClientMockRecorder) UpsertPort(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertPort", reflect.TypeOf((*MockServiceClient)(nil).UpsertPort), varargs...)
}
