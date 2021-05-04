// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/chat (interfaces: ChatRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
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

// GetRestaurant mocks base method
func (m *MockChatRepo) GetRestaurant(arg0 context.Context, arg1 int) (*models.InfoOpponent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRestaurant", arg0, arg1)
	ret0, _ := ret[0].(*models.InfoOpponent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRestaurant indicates an expected call of GetRestaurant
func (mr *MockChatRepoMockRecorder) GetRestaurant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRestaurant", reflect.TypeOf((*MockChatRepo)(nil).GetRestaurant), arg0, arg1)
}

// GetUser mocks base method
func (m *MockChatRepo) GetUser(arg0 context.Context, arg1 int) (*models.InfoOpponent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*models.InfoOpponent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockChatRepoMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockChatRepo)(nil).GetUser), arg0, arg1)
}

// SaveMessageFromRestaurant mocks base method
func (m *MockChatRepo) SaveMessageFromRestaurant(arg0 context.Context, arg1 models.WsMessageForRepo) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMessageFromRestaurant", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveMessageFromRestaurant indicates an expected call of SaveMessageFromRestaurant
func (mr *MockChatRepoMockRecorder) SaveMessageFromRestaurant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMessageFromRestaurant", reflect.TypeOf((*MockChatRepo)(nil).SaveMessageFromRestaurant), arg0, arg1)
}

// SaveMessageFromUser mocks base method
func (m *MockChatRepo) SaveMessageFromUser(arg0 context.Context, arg1 models.WsMessageForRepo) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMessageFromUser", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveMessageFromUser indicates an expected call of SaveMessageFromUser
func (mr *MockChatRepoMockRecorder) SaveMessageFromUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMessageFromUser", reflect.TypeOf((*MockChatRepo)(nil).SaveMessageFromUser), arg0, arg1)
}
