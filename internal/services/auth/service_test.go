package auth

import (
	"context"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/auth/mocks"
	proto "github.com/borscht/backend/services/proto/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_CheckKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	session := proto.SessionValue{
		Session: "session1",
	}
	sessionResult := proto.SessionInfo{
		Id:   1,
		Role: config.RoleUser,
	}

	clientMock.EXPECT().CheckKey(c, &session).Return(&sessionResult, nil)

	output, check, err := authService.CheckKey(c, session.Session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, output.Id, 1)
	require.EqualValues(t, check, true)
}

func TestService_CreateKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	session := models.SessionInfo{
		Id:   1,
		Role: config.RoleUser,
	}
	sessionInfo := proto.SessionInfo{
		Id:   int32(session.Id),
		Role: session.Role,
	}
	sessionVal := proto.SessionValue{
		Session: "session1",
	}

	clientMock.EXPECT().CreateKey(c, &sessionInfo).Return(&sessionVal, nil)

	key, err := authService.CreateKey(c, session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, key, "session1")
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	userToCreate := models.User{
		Name:     "Daria",
		Email:    "d@mail.ru",
		Password: "111111",
		Phone:    "89111111111",
	}
	userRequest := proto.User{
		Email:    userToCreate.Email,
		Phone:    userToCreate.Phone,
		Password: userToCreate.Password,
		Name:     userToCreate.Name,
	}

	clientMock.EXPECT().CreateUser(c, &userRequest).Return(&proto.SuccessUserResponse{}, nil)

	user, err := authService.Create(c, userToCreate)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, user.Role, config.RoleUser)
}

func TestService_CheckUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	userToCheck := models.UserAuth{
		Password: "111111",
		Login:    "89111111111",
	}
	userRequest := proto.UserAuth{
		Login:    userToCheck.Login,
		Password: userToCheck.Password,
	}

	clientMock.EXPECT().CheckUserExists(c, &userRequest).Return(&proto.SuccessUserResponse{}, nil)

	user, err := authService.CheckUserExists(c, userToCheck)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, user.Role, config.RoleUser)
}

func TestService_GetByUid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	uid := 1
	userRequest := proto.UID{
		Uid: int32(uid),
	}

	clientMock.EXPECT().GetByUid(c, &userRequest).Return(&proto.SuccessUserResponse{}, nil)

	user, err := authService.GetByUid(c, uid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, user.Role, config.RoleUser)
}

func TestService_CreateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	restaurantToCreate := models.RestaurantInfo{
		Title:         "rest1",
		AdminEmail:    "d@mail.ru",
		AdminPhone:    "89111111111",
		AdminPassword: "111111",
	}
	restRequest := proto.User{
		Name:     restaurantToCreate.Title,
		Email:    restaurantToCreate.AdminEmail,
		Phone:    restaurantToCreate.AdminPhone,
		Password: restaurantToCreate.AdminPassword,
	}

	clientMock.EXPECT().CreateRestaurant(c, &restRequest).Return(&proto.SuccessRestaurantResponse{}, nil)

	restaurant, err := authService.CreateRestaurant(c, restaurantToCreate)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurant.Role, config.RoleAdmin)
}

func TestService_CheckRestaurantExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	restaurantToCheck := models.RestaurantAuth{
		Login:    "89111111111",
		Password: "111111",
	}
	restRequest := proto.UserAuth{
		Login:    restaurantToCheck.Login,
		Password: restaurantToCheck.Password,
	}

	clientMock.EXPECT().CheckRestaurantExists(c, &restRequest).Return(&proto.SuccessRestaurantResponse{}, nil)

	restaurant, err := authService.CheckRestaurantExists(c, restaurantToCheck)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurant.Role, config.RoleAdmin)
}

func TestService_GetByRid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	rid := 1
	restRequest := proto.RID{
		Rid: int32(rid),
	}

	clientMock.EXPECT().GetByRid(c, &restRequest).Return(&proto.SuccessRestaurantResponse{}, nil)

	restaurant, err := authService.GetByRid(c, rid)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, restaurant.Role, config.RoleAdmin)
}

func TestService_CheckSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	session := "session1"
	sessionRequest := proto.SessionValue{
		Session: session,
	}

	clientMock.EXPECT().CheckSession(c, &sessionRequest).Return(&proto.SessionInfo{}, nil)

	_, check, err := authService.CheckSession(c, session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, check, true)
}

func TestService_CreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	session := models.SessionInfo{
		Id:   1,
		Role: config.RoleUser,
	}
	sessionRequest := proto.SessionInfo{
		Id:   int32(session.Id),
		Role: session.Role,
	}
	sessionResult := proto.SessionValue{
		Session: "session1",
	}
	clientMock.EXPECT().CreateSession(c, &sessionRequest).Return(&sessionResult, nil)

	sessionOutput, err := authService.CreateSession(c, session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, sessionOutput, "session1")
}

func TestService_DeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockAuthClient(ctrl)
	authService := NewService(clientMock)

	c := context.Background()

	session := "session1"
	sessionRequest := proto.SessionValue{
		Session: session,
	}

	clientMock.EXPECT().DeleteSession(c, &sessionRequest).Return(&proto.Error{}, nil)

	err := authService.DeleteSession(c, session)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}
