// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/user (interfaces: UserHandler)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v4 "github.com/labstack/echo/v4"
	reflect "reflect"
)

// MockUserHandler is a mock of UserHandler interface
type MockUserHandler struct {
	ctrl     *gomock.Controller
	recorder *MockUserHandlerMockRecorder
}

// MockUserHandlerMockRecorder is the mock recorder for MockUserHandler
type MockUserHandlerMockRecorder struct {
	mock *MockUserHandler
}

// NewMockUserHandler creates a new mock instance
func NewMockUserHandler(ctrl *gomock.Controller) *MockUserHandler {
	mock := &MockUserHandler{ctrl: ctrl}
	mock.recorder = &MockUserHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserHandler) EXPECT() *MockUserHandlerMockRecorder {
	return m.recorder
}

// CheckAuth mocks base method
func (m *MockUserHandler) CheckAuth(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckAuth indicates an expected call of CheckAuth
func (mr *MockUserHandlerMockRecorder) CheckAuth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockUserHandler)(nil).CheckAuth), arg0)
}

// Create mocks base method
func (m *MockUserHandler) Create(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockUserHandlerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserHandler)(nil).Create), arg0)
}

// EditProfile mocks base method
func (m *MockUserHandler) EditProfile(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProfile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditProfile indicates an expected call of EditProfile
func (mr *MockUserHandlerMockRecorder) EditProfile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProfile", reflect.TypeOf((*MockUserHandler)(nil).EditProfile), arg0)
}

// GetUserData mocks base method
func (m *MockUserHandler) GetUserData(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUserData indicates an expected call of GetUserData
func (mr *MockUserHandlerMockRecorder) GetUserData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockUserHandler)(nil).GetUserData), arg0)
}

// Login mocks base method
func (m *MockUserHandler) Login(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login
func (mr *MockUserHandlerMockRecorder) Login(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserHandler)(nil).Login), arg0)
}

// Logout mocks base method
func (m *MockUserHandler) Logout(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout
func (mr *MockUserHandlerMockRecorder) Logout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockUserHandler)(nil).Logout), arg0)
}