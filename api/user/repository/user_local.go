package repository

import (
	"backend/api/domain"
	"backend/api/domain/user"
	errors "backend/utils"
	"fmt"
	"net/http"
)

type userRepo struct {
}

func NewUserRepo() user.UserRepo {
	return &userRepo{}
}

func (u *userRepo) Create(ctx *domain.CustomContext, newUser user.User) error {
	for _, user := range *ctx.Users {
		if (user.Phone == newUser.Phone) && user.Password == newUser.Password {
			return errors.Create(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
		}
	}

	// записываем нового
	*ctx.Users = append(*ctx.Users, newUser)
	return nil
}

func (u *userRepo) GetByLogin(ctx *domain.CustomContext, check user.UserAuth) (user.User, error) {
	for _, user := range *ctx.Users {
		if (user.Email == check.Login || user.Phone == check.Login) && user.Password == check.Password {
			return user, nil
		}
	}

	return user.User{}, errors.Authorization("not user bd")
}

func (u *userRepo) GetByNumber(ctx *domain.CustomContext, number string) (user.User, error) {
	for _, user := range *ctx.Users {
		if user.Phone == number {
			return user, nil
		}
	}

	return user.User{}, errors.Authorization("user not found")
}

func (u *userRepo) Update(ctx *domain.CustomContext, newUser user.UserData) error {
	for i, user := range *ctx.Users {
		if user.Email == newUser.Email && user.Phone != ctx.User.Phone { // если у кого-то другого уже есть такой email
			return errors.Create(http.StatusBadRequest, "user with this email already exists")
		}
		if user.Phone == newUser.Phone && user.Phone != ctx.User.Phone { // если у кого-то другого уже есть такой телефон
			return errors.Create(http.StatusBadRequest, "user with this number already exists")
		}

		if user.Phone == ctx.User.Phone {
			if newUser.Password != "" {
				if newUser.PasswordOld != user.Password {
					fmt.Println(newUser.PasswordOld, " ", user.Password)
					return errors.Create(http.StatusBadRequest, "invalid old password")
				}
				(*ctx.Users)[i].Password = newUser.Password
			}

			(*ctx.Users)[i].Phone = newUser.Phone
			(*ctx.Users)[i].Email = newUser.Email
			(*ctx.Users)[i].Name = newUser.Name

			fmt.Println(*ctx.Users)
			//return cc.SendOK(profileEdits)
			return nil
		}
	}

	return errors.Authorization("user not found")
}