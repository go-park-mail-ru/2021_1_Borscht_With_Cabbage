// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/restaurantAdmin (interfaces: AdminHandler)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v4 "github.com/labstack/echo/v4"
	reflect "reflect"
)

// MockAdminHandler is a mock of AdminHandler interface
type MockAdminHandler struct {
	ctrl     *gomock.Controller
	recorder *MockAdminHandlerMockRecorder
}

// MockAdminHandlerMockRecorder is the mock recorder for MockAdminHandler
type MockAdminHandlerMockRecorder struct {
	mock *MockAdminHandler
}

// NewMockAdminHandler creates a new mock instance
func NewMockAdminHandler(ctrl *gomock.Controller) *MockAdminHandler {
	mock := &MockAdminHandler{ctrl: ctrl}
	mock.recorder = &MockAdminHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdminHandler) EXPECT() *MockAdminHandlerMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockAdminHandler) Create(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockAdminHandlerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAdminHandler)(nil).Create), arg0)
}

// EditProfile mocks base method
func (m *MockAdminHandler) EditProfile(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProfile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditProfile indicates an expected call of EditProfile
func (mr *MockAdminHandlerMockRecorder) EditProfile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProfile", reflect.TypeOf((*MockAdminHandler)(nil).EditProfile), arg0)
}

// GetUserData mocks base method
func (m *MockAdminHandler) GetUserData(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUserData indicates an expected call of GetUserData
func (mr *MockAdminHandlerMockRecorder) GetUserData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockAdminHandler)(nil).GetUserData), arg0)
}

// Login mocks base method
func (m *MockAdminHandler) Login(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login
func (mr *MockAdminHandlerMockRecorder) Login(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAdminHandler)(nil).Login), arg0)
}