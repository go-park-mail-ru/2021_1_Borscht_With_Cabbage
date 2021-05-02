package internal

//import (
//	"github.com/borscht/backend/internal/image/mocks"
//	"github.com/borscht/backend/internal/models"
//	adminMock "github.com/borscht/backend/internal/restaurantAdmin/mocks"
//	"github.com/borscht/backend/utils/logger"
//	"github.com/borscht/backend/utils/secure"
//	"github.com/golang/mock/gomock"
//	"testing"
//)

//func TestRestaurantUsecase_CreateRestaurant(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
//	imageRepoMock := mocks.NewMockImageRepo(ctrl)
//
//	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
//	c := context.Background()
//	ctx := context.WithValue(c, "request_id", 1)
//
//	logger.InitLogger()
//
//	restaurant := models.RestaurantInfo{
//		AdminEmail:   "dasha@mail.ru",
//		AdminPhone:   "89111111111",
//		Title:        "rest1",
//		Description:  "hey",
//		DeliveryCost: 200,
//	}
//	restaurantWithAvatar := models.RestaurantInfo{
//		AdminEmail:        "dasha@mail.ru",
//		AdminPhone:        "89111111111",
//		Title:             "rest1",
//		Description:       "hey",
//		DeliveryCost:      200,
//		Avatar:            config.DefaultRestaurantImage,
//		AdminHashPassword: secure.HashPassword(ctx, secure.GetSalt(), restaurant.AdminPassword),
//	// ПРОБЛЕМА: соль рандомная
//	}
//
//	restaurantRepoMock.EXPECT().CreateRestaurant(ctx, restaurantWithAvatar).Return(1, nil)
//
//	restaurantResponse, err := restaurantUsecase.CreateRestaurant(ctx, restaurant)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//
//	require.EqualValues(t, restaurantResponse.ID, 1)
//}
//
//func TestRestaurantUsecase_CheckRestaurantExists(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
//	imageRepoMock := mocks.NewMockImageRepo(ctrl)
//
//	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
//	c := context.Background()
//	ctx := context.WithValue(c, "request_id", 1)
//
//	logger.InitLogger()
//
//	user := models.RestaurantAuth{
//		Login:    "dasha@mail.ru",
//		Password: "1111111",
//	}
//	restaurant := models.RestaurantInfo{
//		AdminHashPassword: secure.HashPassword(ctx, secure.GetSalt(), user.Password),
//	}
//
//	restaurantRepoMock.EXPECT().GetByLogin(ctx, user.Login).Return(&restaurant, nil)
//
//	_, err := restaurantUsecase.CheckRestaurantExists(ctx, user)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//}
//
//func TestRestaurantUsecase_GetByRid(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	restaurantRepoMock := adminMock.NewMockAdminRestaurantRepo(ctrl)
//	imageRepoMock := mocks.NewMockImageRepo(ctrl)
//
//	restaurantUsecase := NewRestaurantUsecase(restaurantRepoMock, imageRepoMock)
//	ctx := new(context.Context)
//
//	restaurant := models.RestaurantInfo{}
//	restaurantRepoMock.EXPECT().GetByRid(*ctx, 1).Return(&restaurant, nil)
//
//	_, err := restaurantUsecase.GetByRid(*ctx, 1)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//}
