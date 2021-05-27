package usecase

import (
	"context"
	"testing"

	"github.com/borscht/backend/configProject"
	imageMock "github.com/borscht/backend/internal/image/mocks"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewUserUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	restaurantRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(restaurantRepoMock, imageRepoMock)
	if userUsecase == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestUserUsecase_UpdateMainAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	c := context.Background()

	user := models.User{
		Uid:      1,
		Email:    "dasha@mail.ru",
		Phone:    "89111111111",
		Name:     "111111",
		Password: "1111111",
	}
	ctx := context.WithValue(c, "User", user)

	address := models.Address{
		Name: "address1",
	}

	userRepoMock.EXPECT().DeleteAddress(ctx, user.Uid)
	userRepoMock.EXPECT().AddAddress(ctx, user.Uid, address)

	err := userUsecase.UpdateMainAddress(ctx, address)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestUserUsecase_UpdateMainAddress_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	c := context.Background()

	address := models.Address{
		Name: "address1",
	}

	err := userUsecase.UpdateMainAddress(c, address)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestUserUsecase_GetMainAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	c := context.Background()

	address := models.Address{
		Name: "address1",
	}

	user := models.User{
		Uid:      1,
		Email:    "dasha@mail.ru",
		Phone:    "89111111111",
		Name:     "111111",
		Password: "1111111",
	}
	ctx := context.WithValue(c, "User", user)

	userRepoMock.EXPECT().GetAddress(ctx, user.Uid).Return(&address, nil)

	addressResult, err := userUsecase.GetMainAddress(ctx)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, addressResult.Name, "address1")
}

func TestUserUsecase_GetMainAddress_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	c := context.Background()

	_, err := userUsecase.GetMainAddress(c)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestUserUsecase_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	c := context.Background()

	user := models.User{
		Email:    "dasha@mail.ru",
		Phone:    "89111111111",
		Name:     "111111",
		Password: "1111111",
	}
	ctx := context.WithValue(c, "User", user)

	response := new(models.SuccessUserResponse)
	var err error
	response, err = userUsecase.GetUserData(ctx)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, configProject.RoleUser)
}

func TestUserUsecase_UpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	c := context.Background()
	userId := models.User{
		Uid: 1,
	}
	ctx := context.WithValue(c, "User", userId)

	user := models.UserData{
		Email:    "dasha@mail.ru",
		Phone:    "89111111111",
		Name:     "Daria",
		Password: "1111111",
		Address: models.Address{
			Name:      "yu2",
			Longitude: 213,
			Latitude:  2142,
			Radius:    1000,
		},
	}
	userWithID := user
	userWithID.ID = 1

	userRepoMock.EXPECT().UpdateData(ctx, userWithID).Return(nil)
	userRepoMock.EXPECT().DeleteAddress(ctx, userId.Uid).Return(nil)
	userRepoMock.EXPECT().AddAddress(ctx, userId.Uid, user.Address).Return(nil)

	response := new(models.SuccessUserResponse)
	var err error
	response, err = userUsecase.UpdateData(ctx, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, configProject.RoleUser)
}

func TestUserUsecase_UploadAvatar(t *testing.T) {
	// TODO
}
