package repository

import (
	"backend/api/domain"
	errors "backend/models"
	"net/http"
)

type userRepo struct {
}

func NewUserRepo() domain.UserRepo {
	return &userRepo{}
}

func (u *userRepo) Create(ctx *domain.CustomContext, newUser domain.User) error {
	for _, user := range *ctx.Users {
		if (user.Phone == newUser.Phone) && user.Password == newUser.Password {
			return errors.Create(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
		}
	}

	// записываем нового
	*ctx.Users = append(*ctx.Users, newUser)
	return nil
}

func (u *userRepo) GetByLogin(ctx *domain.CustomContext, check domain.UserAuth) (domain.User, error) {
	for _, user := range *ctx.Users {
		if (user.Email == check.Login || user.Phone == check.Login) && user.Password == check.Password {
			return user, nil
		}
	}

	return domain.User{}, errors.Authorization("not user bd")
}
