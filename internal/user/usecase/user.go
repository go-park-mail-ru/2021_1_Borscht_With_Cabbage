package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/uniq"
)

// TODO: хранить статику в /var/...

const (
	HeadAvatar = "static/avatar/"
)

type userUsecase struct {
	userRepository user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) user.UserUsecase {
	return &userUsecase{
		userRepository: repo,
	}
}

func (u *userUsecase) Create(ctx context.Context, newUser models.User) (*models.User, error) {

	// TODO валидация какая нибудь
	newUser.Avatar = config.DefaultAvatar

	uid, err := u.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.Uid = uid
	return &newUser, nil
}

func (u *userUsecase) CheckUserExists(ctx context.Context, user models.UserAuth) (*models.User, error) {
	return u.userRepository.CheckUserExists(ctx, user)
}

func (u *userUsecase) GetByUid(ctx context.Context, uid int) (models.User, error) {
	return u.userRepository.GetByUid(ctx, uid)
}

func (u *userUsecase) Update(ctx context.Context, newUser models.UserData, uid int) error {
	// TODO валидация

	return u.userRepository.Update(ctx, newUser, uid)
}

func (u *userUsecase) UploadAvatar(ctx context.Context, image *multipart.FileHeader) (string, error) {
	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := uniq.GetUniqFilename(ctx, image.Filename)
	if localErr != nil {
		return "", localErr
	}

	filename := HeadAvatar + uid + expansion
	if err := u.userRepository.UploadAvatar(ctx, image, filename); err != nil {
		return "", err
	}

	return config.Repository + filename, nil
}
