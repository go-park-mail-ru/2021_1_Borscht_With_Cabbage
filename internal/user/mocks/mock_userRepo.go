// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/user (interfaces: UserRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUserRepo is a mock of UserRepo interface
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// AddAddress mocks base method
func (m *MockUserRepo) AddAddress(arg0 context.Context, arg1 int, arg2 models.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAddress", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAddress indicates an expected call of AddAddress
func (mr *MockUserRepoMockRecorder) AddAddress(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAddress", reflect.TypeOf((*MockUserRepo)(nil).AddAddress), arg0, arg1, arg2)
}

// DeleteAddress mocks base method
func (m *MockUserRepo) DeleteAddress(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAddress", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAddress indicates an expected call of DeleteAddress
func (mr *MockUserRepoMockRecorder) DeleteAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAddress", reflect.TypeOf((*MockUserRepo)(nil).DeleteAddress), arg0, arg1)
}

// GetAddress mocks base method
func (m *MockUserRepo) GetAddress(arg0 context.Context, arg1 int) (*models.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress", arg0, arg1)
	ret0, _ := ret[0].(*models.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAddress indicates an expected call of GetAddress
func (mr *MockUserRepoMockRecorder) GetAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockUserRepo)(nil).GetAddress), arg0, arg1)
}

// UpdateAvatar mocks base method
func (m *MockUserRepo) UpdateAvatar(arg0 context.Context, arg1 int, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatar", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAvatar indicates an expected call of UpdateAvatar
func (mr *MockUserRepoMockRecorder) UpdateAvatar(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatar", reflect.TypeOf((*MockUserRepo)(nil).UpdateAvatar), arg0, arg1, arg2)
}

// UpdateData mocks base method
func (m *MockUserRepo) UpdateData(arg0 context.Context, arg1 models.UserData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateData", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateData indicates an expected call of UpdateData
func (mr *MockUserRepoMockRecorder) UpdateData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateData", reflect.TypeOf((*MockUserRepo)(nil).UpdateData), arg0, arg1)
}
