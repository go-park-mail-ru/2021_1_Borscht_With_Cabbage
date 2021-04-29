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
	UpdateMainAddress(c echo.Context) error
	GetMainAddress(c echo.Context) error
}

type UserUsecase interface {
	GetUserData(ctx context.Context) (*models.SuccessUserResponse, error)
	UpdateData(ctx context.Context, newUser models.UserData) (*models.SuccessUserResponse, error)
	UploadAvatar(ctx context.Context, image *multipart.FileHeader) (*models.UserImageResponse, error)
	AddAddress(ctx context.Context, uid int, address models.Address) error
	UpdateMainAddress(ctx context.Context, address models.Address) error
	GetMainAddress(ctx context.Context) (*models.Address, error)
}

type UserRepo interface {
	UpdateData(ctx context.Context, newUser models.UserData) error
	UpdateAvatar(ctx context.Context, idUser int, filename string) error
	AddAddress(ctx context.Context, uid int, address models.Address) error
	DeleteAddress(ctx context.Context, uid int) error
	GetAddress(ctx context.Context, uid int) (*models.Address, error)
}
