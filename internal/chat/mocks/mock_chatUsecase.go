// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/chat (interfaces: ChatUsecase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
	gomock "github.com/golang/mock/gomock"
	websocket "github.com/gorilla/websocket"
	reflect "reflect"
)

// MockChatUsecase is a mock of ChatUsecase interface
type MockChatUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockChatUsecaseMockRecorder
}

// MockChatUsecaseMockRecorder is the mock recorder for MockChatUsecase
type MockChatUsecaseMockRecorder struct {
	mock *MockChatUsecase
}

// NewMockChatUsecase creates a new mock instance
func NewMockChatUsecase(ctrl *gomock.Controller) *MockChatUsecase {
	mock := &MockChatUsecase{ctrl: ctrl}
	mock.recorder = &MockChatUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChatUsecase) EXPECT() *MockChatUsecaseMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockChatUsecase) Connect(arg0 context.Context, arg1 *websocket.Conn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockChatUsecaseMockRecorder) Connect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockChatUsecase)(nil).Connect), arg0, arg1)
}

// GetAllChats mocks base method
func (m *MockChatUsecase) GetAllChats(arg0 context.Context) ([]models.BriefInfoChat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChats", arg0)
	ret0, _ := ret[0].([]models.BriefInfoChat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllChats indicates an expected call of GetAllChats
func (mr *MockChatUsecaseMockRecorder) GetAllChats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChats", reflect.TypeOf((*MockChatUsecase)(nil).GetAllChats), arg0)
}

// GetAllMessages mocks base method
func (m *MockChatUsecase) GetAllMessages(arg0 context.Context, arg1 int) (*models.InfoChat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMessages", arg0, arg1)
	ret0, _ := ret[0].(*models.InfoChat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMessages indicates an expected call of GetAllMessages
func (mr *MockChatUsecaseMockRecorder) GetAllMessages(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMessages", reflect.TypeOf((*MockChatUsecase)(nil).GetAllMessages), arg0, arg1)
}

// ProcessMessage mocks base method
func (m *MockChatUsecase) ProcessMessage(arg0 context.Context, arg1 *websocket.Conn, arg2 models.FromClient) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessMessage indicates an expected call of ProcessMessage
func (mr *MockChatUsecaseMockRecorder) ProcessMessage(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessMessage", reflect.TypeOf((*MockChatUsecase)(nil).ProcessMessage), arg0, arg1, arg2)
}

// UnConnect mocks base method
func (m *MockChatUsecase) UnConnect(arg0 context.Context, arg1 *websocket.Conn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnConnect", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnConnect indicates an expected call of UnConnect
func (mr *MockChatUsecaseMockRecorder) UnConnect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnConnect", reflect.TypeOf((*MockChatUsecase)(nil).UnConnect), arg0, arg1)
}
