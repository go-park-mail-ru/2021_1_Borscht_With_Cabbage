// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/user (interfaces: UserUsecase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
	gomock "github.com/golang/mock/gomock"
	multipart "mime/multipart"
	reflect "reflect"
)

// MockUserUsecase is a mock of UserUsecase interface
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// AddAddress mocks base method
func (m *MockUserUsecase) AddAddress(arg0 context.Context, arg1 int, arg2 models.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAddress", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAddress indicates an expected call of AddAddress
func (mr *MockUserUsecaseMockRecorder) AddAddress(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAddress", reflect.TypeOf((*MockUserUsecase)(nil).AddAddress), arg0, arg1, arg2)
}

// GetAddress mocks base method
func (m *MockUserUsecase) GetAddress(arg0 context.Context, arg1 int) (*models.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress", arg0, arg1)
	ret0, _ := ret[0].(*models.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAddress indicates an expected call of GetAddress
func (mr *MockUserUsecaseMockRecorder) GetAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockUserUsecase)(nil).GetAddress), arg0, arg1)
}

// GetMainAddress mocks base method
func (m *MockUserUsecase) GetMainAddress(arg0 context.Context) (*models.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMainAddress", arg0)
	ret0, _ := ret[0].(*models.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMainAddress indicates an expected call of GetMainAddress
func (mr *MockUserUsecaseMockRecorder) GetMainAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMainAddress", reflect.TypeOf((*MockUserUsecase)(nil).GetMainAddress), arg0)
}

// GetUserData mocks base method
func (m *MockUserUsecase) GetUserData(arg0 context.Context) (*models.SuccessUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", arg0)
	ret0, _ := ret[0].(*models.SuccessUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserData indicates an expected call of GetUserData
func (mr *MockUserUsecaseMockRecorder) GetUserData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockUserUsecase)(nil).GetUserData), arg0)
}

// UpdateData mocks base method
func (m *MockUserUsecase) UpdateData(arg0 context.Context, arg1 models.UserData) (*models.SuccessUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateData", arg0, arg1)
	ret0, _ := ret[0].(*models.SuccessUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateData indicates an expected call of UpdateData
func (mr *MockUserUsecaseMockRecorder) UpdateData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateData", reflect.TypeOf((*MockUserUsecase)(nil).UpdateData), arg0, arg1)
}

// UpdateMainAddress mocks base method
func (m *MockUserUsecase) UpdateMainAddress(arg0 context.Context, arg1 models.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMainAddress", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMainAddress indicates an expected call of UpdateMainAddress
func (mr *MockUserUsecaseMockRecorder) UpdateMainAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMainAddress", reflect.TypeOf((*MockUserUsecase)(nil).UpdateMainAddress), arg0, arg1)
}

// UploadAvatar mocks base method
func (m *MockUserUsecase) UploadAvatar(arg0 context.Context, arg1 *multipart.FileHeader) (*models.UserImageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadAvatar", arg0, arg1)
	ret0, _ := ret[0].(*models.UserImageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadAvatar indicates an expected call of UploadAvatar
func (mr *MockUserUsecaseMockRecorder) UploadAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadAvatar", reflect.TypeOf((*MockUserUsecase)(nil).UploadAvatar), arg0, arg1)
}
