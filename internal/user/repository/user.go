package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/borscht/backend/database/local"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	errors "github.com/borscht/backend/utils"
)

type userRepo struct {
	db local.Database
}

func NewUserRepo() user.UserRepo {
	return &userRepo{
		db: local.GetInstance(),
	}
}

func (u *userRepo) Create(ctx context.Context, newUser models.User) error {
	for _, curUser := range *u.db.GetModels().Users {
		if (curUser.Phone == newUser.Phone) && curUser.Password == newUser.Password {
			return errors.Authorization("User with this number already exists") // такой юзер уже есть
		}
	}

	// записываем нового
	*u.db.GetModels().Users = append(*u.db.GetModels().Users, newUser)
	return nil
}

func (u *userRepo) GetByLogin(ctx context.Context, check models.UserAuth) (models.User, error) {
	for _, curUser := range *u.db.GetModels().Users {
		if (curUser.Email == check.Login || curUser.Phone == check.Login) && curUser.Password == check.Password {
			return curUser, nil
		}
	}

	return models.User{}, errors.Authorization("not curUser bd")
}

func (u *userRepo) GetByNumber(ctx context.Context, number string) (models.User, error) {
	for _, curUser := range *u.db.GetModels().Users {
		if curUser.Phone == number {
			return curUser, nil
		}
	}

	return models.User{}, errors.Authorization("curUser not found")
}

func (u *userRepo) Update(ctx context.Context, newUser models.UserData) error {
	for i, curUser := range *u.db.GetModels().Users {
		if curUser.Email == newUser.Email && curUser.Phone != ctx.Value("User").(models.User).Phone { // если у кого-то другого уже есть такой email
			return errors.BadRequest("curUser with this email already exists")
		}
		if curUser.Phone == newUser.Phone && curUser.Phone != ctx.Value("User").(models.User).Phone { // если у кого-то другого уже есть такой телефон
			return errors.Authorization("User with this number already exists")
		}

		if curUser.Phone == ctx.Value("User").(models.User).Phone {
			if newUser.Password != "" {
				if newUser.PasswordOld != curUser.Password {
					fmt.Println(newUser.PasswordOld, " ", curUser.Password)
					return errors.Authorization("invalid old password")
				}
				(*u.db.GetModels().Users)[i].Password = newUser.Password
			}

			(*u.db.GetModels().Users)[i].Phone = newUser.Phone
			(*u.db.GetModels().Users)[i].Email = newUser.Email
			(*u.db.GetModels().Users)[i].Name = newUser.Name

			fmt.Println(*u.db.GetModels().Users)
			//return cc.SendResponse(profileEdits)
			return nil
		}
	}

	return errors.Authorization("curUser not found")
}

func (u *userRepo) UploadAvatar(image *multipart.FileHeader, filename string) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		return errors.FailServer(err.Error())
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return errors.FailServer(err.Error())
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return errors.FailServer(err.Error())
	}

	return nil
}
