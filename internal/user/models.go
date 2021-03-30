package user

import (
	"github.com/borscht/backend/internal/models"
	"github.com/labstack/echo/v4"
	"mime/multipart"
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
	Create(user models.User) (int, error)
	CheckUserExists(user models.UserAuth) (models.User, error)
	GetByUid(uid int) (models.User, error)
	Update(newUser models.UserData, uid int) error
	UploadAvatar(image *multipart.FileHeader) (string, error)
}

type UserRepo interface {
	Create(user models.User) (int, error)
	CheckUserExists(user models.UserAuth) (models.User, error)
	GetByUid(uid int) (models.User, error)
	Update(newUser models.UserData, uid int) error
	UploadAvatar(image *multipart.FileHeader, filename string) error
}
