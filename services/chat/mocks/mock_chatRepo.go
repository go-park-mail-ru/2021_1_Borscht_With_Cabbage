// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/services/chat/repository (interfaces: ChatRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	chat "github.com/borscht/backend/services/proto/chat"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockChatRepo is a mock of ChatRepo interface
type MockChatRepo struct {
	ctrl     *gomock.Controller
	recorder *MockChatRepoMockRecorder
}

// MockChatRepoMockRecorder is the mock recorder for MockChatRepo
type MockChatRepoMockRecorder struct {
	mock *MockChatRepo
}

// NewMockChatRepo creates a new mock instance
func NewMockChatRepo(ctrl *gomock.Controller) *MockChatRepo {
	mock := &MockChatRepo{ctrl: ctrl}
	mock.recorder = &MockChatRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChatRepo) EXPECT() *MockChatRepoMockRecorder {
	return m.recorder
}

// GetAllChatsRestaurant mocks base method
func (m *MockChatRepo) GetAllChatsRestaurant(arg0 context.Context, arg1 int) ([]*chat.BriefInfoChat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChatsRestaurant", arg0, arg1)
	ret0, _ := ret[0].([]*chat.BriefInfoChat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllChatsRestaurant indicates an expected call of GetAllChatsRestaurant
func (mr *MockChatRepoMockRecorder) GetAllChatsRestaurant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChatsRestaurant", reflect.TypeOf((*MockChatRepo)(nil).GetAllChatsRestaurant), arg0, arg1)
}

// GetAllChatsUser mocks base method
func (m *MockChatRepo) GetAllChatsUser(arg0 context.Context, arg1 int) ([]*chat.BriefInfoChat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChatsUser", arg0, arg1)
	ret0, _ := ret[0].([]*chat.BriefInfoChat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllChatsUser indicates an expected call of GetAllChatsUser
func (mr *MockChatRepoMockRecorder) GetAllChatsUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChatsUser", reflect.TypeOf((*MockChatRepo)(nil).GetAllChatsUser), arg0, arg1)
}

// GetAllMessagesFromRestaurant mocks base method
func (m *MockChatRepo) GetAllMessagesFromRestaurant(arg0 context.Context, arg1, arg2 int) ([]*chat.InfoMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMessagesFromRestaurant", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*chat.InfoMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMessagesFromRestaurant indicates an expected call of GetAllMessagesFromRestaurant
func (mr *MockChatRepoMockRecorder) GetAllMessagesFromRestaurant(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMessagesFromRestaurant", reflect.TypeOf((*MockChatRepo)(nil).GetAllMessagesFromRestaurant), arg0, arg1, arg2)
}

// GetAllMessagesFromUser mocks base method
func (m *MockChatRepo) GetAllMessagesFromUser(arg0 context.Context, arg1, arg2 int) ([]*chat.InfoMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMessagesFromUser", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*chat.InfoMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMessagesFromUser indicates an expected call of GetAllMessagesFromUser
func (mr *MockChatRepoMockRecorder) GetAllMessagesFromUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMessagesFromUser", reflect.TypeOf((*MockChatRepo)(nil).GetAllMessagesFromUser), arg0, arg1, arg2)
}