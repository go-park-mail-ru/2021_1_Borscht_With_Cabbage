package repository

import (
	"backend/api/domain"
	"net/http"
)

type userRepo struct {
}

func NewUserRepo(repo domain.UserRepo) domain.UserRepo {
	return &userRepo{}
}

func (u *userRepo) Create(ctx *domain.CustomContext, newUser domain.User) error {
	for _, user := range *ctx.Users {
		if (user.Phone == newUser.Phone) && user.Password == newUser.Password {
			return ctx.String(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
		}
	}

	// записываем нового
	*ctx.Users = append(*ctx.Users, newUser)
	return nil
}
