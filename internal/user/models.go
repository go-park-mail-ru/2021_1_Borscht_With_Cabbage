package user

import (
	"context"
	"mime/multipart"

	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	GetUserData(c echo.Context) error
	UpdateData(c echo.Context) error
	UploadAvatar(c echo.Context) error
	CheckAuth(c echo.Context) error
	Logout(c echo.Context) error
}

type UserUsecase interface {
	Create(ctx context.Context, user models.User) (*models.SuccessUserResponse, error)
	CheckUserExists(ctx context.Context, user models.UserAuth) (*models.SuccessUserResponse, error)
	GetByUid(ctx context.Context, uid int) (*models.SuccessUserResponse, error)
	GetUserData(ctx context.Context) (*models.SuccessUserResponse, error)
	UpdateData(ctx context.Context, newUser models.UserData) (*models.SuccessUserResponse, error)
	UploadAvatar(ctx context.Context, image *multipart.FileHeader) (*models.UserImageResponse, error)
}

type UserRepo interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetByUid(ctx context.Context, uid int) (models.User, error)
	UpdateData(ctx context.Context, newUser models.UserData) error
	UpdateAvatar(ctx context.Context, idUser int, filename string) error
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateMainAddress(ctx context.Context, uid int, address string) error
}
