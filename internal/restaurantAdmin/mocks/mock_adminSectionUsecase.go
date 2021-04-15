// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/restaurantAdmin (interfaces: AdminSectionUsecase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAdminSectionUsecase is a mock of AdminSectionUsecase interface
type MockAdminSectionUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAdminSectionUsecaseMockRecorder
}

// MockAdminSectionUsecaseMockRecorder is the mock recorder for MockAdminSectionUsecase
type MockAdminSectionUsecaseMockRecorder struct {
	mock *MockAdminSectionUsecase
}

// NewMockAdminSectionUsecase creates a new mock instance
func NewMockAdminSectionUsecase(ctrl *gomock.Controller) *MockAdminSectionUsecase {
	mock := &MockAdminSectionUsecase{ctrl: ctrl}
	mock.recorder = &MockAdminSectionUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdminSectionUsecase) EXPECT() *MockAdminSectionUsecaseMockRecorder {
	return m.recorder
}

// AddSection mocks base method
func (m *MockAdminSectionUsecase) AddSection(arg0 context.Context, arg1 models.Section) (*models.Section, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSection", arg0, arg1)
	ret0, _ := ret[0].(*models.Section)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSection indicates an expected call of AddSection
func (mr *MockAdminSectionUsecaseMockRecorder) AddSection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSection", reflect.TypeOf((*MockAdminSectionUsecase)(nil).AddSection), arg0, arg1)
}

// DeleteSection mocks base method
func (m *MockAdminSectionUsecase) DeleteSection(arg0 context.Context, arg1 int) (*models.DeleteSuccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSection", arg0, arg1)
	ret0, _ := ret[0].(*models.DeleteSuccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSection indicates an expected call of DeleteSection
func (mr *MockAdminSectionUsecaseMockRecorder) DeleteSection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSection", reflect.TypeOf((*MockAdminSectionUsecase)(nil).DeleteSection), arg0, arg1)
}

// UpdateSection mocks base method
func (m *MockAdminSectionUsecase) UpdateSection(arg0 context.Context, arg1 models.Section) (*models.Section, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSection", arg0, arg1)
	ret0, _ := ret[0].(*models.Section)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSection indicates an expected call of UpdateSection
func (mr *MockAdminSectionUsecaseMockRecorder) UpdateSection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSection", reflect.TypeOf((*MockAdminSectionUsecase)(nil).UpdateSection), arg0, arg1)
}
