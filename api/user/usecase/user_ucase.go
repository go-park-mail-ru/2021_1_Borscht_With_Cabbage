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

func (u *userUsecase) GetByLogin(ctx *domain.CustomContext, user domain.UserAuth) (domain.User, error) {
	return u.userRepository.GetByLogin(ctx, user)
}

func (u *userUsecase) GetByNumber(ctx *domain.CustomContext, number string) (domain.User, error) {
	return u.userRepository.GetByNumber(ctx, number)
}

func (u *userUsecase) Update(ctx *domain.CustomContext, newUser domain.UserData) error {
	// TODO валидация

	return u.userRepository.Update(ctx, newUser)
}
