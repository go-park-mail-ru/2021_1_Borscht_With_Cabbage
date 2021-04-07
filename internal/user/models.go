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
	Create(ctx context.Context, user models.User) (*models.User, error)
	CheckUserExists(ctx context.Context, user models.UserAuth) (*models.User, error)
	GetByUid(ctx context.Context, uid int) (models.User, error)
	Update(ctx context.Context, newUser models.UserData, uid int) error
	UploadAvatar(ctx context.Context, image *multipart.FileHeader) (string, error)
}

type UserRepo interface {
	Create(ctx context.Context, user models.User) (int, error)
	CheckUserExists(ctx context.Context, user models.UserAuth) (*models.User, error)
	GetByUid(ctx context.Context, uid int) (models.User, error)
	Update(ctx context.Context, newUser models.UserData, uid int) error
}
