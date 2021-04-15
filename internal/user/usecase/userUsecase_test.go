package usecase

import (
	"context"
	"github.com/borscht/backend/config"
	imageMock "github.com/borscht/backend/internal/image/mocks"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user/mocks"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/secure"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
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

func TestUserUsecase_Create(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//userRepoMock := mocks.NewMockUserRepo(ctrl)
	//imageRepoMock := imageMock.NewMockImageRepo(ctrl)
	//
	//userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
	//c := context.Background()
	//ctx := context.WithValue(c, "request_id", 1)
	//
	//logger.InitLogger()
	//
	//user := models.User{
	//	Email:    "dasha@mail.ru",
	//	Phone: "89111111111",
	//	Name: "111111",
	//	Password: "1111111",
	//}
	//response := models.SuccessUserResponse{
	//	user, config.RoleUser,
	//}
	//
	//userRepoMock.EXPECT().Create(ctx, user).Return(&response, nil)
	//
	//_, err := userUsecase.Create(ctx, user)
	//if err != nil {
	//	t.Errorf("unexpected err: %s", err)
	//	return
	//}
}

func TestUserUsecase_CheckUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
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
		HashPassword: secure.HashPassword(ctx, secure.GetSalt(), "111111"),
	}

	userRepoMock.EXPECT().GetByLogin(ctx, userAuth.Login).Return(&user, nil)

	response := new(models.SuccessUserResponse)
	var err error
	response, err = userUsecase.CheckUserExists(ctx, userAuth)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, config.RoleUser)
}

func TestUserUsecase_GetByUid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mocks.NewMockUserRepo(ctrl)
	imageRepoMock := imageMock.NewMockImageRepo(ctrl)

	userUsecase := NewUserUsecase(userRepoMock, imageRepoMock)
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

	userRepoMock.EXPECT().GetByUid(ctx, 1).Return(user, nil)

	response := new(models.SuccessUserResponse)
	var err error
	response, err = userUsecase.GetByUid(ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, config.RoleUser)
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

	require.EqualValues(t, response.Role, config.RoleUser)
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
	}
	userWithID := user
	userWithID.ID = 1

	userRepoMock.EXPECT().UpdateData(ctx, userWithID).Return(nil)

	response := new(models.SuccessUserResponse)
	var err error
	response, err = userUsecase.UpdateData(ctx, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, response.Role, config.RoleUser)
}

func TestUserUsecase_UploadAvatar(t *testing.T) {
	// TODO
}
