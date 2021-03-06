// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/chat (interfaces: ChatHandler)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v4 "github.com/labstack/echo/v4"
	reflect "reflect"
)

// MockChatHandler is a mock of ChatHandler interface
type MockChatHandler struct {
	ctrl     *gomock.Controller
	recorder *MockChatHandlerMockRecorder
}

// MockChatHandlerMockRecorder is the mock recorder for MockChatHandler
type MockChatHandlerMockRecorder struct {
	mock *MockChatHandler
}

// NewMockChatHandler creates a new mock instance
func NewMockChatHandler(ctrl *gomock.Controller) *MockChatHandler {
	mock := &MockChatHandler{ctrl: ctrl}
	mock.recorder = &MockChatHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChatHandler) EXPECT() *MockChatHandlerMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockChatHandler) Connect(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockChatHandlerMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockChatHandler)(nil).Connect), arg0)
}

// GetAllChats mocks base method
func (m *MockChatHandler) GetAllChats(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChats", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAllChats indicates an expected call of GetAllChats
func (mr *MockChatHandlerMockRecorder) GetAllChats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChats", reflect.TypeOf((*MockChatHandler)(nil).GetAllChats), arg0)
}

// GetAllMessages mocks base method
func (m *MockChatHandler) GetAllMessages(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMessages", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAllMessages indicates an expected call of GetAllMessages
func (mr *MockChatHandlerMockRecorder) GetAllMessages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMessages", reflect.TypeOf((*MockChatHandler)(nil).GetAllMessages), arg0)
}

// GetKey mocks base method
func (m *MockChatHandler) GetKey(arg0 v4.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetKey indicates an expected call of GetKey
func (mr *MockChatHandlerMockRecorder) GetKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockChatHandler)(nil).GetKey), arg0)
}
