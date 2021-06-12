// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/borscht/backend/internal/restaurant (interfaces: RestaurantRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "github.com/borscht/backend/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRestaurantRepo is a mock of RestaurantRepo interface
type MockRestaurantRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRestaurantRepoMockRecorder
}

// MockRestaurantRepoMockRecorder is the mock recorder for MockRestaurantRepo
type MockRestaurantRepoMockRecorder struct {
	mock *MockRestaurantRepo
}

// NewMockRestaurantRepo creates a new mock instance
func NewMockRestaurantRepo(ctrl *gomock.Controller) *MockRestaurantRepo {
	mock := &MockRestaurantRepo{ctrl: ctrl}
	mock.recorder = &MockRestaurantRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRestaurantRepo) EXPECT() *MockRestaurantRepoMockRecorder {
	return m.recorder
}

// GetAllCategories mocks base method
func (m *MockRestaurantRepo) GetAllCategories(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCategories", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCategories indicates an expected call of GetAllCategories
func (mr *MockRestaurantRepoMockRecorder) GetAllCategories(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCategories", reflect.TypeOf((*MockRestaurantRepo)(nil).GetAllCategories), arg0)
}

// GetById mocks base method
func (m *MockRestaurantRepo) GetById(arg0 context.Context, arg1 int, arg2 models.Coordinates) (*models.RestaurantWithDishes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.RestaurantWithDishes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById
func (mr *MockRestaurantRepoMockRecorder) GetById(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockRestaurantRepo)(nil).GetById), arg0, arg1, arg2)
}

// GetRecommendations mocks base method
func (m *MockRestaurantRepo) GetRecommendations(arg0 context.Context, arg1 models.RecommendationsParams) ([]models.RestaurantInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecommendations", arg0, arg1)
	ret0, _ := ret[0].([]models.RestaurantInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecommendations indicates an expected call of GetRecommendations
func (mr *MockRestaurantRepoMockRecorder) GetRecommendations(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecommendations", reflect.TypeOf((*MockRestaurantRepo)(nil).GetRecommendations), arg0, arg1)
}

// GetReviews mocks base method
func (m *MockRestaurantRepo) GetReviews(arg0 context.Context, arg1 int) ([]models.RestaurantReview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReviews", arg0, arg1)
	ret0, _ := ret[0].([]models.RestaurantReview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReviews indicates an expected call of GetReviews
func (mr *MockRestaurantRepoMockRecorder) GetReviews(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReviews", reflect.TypeOf((*MockRestaurantRepo)(nil).GetReviews), arg0, arg1)
}

// GetVendor mocks base method
func (m *MockRestaurantRepo) GetVendor(arg0 context.Context, arg1 models.RestaurantRequest) ([]models.RestaurantInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVendor", arg0, arg1)
	ret0, _ := ret[0].([]models.RestaurantInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVendor indicates an expected call of GetVendor
func (mr *MockRestaurantRepoMockRecorder) GetVendor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVendor", reflect.TypeOf((*MockRestaurantRepo)(nil).GetVendor), arg0, arg1)
}