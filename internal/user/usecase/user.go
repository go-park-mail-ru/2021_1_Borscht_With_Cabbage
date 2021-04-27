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
	"github.com/borscht/backend/utils/secure"
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

func (u *userUsecase) Create(ctx context.Context, newUser models.User) (*models.SuccessUserResponse, error) {
	newUser.Avatar = config.DefaultUserImage

	newUser.HashPassword = secure.HashPassword(ctx, secure.GetSalt(), newUser.Password)

	uid, err := u.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.Uid = uid
	newUser.Password = "" // TODO: подумать как это более аккуратно сделать
	newUser.HashPassword = nil

	err = u.userRepository.AddAddress(ctx, uid, newUser.Address)
	if err != nil {
		// TODO: подумать над тем чтобы удалять пользователя
		return nil, err
	}

	response := &models.SuccessUserResponse{
		User: newUser,
		Role: config.RoleUser,
	}

	return response, nil
}

func (u *userUsecase) CheckUserExists(ctx context.Context, userAuth models.UserAuth) (*models.SuccessUserResponse, error) {
	user, err := u.userRepository.GetByLogin(ctx, userAuth.Login)
	if err != nil {
		return nil, err
	}
	logger.UsecaseLevel().InlineDebugLog(ctx, &user.HashPassword)

	if !secure.CheckPassword(ctx, user.HashPassword, userAuth.Password) {
		err = errors.AuthorizationError("bad password")
		logger.UsecaseLevel().ErrorLog(ctx, err)
		return nil, err
	}
	user.HashPassword = nil
	return &models.SuccessUserResponse{
		User: *user,
		Role: config.RoleUser,
	}, nil
}

func (u *userUsecase) GetByUid(ctx context.Context, uid int) (*models.SuccessUserResponse, error) {
	user, err := u.userRepository.GetByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	address, err := u.userRepository.GetAddress(ctx, uid)
	if err != nil {
		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"address error": err})
		return nil, err
	}
	if address != nil {
		user.Address = *address
	}

	return &models.SuccessUserResponse{
		User: user,
		Role: config.RoleUser,
	}, nil
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
		Role: config.RoleUser,
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
