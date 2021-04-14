package usecase

import (
	"context"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/image/mocks"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRestaurantUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)

	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
	if restaurantUsecase == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestRestaurantUsecase_UpdateRestaurantData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)

	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)

	restaurantAdmin := models.RestaurantInfo{
		ID:       1,
		Rating:   4,
		AvgCheck: 1400,
		Avatar:   config.DefaultAvatar,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurantAdmin)

	restaurant := models.RestaurantUpdateData{
		ID:           1,
		AdminEmail:   "dasha@mail.ru",
		AdminPhone:   "89111111111",
		Title:        "rest1",
		Description:  "hey",
		DeliveryCost: 200,
	}

	restaurantRepoMock.EXPECT().UpdateRestaurantData(ctx, restaurant).Return(nil)

	restaurantResponse, err := restaurantUsecase.UpdateRestaurantData(ctx, restaurant)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurantResponse.ID, 1)
}

func TestRestaurantUsecase_CreateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)

	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
	ctx := new(context.Context)

	restaurant := models.RestaurantInfo{
		AdminEmail:   "dasha@mail.ru",
		AdminPhone:   "89111111111",
		Title:        "rest1",
		Description:  "hey",
		DeliveryCost: 200,
	}
	restaurantWithAvatar := models.RestaurantInfo{
		AdminEmail:   "dasha@mail.ru",
		AdminPhone:   "89111111111",
		Title:        "rest1",
		Description:  "hey",
		DeliveryCost: 200,
		Avatar:       config.DefaultAvatar,
	}

	restaurantRepoMock.EXPECT().CreateRestaurant(*ctx, restaurantWithAvatar).Return(1, nil)

	restaurantResponse, err := restaurantUsecase.CreateRestaurant(*ctx, restaurant)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurantResponse.ID, 1)
}

func TestRestaurantUsecase_CheckRestaurantExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)

	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
	ctx := new(context.Context)

	user := models.RestaurantAuth{
		Login:    "dasha@mail.ru",
		Password: "1111111",
	}
	restaurant := models.RestaurantInfo{}

	restaurantRepoMock.EXPECT().CheckRestaurantExists(*ctx, user).Return(&restaurant, nil)

	_, err := restaurantUsecase.CheckRestaurantExists(*ctx, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestRestaurantUsecase_GetByRid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
	imageRepoMock := mocks.NewMockImageRepo(ctrl)

	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
	ctx := new(context.Context)

	restaurant := models.RestaurantInfo{}
	restaurantRepoMock.EXPECT().GetByRid(*ctx, 1).Return(&restaurant, nil)

	_, err := restaurantUsecase.GetByRid(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestRestaurantUsecase_UploadRestaurantImage(t *testing.T) {
	// todo
}
