package internal

import (
	"context"
	"testing"

	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	authRepoMocks "github.com/borscht/backend/services/auth/mocks"
	protoAuth "github.com/borscht/backend/services/proto/auth"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/secure"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//func TestService_CreateUser(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
//	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
//	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
//	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)
//
//	c := context.Background()
//	ctx := context.WithValue(c, "request_id", 1)
//
//	logger.InitLogger()
//
//	user := models.User{
//		Email:    "dasha@mail.ru",
//		Phone:    "89111111111",
//		Name:     "111111",
//		Password: "1111111",
//		Avatar:   config.DefaultUserImage,
//	}
//	userProto := protoAuth.User{
//		Email:    user.Email,
//		Phone:    user.Phone,
//		Name:     user.Name,
//		Password: user.Password,
//	}
//
//	response := models.SuccessUserResponse{
//		User: user,
//		Role: config.RoleUser,
//	}
//	response.Uid = 1
//
//	userAuthRepoMock.EXPECT().Create(ctx, user).Return(1, nil)
//
//	userResponse, err := authService.CreateUser(ctx, &userProto)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//
//	require.EqualValues(t, userResponse.UID, 1)
//	require.EqualValues(t, userResponse.Email, "dasha@mail.ru")
//}

func TestService_GetByUid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	user := models.User{
		Email:        "dasha@mail.ru",
		Phone:        "89111111111",
		Name:         "111111",
		Password:     "1111111",
		HashPassword: secure.HashPassword(ctx, secure.GetSalt(), "111111"),
	}

	userAuthRepoMock.EXPECT().GetByUid(ctx, 1).Return(user, nil)
	userAuthRepoMock.EXPECT().GetAddress(ctx, 1).Return(&models.Address{}, nil)

	response, err := authService.GetByUid(ctx, &protoAuth.UID{Uid: 1})
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, configProject.RoleUser)
}

func TestService_CheckUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)
	logger.InitLogger()

	userAuth := models.UserAuth{
		Login:    "dasha@mail.ru",
		Password: "111111",
	}

	user := models.User{
		Email:        "dasha@mail.ru",
		Phone:        "89111111111",
		Name:         "111111",
		Password:     "1111111",
		Uid:          1,
		HashPassword: secure.HashPassword(ctx, secure.GetSalt(), "111111"),
	}
	userProto := protoAuth.UserAuth{
		Login:    user.Email,
		Password: user.Password,
	}

	userAuthRepoMock.EXPECT().GetByLogin(ctx, userAuth.Login).Return(&user, nil)
	userAuthRepoMock.EXPECT().GetAddress(ctx, 1).Return(&models.Address{}, nil)

	response, err := authService.CheckUserExists(ctx, &userProto)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, configProject.RoleUser)
}

//func TestRestaurantUsecase_CreateRestaurant(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
//	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
//	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
//	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)
//
//	c := context.Background()
//	ctx := context.WithValue(c, "request_id", 1)
//
//	logger.InitLogger()
//
//	restaurant := models.RestaurantInfo{
//		AdminEmail:    "dasha@mail.ru",
//		AdminPhone:    "89111111111",
//		Title:         "rest1",
//		AdminPassword: "111111",
//	}
//	restaurantWithAvatar := models.RestaurantInfo{
//		AdminEmail:        "dasha@mail.ru",
//		AdminPhone:        "89111111111",
//		Title:             "rest1",
//		Avatar:            config.DefaultRestaurantImage,
//		AdminHashPassword: secure.HashPassword(ctx, secure.GetSalt(), restaurant.AdminPassword),
//		// ПРОБЛЕМА: соль рандомная
//	}
//	restaurantProto := protoAuth.User{
//		Email:    restaurant.AdminEmail,
//		Phone:    restaurant.AdminPhone,
//		Name:     restaurant.Title,
//		Password: restaurant.AdminPassword,
//	}
//
//	restaurantAuthRepoMock.EXPECT().CreateRestaurant(ctx, restaurantWithAvatar).Return(1, nil)
//
//	restaurantResponse, err := authService.CreateRestaurant(ctx, &restaurantProto)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//
//	require.EqualValues(t, restaurantResponse.Email, "dasha@mail.ru")
//}

func TestService_CheckRestaurantExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	user := models.RestaurantAuth{
		Login:    "dasha@mail.ru",
		Password: "1111111",
	}
	restaurant := models.RestaurantInfo{
		AdminEmail:    "dasha@mail.ru",
		AdminPhone:    "81111111111",
		Title:         "rest1",
		AdminPassword: "111111",
	}
	userProto := protoAuth.UserAuth{
		Login:    user.Login,
		Password: user.Password,
	}

	restaurantAuthRepoMock.EXPECT().GetByLogin(ctx, user.Login).Return(&restaurant, nil)

	restaurantResponse, err := authService.CheckRestaurantExists(ctx, &userProto)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurantResponse.Email, "dasha@mail.ru")
}

func TestService_GetByRid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)

	ctx := new(context.Context)

	restaurant := models.RestaurantInfo{}
	restaurantAuthRepoMock.EXPECT().GetByRid(*ctx, 1).Return(&restaurant, nil)
	restaurantAuthRepoMock.EXPECT().GetAddress(*ctx, 1).Return(&models.Address{}, nil)

	_, err := authService.GetByRid(*ctx, &protoAuth.RID{Rid: 1})
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_CheckSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)

	ctx := new(context.Context)

	sessionInfo := models.SessionInfo{
		Id:   1,
		Role: configProject.RoleUser,
	}
	session := protoAuth.SessionValue{
		Session: "session1",
	}

	sessionRepoMock.EXPECT().Check(*ctx, headSession+session.Session).Return(sessionInfo, true, nil)

	sessionResult, err := authService.CheckSession(*ctx, &session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, sessionResult.Id, 1)
}

//func TestService_CreateSession(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
//	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
//	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
//	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)
//
//	ctx := new(context.Context)
//
//	session := protoAuth.SessionInfo{
//		Id: 1,
//		Role: config.RoleUser,
//		Exists: true,
//	}
//
//	// что делать с рандомной сессией?
//	sessionRepoMock.EXPECT().Create().Return()
//	_, err := authService.GetByRid(*ctx, &protoAuth.RID{Rid: 1})
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//}

func TestService_DeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userAuthRepoMock := authRepoMocks.NewMockUserAuthRepo(ctrl)
	restaurantAuthRepoMock := authRepoMocks.NewMockRestaurantAuthRepo(ctrl)
	sessionRepoMock := authRepoMocks.NewMockSessionRepo(ctrl)
	authService := NewService(userAuthRepoMock, restaurantAuthRepoMock, sessionRepoMock)

	ctx := new(context.Context)

	session := protoAuth.SessionValue{
		Session: "session1",
	}

	sessionRepoMock.EXPECT().Delete(*ctx, headSession+session.Session).Return(nil)

	_, err := authService.DeleteSession(*ctx, &session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}
