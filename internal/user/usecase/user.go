package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/image"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/uniq"
)

// TODO: хранить статику в /var/...

const (
	HeadAvatar = "static/user/avatar/"
)

type userUsecase struct {
	userRepository  user.UserRepo
	imageRepository image.ImageRepo
}

func NewUserUsecase(repo user.UserRepo, image image.ImageRepo) user.UserUsecase {
	return &userUsecase{
		userRepository:  repo,
		imageRepository: image,
	}
}

func (u *userUsecase) Create(ctx context.Context, newUser models.User) (*models.SuccessUserResponse, error) {
	newUser.Avatar = config.DefaultUserImage

	uid, err := u.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.Uid = uid
	newUser.Password = "" // TODO: подумать как это более аккуратно сделать

	response := &models.SuccessUserResponse{
		User: newUser,
		Role: config.RoleUser,
	}

	return response, nil
}

func (u *userUsecase) CheckUserExists(ctx context.Context, user models.UserAuth) (*models.User, error) {
	return u.userRepository.CheckUserExists(ctx, user)
}

func (u *userUsecase) GetByUid(ctx context.Context, uid int) (models.User, error) {
	return u.userRepository.GetByUid(ctx, uid)
}

func (u *userUsecase) UpdateData(ctx context.Context, newUser models.UserData) (*models.SuccessUserResponse, error) {
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	newUser.ID = user.Uid
	err := u.userRepository.UpdateData(ctx, newUser)
	if err != nil {
		return nil, err
	}

	responseUser := models.User{
		Uid:         newUser.ID,
		Name:        newUser.Name,
		Email:       newUser.Email,
		Phone:       newUser.Phone,
		Avatar:      user.Avatar,
		MainAddress: user.MainAddress,
	}

	response := &models.SuccessUserResponse{
		User: responseUser,
		Role: config.RoleUser,
	}

	return response, err
}

func (u *userUsecase) UploadAvatar(ctx context.Context, image *multipart.FileHeader) (*models.UserImageResponse, error) {
	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := uniq.GetUniqFilename(ctx, image.Filename)
	if localErr != nil {
		return nil, localErr
	}

	// удаление изображения
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	if user.Avatar != config.DefaultUserImage {
		removeFile := strings.Replace(user.Avatar, config.Repository, "", -1)
		err := u.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	filename := HeadAvatar + uid + expansion
	err := u.imageRepository.UploadImage(ctx, filename, image)
	if err != nil {
		return nil, err
	}

	filename = config.Repository + HeadAvatar + uid + expansion
	err = u.userRepository.UpdateAvatar(ctx, user.Uid, filename)
	if err != nil {
		return nil, err
	}

	response := &models.UserImageResponse{
		Filename: filename,
	}
	return response, nil
}
