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
	EditProfile(c echo.Context) error
	CheckAuth(c echo.Context) error
	Logout(c echo.Context) error
}

type UserUsecase interface {
	Create(ctx context.Context, user models.User) error
	GetByLogin(ctx context.Context, user models.UserAuth) (models.User, error)
	GetByNumber(ctx context.Context, number string) (models.User, error)
	Update(ctx context.Context, newUser models.UserData) error
	UploadAvatar(ctx context.Context, image *multipart.FileHeader) (string, error)
}

type UserRepo interface {
	Create(ctx context.Context, user models.User) error
	GetByLogin(ctx context.Context, user models.UserAuth) (models.User, error)
	GetByNumber(ctx context.Context, number string) (models.User, error)
	Update(ctx context.Context, newUser models.UserData) error
	UploadAvatar(image *multipart.FileHeader, filename string) error
}
