// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/order (interfaces: OrderUsecase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockOrderUsecase is a mock of OrderUsecase interface
type MockOrderUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUsecaseMockRecorder
}

// MockOrderUsecaseMockRecorder is the mock recorder for MockOrderUsecase
type MockOrderUsecaseMockRecorder struct {
	mock *MockOrderUsecase
}

// NewMockOrderUsecase creates a new mock instance
func NewMockOrderUsecase(ctrl *gomock.Controller) *MockOrderUsecase {
	mock := &MockOrderUsecase{ctrl: ctrl}
	mock.recorder = &MockOrderUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOrderUsecase) EXPECT() *MockOrderUsecaseMockRecorder {
	return m.recorder
}

// AddBasket mocks base method
func (m *MockOrderUsecase) AddBasket(arg0 context.Context, arg1 models.BasketForUser) (*models.BasketForUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBasket", arg0, arg1)
	ret0, _ := ret[0].(models.BasketForUser)
	ret1, _ := ret[1].(error)
	return &ret0, ret1
}

// AddToBasket mocks base method
func (m *MockOrderUsecase) AddToBasket(arg0 context.Context, arg1 models.DishToBasket, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBasket", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBasket indicates an expected call of AddToBasket
func (mr *MockOrderUsecaseMockRecorder) AddToBasket(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBasket", reflect.TypeOf((*MockOrderUsecase)(nil).AddToBasket), arg0, arg1, arg2)
}

// Create mocks base method
func (m *MockOrderUsecase) Create(arg0 context.Context, arg1 int, arg2 models.CreateOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockOrderUsecaseMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderUsecase)(nil).Create), arg0, arg1, arg2)
}

// DeleteFromBasket mocks base method
func (m *MockOrderUsecase) DeleteFromBasket(arg0 context.Context, arg1 models.DishToBasket, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromBasket", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromBasket indicates an expected call of DeleteFromBasket
func (mr *MockOrderUsecaseMockRecorder) DeleteFromBasket(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromBasket", reflect.TypeOf((*MockOrderUsecase)(nil).DeleteFromBasket), arg0, arg1, arg2)
}

// GetBasket mocks base method
func (m *MockOrderUsecase) GetBasket(arg0 context.Context, arg1 int) (*models.BasketForUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBasket", arg0, arg1)
	ret0, _ := ret[0].(models.BasketForUser)
	ret1, _ := ret[1].(error)
	return &ret0, ret1
}

// GetBasket indicates an expected call of GetBasket
func (mr *MockOrderUsecaseMockRecorder) GetBasket(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBasket", reflect.TypeOf((*MockOrderUsecase)(nil).GetBasket), arg0, arg1)
}

// GetRestaurantOrders mocks base method
func (m *MockOrderUsecase) GetRestaurantOrders(arg0 context.Context, arg1 string) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRestaurantOrders", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRestaurantOrders indicates an expected call of GetRestaurantOrders
func (mr *MockOrderUsecaseMockRecorder) GetRestaurantOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRestaurantOrders", reflect.TypeOf((*MockOrderUsecase)(nil).GetRestaurantOrders), arg0, arg1)
}

// GetUserOrders mocks base method
func (m *MockOrderUsecase) GetUserOrders(arg0 context.Context, arg1 int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders
func (mr *MockOrderUsecaseMockRecorder) GetUserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockOrderUsecase)(nil).GetUserOrders), arg0, arg1)
}
