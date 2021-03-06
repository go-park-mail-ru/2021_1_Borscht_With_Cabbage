package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/image"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/borscht/backend/utils/uniq"
)

const (
	HeadAvatar = "static/user/avatar/"
)

type userUsecase struct {
	FolderSaveImages string
	userRepository   user.UserRepo
	imageRepository  image.ImageRepo
}

func NewUserUsecase(repo user.UserRepo, image image.ImageRepo) user.UserUsecase {
	return &userUsecase{
		FolderSaveImages: config.ConfigStatic.Folder,

		userRepository:  repo,
		imageRepository: image,
	}
}

func (a userUsecase) AddAddress(ctx context.Context, uid int, address models.Address) error {
	return a.userRepository.AddAddress(ctx, uid, address)
}
func (a userUsecase) GetAddress(ctx context.Context, rid int) (*models.Address, error) {
	return a.userRepository.GetAddress(ctx, rid)
}

func (a userUsecase) UpdateMainAddress(ctx context.Context, address models.Address) error {
	logger.UsecaseLevel().InlineDebugLog(ctx, "update address")
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return failError
	}

	err := a.userRepository.DeleteAddress(ctx, user.Uid)
	if err != nil {
		return err
	}

	return a.userRepository.AddAddress(ctx, user.Uid, address)
}

func (a userUsecase) GetMainAddress(ctx context.Context) (*models.Address, error) {
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	return a.userRepository.GetAddress(ctx, user.Uid)
}

func (u *userUsecase) GetUserData(ctx context.Context) (*models.SuccessUserResponse, error) {
	user := ctx.Value("User")

	responseUser, ok := user.(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	return &models.SuccessUserResponse{
		User: responseUser,
		Role: configProject.RoleUser,
	}, nil
}

func correctUserData(newUser *models.UserData, oldUser *models.User) {
	newUser.ID = oldUser.Uid
	if newUser.Name == "" {
		newUser.Name = oldUser.Name
	}
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}
	if newUser.Phone == "" {
		newUser.Phone = oldUser.Phone
	}
	if newUser.Address.Name == "" {
		newUser.Address.Name = oldUser.Address.Name
		newUser.Address.Latitude = oldUser.Address.Latitude
		newUser.Address.Longitude = oldUser.Address.Longitude
	}
}

func (u *userUsecase) UpdateData(ctx context.Context, newUser models.UserData) (*models.SuccessUserResponse, error) {
	user, ok := ctx.Value("User").(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	correctUserData(&newUser, &user)
	err := u.userRepository.UpdateData(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = u.userRepository.DeleteAddress(ctx, newUser.ID)
	if err != nil {
		return nil, err
	}

	err = u.userRepository.AddAddress(ctx, newUser.ID, newUser.Address)
	if err != nil {
		return nil, err
	}

	responseUser := models.User{
		Uid:     newUser.ID,
		Name:    newUser.Name,
		Email:   newUser.Email,
		Phone:   newUser.Phone,
		Avatar:  user.Avatar,
		Address: user.Address,
	}

	response := &models.SuccessUserResponse{
		User: responseUser,
		Role: configProject.RoleUser,
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
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, failError
	}

	if user.Avatar != config.ConfigStatic.DefaultUserImage {
		removeFile := u.FolderSaveImages + strings.Replace(user.Avatar, config.ConfigStatic.Repository, "", -1)
		err := u.imageRepository.DeleteImage(ctx, removeFile)
		if err != nil {
			return nil, err
		}
	}

	filename := u.FolderSaveImages + HeadAvatar + uid + expansion
	err := u.imageRepository.UploadImage(ctx, filename, image)
	if err != nil {
		return nil, err
	}

	filename = config.ConfigStatic.Repository + HeadAvatar + uid + expansion
	err = u.userRepository.UpdateAvatar(ctx, user.Uid, filename)
	if err != nil {
		return nil, err
	}

	response := &models.UserImageResponse{
		Filename: filename,
	}
	return response, nil
}
