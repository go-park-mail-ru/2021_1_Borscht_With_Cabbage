package usecase

import (
	"context"
	"database/sql"
	"github.com/borscht/backend/internal/models"
	mocks2 "github.com/borscht/backend/internal/restaurant/mocks"
	"github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewRestaurantUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := mocks2.NewMockRestaurantRepo(ctrl)
	repoAdminMock := mocks.NewMockAdminRestaurantRepo(ctrl)

	userUsecase := NewRestaurantUsecase(restaurantRepoMock, repoAdminMock)
	if userUsecase == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantUsecase_GetVendor(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//restaurantRepoMock := mocks2.NewMockRestaurantRepo(ctrl)
	//repoAdminMock := mocks.NewMockAdminRestaurantRepo(ctrl)
	//
	//userUsecase := NewRestaurantUsecase(restaurantRepoMock, repoAdminMock)
	//c := context.Background()
	//
	//request := models.RestaurantRequest{}
	//request2 := models.RestaurantRequest{Receipt: 100000, Time: 1000}
	//restaurants := []models.RestaurantInfo{{ID: 1}}
	//restaurantRepoMock.EXPECT().GetAllCategories(c).Return(make([]string, 0), nil)
	//restaurantRepoMock.EXPECT().GetVendor(c, request2).Return(restaurants, nil)
	//repoAdminMock.EXPECT().GetAddress(c, 1).Return(&models.Address{}, nil)
	//
	//_, err := userUsecase.GetVendor(c, request)
	//if err != nil {
	//	t.Errorf("unexpected err: %s", err)
	//	return
	//}
}

func TestRestaurantUsecase_GetVendor_CategoriesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := mocks2.NewMockRestaurantRepo(ctrl)
	repoAdminMock := mocks.NewMockAdminRestaurantRepo(ctrl)

	userUsecase := NewRestaurantUsecase(restaurantRepoMock, repoAdminMock)
	c := context.Background()

	request := models.RestaurantRequest{}

	restaurantRepoMock.EXPECT().GetAllCategories(c).Return(make([]string, 0), sql.ErrNoRows)

	_, err := userUsecase.GetVendor(c, request)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestRestaurantUsecase_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := mocks2.NewMockRestaurantRepo(ctrl)
	repoAdminMock := mocks.NewMockAdminRestaurantRepo(ctrl)

	userUsecase := NewRestaurantUsecase(restaurantRepoMock, repoAdminMock)
	c := context.Background()

	restaurantRepoMock.EXPECT().GetById(c, 1, models.Coordinates{}).Return(&models.RestaurantWithDishes{}, nil)
	repoAdminMock.EXPECT().GetAddress(c, 0).Return(&models.Address{}, nil)

	_, err := userUsecase.GetById(c, models.Coordinates{}, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}
