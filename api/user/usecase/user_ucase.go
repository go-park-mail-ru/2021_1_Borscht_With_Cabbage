package usecase

import (
	"backend/api/domain"
	"backend/api/domain/user"
)

type userUsecase struct {
	userRepository user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) user.UserUsecase {
	return &userUsecase{
		userRepository: repo,
	}
}

func (u *userUsecase) Create(ctx *domain.CustomContext, newUser user.User) error {

	// TODO валидация какая нибудь

	return u.userRepository.Create(ctx, newUser)
}

func (u *userUsecase) GetByLogin(ctx *domain.CustomContext, user user.UserAuth) (user.User, error) {
	return u.userRepository.GetByLogin(ctx, user)
}

func (u *userUsecase) GetByNumber(ctx *domain.CustomContext, number string) (user.User, error) {
	return u.userRepository.GetByNumber(ctx, number)
}

func (u *userUsecase) Update(ctx *domain.CustomContext, newUser user.UserData) error {
	// TODO валидация

	return u.userRepository.Update(ctx, newUser)
}
