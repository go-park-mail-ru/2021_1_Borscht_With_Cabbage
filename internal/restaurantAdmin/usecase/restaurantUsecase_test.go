package usecase

import (
	"context"
	"testing"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/image/mocks"
	"github.com/borscht/backend/internal/models"
	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
		Avatar:   config.DefaultRestaurantImage,
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
	restaurantRepoMock.EXPECT().UpdateAddress(ctx, restaurant.ID, restaurant.Address)

	restaurantResponse, err := restaurantUsecase.UpdateRestaurantData(ctx, restaurant)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurantResponse.ID, 1)
}

func TestRestaurantUsecase_UploadRestaurantImage(t *testing.T) {
	// todo
}
