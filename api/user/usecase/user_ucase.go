package usecase

import (
	"backend/api/domain"
)

type userUsecase struct {
	userRepository domain.UserRepo
}

func NewUserUsecase(repo domain.UserRepo) domain.UserUsecase {
	return &userUsecase{
		userRepository: repo,
	}
}

func (u *userUsecase) Create(ctx *domain.CustomContext, newUser domain.User) error {

	// TODO валидация какая нибудь

	return u.userRepository.Create(ctx, newUser)
}
