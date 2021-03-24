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

func (u *userUsecase) Create(newUser domain.User) error {

	// TODO валидация какая нибудь

	return u.userRepository.Create(newUser)
}

func (u *userUsecase) LoginCheck(user domain.UserAuth) (domain.User, error) {
	return u.userRepository.CheckUserExists(user)
}

func (u *userUsecase) GetByUid(uid int32) (domain.User, error) {
	return u.userRepository.GetByUid(uid)
}

func (u *userUsecase) Update(newUser domain.UserData) error {
	// TODO валидация

	return u.userRepository.Update(newUser)
}
